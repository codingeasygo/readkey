package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"

	"github.com/sutils/readkey"

	"golang.org/x/net/websocket"
)

var CharTerm = []byte{3}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: rkeynet <remote uri>\n")
		return
	}
	var uri string
	uri = os.Args[1]
	var err error
	var conn io.ReadWriteCloser
	rurl, err := url.Parse(uri)
	if err == nil {
		switch rurl.Scheme {
		case "tcp":
			conn, err = net.Dial("tcp", rurl.Host)
		case "ws":
			conn, err = websocket.Dial(uri, "", "https://"+rurl.Host)
		case "wss":
			conn, err = websocket.Dial(uri, "", "https://"+rurl.Host)
		case "unix":
			conn, err = net.Dial("unix", rurl.Host)
		default:
			err = fmt.Errorf("scheme(%v) is not suppored", rurl.Scheme)
		}
	}
	if err != nil {
		fmt.Printf("connect to %v fail with %v\n", uri, err)
		os.Exit(1)
	}
	readkey.Open()
	defer func() {
		conn.Close()
		readkey.Close()
		os.Exit(1)
	}()
	go func() {
		io.Copy(os.Stdout, conn)
		fmt.Printf("connection is closed\n")
		readkey.Close()
		os.Exit(1)
	}()
	stopc := 0
	for {
		key, err := readkey.Read()
		if err != nil {
			break
		}
		if bytes.Equal(key, CharTerm) {
			stopc++
			if stopc > 5 {
				break
			}
		} else {
			stopc = 0
		}
		_, err = conn.Write(key)
		if err != nil {
			break
		}
	}
}
