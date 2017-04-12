package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/a-h/round"
	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/parnurzeal/gorequest"
)

type getWithproxy struct {
	proxy     string
	url       string
	serverURL string
}

type proxy struct {
	IP      string
	Port    int
	Country string
	Respone float64
	Status  bool
}

func (g *getWithproxy) getproxy() {
	httpProxy := fmt.Sprintf("https://%s", g.proxy)
	str := strings.Split(g.proxy, ":")
	ip := str[0]
	port, _ := strconv.Atoi(str[1])

	request := gorequest.New().Proxy(httpProxy).Timeout(2 * time.Second)
	timeStart := time.Now()
	resp, _, err := request.Get(g.url).Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusRequestTimeout).End()
	if err == nil && resp.StatusCode == 200 {
		fmt.Println("GOOD: ", g.proxy)
		country := ipToCountry(ip)
		respone := round.ToEven(time.Since(timeStart).Seconds(), 3)
		u := proxy{Country: country, IP: ip, Port: port, Respone: respone, Status: true}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(u)
		res, _ := http.Post(g.serverURL, "application/json; charset=utf-8", b)
		io.Copy(os.Stdout, res.Body)
	} else {
		fmt.Println("BAD: ", g.proxy)
	}
}

func ipToCountry(ip string) string {
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		fmt.Printf("Could not open GeoIP database\n")
		os.Exit(1)
	}
	defer db.Close()
	country, _ := db.Country(net.ParseIP(ip))
	return country.Country.IsoCode
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var (
		url       = flag.String("url", "https://m.vk.com", "")
		fileIn    = flag.String("in", "proxylist.txt", "full path to proxy file")
		serverURL = flag.String("apiurl", "", "API url")
		treds     = flag.Int("treds", 50, "number of treds")
	)

	flag.Parse()

	for {
		content, _ := ioutil.ReadFile(*fileIn)
		if len(content) == 0 {
			time.Sleep(5 * time.Second)
		} else {
			ioutil.WriteFile(*fileIn, []byte(""), 0644)
			proxys := strings.Split(string(content), "\n")

			workers := *treds

			wg := new(sync.WaitGroup)
			in := make(chan string, 2*workers)

			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for proxy := range in {
						gp := getWithproxy{
							proxy:     proxy,
							url:       *url,
							serverURL: *serverURL,
						}
						gp.getproxy()
					}
				}()
			}

			for _, proxy := range proxys {
				if proxy != "" {
					in <- proxy
				}
			}
			close(in)
			wg.Wait()
		}
	}
}
