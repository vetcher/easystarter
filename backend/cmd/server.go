package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vetcher/easystarter/backend/services"
)

const (
	VERSION    = "0.1"
	WelcomeTip = "Easy Starter Server " + VERSION
)

var addr = flag.String("addr", ":8999", "Server's listen port")

func handleSignals() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-sigChan
	log.Print("Stop all services")
	err := <-services.ServeStopServices(<-services.ServeAllServicesNames()...)
	if err != nil {
		log.Println(err)
	}
	log.Print("Terminate")
	os.Exit(int(sig.(syscall.Signal)))
}

func main() {
	flag.Parse()
	go handleSignals()
	http.HandleFunc("/cfg/loadall", services.NewHTTPWrapper(services.HandleLoadAllServices))
	http.HandleFunc("/cfg/load", services.NewHTTPWrapper(services.HandleLoadServices))
	http.HandleFunc("/env/get", services.NewHTTPWrapper(services.HandleGetEnv))
	http.HandleFunc("/env/reload", services.NewHTTPWrapper(services.HandleReloadEnv))
	http.HandleFunc("/svc/names", services.NewHTTPWrapper(services.HandleAllServicesNames))
	http.HandleFunc("/svc/info", services.NewHTTPWrapper(services.HandleServicesInfo))
	http.HandleFunc("/svc/start", services.NewHTTPWrapper(services.HandleStartServices))
	http.HandleFunc("/svc/stop", services.NewHTTPWrapper(services.HandleStopServices))
	http.HandleFunc("/svc/restart", services.NewHTTPWrapper(services.HandleRestartServices))
	http.HandleFunc("/svc/kill", services.NewHTTPWrapper(services.HandleKillServices))
	http.HandleFunc("/stop", services.StopServer)
	log.Println(WelcomeTip)
	log.Println("Serve on", *addr)
	log.Println(http.ListenAndServe(*addr, nil))
}
