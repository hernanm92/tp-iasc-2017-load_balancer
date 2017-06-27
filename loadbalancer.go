package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"tp-iasc-2017-load_balancer/httpclient"
	"tp-iasc-2017-load_balancer/libs"
	"tp-iasc-2017-load_balancer/scheduler"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends             []string `json:"backends"`
	WaitTimeSeconds      int      `json:"wait_time_seconds"`
	ExpiredTimeInMinutes int      `json:"expired_time_in_minutes"`
}

var config Config
var cacheClient = cache.CacheClient{cache.NewClient()}
var schedulerClient = scheduler.ServerScheduler{}
var servers []scheduler.ServerData

var httpClient = httpclient.HttpClient{config.WaitTimeSeconds}

func main() {
	LoadConfigFile("config.json")

	servers = schedulerClient.InitServers(config.Backends)

	fmt.Println("Starting Server...")

	router := gin.Default()

	router.Any("*path", ReverseProxy)

	router.Run(":8080")
}

func ReverseProxy(context *gin.Context) {
	if cacheClient.IsCacheble(context.Request) {
		//cacheoss
		cacheRequest(context)
	} else {
		fmt.Println("-----el request no es cacheable, no uso cache------")
		MakeRequest(context)
	}
}

func MakeRequest(context *gin.Context) (string, int) {
	server, errorCode := schedulerClient.RandomServer(servers)
	if errorCode == -1 {
		context.String(200, "En estos momentos no se puede antender esta solicitud")
		return "", -1
	}

	bodystring, err := httpClient.DoRequest(context.Request, server.Url)

	if err != nil {
		code := checkError(err)
		if code == 408 || code == -1 {
			SetUnAvailableCurrentServerBy(server.Url)
			MakeRequest(context)
		} else {
			fmt.Println(err)
			context.String(500, "Error en el servidor")
		}
		return "", -1
	}

	context.String(200, bodystring)
	return bodystring, 0
}

func cacheRequest(c *gin.Context) {
	if cacheClient.ExistsOrNotExpiredKey(c.Request) {
		fmt.Println("----existe el request en cache, traigo y respondo al cliente----------")
		data := cacheClient.GetRequestValue(c.Request)
		c.String(200, data)

	} else {
		fmt.Println("-----No existe en cache, hago el request y guardo----")
		bodystring, code := MakeRequest(c)
		if code != -1 {
			cacheClient.SetRequest(c.Request, bodystring, config.ExpiredTimeInMinutes)
		}
	}
}

func SetUnAvailableCurrentServerBy(url string) {
	for index := 0; index < len(servers); index++ {
		server := servers[index]
		if strings.EqualFold(server.Url, url) {
			server.UnAvailableTime = 2
			servers[index] = server
			break
		}
	}
	fmt.Println(servers)
}

func checkError(err error) int {
	//aca se puede validar mas erroress
	timeout, _ := regexp.MatchString("Timeout", err.Error())
	errorConnection, _ := regexp.MatchString("No se puede establecer una conexión", err.Error())
	if timeout {
		return 408
	}

	if errorConnection {
		return -1
	}
	return 500
}

func LoadConfigFile(filename string) {
	configFile, _ := os.Open(filename)
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
