package extperipheral

import (
	"arhat.dev/arhat-proto/arhatgopb"
)

type Peripheral interface {
	Operate(params map[string]string, data []byte) ([][]byte, error)
	CollectMetrics(params map[string]string) ([]*arhatgopb.PeripheralMetricsMsg_Value, error)
	Close()
}

type PeripheralConnector interface {
	Connect(target string, params map[string]string, tlsConfig *arhatgopb.TLSConfig) (Peripheral, error)
}
