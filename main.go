package main

import (
"avtoSell/db"
"avtoSell/router"
"log"
"net/http"
)

func main() {
	connection, err := db.Connect()
	if err != nil{
		return
	}

	log.Fatal(http.ListenAndServe(":8090", router.Init(connection)))
}
