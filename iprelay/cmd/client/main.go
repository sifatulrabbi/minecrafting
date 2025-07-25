package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	doRegisterReq()
evLoop:
	for {
		select {
		case <-ctx.Done():
			break evLoop
		case <-time.Tick(time.Second * 5):
			doRegisterReq()
		}
	}
	fmt.Println("Shutting down the application...")
}

func doRegisterReq() {
	// url := "http://localhost:8080/register-ip"
	url := "https://sifatulrabbi.com/minecraft/register-ip"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))
}
