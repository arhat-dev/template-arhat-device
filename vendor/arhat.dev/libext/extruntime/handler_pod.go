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

	"arhat.dev/aranya-proto/aranyagopb/runtimepb"
)

func (h *Handler) handlePodList(ctx context.Context, _ uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(runtimepb.PodListCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	msg, err := h.impl.ListPods(ctx, opts)
	if err != nil {
		return nil, err
	}

	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	return &runtimepb.Packet{
		Kind:    runtimepb.MSG_POD_STATUS_LIST,
		Payload: data,
	}, nil
}

func (h *Handler) handlePodEnsure(ctx context.Context, _ uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(runtimepb.PodEnsureCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	msg, err := h.impl.EnsurePod(ctx, opts)
	if err != nil {
		return nil, err
	}

	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	return &runtimepb.Packet{
		Kind:    runtimepb.MSG_POD_STATUS,
		Payload: data,
	}, nil
}

func (h *Handler) handlePodDelete(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error) {
	opts := new(runtimepb.PodDeleteCmd)
	err := opts.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	msg, err := h.impl.DeletePod(ctx, opts)
	if err != nil {
		return nil, err
	}

	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	return &runtimepb.Packet{
		Kind:    runtimepb.MSG_POD_STATUS,
		Payload: data,
	}, nil
}
