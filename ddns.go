package main

import (
	"context"
	"encoding/json"
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

type Domains struct {
	Domains []Domain `json:"domains"`
}

type Domain struct {
	Domain      string `json:"domain"`
	Provider    string `json:"provider"`
	IsSubdomain bool   `json:"is_subdomain"`
	Options     Options
}

type Options struct {
	ApiToken string `json:"api_token"`
}

type Config struct {
	Configs struct {
		Period int `json:"period"`
	} `json:"configs"`
}

func main() {
	//appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//checkError(err)
	//configPath := path.Join(appDir, "config.json")
	data, err := ioutil.ReadFile("config.json")
	checkError(err)

	var domains Domains
	var period Config

	err = json.Unmarshal(data, &period)
	checkError(err)
	err = json.Unmarshal(data, &domains)
	checkError(err)

	runPeriod := period.Configs.Period
	if runPeriod > 0 {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go runTimer(runPeriod, domains, &wg)
		wg.Wait()
	} else {
		fmt.Println("Period must be greater than 0")
	}
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

func getDomainFromSubdomain(subdomain string) string {
	split := strings.Split(subdomain, ".")
	if len(split) > 2 {
		var rootDomain []string
		for i, d := range split {
			if i > 0 {
				rootDomain = append(rootDomain, d)
			}
		}
		return strings.Join(rootDomain, ".")
	} else {
		return subdomain
	}
}

func runTimer(interval int, domains Domains, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				curIP := getCurrentIP()
				fmt.Println("Current IP: ", curIP)
				if onStart {
					prevIP = curIP
					onStart = false
				} else {
					if curIP != prevIP {
						for _, d := range domains.Domains {

							domain := d.Domain
							token := d.Options.ApiToken
							if !d.IsSubdomain {
								api, err := cloudflare.NewWithAPIToken(token)
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
							} else {
								rootDomain := getDomainFromSubdomain(domain)

								api, err := cloudflare.NewWithAPIToken(token)
								checkError(err)
								ctx := context.Background()
								zoneID, err := api.ZoneIDByName(rootDomain)
								checkError(err)

								dnsRecord, err := api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{Name: domain})
								checkError(err)

								record := dnsRecord[0]
								record.Content = curIP

								err = api.UpdateDNSRecord(ctx, zoneID, record.ID, record)
								checkError(err)
							}
						}
						prevIP = curIP
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
