package main

import (
	"fmt"
	"os"
	"log"

	"github.com/Ahmed0427/rssy/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not in the environment")	
	}
	
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	fmt.Println(cfg)
}
