package main

import (
	"encoding/hex"
	"log"
)

func dataStore(hash string) string {
	var Payload string

	hexStr, err := CACHE.Get(hash).Result()
	if err != nil {
		log.Print(err)
	}
	log.Print(hexStr)

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS demoTable (id INT NOT NULL AUTO_INCREMENT, token VARCHAR(100), text TEXT, PRIMARY KEY(id))")
	_, err = DB.Exec("insert into demoTable values(null,?,?)", hash, hexStr)

	err = DB.QueryRow("SELECT text FROM demoTable WHERE token = ?", hash).Scan(&Payload) // WHERE number = 13
	if err != nil {
		log.Print(err) // proper error handling instead of panic in your app
	}
	decoded, err := hex.DecodeString(Payload)
	return string(decoded)
}
