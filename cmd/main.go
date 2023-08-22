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
var API_VERSION = "no-version"
var POD_NAME = "pod no-name"
var POD_PATH = "/pod-b"
	
type Secret struct {
	Username		string	`json:username`
	Password		string	`json:"password"`
}

type InfoPod struct {
	POD_NAME			string `json:"pod_name"`
	API_VERSION			string `json:"version"`
	OSPID				string `json:"os_pid"`
	IpAdress			string `json:"ip_address"`
}

func methodA(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> methodA => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>methodA : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
}

func methodB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> methodB => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>methodB : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> index => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>Hello World Go Web : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> check => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>Health check : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
}

func live(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> live => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>Live check: " + POD_NAME + " ver: " + API_VERSION + " </h1>")
}

func infoPod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> infoPod => " + r.Method + " => path:  " + r.URL.Path)
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

	if os.Getenv("POD_NAME") !=  "" {
		POD_NAME = os.Getenv("POD_NAME")
	}
	my_info_pod.POD_NAME = POD_NAME

	if os.Getenv("POD_PATH") !=  "" {
		POD_PATH = os.Getenv("POD_PATH")
	}
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

	http.HandleFunc(POD_PATH + "/", index)
	http.HandleFunc("/health", check)
	http.HandleFunc("/live", live)
	http.HandleFunc(POD_PATH + "/info", infoPod)
	http.HandleFunc(POD_PATH + "/a", methodA)
	http.HandleFunc(POD_PATH + "/b", methodB)

	fmt.Println("Server starting...", PORT)
	fmt.Println("my_secret_loaded_from_volume : ", my_secret_loaded_from_volume)
	fmt.Println("my_info_pod : ", my_info_pod)

	http.ListenAndServe(fmt.Sprintf("%s%d", ":", PORT), nil)
}