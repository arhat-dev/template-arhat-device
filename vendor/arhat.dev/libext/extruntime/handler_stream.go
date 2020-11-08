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
	"sync"
	"sync/atomic"
	"time"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/iohelper"
)

func (h *Handler) handleStreamOperation(
	ctx context.Context,
	sid uint64,
	useStdin, useStdout, useStderr, useTty bool,
	preRun func() error,
	run func(stdout, stderr io.WriteCloser) *aranyagopb.ErrorMsg,
) {
	var (
		seq       uint64
		preRunErr error
		err       *aranyagopb.ErrorMsg
		wg        = new(sync.WaitGroup)
	)

	defer func() {
		wg.Wait()

		lastMsg := &arhatgopb.Msg{
			Kind:    arhatgopb.MSG_DATA_OUTPUT,
			Id:      sid,
			Ack:     nextSeq(&seq),
			Payload: nil,
		}
		if err != nil {
			data, _ := err.Marshal()
			lastMsg.Payload = data
			lastMsg.Kind = arhatgopb.MSG_ERROR
		}

		_ = h.SendMsg(lastMsg)

		h.streams.Del(sid)
	}()

	if preRun != nil {
		preRunErr = preRun()
		if preRunErr != nil {
			return
		}
	}

	stdout, stderr, closeStream := h.createTerminalStream(ctx, sid, useStdout, useStderr, useStdin && useTty, &seq, wg)
	defer closeStream()

	err = run(stdout, stderr)
}

func (h *Handler) createTerminalStream(
	ctx context.Context,
	sid uint64,
	useStdout, useStderr, interactive bool,
	pSeq *uint64,
	wg *sync.WaitGroup,
) (stdout, stderr io.WriteCloser, close func()) {
	var (
		readStdout  io.ReadCloser
		readStderr  io.ReadCloser
		readTimeout = 100 * time.Millisecond
		seqMu       *sync.Mutex
	)

	if interactive {
		readTimeout = 20 * time.Millisecond
	}

	if useStdout && useStderr {
		seqMu = new(sync.Mutex)
	}

	if useStdout {
		readStdout, stdout = iohelper.Pipe()
		wg.Add(1)
		go func() {
			defer func() {
				_ = readStdout.Close()
				wg.Done()
			}()

			h.uploadDataOutput(ctx, sid, readStdout, arhatgopb.MSG_DATA_OUTPUT, readTimeout, pSeq, seqMu)
		}()
	}

	if useStderr {
		readStderr, stderr = iohelper.Pipe()
		wg.Add(1)
		go func() {
			defer func() {
				_ = readStderr.Close()
				wg.Done()
			}()

			h.uploadDataOutput(ctx, sid, readStderr, arhatgopb.MSG_RUNTIME_DATA_STDERR, readTimeout, pSeq, seqMu)
		}()
	}

	return stdout, stderr, func() {
		if stdout != nil {
			_ = stdout.Close()
		}

		if stderr != nil {
			_ = stderr.Close()
		}
	}
}

func (h *Handler) uploadDataOutput(
	ctx context.Context,
	sid uint64,
	rd io.Reader,
	kind arhatgopb.MsgType,
	readTimeout time.Duration,
	pSeq *uint64,
	seqMu *sync.Mutex,
) {
	r := iohelper.NewTimeoutReader(rd, 4096)
	go r.StartBackgroundReading()

	stopSig := ctx.Done()
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	for r.WaitUntilHasData(stopSig) {
		timer.Reset(readTimeout)
		data, isTimeout := r.ReadUntilTimeout(timer.C)
		if !isTimeout && !timer.Stop() {
			<-timer.C
		}

		if seqMu != nil {
			seqMu.Lock()
		}

		err := h.SendMsg(&arhatgopb.Msg{
			Kind:    kind,
			Id:      sid,
			Ack:     nextSeq(pSeq),
			Payload: data,
		})

		if seqMu != nil {
			seqMu.Unlock()
		}

		if err != nil {
			return
		}
	}
}

func nextSeq(p *uint64) uint64 {
	seq := atomic.LoadUint64(p)
	for !atomic.CompareAndSwapUint64(p, seq, seq+1) {
		seq++
	}

	return seq
}
