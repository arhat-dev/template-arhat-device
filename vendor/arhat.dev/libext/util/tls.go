package util

import (
	"crypto/tls"

	"github.com/pion/dtls/v2"
)

func ConvertTLSConfigToDTLSConfig(tlsConfig *tls.Config) *dtls.Config {
	if tlsConfig == nil {
		return nil
	}

	var cs []dtls.CipherSuiteID
	for i := range tlsConfig.CipherSuites {
		cs = append(cs, dtls.CipherSuiteID(tlsConfig.CipherSuites[i]))
	}

	return &dtls.Config{
		Certificates:          tlsConfig.Certificates,
		CipherSuites:          cs,
		InsecureSkipVerify:    tlsConfig.InsecureSkipVerify,
		VerifyPeerCertificate: tlsConfig.VerifyPeerCertificate,
		RootCAs:               tlsConfig.RootCAs,
		ClientCAs:             tlsConfig.ClientCAs,
		ServerName:            tlsConfig.ServerName,

		// TODO: support more dTLS options
		SignatureSchemes:       nil,
		SRTPProtectionProfiles: nil,
		ClientAuth:             0,
		ExtendedMasterSecret:   0,
		FlightInterval:         0,
		PSK:                    nil,
		PSKIdentityHint:        nil,
		InsecureHashes:         false,
		LoggerFactory:          nil,
		MTU:                    0,
		ReplayProtectionWindow: 0,
	}
}
