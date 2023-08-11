package main

import (
	"fmt"
	"net"
	"strconv"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"log"
	"os"
)

var my_secret_loaded_from_volume Secret
var my_info_pod	InfoPod
var PORT = 3000
var API_VERSION = "0.0"
	
type Secret struct {
	Username		string	`json:username`
	Password		string	`json:"password"`
}

type InfoPod struct {
	API_VERSION			string `json:"version"`
	OSPID				string `json:"os_pid"`
	IpAdress			string `json:"ip_address"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> index")
	fmt.Fprintf(w, "<h1>Hello World Go Web 2.0</h1>")
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> check")
	fmt.Fprintf(w, "<h1>Health check</h1>")
}

func live(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> live")
	fmt.Fprintf(w, "<h1>Live check</h1>")
}

func infoPod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> infoPod")
	json.NewEncoder(w).Encode(my_info_pod)
}

func setInfoPod(){
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error to get the POD IP address !!!", err)
		os.Exit(3)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				my_info_pod.IpAdress = ipnet.IP.String()
			}
		}
	}
	my_info_pod.OSPID = strconv.Itoa(os.Getpid())

	if os.Getenv("API_VERSION") !=  "" {
		API_VERSION = os.Getenv("API_VERSION")
	}
	my_info_pod.API_VERSION = API_VERSION
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

	if os.Getenv("PORT") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("PORT"))
		PORT = intVar
	}

	setInfoPod()
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/health", check)
	http.HandleFunc("/live", live)
	http.HandleFunc("/info", infoPod)

	fmt.Println("Server starting...", PORT)
	fmt.Println("my_secret_loaded_from_volume : ", my_secret_loaded_from_volume)
	fmt.Println("my_info_pod : ", my_info_pod)

	http.ListenAndServe(fmt.Sprintf("%s%d", ":", PORT), nil)
}