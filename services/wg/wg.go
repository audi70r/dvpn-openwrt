package wg

import (
	"github.com/solarlabsteam/dvpn-openwrt/services/node"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"
	"io/ioutil"
	"os/exec"
	"time"
)

func InitWG() {
	// find wireguard.toml config or initiate it if it's not found
	wgConfigPath := appconf.Paths.WireGuardConfigFullPath()
	_, readErr := ioutil.ReadFile(wgConfigPath)

	if readErr != nil {
		// delay to make sure wg is loaded
		time.Sleep(time.Second * 30)
		if wireguardErr := createWireguardConfig(); wireguardErr != nil {
			panic(wireguardErr)
		}
	}
}

func createWireguardConfig() (err error) {
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeWireguard, node.DVPNNodeConfig, node.DVPNNodeInit, appconf.DVPNNodeHomeDirParam, appconf.Paths.SentinelPath())

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
