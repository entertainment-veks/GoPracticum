package os

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ListenShutdownSignals(ctx context.Context, server *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigs
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("error during shutdown server", err)
	}
}
