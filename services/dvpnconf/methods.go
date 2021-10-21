package dvpnconf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/solarlabsteam/dvpn-openwrt/services/node"
	appconf "github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"

	"github.com/pelletier/go-toml"
)

func GetConfigs() (config []byte, err error) {
	var dVPNConfig dVPNConfig

	configPath := appconf.Paths.DVPNConfigFullPath()
	confBytes, readErr := ioutil.ReadFile(configPath)

	// if config could not be read, attempt to init config
	if readErr != nil {
		return initConfig()
	}

	wgConfigPath := appconf.Paths.WireGuardConfigFullPath()
	_, readErr = ioutil.ReadFile(wgConfigPath)

	if readErr != nil {
		return initWireguardConfig()
	}

	tlsCertPath := appconf.Paths.CertificateFullPath()
	_, readErr = ioutil.ReadFile(tlsCertPath)

	if readErr != nil {
		return generateCertificate()
	}

	if err = toml.Unmarshal(confBytes, &dVPNConfig); err != nil {
		return config, err
	}

	config, _ = json.Marshal(dVPNConfig)

	return config, err
}

func PostConfig(config dVPNConfig) (resp []byte, err error) {
	configPath := appconf.Paths.DVPNConfigFullPath()

	configBytes, err := toml.Marshal(config)

	if err != nil {
		return resp, err
	}

	if err = ioutil.WriteFile(configPath, configBytes, 0644); err != nil {
		return resp, err
	}

	resp, err = json.Marshal(config)

	if err != nil {
		return resp, err
	}

	return resp, err
}

func initConfig() (config []byte, err error) {
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeConfig, node.DVPNNodeInit)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}

func initWireguardConfig() (config []byte, err error) {
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeWireguard, node.DVPNNodeConfig, node.DVPNNodeInit)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}

func generateCertificate() (config []byte, err error) {
	cmd := exec.Command(node.BinSH, os.Getenv("HOME")+node.SHGenerateCertPath)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}
