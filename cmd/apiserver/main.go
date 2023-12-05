package main

import (
	"flag"
	"log"
	"mjcomparer/internal/app/apiserver"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
