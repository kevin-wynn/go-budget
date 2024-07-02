package main

import (
	"encoding/json"
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
	Amount     int
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

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accounts = []Account{}

	db.Find(&accounts)
	json.NewEncoder(w).Encode(&accounts)
	a, err := json.Marshal(&accounts)
	if err != nil {
		log.Fatalf("failed to marshal json for accounts %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(a)
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var categories = []Category{}

	db.Find(&categories)
	json.NewEncoder(w).Encode(&categories)
	c, err := json.Marshal(&categories)
	if err != nil {
		log.Fatalf("failed to marshal json for categories %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(c)
}

func main() {
	var b Budget
	b.GetBudget()
	SetUpDatabase(b)

	r := mux.NewRouter()
	r.HandleFunc("/accounts", AccountsHandler)
	r.HandleFunc("/categories", CategoriesHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
