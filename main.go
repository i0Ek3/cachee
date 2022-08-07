package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/i0Ek3/cachee/cachee"
)

var db = map[string]string{
	"A": "1",
	"B": "2",
	"C": "3",
}

func main() {
	cachee.NewGroup("cache", 2<<10, cachee.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slow] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := cachee.NewHTTPPool(addr)
	log.Println("cachee is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
