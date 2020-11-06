/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package libext

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"runtime"
	"strings"
	"time"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/pipenet"
	"github.com/pion/dtls/v2"
	"golang.org/x/sync/errgroup"

	"arhat.dev/libext/codec"
	"arhat.dev/libext/types"
	"arhat.dev/libext/util"
)

type (
	connectFunc func() (net.Conn, error)
)

func NewClient(
	ctx context.Context,
	kind arhatgopb.ExtensionType,
	name string,
	c types.Codec,

	// connection management
	dialer *net.Dialer,
	endpointURL string,
	tlsConfig *tls.Config,
) (*Client, error) {
	u, err := url.Parse(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint url: %w", err)
	}

	jsonCodec := codec.GetCodec(arhatgopb.CODEC_JSON)

	regMsg, err := util.NewMsg(jsonCodec.Marshal, arhatgopb.MSG_REGISTER, 0, 0, &arhatgopb.RegisterMsg{
		Name:          name,
		ExtensionType: kind,
		Codec:         c.Type(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create register message: %w", err)
	}

	regMsgBuf := new(bytes.Buffer)
	err = jsonCodec.NewEncoder(regMsgBuf).Encode(regMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal register message: %w", err)
	}

	if dialer == nil {
		dialer = &net.Dialer{
			Timeout:       0,
			Deadline:      time.Time{},
			LocalAddr:     nil,
			FallbackDelay: 0,
			KeepAlive:     0,
			Resolver:      nil,
			Control:       nil,
		}
	}

	var (
		connector connectFunc
	)
	switch s := strings.ToLower(u.Scheme); s {
	case "tcp", "tcp4", "tcp6": // nolint:goconst
		_, err = net.ResolveTCPAddr(s, u.Host)
		connector = func() (net.Conn, error) {
			return dialer.DialContext(ctx, s, u.Host)
		}
	case "udp", "udp4", "udp6": // nolint:goconst
		_, err = net.ResolveUDPAddr(s, u.Host)
		connector = func() (net.Conn, error) {
			return dialer.DialContext(ctx, s, u.Host)
		}
	case "unix": // nolint:goconst
		_, err = net.ResolveUnixAddr(s, u.Path)
		connector = func() (net.Conn, error) {
			return dialer.DialContext(ctx, s, u.Path)
		}
	case "pipe":
		connector = func() (net.Conn, error) {
			if runtime.GOOS == "windows" {
				host := u.Host
				path := u.Path
				if path == "" {
					host = "."
					path = u.Host
				}

				return pipenet.DialContext(ctx, fmt.Sprintf(`\\%s\pipe\%s`, host, path))
			}

			return pipenet.DialContext(ctx, u.Path)
		}
	default:
		return nil, fmt.Errorf("unsupported endpoint protocol %q", u.Scheme)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to resolve %s address: %w", u.Scheme, err)
	}

	return &Client{
		ctx: ctx,

		codec:  c,
		regMsg: regMsgBuf.Bytes(),

		createConnection: func() (_ net.Conn, err error) {
			var innerConn net.Conn

			innerConn, err = connector()
			if err != nil {
				return nil, err
			}
			if tlsConfig == nil {
				return innerConn, nil
			}

			defer func() {
				if err != nil {
					_ = innerConn.Close()
				}
			}()

			// tls enabled
			_, isUDPConn := innerConn.(*net.UDPConn)
			if isUDPConn {
				// currently only udp supports dtls
				dtlsConfig := util.ConvertTLSConfigToDTLSConfig(tlsConfig)
				dtlsConfig.ConnectContextMaker = func() (context.Context, func()) {
					return context.WithCancel(ctx)
				}
				return dtls.ClientWithContext(ctx, innerConn, dtlsConfig)
			}

			return tls.Client(innerConn, tlsConfig), nil
		},
	}, nil
}

type Client struct {
	ctx context.Context

	codec  types.Codec
	regMsg []byte

	createConnection connectFunc
}

// ProcessNewStream creates a new connection and handles message stream until connection lost
// or msgCh closed
// the provided `cmdCh` and `msgCh` are expected to be freshly created
// usually this function is used in conjunction with Controller.RefreshChannels
func (c *Client) ProcessNewStream(
	cmdCh chan<- *arhatgopb.Cmd,
	msgCh <-chan *arhatgopb.Msg,
) error {
	conn, err := c.createConnection()
	if err != nil {
		return fmt.Errorf("failed to dial endpoint: %w", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	_, err = conn.Write(c.regMsg)
	if err != nil {
		return fmt.Errorf("failed to register myself: %w", err)
	}

	wg, ctx := errgroup.WithContext(c.ctx)

	keepaliveCh := make(chan struct{})
	wg.Go(func() error {
		defer func() {
			_ = conn.Close()
		}()

		enc := c.codec.NewEncoder(conn)
		for {
			select {
			case _, more := <-keepaliveCh:
				if !more {
					return io.EOF
				}
				err2 := enc.Encode(&arhatgopb.Msg{
					Kind: arhatgopb.MSG_PONG,
				})
				if err2 != nil {
					return fmt.Errorf("failed to encode pong message: %w", err2)
				}
			case msg, more := <-msgCh:
				if !more {
					return io.EOF
				}
				err2 := enc.Encode(msg)
				if err2 != nil {
					return fmt.Errorf("failed to encode message: %w", err2)
				}
			case <-ctx.Done():
				return io.EOF
			}
		}
	})

	wg.Go(func() error {
		defer func() {
			close(keepaliveCh)
			close(cmdCh)
		}()

		dec := c.codec.NewDecoder(conn)
		for {
			cmd := new(arhatgopb.Cmd)
			err2 := checkNetworkReadErr(dec.Decode(cmd))
			if err2 != nil {
				return err2
			}

			if cmd.Kind == arhatgopb.CMD_PING {
				select {
				case keepaliveCh <- struct{}{}:
				case <-ctx.Done():
					return io.EOF
				}
				continue
			}

			select {
			case cmdCh <- cmd:
			case <-ctx.Done():
				return io.EOF
			}
		}
	})

	err = wg.Wait()

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func checkNetworkReadErr(err error) error {
	if err == nil {
		return nil
	}

	switch t := err.(type) {
	case *net.OpError:
		if t.Err.Error() == "use of closed network connection" {
			return io.EOF
		}
	default:
		if strings.Contains(err.Error(), "closed") {
			return io.EOF
		}

		return t
	}

	return err
}
