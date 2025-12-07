package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for i := 0; i < 1000; i++ {
		go func(i int) {
			resp, err := http.Get("http://localhost:8080/boards")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("запрос %d отправлен\n", i)
			defer resp.Body.Close()
		}(i)
	}

	time.Sleep(1 * time.Hour)
}
