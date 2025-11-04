package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"productsParser/internal/config"
	"productsParser/internal/runner"
	"syscall"
	"time"
)

func main() {
	token := flag.String("token", "", "tg token")
	timeout := flag.Duration("httpTimeout", time.Minute*2, "Таймаут HTTP соединений")

	flag.Parse()

	cfg := config.NewGlobalConfig(true, *token, *timeout, true)

	gatewayApp, err := runner.NewApp(context.TODO(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := gatewayApp.Ctx()

	ctxNotify, stop := signal.NotifyContext(
		ctx,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTSTP,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer stop()

	if err := gatewayApp.Run(ctxNotify); err != nil {
		log.Fatal(err)
	}
}
