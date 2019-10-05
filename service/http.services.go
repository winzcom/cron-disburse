package service

import (
	"cron-disburse/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

type PayStackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		RecipientCode string `json:"recipient_code"`
	} `json:"data"`
}

func TransferRecipient(user model.User, l *model.Loan, db *gorm.DB) (error, PayStackResponse) {
	form := url.Values{
		"type":           {"nuban"},
		"name":           {user.FirstName},
		"account_number": {l.BankAccountNo},
		"bank_code":      {l.BankCode},
		"currency":       {"NGN"},
	}
	tR, _ := http.NewRequest("POST", "https://api.paystack.co/transferrecipient", strings.NewReader(form.Encode()))
	client := &http.Client{}
	tR.Header.Set("content-type", "application/x-www-form-urlencoded")
	tR.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	resp, err := client.Do(tR)

	if err != nil {
		return err, PayStackResponse{}
	}

	defer resp.Body.Close()

	p := make([]byte, 9999)

	var result PayStackResponse

	n, _ := resp.Body.Read(p)
	er := json.Unmarshal(p[:n], &result)

	if er != nil {
		fmt.Println("Printing the roor ", er)
	}

	if result.Status == false {
		return errors.New("Failed to do stuff"), result
	}
	return nil, result
}

func Transfer(recipientCode string, loan *model.Loan, user model.User, db *gorm.DB) (interface{}, interface{}) {
	form := url.Values{
		"source":    {"balance"},
		"reason":    {"disbursement"},
		"amount":    {strconv.Itoa(loan.Amount * 100)},
		"currency":  {"NGN"},
		"recipient": {recipientCode},
	}
	fmt.Println("transfer ", form.Encode())
	tR, _ := http.NewRequest("POST", "https://api.paystack.co/transfer", strings.NewReader(form.Encode()))
	client := &http.Client{}
	tR.Header.Set("content-type", "application/x-www-form-urlencoded")
	tR.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	resp, err := client.Do(tR)

	if err != nil {
		fmt.Println("error ooos", err)
		return err, PayStackResponse{}
	}

	defer resp.Body.Close()

	p := make([]byte, 9999)

	var result PayStackResponse

	n, _ := resp.Body.Read(p)
	er := json.Unmarshal(p[:n], &result)

	if er != nil {
		fmt.Println("Printing the roor ", er)
	}

	if result.Status == false {
		fmt.Println(result)
		return errors.New("Failed to do stuff"), result
	}

	loan.ApprovalStatus = 1
	loan.LoanDesc = "running"
	loan.UpdateLoan(db)

	return nil, nil
}

func StartTransfer(loan model.Loan, user model.User, db *gorm.DB) {
	_, f := TransferRecipient(user, &loan, db)

	if f.Status == true {
		Transfer(f.Data.RecipientCode, &loan, user, db)
	} else {
		fmt.Println("failed to tr")
	}
}
