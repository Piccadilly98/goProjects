package data_structs

import (
	"fmt"
)

type VPNConfig struct {
	RemoteHost           string
	RemotePort           int
	CaFilename           string
	CaInbuilt            string
	CertFilename         string
	CertInbuilt          string
	KeyFileName          string
	KeyInbuilt           string
	TlsAuth              string
	SecretFilename       string
	AuthUserPass         bool
	AuthUserPassFilename string
	Proto                string
}

func (v VPNConfig) String() string {
	return fmt.Sprintf("Host: %s\nPort: %d\nca_file: %s\nCaInbuilt: %v\ncert_file: %s\nCertInbuilt: %v\nkey_file: %s\nKeyInbuilt: %v\nsecret_file: %v\nsecret_tls: %v\nauth: %v\nProto: %s", v.RemoteHost, v.RemotePort, v.CaFilename, v.CaInbuilt, v.CertFilename, v.CertInbuilt, v.KeyFileName, v.KeyInbuilt, v.SecretFilename, v.TlsAuth, v.AuthUserPass, v.Proto)
}
