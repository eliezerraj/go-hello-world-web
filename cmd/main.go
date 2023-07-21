package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"os"
)

var my_secret_loaded_from_volume Secret

type Secret struct {
	Username		string	`json:username`
	Password		string	`json:"password"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Health check</h1>")
}

func live(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Live check</h1>")
}

func init() {
	file_user, err := ioutil.ReadFile("/var/go-hello-world-web/secret/username")
    if err != nil {
        log.Fatal(err)
		os.Exit(3)
    }
	file_pass, err := ioutil.ReadFile("/var/go-hello-world-web/secret/password")
    if err != nil {
        log.Fatal(err)
		os.Exit(3)
    }
	my_secret_loaded_from_volume.Username = string(file_user)
	my_secret_loaded_from_volume.Password = string(file_pass)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/health", check)
	http.HandleFunc("/live", live)
	fmt.Println("Server starting...")

	fmt.Println("my_secret_loaded_from_volume : ", my_secret_loaded_from_volume)

	http.ListenAndServe(":3000", nil)
}