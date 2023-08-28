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
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

)

var my_secret_loaded_from_volume Secret
var my_info_pod	InfoPod
var PORT = 3000
var API_VERSION = "no-version"
var POD_NAME = "pod no-name"
var POD_PATH = ""
	
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
	return
}

func methodB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> methodB => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>methodB : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> index => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>Hello World Go Web : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
	return
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> check => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>Health check : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
	return
}

func live(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> live => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>Live check: " + POD_NAME + " ver: " + API_VERSION + " </h1>")
	return
}

func infoPod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> infoPod => " + r.Method + " => path:  " + r.URL.Path)
	json.NewEncoder(w).Encode(my_info_pod)
	return
}

func sum(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> sum => " + r.Method + " => path:  " + r.URL.Path)
	
	vars := mux.Vars(r)
	intVar, err := strconv.Atoi(vars["data"])
	if err != nil{
		log.Printf("ERRO => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("erro.ErrConvertion")
		return
	}
	intVar++
	json.NewEncoder(w).Encode(intVar)
	return
}

func version(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> version => " + r.Method + " => path:  " + r.URL.Path)
	json.NewEncoder(w).Encode(API_VERSION)
	return
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
	log.Printf("Starting !")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc(POD_PATH +"/index", index )
	myRouter.HandleFunc(POD_PATH +"/", index )
	myRouter.HandleFunc(POD_PATH +"/version", version )
	myRouter.HandleFunc(POD_PATH +"/info", infoPod )
	myRouter.HandleFunc("/health", check)
	myRouter.HandleFunc("/live", live)
	myRouter.HandleFunc(POD_PATH + "/a", methodA)
	myRouter.HandleFunc(POD_PATH + "/b", methodB)
	myRouter.HandleFunc(POD_PATH + "/sum/{data}", sum)

	fmt.Println("Server starting...", PORT)
	fmt.Println("my_secret_loaded_from_volume : ", my_secret_loaded_from_volume)
	fmt.Println("my_info_pod : ", my_info_pod)

	http_srv := http.Server{
		Addr:		":" +  strconv.Itoa(PORT),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(5) * time.Second,   
		WriteTimeout: time.Duration(5) * time.Second,  
		IdleTimeout:  time.Duration(5) * time.Second, 
	}

	go func() {
		err := http_srv.ListenAndServe()
		if err != nil {
			log.Print("message ==> ", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	log.Printf("Stopping Server")
	ctx , cancel := context.WithTimeout(context.Background(), time.Duration(5) * time.Second)
	defer cancel()

	if err := http_srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Print("WARNING Dirty Shutdown", err)
		return
	}

	log.Printf("Stop Done !")
}