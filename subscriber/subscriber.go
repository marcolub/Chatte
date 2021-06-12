package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main(){
	// create redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "mysecret",
		DB:0,
	})
	// ping the redis server and check for errors
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	ctx := context.Background()
	// subscribe to given topic
	topic := redisClient.Subscribe(ctx,"new_users")
	// get the channel to use
	channel := topic.Channel()
	// messages sent to the channel
	for msg := range channel {
		u := &User{}
		// unmarshal data
		err = u.UnmarshalBinary([]byte(msg.Payload))
		if err != nil{
			panic(err)
		}
		fmt.Println(u)
	}
} 

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

func (u *User) String() string {
	return "User: " + u.Username 
}