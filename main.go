package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var counter = -1
var mutex = &sync.Mutex{}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	logger := log.New(os.Stdout, "counter ", log.LstdFlags)
	logger.Println(string(counter))

	if counter == -1 {
		var saved, _ = ioutil.ReadFile("/tmp/nowthis_counter.log")
		counter, _ = strconv.Atoi(string(saved))
	}
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
	ioutil.WriteFile("/tmp/nowthis_counter.log", []byte(fmt.Sprintf("%d\n", counter)), 0644)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger := log.New(os.Stdout, "web", log.LstdFlags)
		logger.Println(r.URL.Path)

		var file = "output/home.html"
		if string(r.URL.Path) != "/" {
			file = "output/" + r.URL.Path
		}
		http.ServeFile(w, r, file)
		//		incrementCounter(w, r)

		//		content, err := ioutil.ReadFile("output/home.html")
		//		if err == nil {
		//			fmt.Fprintf(w, string(content))
		//		}
	})

	http.HandleFunc("/visits", incrementCounter)
	//	fs := http.FileServer(http.Dir("output/"))
	//	http.Handle("output/", http.StripPrefix("output/", fs))

	http.ListenAndServe(":8080", nil)

	//	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/home/dbajet/go/src/nowthisio/output"))))

}
