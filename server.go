package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"encoding/json"
	"sync"
)

type Response struct {
	Message string `json:"message"`
	Status int `json:"status"`
	IsValid bool `json:"isvalid"`
}

var Users = struct{
	m map [string] User
	sync.RWMutex
}{m: make(map[string] User)}

type User struct{
	username string
	WebSocket *websocket.Conn
}

func ToArrayByte(value string)[]byte {
	return []byte(value)
}

func ConcatMessage(username string, array[] byte)string {
	return username + " : " + string (array[:])
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte ("Hello World!"))
}

func HelloJson(w http.ResponseWriter, r *http.Request) {
	response := CreateResponse("This is json",200, true)
	json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string, status int, valid bool) Response{
	return Response { message , status , valid}
}

func LoadStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r, "./view/index.html")
}

func Validate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	response := Response{}

	if UserExists(username){
		response.IsValid = false
	}else{
		response.IsValid = true
	}
	json.NewEncoder(w).Encode(response)
}

func UserExists (username string) bool{
	Users.RLock()
	defer Users.RUnlock()

	if _, ok := Users.m[username]; ok {
		return true;
	}
	return false;
}

func CreateUser(username string, ws *websocket.Conn) User {
	return User{ username, ws}
}

func AddUser (user User) {
	Users.Lock()
	defer Users.Unlock()

	Users.m[user.username] = user
}

func RemoveUser(username string) {
	Users.Lock()
	defer Users.Unlock()
	delete(Users.m, username)
	log.Println("The user "+ username + " has left.")
}

func SendMessage(type_message int, message []byte) {
	Users.RLock()
	defer Users.RUnlock()

	for _,user := range Users.m{
		if err := user.WebSocket.WriteMessage(type_message, message); err != nil{
			return
		}
	}
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	ws, err := websocket.Upgrade(w,r,nil,1024,1024)
	if err != nil{
		log.Println(err)
		return
	}

	current_user := CreateUser(username, ws)
	AddUser(current_user)
	log.Println("New user added.")

	for{
		type_message, message, err := ws.ReadMessage()
		if err != nil{
			RemoveUser(username)
			return
		}
		finalMessage := ConcatMessage(username, message)
		SendMessage(type_message, ToArrayByte(finalMessage))
	}
}

func main() {

	cssHandle := http.FileServer(http.Dir("./view/css/"))
	jsHandle := http.FileServer(http.Dir("./view/js/"))

	mux := mux.NewRouter()
	mux.HandleFunc("/Hello", HelloWorld).Methods("GET")
	mux.HandleFunc("/HelloJson", HelloJson).Methods("GET")
	mux.HandleFunc("/Static", LoadStatic).Methods("GET")
	mux.HandleFunc("/Validate", Validate).Methods("POST")
	mux.HandleFunc("/Chat/{username}", WebSocket).Methods("GET")

	http.Handle("/", mux)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandle))
	log.Println("Server is on the port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}