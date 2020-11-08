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

package runtime

import (
	"context"
	"io"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/runtimepb"
	"arhat.dev/pkg/wellknownerrors"

	"arhat.dev/libext/extruntime"
	"arhat.dev/libext/types"
)

var _ extruntime.RuntimeEngine = (*SampleRuntime)(nil)

type SampleRuntime struct{}

func (r *SampleRuntime) Name() string          { return "" }
func (r *SampleRuntime) Version() string       { return "" }
func (r *SampleRuntime) OS() string            { return "" }
func (r *SampleRuntime) OSImage() string       { return "" }
func (r *SampleRuntime) Arch() string          { return "" }
func (r *SampleRuntime) KernelVersion() string { return "" }

func (r *SampleRuntime) Exec(
	_ context.Context,
	podUID, container string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	command []string,
	tty bool,
	errCh chan<- *aranyagopb.ErrorMsg,
) (doResize types.ResizeHandleFunc, err error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) Attach(
	ctx context.Context,
	podUID, container string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	errCh chan<- *aranyagopb.ErrorMsg,
) (doResize types.ResizeHandleFunc, err error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) Logs(
	ctx context.Context,
	options *aranyagopb.LogsCmd,
	stdout, stderr io.Writer,
) error {
	return wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) PortForward(
	ctx context.Context,
	podUID string,
	protocol string,
	port int32,
	upstream io.Reader,
	downstream io.Writer,
) error {
	return wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) EnsurePod(
	ctx context.Context, options *runtimepb.PodEnsureCmd,
) (*runtimepb.PodStatusMsg, error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) DeletePod(
	ctx context.Context, options *runtimepb.PodDeleteCmd,
) (*runtimepb.PodStatusMsg, error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) ListPods(
	ctx context.Context, options *runtimepb.PodListCmd,
) (*runtimepb.PodStatusListMsg, error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) EnsureImages(
	ctx context.Context, options *runtimepb.ImageEnsureCmd,
) (*runtimepb.ImageStatusListMsg, error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) DeleteImages(
	ctx context.Context, options *runtimepb.ImageDeleteCmd,
) (*runtimepb.ImageStatusListMsg, error) {
	return nil, wellknownerrors.ErrNotSupported
}

func (r *SampleRuntime) ListImages(
	ctx context.Context, options *runtimepb.ImageListCmd,
) (*runtimepb.ImageStatusListMsg, error) {
	return nil, wellknownerrors.ErrNotSupported
}
