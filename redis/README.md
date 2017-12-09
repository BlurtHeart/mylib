RedisClient
==============

RedisClient is a redis client written in golang. It can help you operate the redis database.

## How to install?

	go get github.com/blurty/mylib/redis
	
## How to use it?
First you must import it

	import (
		"github.com/blurty/mylib/redis"
	)
	
Then in your code you can type it like this:

    package main

    import (
        "fmt"
        "time"
        "github.com/blurty/mylib/redis"
    )

    func main() {
            client := &redis.RedisClient{
            Addr:     "127.0.0.1:6379",
            Password: "",
            DB:       0,
        }
        client.Connect()

        // value is list
        key := "user13"
        value := "test1"
        value2 := "test2"

        client.SetList(key, value)
        client.SetList(key, value2)
        duration, _ := time.ParseDuration("10s")
        client.SetExpires(key, duration)

        // fuzzy query
        keys, err := client.GetKeys("user*")
        if err != nil {
            panic(err)
        }
        for _, k := range keys {
            v, err := client.GetList(k, 0, -1)
            if err != nil {
                panic(err)
            }
        fmt.Println("key:", k, "value:", v)
        }
    }