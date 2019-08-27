package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Invisibi-nd/slack-bot/appruntime"
	"github.com/Invisibi-nd/slack-bot/service"
)

func main() {

	go service.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 60 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit
	appruntime.Logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	done := make(chan bool)
	go func() {
		service.Shutdown()
		done <- true
	}()

	wait := make(chan bool)
	go func() {
		select {
		case <-ctx.Done():
			appruntime.Logger.Fatal("Server Shutdown Timeout ")
			wait <- true
			break
		case <-done:
			wait <- true
			break
		}
	}()
	<-wait
	appruntime.Logger.Info("Server exiting")

}
