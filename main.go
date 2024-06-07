package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{Addr: ":3000"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.Write([]byte("waiting finished \n"))
	})

	go func() {
		fmt.Println("server Running at localhost:3000")
		if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			panic("server down: " + err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Println("shutdown down server ...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("could not graceful shutdown server")
	}
	fmt.Println("server graceful shutdown finished")
}
