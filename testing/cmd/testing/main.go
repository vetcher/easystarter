package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	dur := flag.Int("duration", 5, "Time to sleep")
	durX := flag.Int("x", 1, "Multiple")
	flag.Parse()
	go func() {
		syschan := make(chan os.Signal)
		signal.Notify(syschan, syscall.SIGTERM, syscall.SIGINT)
		println("WAIT FOR SIGNAL")
		x := <-syschan
		println("STOP", x)
		os.Exit(0)
	}()
	println("Start test service, sleep", *dur**durX, "secs")
	time.Sleep(time.Second * time.Duration(*dur) * time.Duration(*durX))
	moredurstr := os.Getenv("SLEEP_MORE")
	moredur, err := strconv.Atoi(moredurstr)
	if err != nil {
		println(err.Error())
	} else {
		if moredur != 0 {
			println("Test service want to sleep", moredur, "seconds more")
			time.Sleep(time.Second * time.Duration(moredur))
		}
	}
	println("Test service wakes up and closing")
}
