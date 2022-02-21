package appconf

import (
	"fmt"
	"path"
	"time"
)

const DVPNNodeHomeDirParam = "--home"

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
	ShadowPath         string
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
		HomeDir:            "",
		SentinelDir:        "/.sentinelnode",
		DVPNConfigPath:     "/config.toml",
		WgPath:             "/wireguard.toml",
		CertificatePath:    "/tls.crt",
		SHGenerateCertPath: "/generatecert.sh",
		ShadowPath:         "/etc/shadow",
	}
}

func LoadTestConf() {
	LoadConf()

	Paths.HomeDir = path.Clean(fmt.Sprintf("./temp"))
	Paths.ShadowPath = path.Clean(fmt.Sprintf("./temp/shadow"))

	fmt.Println(Paths.HomeDir)
	fmt.Println(Paths.ShadowPath)
}

func (p *PathsConf) SentinelPath() string {
	return fmt.Sprintf("%v%v", p.HomeDir, p.SentinelDir)
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
