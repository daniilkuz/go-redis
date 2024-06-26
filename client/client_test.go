package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient2(t *testing.T) {

	client, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Set(context.TODO(), "foo", "1"); err != nil {
		log.Fatal(err)
	}
	val, err := client.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET => ", val)

}

func TestNewClient(t *testing.T) {

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
