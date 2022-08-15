package main

import (
	"flag"
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

func createGroup() *cachee.Group {
	return cachee.NewGroup("cache", 2<<10, cachee.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slow] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, group *cachee.Group) {
	peers := cachee.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("cachee is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *cachee.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			key := req.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, err = w.Write(view.ByteSlice())
			if err != nil {
				return 
			}
		}))
	log.Println("frontend is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var (
		port int
		api bool
		addrs []string
	)
	flag.IntVar(&port, "port", 8001, "cachee server port")
	flag.BoolVar(&api, "api", false, "start an api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	g := createGroup()
	if api {
		go startAPIServer(apiAddr, g)
	}
	startCacheServer(addrMap[port], []string(addrs), g)
}
