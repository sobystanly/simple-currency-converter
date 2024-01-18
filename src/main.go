package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"platform-sre-interview-excercise-master/cache"
	"platform-sre-interview-excercise-master/controller"
	"platform-sre-interview-excercise-master/handler"
	"syscall"
	"time"
)

func main() {
	configureLogger()
	client := &http.Client{}
	ratesCache := cache.NewCache()
	converter := controller.NewConverter(client)
	exchangeRate := handler.NewExchangeRate(converter, &ratesCache)

	router := exchangeRate.GetRoutes()

	logrus.Debug("starting exchange-rates service")

	//listen for sigint/term from OS to trigger graceful shutdown
	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8001"),
		Handler: router,
	}

	go serve(srv)

	sig := <-terminationChannel

	logrus.Infof("Termination signal '%s' received, initiating graceful shutdown...", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(25)*time.Second)
	defer cancel()

	//shutdown http server
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to gracefully shutdown HTTP server: %s", err.Error())
	} else {
		logrus.Infof("Successfully shut down http server gracefully.")
	}
}

func serve(srv *http.Server) {
	//start http server
	err := srv.ListenAndServe()
	if err != nil {
		logrus.Panic("can't start http server", err)
	}
	logrus.Debugf("successfully started exchange-rates service on port: %d", 8001)
}

func configureLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
