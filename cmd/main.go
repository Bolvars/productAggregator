package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"productsParser/internal/config"
	"productsParser/internal/runner"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	//token := flag.String("token", "", "token")
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, используются системные переменные")
	}
	isTg := flag.Bool("isTg", false, "is tgbot or maxbot")
	isWebhook := flag.Bool("isWebhook", true, "is tgbot or maxbot")
	timeout := flag.Duration("httpTimeout", time.Minute*2, "Таймаут HTTP соединений")

	flag.Parse()
	token := os.Getenv("TOKEN")
	secret := os.Getenv("SECRET")
	host := os.Getenv("HOST")
	defaultURL := os.Getenv("URLs")

	cfg := config.NewGlobalConfig(*isTg, *isWebhook, token, secret, host, defaultURL, *timeout, true)

	gatewayApp, err := runner.NewApp(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := gatewayApp.Ctx()

	ctxNotify, stop := signal.NotifyContext(
		ctx,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGTSTP,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer stop()

	if err := gatewayApp.Run(ctxNotify); err != nil {
		log.Fatal(err)
	}
}
