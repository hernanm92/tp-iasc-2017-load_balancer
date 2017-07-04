package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"tp-iasc-2017-load_balancer/constants"
	"tp-iasc-2017-load_balancer/httpclient"
	"tp-iasc-2017-load_balancer/libs"
	"tp-iasc-2017-load_balancer/scheduler"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends              []string `json:"backends"`
	RequestTimeout        int      `json:"request_timeout"`
	CacheExpiredTime      int      `json:"cache_expired_time"`
	ServerNoAvailableTime int      `json:"server_no_available_time"`
}

var config Config
var cacheClient = cache.CacheClient{cache.NewClient()}
var schedulerClient = scheduler.ServerScheduler{}
var servers []scheduler.ServerData

var httpClient httpclient.HttpClient

func main() {
	LoadConfigFile("config.json")

	servers = schedulerClient.InitServers(config.Backends)
	httpClient = httpclient.HttpClient{config.RequestTimeout}

	fmt.Println("Starting Server...")

	router := gin.Default()

	router.Any("*path", ReverseProxy)

	router.Run(":8080")
}

func ReverseProxy(context *gin.Context) {
	if cacheClient.IsCacheble(context.Request) {
		cacheRequest(context)
	} else {
		fmt.Println("Don't search request in cache, call the backend but don't save it")
		MakeRequest(context)
	}
}

func MakeRequest(context *gin.Context) (string, int) {
	server, errorCode := schedulerClient.GetRandomAvailableServer(servers)
	if errorCode == constants.UNAVAILABLE_SERVER_CODE {
		context.String(constants.UNAVAILABLE_SERVER_CODE, "In this moment we can not attend your request")
		return "", constants.UNAVAILABLE_SERVER_CODE
	}

	bodystring, err := httpClient.DoRequest(context.Request, server.Url)

	if err != nil {
		code := checkError(err)
		if code == constants.TIMEOUT_ERROR_CODE || code == constants.NO_CONNECTION_SERVER {
			SetFutureAvailableTime(server)
			MakeRequest(context)
		} else {
			fmt.Println(err)
			context.String(500, "Server Error")
		}
		return "", constants.ERROR_REQUEST_CODE
	}

	context.String(200, bodystring)
	return bodystring, constants.NO_ERROR_CODE
}

func cacheRequest(c *gin.Context) {
	if cacheClient.ExistsOrNotExpiredKey(c.Request) {
		fmt.Println("The request exists in cache, don't call the backend")
		data := cacheClient.GetRequestValue(c.Request)
		c.String(200, data)

	} else {
		fmt.Println("The request don't exists in cache, call the backend and save it")
		bodystring, code := MakeRequest(c)
		if code != constants.ERROR_REQUEST_CODE {
			cacheClient.SetRequest(c.Request, bodystring, config.CacheExpiredTime)
		}
	}
}

func SetFutureAvailableTime(serverToUpdate scheduler.ServerData) {
	for index := 0; index < len(servers); index++ {
		server := servers[index]
		if server.Id == serverToUpdate.Id {
			serverToUpdate.EnabledFrom = time.Now().Add(time.Duration(config.ServerNoAvailableTime) * time.Minute)
			fmt.Println(server)
			servers[index] = serverToUpdate
			break
		}
	}
	fmt.Println(servers)
}

func checkError(err error) int {
	timeout, _ := regexp.MatchString("Timeout", err.Error())
	if timeout {
		return constants.TIMEOUT_ERROR_CODE
	}

	errorConnection, _ := regexp.MatchString("No se puede establecer una conexión ", err.Error())
	if errorConnection {
		return constants.NO_CONNECTION_SERVER
	}

	return 500
}

func LoadConfigFile(filename string) {
	configFile, _ := os.Open(filename)
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
