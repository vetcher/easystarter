package main

import (
	"net/http"

	"log"

	"github.com/vetcher/easystarter/backend/services"
)

func main() {
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleAllServicesNames))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleGetEnv))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleKillServices))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleLoadAllServices))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleLoadServices))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleReloadEnv))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleRestartServices))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleServicesInfo))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleStartServices))
	http.HandleFunc("", services.NewHTTPWrapper(services.HandleStopServices))
	log.Println(http.ListenAndServe(":8999", nil))
}
