package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	dbFile        string
	todos         Todos
	doneChannel   chan string
	signalChannel chan os.Signal
	removePeriod  time.Duration
)

func startMain() {
	var wait time.Duration

	doneChannel = make(chan string, 1)
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()
	todos = CreateTodos()
	initRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	for {
		select {
		case <-signalChannel:
			todos.writeToFile()
			ctx, cancel := context.WithTimeout(context.Background(), wait)
			defer cancel()
			srv.Shutdown(ctx)
			log.Println("shutting down")
			os.Exit(0)
		case id := <-doneChannel:
			mutex.Lock()
			if err := todos.deleteById(id); err == nil {
				log.Printf("Todo item with %s id was deleted due it was done for a certain time now.\n", id)
			} else {
				log.Printf("Failed to delete todo item with %s id, error: %v.\n", id, err)
			}
			mutex.Unlock()
			continue
		}
	}
}

func main() {
	dbFile = "todo-database.json"
	removePeriod = 5 * time.Minute
	startMain()
}
