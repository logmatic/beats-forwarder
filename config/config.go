// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "github.com/elastic/go-ucfg"

type Config struct {
	ucfg.Config
	Input  InputConfig        `config:"input"`
	Output OutputConfig       `config:"output"`
}

type TlsConfig struct {
	Enable   *bool   `config:"enable"`
	CaPath   *string `config:"ca_path"`
	CertPath *string `config:"cert_path"`
	KeyPath  *string `config:"key_path"`
}

type LJConfig struct {
	V1 *bool   `config:"V1"`
	V2 *bool `config:"V2"`
}

type InputConfig struct {
	Host      *string
	Port      *int
	Keepalive *int
	Timeout   *int
	LJ        LJConfig
	TlsConfig TlsConfig     `config:"tls"`
}

type OutputUdpTcpConfig struct {
	Network   *string        `config:"network"`
	Raddr     *string        `config:"raddr"`
	TlsConfig *TlsConfig     `config:"tls"`
}

type OutputSyslogConfig struct {
	Tag      *string         `config:"tag"`
	Hostname *string         `config:"hostname"`
	Network  *string         `config:"network"`
	Raddr    *string         `config:"raddr"`
}

type OutputLogmaticConfig struct {
	Key   *string
	Network *string
	Raddr   *string
}

type OutputConfig struct {
	Type 	 *string
	UDPTCP   OutputUdpTcpConfig      `config:"udp_tcp"`
	Syslog   OutputSyslogConfig        `config:"syslog"`
	Logmatic OutputLogmaticConfig  `config:"logmatic"`
}