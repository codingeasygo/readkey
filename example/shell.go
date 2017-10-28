package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"github.com/sutils/readkey"
)

func main() {
	cmd := exec.Command("bash")
	pipe, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, pipe)
	defer readkey.Close()
	var key []byte
	for err == nil {
		key, err = readkey.ReadKey()
		if err != nil {
			break
		}
		if key[0] == 27 { //esc
			break
		}
		_, err = pipe.Write(key)
		if err != nil {
			break
		}
	}
}
