package main

/*
$ curl http://localhost:9999/_simple_cache/scores/Tom
630
$ curl http://localhost:9999/_simple_cache/scores/kkk
kkk not exist
*/

import (
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

func main() {
	simpleCache.NewGroup("scores", 2<<10, simpleCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := simpleCache.NewHTTPPool(addr)
	log.Println("simple-cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
