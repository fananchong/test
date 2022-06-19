package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func httpGet(client *http.Client, url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(response)
}

func httpPut(client *http.Client, url string, data string) string {
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(response)
}

func main() {

	var clients []*http.Client

	clientNum := 100
	for i := 0; i < clientNum; i++ {
		clients = append(clients, &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					KeepAlive: 60 * time.Second,
					Timeout:   60 * time.Second,
				}).DialContext,
			},
		})
	}

	urls := []string{"http://127.0.0.1:12380", "http://127.0.0.1:22380", "http://127.0.0.1:32380"}

	var putCounter int32
	var putNum = 100000
	var putUrlCounter = 0
	var putWG sync.WaitGroup
	for i := 0; i < clientNum; i++ {
		putWG.Add(1)
		go func() {
			for i := 0; i < putNum; i++ {
				key := fmt.Sprintf("key%v", i)
				value := fmt.Sprintf("value%v", i)
				url := fmt.Sprintf("%s/%s", urls[(putUrlCounter+1)%len(urls)], key)
				httpPut(clients[i%len(clients)], url, value)
				atomic.AddInt32(&putCounter, 1)
			}
			putWG.Done()
		}()
	}
	putTick := time.NewTicker(1 * time.Second)
	go func() {
		for range putTick.C {
			c := atomic.SwapInt32(&putCounter, 0)
			fmt.Printf("put %v/s\n", c)
		}
	}()
	putWG.Wait()
	putTick.Stop()

	var getCounter int32
	var getNum = 100000
	var getUrlCounter = 0
	var getWG sync.WaitGroup
	for i := 0; i < clientNum; i++ {
		getWG.Add(1)
		go func() {
			for i := 0; i < getNum; i++ {
				key := fmt.Sprintf("key%v", i)
				url := fmt.Sprintf("%s/%s", urls[(getUrlCounter+1)%len(urls)], key)
				httpGet(clients[i%len(clients)], url)
				atomic.AddInt32(&getCounter, 1)
			}
			getWG.Done()
		}()
	}
	getTick := time.NewTicker(1 * time.Second)
	go func() {
		for range getTick.C {
			c := atomic.SwapInt32(&getCounter, 0)
			fmt.Printf("get %v/s\n", c)
		}
	}()
	getWG.Wait()
	getTick.Stop()
}
