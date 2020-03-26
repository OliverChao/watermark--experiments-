package main

import (
	"WaterMasking/controller"
	"WaterMasking/model"
	"WaterMasking/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	model.FlagConfInit()
}

func main() {
	service.ConnectDB()
	service.InitSourceCacheData()

	router := controller.MapRoutes()
	server := &http.Server{
		Addr:    "0.0.0.0:" + "8080",
		Handler: router,
	}
	handleSignal(server)

	if err := server.ListenAndServe(); nil != err {
		logger.Fatalf("listen and serve failed: " + err.Error())
	}

}

func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Infof("got signal [%s], exiting pipe now", s)
		if err := server.Close(); nil != err {
			logger.Errorf("server close failed: " + err.Error())
		}

		service.DisconnectDB()

		logger.Infof("Pipe exited")
		os.Exit(0)
	}()
}
