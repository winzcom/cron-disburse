package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Loan struct {
	//gorm.Model
	ID                        int `gorm:"primary_key"`
	Amount                    int
	LoanTenor                 int
	LoanDesc                  string `gorm:"type:varchar(20)"`
	ApprovalStatus            int
	ToRepay                   int
	InterestRate              float32
	LoanPurpose               string `gorm:"type:varchar(30)"`
	BankName                  string `gorm:"type:varchar(30)"`
	BankCode                  string
	BankAccountNo             string  `gorm:"type:varchar(30)"`
	BankAccountValidated      uint8   `gorm:"type:int(1)"`
	IsBankAccountCorrect      uint8   `gorm:"type:int(1)"`
	CorrectBankAccountName    string  `gorm:"type:varchar(30)"`
	MonthlyRepayment          float32 `gorm:"type:float(10)"`
	TotalRepayment            float32 `gorm:"type:float(10)"`
	PrincipalMonthlyRepayment float32 `gorm:"type:float(10)"`
	MonthlyInterest           float32 `gorm:"type:float(10)"`
	PayDay                    uint8   `gorm:"type:int(2);column:payday"`
	LoanStarts                []uint8 `gorm:"type:timestamp"`
	LoanEnds                  []uint8 `gorm:"type:timestamp"`
	NextDueDate               []uint8 `gorm:"type:timestamp"`
	RepaymentMode             string  `gorm:"type:varchar(30)"`
	ApprovalOrRejectionNote   string  `gorm:"type:varchar(100)"`
	CreatedAt                 string  `gorm:"type:timestamp"`
	UpdatedAt                 string  `gorm:"type:timestamp"`
	User                      User
	UserID                    int
}

func (Loan) TableName() string {
	return "loan_requests"
}

func (l *Loan) UpdateLoan(db *gorm.DB) {
	fmt.Println("approval ", l.ApprovalStatus)
	db.Save(l)
}
