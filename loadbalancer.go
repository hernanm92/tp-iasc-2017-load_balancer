package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"tp-iasc-2017-load_balancer/libs"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends             []string `json:"backends"`
	WaitTimeSeconds      int      `json:"wait_time_seconds"`
	ExpiredTimeInMinutes int      `json:"expired_time_in_minutes"`
}

type IterceptorTransport struct {
}

type Msg struct {
	message string
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

	//esta logica tine q estar en otro lado
	if !cacheClient.IsCacheble(c.Request) {
		//no usa cache, actuo normal
		resp := DoRequest("GET", target+c.Request.RequestURI, nil, nil)
		r, _ := json.Marshal(resp)
		c.Writer.Write(r)

	} else {
		//cacheo
		if cacheClient.ExistsOrNotExpiredKey(c.Request) {
			//existe, traigo de la redis y respondo
			msg("existe request")
			//data := cacheClient.GetRequestValue(c.Request)

		} else {
			//hago el llamado, guardo en el roundtrip y envio
			msg("-----No existe hago el llamado----")
		}

	}
}

var client = &http.Client{}

func DoRequest(method string, url string, header http.Header, body io.ReadCloser) *http.Response {

	//fijarse como meter el body
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)
	return resp
}

func msg(msg string) {
	fmt.Println("----------------- " + msg + " --------------")
}

func RandomServer() string {
	//buscar los q no estaninhabilitados sino devovler al cliente que no
	//esta disponible el reques pot falta de servidores.
	n := rand.Intn(100) % len(config.Backends)

	return config.Backends[n]
}

func LoadConfigFile(filename string) {
	configFile, _ := os.Open(filename)
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
