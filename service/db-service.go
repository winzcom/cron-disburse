package service

import (
	"cron-disburse/model"
	"fmt"
	"os"
	"sync"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open("mysql", "acreditbb:N0tju5tm3*&^@(stagingdb.cluster-ce681ga45xey.eu-west-1.rds.amazonaws.com)/aella_money?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("error connecting to the DB")
	}

	db.LogMode(true)

	return db
	// var (
	// 	loans []model.Loan
	// 	user  model.User
	// )
	//query := db.Table("loan_requests").Model(&loans)

	//fmt.Println(query)
}

func FindUser(loan model.Loan, db *gorm.DB) model.User {
	var user model.User
	db.Find(&user, "id=?", loan.UserID)
	return user
}

func CardReuseable(cards []model.Card) bool {
	for _, card := range cards {
		date := time.Date(card.ExpYear, time.Month(card.ExpMonth), 30, 00, 00, 00, 00, time.UTC)
		diff := date.Sub(time.Now())
		hours := diff.Hours()
		fmt.Println("hours ", hours)
		if card.Reusable == 1 && hours > 1460 {
			return true
		}
	}

	return false
}

func FindCards(userId int, db *gorm.DB) []model.Card {
	var cards []model.Card
	db.Find(&cards, "user_id=? and reusable=?", userId, 1)
	return cards
}

func ApprovedLoan(db *gorm.DB, c chan []model.Loan, allLoans chan bool, more chan bool) {
	var loans []model.Loan
	fmt.Println("paystack key ", os.Getenv("PAYSTACK_SECRET_KEY"))
	t := time.Now()
	fmt.Println("second ", t.Second())
	t1 := t.Add(time.Hour * 24 * -3)
	for i := 0; i < 2; i++ {
		format := t1.Format("2006-01-02")
		db.Find(&loans, "approval_status=? and created_at >= ?", 2, format).Limit(2)
		c <- loans
		<-more
		time.Sleep(1 * time.Second)
	}
	fmt.Println("finished ", time.Now().Second())
	allLoans <- true
}

func Run(db *gorm.DB, done chan<- bool, cha chan []model.Loan) {
	allLoans := make(chan bool, 1)
	more := make(chan bool, 1)
	go ApprovedLoan(db, cha, allLoans, more)
	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case ls := <-cha:
				for _, l := range ls {
					wg.Add(1)
					go func(l model.Loan) {
						fmt.Println(l.LoanDesc)
						user := FindUser(l, db)
						cards := FindCards(user.ID, db)
						if !CardReuseable(cards) {
							return
						}
						StartTransfer(l, user, db)
						wg.Done()
					}(l)
				}
				wg.Wait()
				more <- true
			case <-allLoans:
				fmt.Println("all done returning")
				done <- true
				db.Close()
				return
			}
		}
	}()

}
