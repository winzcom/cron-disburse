package main

import (
	"cron-disburse/model"
	"cron-disburse/service"
	"net/http"
)

func main() {

	http.HandleFunc("/disburse", func(w http.ResponseWriter, r *http.Request) {

		db := service.Connect()

		cha := make(chan []model.Loan)

		done := make(chan bool, 1)

		w.Write([]byte("Running"))

		//defer db.Close()

		defer func() {
			go service.Run(db, done, cha)
		}()
		return
	})

	http.ListenAndServe(":9090", nil)
}
