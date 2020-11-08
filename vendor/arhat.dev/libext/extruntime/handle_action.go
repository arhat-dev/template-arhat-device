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
	"io"
	"net"
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
			errCh := make(chan *aranyagopb.ErrorMsg)
			if opts.Stdin {
				err = h.streams.Add(sid, func() (_ io.WriteCloser, _ types.ResizeHandleFunc, err error) {
					var (
						pr io.ReadCloser
						pw io.WriteCloser
					)

					if opts.Stdin {
						pr, pw = iohelper.Pipe()
						defer func() {
							if err != nil {
								_ = pw.Close()
								_ = pr.Close()
							}
						}()
					}

					resize, err := h.impl.Exec(
						ctx, opts.PodUid, opts.Container, pr, stdout, stderr, opts.Command, opts.Tty, errCh,
					)
					if err != nil {
						return nil, nil, err
					}

					return pw, resize, err
				})
			} else {
				_, err = h.impl.Exec(
					ctx, opts.PodUid, opts.Container, nil, stdout, stderr, opts.Command, opts.Tty, errCh,
				)
			}
			if err != nil {
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: err.Error(),
					Code:        0,
				}
			}

			select {
			case <-ctx.Done():
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: ctx.Err().Error(),
					Code:        0,
				}
			case errMsg := <-errCh:
				return errMsg
			}
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
			errCh := make(chan *aranyagopb.ErrorMsg)

			if opts.Stdin {
				err = h.streams.Add(sid, func() (_ io.WriteCloser, _ types.ResizeHandleFunc, err error) {
					var (
						stdin io.ReadCloser
						pw    io.WriteCloser
					)

					if opts.Stdin {
						stdin, pw = iohelper.Pipe()
						defer func() {
							if err != nil {
								_ = pw.Close()
								_ = stdin.Close()
							}
						}()
					}

					resize, err := h.impl.Attach(ctx, opts.PodUid, opts.Container, stdin, stdout, stderr, errCh)
					if err != nil {
						return nil, nil, err
					}

					return pw, resize, err
				})
			} else {
				_, err = h.impl.Attach(ctx, opts.PodUid, opts.Container, nil, stdout, stderr, errCh)
			}
			if err != nil {
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: err.Error(),
					Code:        0,
				}
			}

			select {
			case <-ctx.Done():
				return &aranyagopb.ErrorMsg{
					Kind:        aranyagopb.ERR_COMMON,
					Description: ctx.Err().Error(),
					Code:        0,
				}
			case errMsg := <-errCh:
				return errMsg
			}
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

func (h *Handler) handlePortForward(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(aranyagopb.PortForwardCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	var (
		seq uint64

		upstream, downstream = net.Pipe()
	)

	defer func() {
		_ = upstream.Close()

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

	err = h.streams.Add(sid, func() (io.WriteCloser, types.ResizeHandleFunc, error) {
		err = h.impl.PortForward(ctx, opts.PodUid, opts.Protocol, opts.Port, upstream, downstream)
		if err != nil {
			return nil, nil, err
		}

		return upstream, nil, nil
	})
	if err != nil {
		return nil, err
	}

	// pipe received data to kubectl
	h.uploadDataOutput(
		ctx,
		sid, downstream,
		arhatgopb.MSG_DATA_OUTPUT,
		20*time.Millisecond,
		&seq, nil,
	)

	return nil, nil
}
