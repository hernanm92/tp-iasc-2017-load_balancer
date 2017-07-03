package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"

	"tp-iasc-2017-load_balancer/constants"
	"tp-iasc-2017-load_balancer/httpclient"
	"tp-iasc-2017-load_balancer/libs"
	"tp-iasc-2017-load_balancer/scheduler"

	"time"

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
		//cacheoss
		cacheRequest(context)
	} else {
		fmt.Println("-----el request no es cacheable, no uso cache------")
		MakeRequest(context)
	}
}

func MakeRequest(context *gin.Context) (string, int) {
	server, errorCode := schedulerClient.GetFirstAvailable(servers)
	if errorCode == constants.NOAVAILABLESERVERCODE {
		context.String(200, "En estos momentos no se puede antender esta solicitud")
		return "", constants.NOAVAILABLESERVERCODE
	}

	bodystring, err := httpClient.DoRequest(context.Request, server.Url)

	if err != nil {
		code := checkError(err)
		if code == constants.TIMEOUTERRORCODE || code == constants.NOCONNECTIONSERVER {
			SetFutureAvailableTime(server)
			MakeRequest(context)
		} else {
			fmt.Println(err)
			context.String(500, "Error en el servidor")
		}
		return "", constants.ERRORREQUESTOCODE
	}

	context.String(200, bodystring)
	return bodystring, constants.NOERRORCODE
}

func cacheRequest(c *gin.Context) {
	if cacheClient.ExistsOrNotExpiredKey(c.Request) {
		fmt.Println("----existe el request en cache, traigo y respondo al cliente----------")
		data := cacheClient.GetRequestValue(c.Request)
		c.String(200, data)

	} else {
		fmt.Println("-----No existe en cache, hago el request y guardo-----")
		bodystring, code := MakeRequest(c)
		if code != constants.ERRORREQUESTOCODE {
			cacheClient.SetRequest(c.Request, bodystring, config.CacheExpiredTime)
		}
	}
}

//pasar el time enable a config
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
	//aca se puede validar mas erroress
	timeout, _ := regexp.MatchString("Timeout", err.Error())
	errorConnection, _ := regexp.MatchString("No se puede establecer una conexión ", err.Error())
	if timeout {
		return constants.TIMEOUTERRORCODE
	}

	if errorConnection {
		return constants.NOCONNECTIONSERVER
	}
	return 500
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
