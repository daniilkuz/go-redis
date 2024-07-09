package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNewClientRedisClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(context.TODO(), "key").Result()
	if err != nil {
		panic(err)
	}
}

func TestNewClient(t *testing.T) {

	client, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("SET => ", "bar")
	if err := client.Set(context.TODO(), "foo", "bar"); err != nil {
		log.Fatal(err)
	}
	val, err := client.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET => ", val)
}

func TestNewClients(t *testing.T) {

	client, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {

		fmt.Println("SET => ", fmt.Sprintf("bar %d", i))
		if err := client.Set(context.TODO(), fmt.Sprint(i), fmt.Sprintf("bar %d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := client.Get(context.TODO(), fmt.Sprint(i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET => ", val)
	}
	time.Sleep(time.Second)
}
