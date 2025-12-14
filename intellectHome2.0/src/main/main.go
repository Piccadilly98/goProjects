package main

import (
	"log"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	_ "github.com/lib/pq"
)

func main() {
	serv, err := server.NewServer(false, 10, 150, true)
	if err != nil {
		log.Fatal(err)
	}
	serv.Start("localhost:8080")
	err = <-serv.ErrorServerChan
	log.Fatal(err)

}
