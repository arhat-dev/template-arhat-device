package libext

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/iohelper"
	"golang.org/x/sync/errgroup"
)

type ExtensionType string

const (
	ExtensionPeripheral ExtensionType = "/peripherals"
)

func NewClient(
	ctx context.Context,
	endpointURL string,
	tlsConfig *tls.Config,
	kind ExtensionType,
	codec Codec,
) (*Client, error) {
	u, err := url.Parse(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint url: %w", err)
	}

	switch strings.ToLower(u.Scheme) {
	case "tcp", "tcp4", "tcp6": // nolint:goconst
	case "unix":
	case "http":
		tlsConfig = nil
		u.Scheme = "tcp"
	case "https":
		if tlsConfig == nil {
			return nil, fmt.Errorf("no tls config provided for https endpoint")
		}
		u.Scheme = "tcp"
	default:
		return nil, fmt.Errorf("unsupported endpoint scheme %s", u.Scheme)
	}

	return &Client{
		ctx:     ctx,
		network: u.Scheme,
		addr:    u.Host + u.Path,

		tlsConfig: tlsConfig,
		codec:     codec,
		endpoint:  kind,
	}, nil
}

type Client struct {
	ctx     context.Context
	network string
	addr    string

	tlsConfig *tls.Config
	codec     Codec
	endpoint  ExtensionType
}

func (c *Client) ProcessNewStream(
	cmdCh chan<- *arhatgopb.Cmd,
	msgCh <-chan *arhatgopb.Msg,
) error {
	client, cleanup, err := c.createHTTPClient()
	if err != nil {
		return fmt.Errorf("failed to dial endpoint: %w", err)
	}

	defer cleanup()

	pr, pw := iohelper.Pipe()

	req, err := http.NewRequest(http.MethodPost, "", pr)
	if err != nil {
		return fmt.Errorf("failed to create post request")
	}

	req.Host = c.addr
	req.URL.Path = string(c.endpoint)
	req.URL.Host = c.addr
	req.URL.Scheme = "http"
	if c.tlsConfig != nil {
		req.URL.Scheme = "https"
	}

	req.Header.Set("Content-Type", c.codec.ContentType())

	wg, ctx := errgroup.WithContext(c.ctx)

	wg.Go(func() error {
		enc := c.codec.NewEncoder(pw)

		defer func() {
			_ = pw.Close()
		}()

		for msg := range msgCh {
			err2 := enc.Encode(msg)
			if err2 != nil {
				return fmt.Errorf("failed to marshal and send msg: %w", err2)
			}
		}

		return io.EOF
	})

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start sync loop for the first time: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	wg.Go(func() error {
		defer func() {
			close(cmdCh)
		}()

		dec := c.codec.NewDecoder(resp.Body)
		for {
			cmd := new(arhatgopb.Cmd)
			err2 := dec.Decode(cmd)
			if err2 != nil {
				if err2 != io.EOF {
					return fmt.Errorf("failed to decode cmd: %w", err2)
				}
				return nil
			}

			select {
			case cmdCh <- cmd:
			case <-ctx.Done():
				return nil
			}
		}
	})

	err = wg.Wait()
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (c *Client) createHTTPClient() (_ *http.Client, cleanup func(), _ error) {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(c.ctx, c.network, c.addr)
	if err != nil {
		return nil, nil, err
	}

	return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return conn, err
				},
				TLSClientConfig:   c.tlsConfig,
				ForceAttemptHTTP2: true,
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}, func() {
			_ = conn.Close()
		}, nil
}
