package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type User struct {
	Username string
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:"redis:6379",
		Password:"mysecret",
		DB:0,
	})
	// ping redis client and check for errors
	err := redisClient.Ping(context.Background()).Err()
	if err != nil{
		time.Sleep(3*time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil{
			panic(err)
		}
	}
	// ctx := context.Background()
	// for {
	// 	// publish generated user to the new_users channel
	// 	err := redisClient.Publish(ctx, "new_users", RandomUser()).Err()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	rand.Seed(time.Now().UnixNano())
	// 	n := rand.Intn(4)
	// 	time.Sleep(time.Duration(n) * time.Second)
	// }
	r := mux.NewRouter()
	r.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		fmt.Println("Hello World!")
	})
	log.Print("Server starting at localhost:4444")
	http.ListenAndServe(":4444",r)
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