package node

import (
	"fmt"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"
	"io"
	"os/exec"
	"time"
)

// StartNodeStreamOutputToSocket will run the node as a child process and stream the std output to a socket connection provided as an argument
func (n *Node) StartNodeStreamOutputToSocket() (err error) {
	cmd := exec.Command(DVPNNodeExec, DVPNNodeStart, fmt.Sprintf("%s=%s", appconf.DVPNNodeHomeDirParam, appconf.Paths.SentinelPath()))
	n.StdOut, _ = cmd.StdoutPipe()
	n.StdErr, _ = cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		return err
	}

	n.Online = true
	n.StartTime = time.Now()
	n.OSProcess = cmd.Process
	n.Pid = cmd.Process.Pid

	go stdOutToSocketBridge(n.StdOut, n.SocketConn)
	go stdOutToSocketBridge(n.StdErr, n.SocketConn)

	go func() {
		cmd.Wait()
		n.resetToDefaults()
	}()

	return nil
}

// Kill will destroy the node child process.
func (n *Node) Kill() (err error) {
	// node is already dead
	if n.OSProcess == nil {
		return nil
	}

	if err = n.OSProcess.Kill(); err != nil {
		return err
	}

	n.resetToDefaults()

	return nil
}

// resetToDefaults will reset the node status values
func (n *Node) resetToDefaults() {
	n.Online = false
	n.StartTime = time.Time{}
	n.OSProcess = nil
	n.Pid = 0
}

// stdOutToSocketBridge use the io and send its output to the provided websocket connection
func stdOutToSocketBridge(r io.Reader, s *Connection) {
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			s.Send(d)
			if err != nil {
				break
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			break
		}
	}

	return
}
