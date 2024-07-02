package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

type Budget struct {
	Categories []Category
	Accounts   []Account
}

type Account struct {
	gorm.Model
	ID   int
	Name string `gorm:"unique"`
	Type string
}

type Category struct {
	gorm.Model
	ID   int
	Name string `gorm:"unique"`
	Due  int
}

type Payee struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Transaction struct {
	gorm.Model
	Date       time.Time
	AccountID  int
	Account    Account
	CategoryID int
	Category   Category
	Amount     float32
	PayeeID    int
	Payee      Payee
}

func (b *Budget) GetBudget() *Budget {
	y, err := os.ReadFile("budget.yaml")
	if err != nil {
		log.Fatalf("error reading budget yaml file %v ", err)
	}
	err = yaml.Unmarshal(y, b)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return b
}

func SetUpDatabase(b Budget) {
	var err error
	db, err = gorm.Open(sqlite.Open("gb.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	db.AutoMigrate(&Account{}, &Category{}, &Payee{}, &Transaction{})

	accounts := []Account{}

	for i := 0; i < len(b.GetBudget().Accounts); i++ {
		accounts = append(accounts, b.GetBudget().Accounts[i])
	}

	categories := []Category{}

	for i := 0; i < len(b.GetBudget().Categories); i++ {
		categories = append(categories, b.GetBudget().Categories[i])
	}

	ar := db.Clauses(clause.OnConflict{UpdateAll: true, Columns: []clause.Column{{Name: "name"}}}).
		Create(accounts)
	if ar.Error != nil {
		log.Fatalf("failed to create initial accounts %v", ar.Error)
	}

	cr := db.Clauses(clause.OnConflict{UpdateAll: true, Columns: []clause.Column{{Name: "name"}}}).
		Create(categories)
	if cr.Error != nil {
		log.Fatalf("failed to create initial categories %v", cr.Error)
	}
}

func ReturnJSON(w http.ResponseWriter, r []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(r)
}

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accounts = []Account{}

	db.Find(&accounts)
	a, err := json.Marshal(&accounts)
	if err != nil {
		log.Fatalf("failed to marshal json for accounts %v", err)
	}

	ReturnJSON(w, a)
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var categories = []Category{}

	db.Find(&categories)
	c, err := json.Marshal(&categories)
	if err != nil {
		log.Fatalf("failed to marshal json for categories %v", err)
	}

	ReturnJSON(w, c)
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		var tt struct {
			Payee    string
			Amount   float32
			Category string
			Account  string
			Date     time.Time
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("failed to read json body for transactions %v", err)
		}
		err = json.Unmarshal(body, &tt)
		if err != nil {
			log.Fatalf("failed to unmarshal json body for transactions %v", err)
		}

		// look up payee to assign
		var p = Payee{Name: tt.Payee}
		var c = Category{Name: tt.Category}
		var a = Account{Name: tt.Account}

		db.FirstOrCreate(&p)
		db.First(&c)
		db.First(&a)

		var t = Transaction{
			Account:  a,
			Category: c,
			Payee:    p,
			Amount:   tt.Amount,
			Date:     tt.Date,
		}

		result := db.Create(&t)
		if result.Error != nil {
			log.Fatalf("failed to create new transaction %v", result.Error)
		}

		rt, err := json.Marshal(&t)
		if err != nil {
			log.Fatalf("failed to marshal json for transactions %v", err)
		}

		ReturnJSON(w, rt)

	case "GET":
		var transactions = []Transaction{}

		db.Joins("Account").Joins("Category").Joins("Payee").Find(&transactions)
		t, err := json.Marshal(&transactions)
		if err != nil {
			log.Fatalf("failed to marshal json for transactions %v", err)
		}

		ReturnJSON(w, t)
	}
}

func main() {
	var b Budget
	b.GetBudget()
	SetUpDatabase(b)

	r := mux.NewRouter()
	r.HandleFunc("/accounts", AccountsHandler)
	r.HandleFunc("/categories", CategoriesHandler)
	r.HandleFunc("/transactions", TransactionsHandler).Methods("GET", "POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
