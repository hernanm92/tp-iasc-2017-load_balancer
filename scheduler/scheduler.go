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
	//obtengo los disponibles, yo supongo q la lista esta actualizada
	n := rand.Intn(100) % len(availableServers)
	//obtengo los random
	return availableServers[n], constants.NO_ERROR_CODE
}

//aca se podria usar un sistemas de listas de disp y no disps
func (scheduler ServerScheduler) GetFirstAvailable(servers []ServerData) (ServerData, int) {
	//n := rand.Intn(100) % len(config.Backends)
	unavailableServers := 0
	availableServer := 0
	for index := 0; index < len(servers); index++ {
		server := servers[index]
		if IsAvailable(server) {
			availableServer = index
			break
		} else {
			unavailableServers++

		}
	}

	if unavailableServers == len(servers) {
		return ServerData{}, constants.NO_AVAILABLE_SERVER_CODE
	}

	return servers[availableServer], constants.NO_ERROR_CODE
}

func IsAvailable(server ServerData) bool {
	now := time.Now()
	result := server.EnabledFrom.Before(now)
	return result
}

//esto pasarlo a una propiedad interna
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
