package main

import (
	"context"
	"fmt"
	"go-redis/client"
	"log"
	"sync"
	"testing"
	"time"
)

func TestServerWithMultiClients(t *testing.T) {
	go func() {
		server := NewServer(Config{})
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			client, err := client.New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			defer client.Close()
			key := fmt.Sprintf("client_foo_%d", it)
			value := fmt.Sprintf("client_bar_%d", it)

			if err := client.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}
			val, err := client.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %d got back  => %s\n", it, val)
			wg.Done()
		}(i)
	}
	wg.Wait()
}