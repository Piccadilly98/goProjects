package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

// type Person struct {
// 	Name string
// 	Age  int
// }

// var people []Person = []Person{
// 	{Name: "Тестовый", Age: 25},
// }

// func main() {
// 	http.HandleFunc("/people", peopleHandler)
// 	http.HandleFunc("/health", healthCheckHandler)
// 	log.Println("server start and listen in localhost, port: 8080")
// 	err := http.ListenAndServe("localhost:8080", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func peopleHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		getPeople(w, r)
// 	case http.MethodPost:
// 		postPerson(w, r)
// 	default:
// 		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
// 	}
// }

// func getPeople(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(people)
// }

// func postPerson(w http.ResponseWriter, r *http.Request) {
// 	var person Person
// 	err := json.NewDecoder(r.Body).Decode(&person)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	people = append(people, person)
// 	fmt.Println(w, "post new person add")
// }

// func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode("good job")
// }

// func hello(w http.ResponseWriter, r *http.Request) {
// 	str := "Hello, dear developer"
// 	b := []byte(str)
// 	_, err := w.Write(b)
// 	if err != nil {
// 		fmt.Println("Error:", err.Error())
// 	} else {
// 		fmt.Println("Correct request processing")
// 	}
// }

// func payCancel(w http.ResponseWriter, r *http.Request) {
// 	str := "Payed canceled"
// 	b := []byte(str)
// 	_, err := w.Write(b)
// 	if err != nil {
// 		fmt.Println("Error:", err.Error())
// 	} else {
// 		fmt.Println("Correct cancel operation")
// 	}
// }

// func pay(w http.ResponseWriter, r *http.Request) {
// 	str := "Pay completed"
// 	b := []byte(str)
// 	_, err := w.Write(b)
// 	if err != nil {
// 		fmt.Println("Error:", err.Error())
// 	} else {
// 		fmt.Println("Correct pay")
// 	}
// }

var money atomic.Int64
var bank atomic.Int64
var mtx sync.Mutex

func payHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error:", err)
		w.Write([]byte(err.Error()))
		return
	}

	str := string(body)
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Error processing request data"))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		return
	} else {
		fmt.Println("Get num:", num, "\nStart processing")
	}
	mtx.Lock()
	if money.Load() < int64(num) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Error!You balance < get number"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("I can't processing balance")
		return
	}
	money.Add(int64(-num))
	helloHandler(w, r)
	mtx.Unlock()
}

func bankHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	str := string(body)
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error processing request data"))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		return
	} else {
		fmt.Println("Get num:", num, "\nStart processing")
	}
	mtx.Lock()
	if money.Load() < int64(num) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Error!You balance < get number"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("I can't processing balance")
		return
	}
	money.Add(int64(-num))
	bank.Add(int64(num))
	helloHandler(w, r)
	mtx.Unlock()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Printf("k: %s --- v: %s\n", k, v)
	}
	fmt.Println(r.URL)
	b := []byte(fmt.Sprintf("Money: %v\nBank: %v\n", money.Load(), bank.Load()))
	for k, v := range r.Header {
		fmt.Printf("k: %s -- v: %s\n", k, v)
	}
	_, err := w.Write(b)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	} else {
		fmt.Println("good")
	}
}
func main() {
	money.Add(50)
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/bank", bankHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatal(err)
	}

}

//localhost:9091/default?foo=x&boo=y

// type fooBoo struct {
// 	Foo string `json:"foo"`
// 	Boo string `json:"boo"`
// }

// func defaultHendler(w http.ResponseWriter, r *http.Request) {
// 	var str fooBoo = fooBoo{}
// 	json.NewEncoder().Encode(&str)
// 	fmt.Println(str)
// }
// func main() {
// 	http.HandleFunc("/default", defaultHendler)
// 	err := http.ListenAndServe(":9091", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
