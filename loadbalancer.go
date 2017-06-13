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
	Backends []string `json:"backends"`
	WaitTime int      `json:"wait_time"`
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
	//si tengo q guardar en cache uso interceptorTransport para tomar el reponse del servidor
	proxy.Transport = &IterceptorTransport{}

	proxy.ServeHTTP(c.Writer, c.Request)
}

//override del roundtrip default de transport
func (t *IterceptorTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	//validar error
	body, err := httputil.DumpResponse(response, true)
	//validar error
	requestString := CreateRequesString(request)
	cacheClient.SetRequest(requestString, string(body))
	data := cacheClient.GetRequest(requestString)
	fmt.Println(data)

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

//falta el body, si es post, y los params si es query
//pasarlo a otra lib y  puede ser mejorable esta funcion
func CreateRequesString(request *http.Request) (req string) {

	accept := request.Header.Get("Accept")
	aceeptEncoding := request.Header.Get("Accept-Encoding")
	acceptLenguage := request.Header.Get("Accept-Language")
	//cacheControl := request.Header.Get("Cache-Control")
	connection := request.Header.Get("Connection")
	cookie := request.Header.Get("Cookie")
	userAgent := request.Header.Get("User-Agent")

	headerString := accept + ";" + aceeptEncoding + ";" + acceptLenguage + ";" + connection + ";" + cookie + ";" + userAgent

	generalRequestString := request.Host + ";" + request.RequestURI + ";" + request.Method + ";" + request.Proto + ";"

	requestString := generalRequestString + headerString

	return requestString
}
