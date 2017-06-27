package scheduler

import (
	"fmt"
)

type ServerData struct {
	Id              int
	Url             string
	UnAvailableTime int
}

type ServerScheduler struct {
}

//cambiar nombre por "primer disponibe"
func (scheduler ServerScheduler) RandomServer(servers []ServerData) (ServerData, int) {
	//buscar los q no estaninhabilitados sino devovler al cliente que no
	//esta disponible el reques pot falta de servidoress
	//n := rand.Intn(100) % len(config.Backends)
	count := 0
	availableServer := 0
	for index := 0; index < len(servers); index++ {
		server := servers[index]
		if server.UnAvailableTime > 0 {

			count++
		} else {
			availableServer = index

			break
		}
	}

	if count == 3 {
		return ServerData{}, -1
	}

	return servers[availableServer], 0
}

func (scheduler ServerScheduler) InitServers(urlArray []string) []ServerData {
	fmt.Println(urlArray)
	servers := []ServerData{
		{
			Id:              1,
			Url:             urlArray[0],
			UnAvailableTime: 0,
		},
		{
			Id:              2,
			Url:             urlArray[1],
			UnAvailableTime: 0,
		},
		{
			Id:              3,
			Url:             urlArray[2],
			UnAvailableTime: 0,
		},
	}
	return servers
}
