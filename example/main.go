package main

/*
$ curl http://localhost:9999/_simple_cache/scores/Tom
630
$ curl http://localhost:9999/_simple_cache/scores/kkk
kkk not exist
*/

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	simpleCache "simple-cache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *simpleCache.Group {
	return simpleCache.NewGroup("scores", 2<<10, simpleCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, adds []string, simple *simpleCache.Group) {
	peers := simpleCache.NewHTTPPool(addr)
	peers.Set(adds...)
	simple.RegisterPeers(peers)
	log.Println("simple cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, simple *simpleCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := simple.Get(key)
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
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Simple Cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var adds []string
	for _, v := range addrMap {
		adds = append(adds, v)
	}

	simple := createGroup()
	if api {
		go startAPIServer(apiAddr, simple)
	}
	startCacheServer(addrMap[port], adds, simple)
}
