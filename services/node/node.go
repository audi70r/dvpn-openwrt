package node

import (
	"io"
	"os"
	"strings"
	"time"
)

type Node struct {
	Online     bool
	Pid        int
	StartTime  time.Time
	SocketConn *Connection   `json:"-"`
	StdOut     io.ReadCloser `json:"-"`
	StdErr     io.ReadCloser `json:"-"`
	OSProcess  *os.Process   `json:"-"`
}

func New() Node {
	var node Node

	outReader := strings.NewReader("out reader")
	node.StdOut = io.NopCloser(outReader)
	errReader := strings.NewReader("err reader")
	node.StdErr = io.NopCloser(errReader)
	node.SocketConn = &Connection{}

	return node
}
