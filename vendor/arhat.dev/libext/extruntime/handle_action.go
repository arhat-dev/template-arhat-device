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

package extruntime

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/runtimepb"
	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/iohelper"

	"arhat.dev/libext/types"
)

func (h *Handler) handleExec(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(aranyagopb.ExecOrAttachCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	h.handleStreamOperation(ctx, sid, opts.Stdin, opts.Stdout, opts.Stderr, opts.Tty,
		// preRun check
		nil,
		// run
		func(stdout, stderr io.WriteCloser) *aranyagopb.ErrorMsg {
			var errCh <-chan *aranyagopb.ErrorMsg
			if opts.Stdin {
				err = h.streams.Add(sid, func() (_ io.WriteCloser, resizeFunc types.ResizeHandleFunc, err error) {
					pr, pw := iohelper.Pipe()
					resizeFunc, errCh, err = h.impl.Exec(
						ctx, opts.PodUid, opts.Container, pr, stdout, stderr, opts.Command, opts.Tty,
					)
					if err != nil {
						_ = pw.Close()
						_ = pr.Close()
						return nil, nil, err
					}

					return pw, resizeFunc, err
				})
			} else {
				_, errCh, err = h.impl.Exec(
					ctx, opts.PodUid, opts.Container, nil, stdout, stderr, opts.Command, opts.Tty,
				)
			}
			if err != nil {
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: err.Error(),
					Code:        0,
				}
			}

			return collectStreamErrors(ctx, errCh)
		},
	)

	return nil, nil
}

func (h *Handler) handleAttach(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(aranyagopb.ExecOrAttachCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	h.handleStreamOperation(ctx, sid, opts.Stdin, opts.Stdout, opts.Stderr, opts.Tty,
		// preRun check
		nil,
		// run
		func(stdout, stderr io.WriteCloser) *aranyagopb.ErrorMsg {
			var errCh <-chan *aranyagopb.ErrorMsg

			if opts.Stdin {
				err = h.streams.Add(sid, func() (_ io.WriteCloser, resizeFunc types.ResizeHandleFunc, err error) {
					stdin, pw := iohelper.Pipe()
					resizeFunc, errCh, err = h.impl.Attach(
						ctx, opts.PodUid, opts.Container, stdin, stdout, stderr,
					)
					if err != nil {
						_ = pw.Close()
						_ = stdin.Close()
						return nil, nil, err
					}

					return pw, resizeFunc, err
				})
			} else {
				_, errCh, err = h.impl.Attach(
					ctx, opts.PodUid, opts.Container, nil, stdout, stderr,
				)
			}
			if err != nil {
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: err.Error(),
					Code:        0,
				}
			}

			return collectStreamErrors(ctx, errCh)
		},
	)

	return nil, nil
}

func (h *Handler) handleLogs(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(aranyagopb.LogsCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	h.handleStreamOperation(ctx, sid, false, true, true, false,
		// preRun check
		nil,
		// run
		func(stdout, stderr io.WriteCloser) *aranyagopb.ErrorMsg {
			err := h.impl.Logs(ctx, opts, stdout, stderr)
			if err != nil {
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: err.Error(),
					Code:        0,
				}
			}

			return nil
		},
	)

	return nil, nil
}

func (h *Handler) handleTtyResize(_ context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(aranyagopb.TerminalResizeCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	h.streams.Resize(sid, opts.Cols, opts.Rows)

	return nil, nil
}

type flexWriteCloser struct {
	writeFunc func([]byte) (int, error)
	closeFunc func() error
}

func (a *flexWriteCloser) Write(p []byte) (n int, err error) {
	return a.writeFunc(p)
}

func (a *flexWriteCloser) Close() error {
	return a.closeFunc()
}

func (h *Handler) handlePortForward(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(aranyagopb.PortForwardCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	var (
		seq uint64

		pr, pw = iohelper.Pipe()
	)

	defer func() {
		_ = pw.Close()

		lastMsg := &arhatgopb.Msg{
			Kind:    arhatgopb.MSG_DATA_OUTPUT,
			Id:      sid,
			Ack:     nextSeq(&seq),
			Payload: nil,
		}

		// send fin msg to close input in aranya
		if err != nil {
			lastMsg.Kind = arhatgopb.MSG_ERROR
			lastMsg.Payload, _ = (&aranyagopb.ErrorMsg{
				Kind:        aranyagopb.ERR_COMMON,
				Description: err.Error(),
				Code:        0,
			}).Marshal()
		}

		// best effort
		_ = h.SendMsg(lastMsg)

		// close this session locally (no more input data should be delivered to this session)
		h.streams.Del(sid)
	}()

	var (
		downstream io.ReadCloser
		closeWrite func()
		errCh      <-chan error
	)

	err = h.streams.Add(sid, func() (io.WriteCloser, types.ResizeHandleFunc, error) {
		downstream, closeWrite, errCh, err = h.impl.PortForward(
			ctx, opts.PodUid, opts.Protocol, opts.Port, pr,
		)
		if err != nil {
			return nil, nil, err
		}

		if downstream == nil || closeWrite == nil || errCh == nil {
			return nil, nil, fmt.Errorf("bad port-forward implementation, missing required return values")
		}

		return &flexWriteCloser{
			writeFunc: pw.Write,
			closeFunc: func() error {
				f, ok := pr.(*os.File)
				if !ok {
					go func() {
						time.Sleep(5 * time.Second)
						_ = pw.Close()
						closeWrite()
						_ = downstream.Close()
					}()
					return nil
				}

				n, _ := iohelper.CheckBytesToRead(f.Fd())

				// assume 100 KB/s
				wait := time.Duration(n) / 1024 / 100 * time.Second
				if wait > 0 {
					go func() {
						time.Sleep(wait)
						_ = pw.Close()
						closeWrite()
						_ = downstream.Close()
					}()
				}

				return nil
			},
		}, nil, nil
	})
	if err != nil {
		return nil, err
	}

	// pipe received data to extension hub
	h.uploadDataOutput(
		ctx,
		sid,
		downstream,
		arhatgopb.MSG_DATA_OUTPUT,
		20*time.Millisecond,
		&seq,
	)

	// downstream read exited

	for {
		// drain errCh
		select {
		case <-ctx.Done():
			return nil, nil
		case e, more := <-errCh:
			if e != nil {
				if err == nil {
					err = e
				} else {
					err = fmt.Errorf("%v; %w", err, e)
				}
			}

			if !more {
				return nil, nil
			}
		}
	}
}
