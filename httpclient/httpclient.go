package httpclient

import (
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	WaitTimeSeconds int
}

var client = &http.Client{}

func (httpClient HttpClient) DoRequest(request *http.Request, url string) (string, error) {

	req, _ := http.NewRequest(request.Method, url+request.RequestURI, request.Body)
	req.Header = request.Header
	client.Timeout = time.Duration(httpClient.WaitTimeSeconds) * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	return getStringBodyFrom(resp), nil
}

func getStringBodyFrom(response *http.Response) string {
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	stringBody := string(body)

	return stringBody
}
