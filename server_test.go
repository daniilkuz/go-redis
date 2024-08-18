package main

import (
	"bytes"
	"context"
	"fmt"
	"go-redis/client"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tidwall/resp"
)

func TestOfficialRedisClient(t *testing.T) {
	listenAddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenAddr,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", listenAddr),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// fmt.Println(rdb)
	// fmt.Println("this is working :)")

	// key := "foo"
	// val := "bar"

	testCases := map[string]string{
		"foo":  "bar",
		"I":    "would",
		"know": "if",
		"you":  "told me",
	}

	for key, val := range testCases {
		err := rdb.Set(context.Background(), key, val, 0).Err()
		if err != nil {
			// panic(err)
			t.Fatal(err)
		}

		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}

		if val != newVal {
			t.Fatalf("expected %v, but got %v instead", val, newVal)
		}
	}

	// fmt.Printf("got this value %v\n", val)

	// val, err := rdb.Get(context.TODO(), "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(val)
}

func TestFooBar(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteString("OK")
	fmt.Println(buf.String())

	in := map[string]string{
		"server":  "redis",
		"version": "6.0",
	}
	out := respWriteMap(in)
	fmt.Println(string(out))
}

func TestServerWithMultiClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
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
	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers, but got %d", len(server.peers))
	}
}
