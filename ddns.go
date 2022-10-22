package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var onStart = true
var prevIP = "0.0.0.0"

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go runTimer(5, &wg)
	wg.Wait()

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getCurrentIP() string {
	resp, err := http.Get("http://checkip.amazonaws.com/")
	checkError(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		checkError(err)
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	clearIP := strings.ReplaceAll(string(body), "\n", "")
	return clearIP
}

func runTimer(interval int, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				curIP := getCurrentIP()
				if onStart {
					prevIP = curIP
					onStart = false
				} else {
					if curIP != prevIP {
						domain := "db.rv.ua"

						// api, err := cloudflare.NewWithAPIToken("APITOKEN")
						api, err := cloudflare.New("APIKEY", "EMAIL")
						checkError(err)
						ctx := context.Background()
						zoneID, err := api.ZoneIDByName(domain)
						checkError(err)

						dnsRecord, err := api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{Name: domain})
						checkError(err)
						record := dnsRecord[0]
						record.Content = curIP

						err = api.UpdateDNSRecord(ctx, zoneID, record.ID, record)
						checkError(err)
					}
				}
			case <-done:
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()
}
