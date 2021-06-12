package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	rdb *redis.Client
)

var clients = make(map[*websocket.Conn]bool) // list of current active clients (opened WebSockets) 
var broadcaster = make(chan ChatMessage) // channel to send/receive ChatMessage data structure
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request)bool{
		return true
	},
}

type User struct {
	Username string
}

type ChatMessage struct {
	Username string `json:"username"`
	Text string `json:"text"`
}

// marshal and unmarshal data 
func (u *User) MarshalBinary()([]byte, error){
	return json.Marshal(u)
}
func (u *User) UnmarshalBinary(data []byte) error{
	if err := json.Unmarshal(data,&u);err != nil{
		return err
	}
	return nil
}

// random names
var Names []string = []string{"Jasper", "Johan", "Edward", "Niel", "Percy", "Adam", "Grape", "Sam", "Redis", "Jennifer", "Jessica", "Angelica", "Amber", "Watch"}
var SirNames []string = []string{"Ericsson", "Redisson", "Edisson", "Tesla", "Bolmer", "Andersson", "Sword", "Fish", "Coder"}

func main(){
	// generate a new redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:"redis:6379",
		Password:"mysecret",
		DB:0,
	})
	// ping redis client and check for errors
	err := rdb.Ping(context.Background()).Err()
	if err != nil{
		time.Sleep(3*time.Second)
		err := rdb.Ping(context.Background()).Err()
		if err != nil{
			panic(err)
		}
	}
	// serve the static page and handle messages
	http.Handle("/",http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket",HandleConnections)
	go HandleMessages()
	log.Print("Server listening on localhost:4444")
	if err := http.ListenAndServe(":4444",nil); err != nil {
		log.Fatal(err)
	}
}

// function for random user
func RandomUser() *User{
	rand.Seed(time.Now().UnixNano())
	nameMax := len(Names)
	sirNameMax := len(SirNames)

	nameIndex:= rand.Intn(nameMax-1)+1
	sirNameIndex:= rand.Intn(sirNameMax-1)+1

	return &User{
		Username: Names[nameIndex]+" "+SirNames[sirNameIndex],
	}
}

func HandleConnections(w http.ResponseWriter,r *http.Request){
	//--- set up receiving messages from other clients
	ws, err := upgrader.Upgrade(w,r,nil)
	if err != nil{
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true

	// generate random username
	myusername := RandomUser().Username
	publish(myusername)
	
	//--- handle sending messages
	for{
		var msg ChatMessage
		// read a new message as JSON and map it to a ChatMessage object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients,ws)
			break
		}
		msg.Username = myusername
		// send a new message to the channel
		broadcaster <- msg
	}
}

func publish(username string){
	ctx := context.Background()
	// publish generated user to the new_users channel
	err := rdb.Publish(ctx, "new_users", username).Err()
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(4)
	time.Sleep(time.Duration(n) * time.Second)
	
}

func HandleMessages(){
	for{
		// grab next message from channel
		msg := <- broadcaster

		storeInRedis(msg)
		messageClients(msg)
	}
}

func storeInRedis(msg ChatMessage){
	json,err := json.Marshal(msg)
	if err != nil{
		panic(err)
	}
	ctx := context.Background()
	if err := rdb.RPush(ctx,"chat_messages",json).Err(); err != nil{
		panic(err)
	}
}

func messageClients(msg ChatMessage){
	// send to every connected client
	for client := range clients{
		messageClient(client,msg)
	}
}

func messageClient(client *websocket.Conn,msg ChatMessage){
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v",err)
		client.Close()
		delete(clients,client)
	}
}