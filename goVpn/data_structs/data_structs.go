package data_structs

import (
	"fmt"
	"time"
)

type VPNConfig struct {
	IsClient             bool
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
	return fmt.Sprintf("\nClient: %v\nHost: %s\nPort: %d\nca_file: %s\nCaInbuilt: %v\ncert_file: %s\nCertInbuilt: %v\nkey_file: %s\nKeyInbuilt: %v\nsecret_file: %v\nsecret_tls: %v\nauth: %v\nauth_file: %v\nProto: %s", v.IsClient, v.RemoteHost, v.RemotePort, v.CaFilename, v.CaInbuilt, v.CertFilename, v.CertInbuilt, v.KeyFileName, v.KeyInbuilt, v.SecretFilename, v.TlsAuth, v.AuthUserPass, v.AuthUserPassFilename, v.Proto)
}

type InitInfo struct {
	TimeInit       time.Time
	Name           string
	Password       string
	Attention      []error
	ConfigFilePath string
	Filename       string
}

func (i InitInfo) String() string {
	return fmt.Sprintf("ConfigFileName: %s\nFileName: %s\nTime start: %v\nUsername: %s\nPassword: %s\nAttentions: %s", i.ConfigFilePath, i.Filename, i.TimeInit, i.Name, i.Password, i.Attention)
}

func NewInitInfo() *InitInfo {
	return &InitInfo{TimeInit: time.Now(), Attention: make([]error, 0)}
}
