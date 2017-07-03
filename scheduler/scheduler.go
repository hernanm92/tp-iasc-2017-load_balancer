package scheduler

import (
	"math/rand"
	"time"
	"tp-iasc-2017-load_balancer/constants"
)

type ServerData struct {
	Id          int
	Url         string
	EnabledFrom time.Time
}

type ServerScheduler struct {
}

func (scheduler ServerScheduler) GetRandomAvailableServer(servers []ServerData) (ServerData, int) {
	//actualizo lista de servidore disponibless
	availableServers := []ServerData{}
	for index := 0; index < len(servers); index++ {
		server := servers[index]
		if IsAvailable(server) {
			availableServers = append(availableServers, server)
		}
	}
	// valido si tengo disponibles
	if len(availableServers) == 0 {
		return ServerData{}, constants.NO_AVAILABLE_SERVER_CODE
	}
	// obtengo un server random dentro de los disponibles
	n := rand.Intn(len(availableServers))

	return availableServers[n], constants.NO_ERROR_CODE
}

func IsAvailable(server ServerData) bool {
	now := time.Now()
	result := server.EnabledFrom.Before(now)
	return result
}

func (scheduler ServerScheduler) InitServers(urlArray []string) []ServerData {
	servers := make([]ServerData, len(urlArray))
	for index := 0; index < len(urlArray); index++ {
		urlServer := urlArray[index]
		servers[index] = ServerData{
			Id:          index,
			Url:         urlServer,
			EnabledFrom: time.Now(),
		}
	}

	return servers
}
