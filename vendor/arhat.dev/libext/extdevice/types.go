package extdevice

import "arhat.dev/arhat-proto/arhatgopb"

type Device interface {
	Operate(params map[string]string, data []byte) ([][]byte, error)
	CollectMetrics(params map[string]string) ([]*arhatgopb.DeviceMetricsMsg_Value, error)
	Close()
}

type DeviceConnector interface {
	Connect(target string, params map[string]string, tlsConfig *arhatgopb.TLSConfig) (Device, error)
}
