package init_test_server

import (
	"log"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	_ "github.com/lib/pq"
)

func InitTestServer(workers bool) (*server.Server, error) {

	serv, err := server.NewServer(true, 30, 150, workers, 0, 0)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return serv, nil
}

func GetPtrStr(str string) *string {
	return &str
}
