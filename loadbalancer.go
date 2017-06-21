package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"tp-iasc-2017-load_balancer/libs"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends             []string `json:"backends"`
	WaitTime             int      `json:"wait_time"`
	ExpiredTimeInMinutes int      `json:"expired_time_in_minutes"`
}

type IterceptorTransport struct {
}

var config Config

var cacheClient = cache.CacheClient{cache.NewClient()}

func main() {
	LoadConfigFile("config.json")

	fmt.Println("Starting Server...")

	router := gin.Default()

	router.Any("*path", ReverseProxy)

	router.Run(":8080")
}

func ReverseProxy(c *gin.Context) {
	target := RandomServer()

	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	//esta logica tine q estar en otro lado
	if !cacheClient.IsCacheble(c.Request) {
		//no usa cache, actuo normal
		proxy.ServeHTTP(c.Writer, c.Request)
		msg("no es cacheable")
	} else {
		//cacheo
		if cacheClient.ExistsOrNotExpiredKey(c.Request) {
			//existe, traigo de la redis y respondo
			msg("existe request")
			data := cacheClient.GetRequestValue(c.Request)
			c.JSON(200, gin.H{"data": data})
		} else {
			//hago el llamado, guardo en el roundtrip y envio
			msg("-----No existe hago el llamado----")
			proxy.Transport = &IterceptorTransport{}
			proxy.ServeHTTP(c.Writer, c.Request)
		}

	}
}

func msg(msg string) {
	fmt.Println("----------------- " + msg + " --------------")
}

//override del roundtrip default de transport
func (t *IterceptorTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	a, _ := response.Body.Read(b)
	//validar error
	cacheClient.SetRequest(request, string(a)) //agregar expiracion

	return response, err
}

func RandomServer() string {
	n := rand.Intn(100) % len(config.Backends)

	return config.Backends[n]
}

func LoadConfigFile(filename string) {
	configFile, _ := os.Open(filename)
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
