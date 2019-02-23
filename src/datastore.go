package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func dataStore(hash string) string {
	var Payload string

	db, err := sql.Open("mysql", AppDb)
	if err != nil {
		log.Print("Open db err: ")
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Print("Ping db err")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppDbNoSQL, AppDbNoSQLPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	hexStr, err := client.Get(hash).Result()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS demoTable (id INT NOT NULL AUTO_INCREMENT, token VARCHAR(100), text TEXT, PRIMARY KEY(id))")
	_, err = db.Exec("insert into demoTable values(null,?,?)", hash, hexStr)

	err = db.QueryRow("SELECT text FROM demoTable WHERE token = ?", hash).Scan(&Payload) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	decoded, err := hex.DecodeString(Payload)
	return string(decoded)
}
