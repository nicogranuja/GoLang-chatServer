package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"encoding/json"
)

type Response struct {
	Message string `json:"message"`
	Status bool `json:"status"`
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte ("Hello World!"))
}

func HelloJson(w http.ResponseWriter, r *http.Request) {
	response := CreateResponse()
	json.NewEncoder(w).Encode(response)
}

func CreateResponse() Response{
	return Response { "This is a json format" , true}
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/Hello", HelloWorld).Methods("GET")
	mux.HandleFunc("/HelloJson", HelloJson).Methods("GET")

	http.Handle("/", mux)
	log.Println("Server is on the port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}