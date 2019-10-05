package main

import (
	"cron-disburse/model"
	"cron-disburse/service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/disburse", func(w http.ResponseWriter, r *http.Request) {

		db := service.Connect()

		cha := make(chan []model.Loan)

		done := make(chan bool, 1)

		w.Write([]byte("Running"))
		go service.Run(db, done, cha)
		return
	})

	http.ListenAndServe(":9090", nil)
}
