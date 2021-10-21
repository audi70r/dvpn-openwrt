package appconf

import (
	"fmt"
	"time"
)

type ServerConf struct {
	Addr                              string
	Port                              string
	HttpServerGracefulShutdownTimeout time.Duration
	WriteTimeout                      time.Duration
	ReadTimeout                       time.Duration
	IdleTimeout                       time.Duration
}

type PathsConf struct {
	BinDir             string
	HomeDir            string
	SentinelDir        string
	DVPNConfigPath     string
	WgPath             string
	CertificatePath    string
	SHGenerateCertPath string
}

var Server ServerConf

var Paths PathsConf

func LoadConf() {
	Server = ServerConf{
		Addr:                              "0.0.0.0",
		Port:                              "9000",
		HttpServerGracefulShutdownTimeout: time.Second * 30,
		WriteTimeout:                      time.Second * 15,
		ReadTimeout:                       time.Second * 15,
		IdleTimeout:                       time.Second * 60,
	}

	Paths = PathsConf{
		BinDir:             "/usr/sbin:/usr/bin:/sbin:/bin:",
		HomeDir:            "/root",
		SentinelDir:        "/.sentinelnode",
		DVPNConfigPath:     "/config.toml",
		WgPath:             "/wireguard.toml",
		CertificatePath:    "/tls.crt",
		SHGenerateCertPath: "/generatecert.sh",
	}
}

func (p *PathsConf) DVPNConfigFullPath() string {
	return fmt.Sprintf("%v%v%v", p.HomeDir, p.SentinelDir, p.DVPNConfigPath)
}

func (p *PathsConf) WireGuardConfigFullPath() string {
	return fmt.Sprintf("%v%v%v", p.HomeDir, p.SentinelDir, p.WgPath)
}

func (p *PathsConf) CertificateFullPath() string {
	return fmt.Sprintf("%v%v%v", p.HomeDir, p.SentinelDir, p.CertificatePath)
}
