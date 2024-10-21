package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Price struct {
	Bid string
}

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		log.Printf("Error creating API request: %v", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making API request: %v", err)
		return
	}
	defer res.Body.Close()

	var data string
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Printf("Error decoding API response: %v", err)
		return
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Printf("Error generating file: %v", err)
		return
	}

	_, err = f.Write([]byte("Dólar: " + data))
	if err != nil {
		log.Printf("Error writing file: %v", err)
		return
	}

	file, err := os.ReadFile("cotacao.txt")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}
	fmt.Println(string(file))
}
