package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

var UNIX_SOCKET_FILE = filepath.Join(os.TempDir(), "easystarter.sock")

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			println(err.Error())
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func main() {
	c, err := net.Dial("unix", UNIX_SOCKET_FILE)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c)
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		text := stdin.Text()
		_, err := c.Write([]byte(text))
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
	}
}
