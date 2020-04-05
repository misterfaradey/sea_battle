package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"local/sea_fight/internal/logger"
	"local/sea_fight/internal/server"
	"local/sea_fight/internal/server/controllers"
	"local/sea_fight/internal/service"
	"log"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const mode = gin.DebugMode //gin.DebugMode gin.ReleaseMode

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger.DefaultLogger()

	srv := service.NewSeaBattleService()

	httpServer := server.NewServer(&server.ServerConf{
		GinMode:        mode,
		Address:        ":8080",
		ReadTimeout:    time.Second * 3,
		WriteTimeout:   time.Second * 3,
		MaxHeaderBytes: 10240,
	}, controllers.NewSeaBattleSquareController(srv))

	go runServer(cancel, httpServer, logger)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	logger.Println("Running ...")

	select {
	case <-signalChan:
	case <-ctx.Done():
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logger.Println("Closing application ...")

	closeServer(ctx, httpServer, logger)

}

func runServer(cancel context.CancelFunc, server server.Server, logger *log.Logger) {

	err := server.Run()
	if err == http.ErrServerClosed {
		logger.Println("server closed")
		return
	}

	if err != nil {
		logger.Println(err)
		cancel()
	}
}

func closeServer(ctx context.Context, server server.Server, logger *log.Logger) {

	err := server.Shutdown(ctx)
	if err != nil {
		logger.Println(err)
	}
}
