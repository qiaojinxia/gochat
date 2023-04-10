package server

type Server interface {
	Init() error
	Close() error
}

var serverManager = make([]Server, 0)

func RegisterServer(sv Server) {
	serverManager = append(serverManager, sv)
}

func InitServers() {
	listener := NewFileListener()
	go listener.Init()
	RegisterServer(listener)
	redisPool = NewRedisServer()
	redisPool.Init()
	RegisterServer(redisPool)
}

func CloseAllServer() {
	for _, v := range serverManager {
		v.Close()
	}
}
