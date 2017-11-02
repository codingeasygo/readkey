package main

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kr/pty"
	"github.com/sutils/readkey"
)

func main() {
	c := exec.Command("bash")
	vpty, tty, err := pty.Open()
	if err != nil {
		return
	}
	w, h := readkey.GetSize()
	err = readkey.SetSize(vpty.Fd(), w, h)
	if err != nil {
		return
	}
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setctty = true
	c.SysProcAttr.Setsid = true
	err = c.Start()
	if err != nil {
		vpty.Close()
		return
	}
	tty.Close()
	go io.Copy(os.Stdout, vpty)
	time.Sleep(100 * time.Millisecond)
	readkey.Open()
	defer readkey.Close()
	var key []byte
	for err == nil {
		key, err = readkey.Read()
		if err != nil {
			break
		}
		if key[0] == 27 || key[0] == 3 { //esc
			break
		}
		_, err = vpty.Write(key)
		if err != nil {
			break
		}
	}
}
