package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
	"tp-iasc-2017-load_balancer/libs"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends             []string `json:"backends"`
	WaitTimeSeconds      int      `json:"wait_time_seconds"`
	ExpiredTimeInMinutes int      `json:"expired_time_in_minutes"`
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
	//esta logica tiene q estar en otro lado
	if cacheClient.IsCacheble(c.Request) {
		//cacheo
		if cacheClient.ExistsOrNotExpiredKey(c.Request) {
			//existe, traigo de la redis y respondo
			fmt.Println("existe request")
			data := cacheClient.GetRequestValue(c.Request)
			c.String(200, data)

		} else {
			//hago el llamado, guardo
			bodystring := DoRequest(c.Request, target)
			cacheClient.SetRequest(c.Request, bodystring, config.ExpiredTimeInMinutes)
			c.String(200, bodystring)
			fmt.Println("-----No existe hago el llamado----")
		}
	} else {
		//no usa cache, actuo normal
		bodystring := DoRequest(c.Request, target)
		c.String(200, bodystring)
		fmt.Println("No es cacheable")

	}
}

//Aca setear el timeout
var client = &http.Client{}

func DoRequest(request *http.Request, url string) (bodystring string) {

	req, _ := http.NewRequest(request.Method, url+request.RequestURI, request.Body)
	req.Header = request.Header
	client.Timeout = time.Duration(config.WaitTimeSeconds) * time.Second
	resp, _ := client.Do(req)
	//code := checkError(err.Error())

	//buscar servidor nuevo y enviar
	/*if code == 408 {
		target := RandomServer() // si target es cero -> enviar un msenaje dicieidno q no esta disponilbe el recurso
		DoRequest(request, target)
	}*/

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(body)

	return stringBody
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

func checkError(msg string) int {
	timeout, _ := regexp.MatchString("Timeout", msg)

	if timeout {
		return 408
	}

	return 500
}
