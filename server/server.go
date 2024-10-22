package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Price struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/", PriceHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func PriceHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancelContextApi := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelContextApi()

	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return
	}
	defer db.Close()

	req, err := callApi(ctx)
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

	var data Price
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Printf("Error decoding API response: %v", err)
		return
	}

	err = CreatePrice(r.Context(), db)
	if err != nil {
		log.Printf("Error creating table in database: %v", err)
		return
	}

	err = InsertPrice(r.Context(), db, data)
	if err != nil {
		log.Printf("Error saving data to the database: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data.USDBRL.Bid)
}

func InsertPrice(ctx context.Context, db *sql.DB, data Price) error {
	ctx, cancelInsert := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancelInsert()

	_, err := db.ExecContext(ctx, `INSERT INTO Price (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		data.USDBRL.Code, data.USDBRL.Codein, data.USDBRL.Name, data.USDBRL.High, data.USDBRL.Low, data.USDBRL.VarBid, data.USDBRL.PctChange, data.USDBRL.Bid, data.USDBRL.Ask, data.USDBRL.Timestamp, data.USDBRL.CreateDate)

	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return err
	}

	return nil
}

func CreatePrice(ctx context.Context, db *sql.DB) error {
	ctx, cancelCreate := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancelCreate()

	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Price (
    	Code        VARCHAR(10),
    	Codein      VARCHAR(10),
    	Name        VARCHAR(100),
    	High        VARCHAR(20),
    	Low         VARCHAR(20),
    	VarBid      VARCHAR(20),
    	PctChange   VARCHAR(10),
    	Bid         VARCHAR(20),
    	Ask         VARCHAR(20),
    	Timestamp   VARCHAR(20),
    	create_date  VARCHAR(20)
	);`)

	if err != nil {
		log.Printf("Error create table: %v", err)
		return err
	}

	return nil
}

func callApi(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
