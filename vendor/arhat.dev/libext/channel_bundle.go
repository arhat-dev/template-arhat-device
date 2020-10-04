package libext

import (
	"arhat.dev/arhat-proto/arhatgopb"
)

func newChannelBundle() *channelBundle {
	return &channelBundle{
		cmdCh:  make(chan *arhatgopb.Cmd, 1),
		msgCh:  make(chan *arhatgopb.Msg, 1),
		closed: make(chan struct{}),
	}
}

type channelBundle struct {
	cmdCh  chan *arhatgopb.Cmd
	msgCh  chan *arhatgopb.Msg
	closed chan struct{}
}

func (cb *channelBundle) Close() {
	close(cb.closed)
	close(cb.msgCh)
}
