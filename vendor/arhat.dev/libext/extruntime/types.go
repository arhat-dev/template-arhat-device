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

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/runtimepb"

	"arhat.dev/libext/types"
)

type RuntimeEngine interface {
	// Name of the runtime engine
	Name() string

	// Version of the runtime engine
	Version() string

	// OS the kernel name of the runtime environment
	OS() string

	// OSImage the os distro name
	OSImage() string

	// Arch the cpu arch of the runtime environment
	Arch() string

	// KernelVersion of the OS
	KernelVersion() string

	/*

		Container interactions

	*/

	// Exec execute a command in a running container
	Exec(
		ctx context.Context,
		podUID, container string,
		stdin io.Reader,
		stdout, stderr io.Writer,
		command []string,
		tty bool,
	) (
		doResize types.ResizeHandleFunc,
		errCh <-chan *aranyagopb.ErrorMsg,
		err error,
	)

	// Attach a running container's stdin/stdout/stderr
	Attach(
		ctx context.Context,
		podUID, container string,
		stdin io.Reader,
		stdout, stderr io.Writer,
	) (
		doResize types.ResizeHandleFunc,
		errCh <-chan *aranyagopb.ErrorMsg,
		err error,
	)

	// Logs retrieve
	Logs(
		ctx context.Context,
		options *aranyagopb.LogsCmd,
		stdout, stderr io.Writer,
	) error

	// PortForward establishes a temporary reverse proxy to cloud
	PortForward(
		ctx context.Context,
		podUID string,
		protocol string,
		port int32,
		upstream io.Reader,
	) (
		downstream io.ReadCloser,
		closeWrite func(),
		readErrCh <-chan error,
		err error,
	)

	/*

		Pod operations

	*/

	// EnsurePod creates containers
	EnsurePod(ctx context.Context, options *runtimepb.PodEnsureCmd) (*runtimepb.PodStatusMsg, error)

	// DeletePod kills all containers and delete pod related volume data
	DeletePod(ctx context.Context, options *runtimepb.PodDeleteCmd) (*runtimepb.PodStatusMsg, error)

	// ListPods show (all) pods we are managing
	ListPods(ctx context.Context, options *runtimepb.PodListCmd) (*runtimepb.PodStatusListMsg, error)

	/*

		Image operations

	*/

	// EnsureImages ensure container images
	EnsureImages(ctx context.Context, options *runtimepb.ImageEnsureCmd) (*runtimepb.ImageStatusListMsg, error)

	// DeleteImages deletes images with specified references
	DeleteImages(ctx context.Context, options *runtimepb.ImageDeleteCmd) (*runtimepb.ImageStatusListMsg, error)

	// DeleteImages lists images with specified references or all images
	ListImages(ctx context.Context, options *runtimepb.ImageListCmd) (*runtimepb.ImageStatusListMsg, error)
}
