package scheduler

import (
	"time"
)

type ServerData struct {
	Id          int
	Url         string
	EnabledFrom time.Time
}

type ServerScheduler struct {
}

//cambiar nombre por "primer disponibe"
func (scheduler ServerScheduler) GetFirstAvailable(servers []ServerData) (ServerData, int) {
	//n := rand.Intn(100) % len(config.Backends)
	count := 0
	availableServer := 0
	for index := 0; index < len(servers); index++ {
		server := servers[index]
		if IsAvailable(server) {
			availableServer = index
			break
		} else {
			count++
		}
	}

	if count == 3 {
		return ServerData{}, -1
	}

	return servers[availableServer], 0
}

func IsAvailable(server ServerData) bool {
	now := time.Now()
	result := server.EnabledFrom.Before(now)
	return result
}

//esto pasarlo a una propiedad interna
func (scheduler ServerScheduler) InitServers(urlArray []string) []ServerData {
	servers := []ServerData{
		{
			Id:          1,
			Url:         urlArray[0],
			EnabledFrom: time.Now(),
		},
		{
			Id:          2,
			Url:         urlArray[1],
			EnabledFrom: time.Now(),
		},
		{
			Id:          3,
			Url:         urlArray[2],
			EnabledFrom: time.Now(),
		},
	}
	return servers
}
