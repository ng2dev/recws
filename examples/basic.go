package main

import (
	"context"
	"github.com/ng2dev/recws"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wss := recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	wss.Dial("wss://echo.websocket.org", nil)

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			go wss.Close()
			log.Printf("Websocket closed %s", wss.GetURL())
			return
		default:
			if !wss.IsConnected() {
				log.Printf("Websocket disconnected %s", wss.GetURL())
				continue
			}

			if err := wss.WriteMessage(1, []byte("Incoming")); err != nil {
				log.Printf("Error: WriteMessage %s", wss.GetURL())
				return
			}

			_, message, err := wss.ReadMessage()
			if err != nil {
				log.Printf("Error: ReadMessage %s", wss.GetURL())
				return
			}

			log.Printf("Success: %s", message)
		}
	}
}
