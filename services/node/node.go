package node

import (
	"github.com/solarlabsteam/dvpn-openwrt/services/socket"
	"io"
	"os"
	"strings"
	"time"
)

type Node struct {
	Online     bool
	Pid        int
	StartTime  time.Time
	SocketConn *socket.Connection `json:"-"`
	StdOut     io.ReadCloser      `json:"-"`
	StdErr     io.ReadCloser      `json:"-"`
	OSProcess  *os.Process        `json:"-"`
}

func New() Node {
	var node Node

	outReader := strings.NewReader("out reader")
	node.StdOut = io.NopCloser(outReader)
	errReader := strings.NewReader("err reader")
	node.StdErr = io.NopCloser(errReader)
	node.SocketConn = &socket.Conn

	return node
}
