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
	"strings"

	"github.com/go-hello-world-web/internal/domain"
	"github.com/go-hello-world-web/internal/erro"

	"github.com/gorilla/mux"
	"github.com/golang-jwt/jwt/v4"

)

var (
	my_secret_loaded_from_volume Secret
	my_info_pod	InfoPod
	PORT = 3000
	API_VERSION = "no-version"
	POD_NAME = "pod no-name"
	POD_PATH = ""
	jwtKey = []byte("my_secret_key")
)

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

func MethodToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> methodToken => " + r.Method + " => path:  " + r.URL.Path)
	fmt.Fprintf(w, "<h1>method com JWT Token  : " + POD_NAME + " ver: " + API_VERSION + " </h1>")
	return
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

func header(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> header => " + r.Method + " => path:  " + r.URL.Path)
	//Iterate over all header fields
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
	}
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
	fmt.Println("==> setInfoPod")
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

func MiddleWareHandlerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-------------------------------------------- \n")
		log.Println("MiddleWareHandlerToken (INICIO)")

		token := r.Header.Get("Authorization")
		tokenSlice := strings.Split(token, " ")
		var bearerToken string

		if len(tokenSlice) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(erro.ErrStatusUnauthorized.Error())
			return
		}
		
		bearerToken = tokenSlice[len(tokenSlice)-1]
		res, err := ScopeValidation(bearerToken,"","")
		if err != nil || res != true {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		log.Println("MiddleWareHandlerToken (FIM)")
		log.Printf("-------------------------------------------- \n")
		
		next.ServeHTTP(w, r)
	})
}

func ScopeValidation(token string, path string, method string) (bool, error){
	fmt.Println("==> ScopeValidation :" , token )

	claims := &domain.JwtData{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, erro.ErrStatusUnauthorized
		}
		return false, erro.ErrTokenExpired
	}

	if !tkn.Valid {
		return false, erro.ErrStatusUnauthorized
	}

	return true, nil
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

	if os.Getenv("JWTKEY") !=  "" {
		jwtKey = []byte(os.Getenv("JWTKEY"))
	}

	setInfoPod()
}

func main() {
	log.Printf("Starting !")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc(POD_PATH +"/index", index)
	myRouter.HandleFunc(POD_PATH +"/", index)
	
	methodToken := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	methodToken.HandleFunc(POD_PATH +"/methodToken", MethodToken)
	methodToken.Use(MiddleWareHandlerToken)

	myRouter.HandleFunc(POD_PATH +"/header", header)
	myRouter.HandleFunc(POD_PATH +"/version", version )
	myRouter.HandleFunc(POD_PATH +"/info", infoPod )
	myRouter.HandleFunc("/health", check)
	myRouter.HandleFunc("/live", live)
	myRouter.HandleFunc(POD_PATH + "/a", methodA)
	myRouter.HandleFunc(POD_PATH + "/b", methodB)
	myRouter.HandleFunc(POD_PATH + "/sum/{data}", sum)

	fmt.Println("Server starting...", PORT)

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