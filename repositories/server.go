package repositories

type Middleware interface {
	Init() error
	Close() error
}

var serverManager =  make([]Middleware,0)

func RegisterServer(sv Middleware){
	serverManager = append(serverManager, sv)
}

func CloseAllServer(){
	for _,v := range serverManager{
		v.Close()
	}
}