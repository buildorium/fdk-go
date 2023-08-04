package fdk

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const socketPath = "/tmp/func.sock"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func HandleFunc(handlerFunc HandlerFunc) {
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	if err := os.Chmod("/tmp/func.sock", 0777); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Remove(socketPath)
		os.Exit(1)
	}()

	m := http.NewServeMux()
	m.HandleFunc("/", handlerFunc)

	server := http.Server{
		Handler: m,
	}
	log.Println("Connect unix socket successful")
	if err := server.Serve(socket); err != nil {
		log.Fatal(err)
	}
}
