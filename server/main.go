package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	ID      int
	Name    string `gorm:"unique"`
	Type    string
	Balance float32
}

type Category struct {
	gorm.Model
	ID        int
	Name      string `gorm:"unique"`
	Due       int
	Assigned  int
	Available int
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
	y, err := os.ReadFile("../data/budget.yaml")
	if err != nil {
		log.Fatalf("error reading budget yaml file %v ", err)
	}
	err = yaml.Unmarshal(y, b)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return b
}

func InitDB(b Budget) {
	var err error
	db, err = gorm.Open(sqlite.Open("../data/gb.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	db.AutoMigrate(&Account{}, &Category{}, &Payee{}, &Transaction{})

	accounts := []Account{{Balance: 0}}

	for i := 0; i < len(b.GetBudget().Accounts); i++ {
		accounts = append(accounts, b.GetBudget().Accounts[i])
	}

	categories := []Category{{Assigned: 0, Available: 0}}

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

func GetAccounts(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var accounts = []Account{}

	db.Find(&accounts)
	a, err := json.Marshal(&accounts)
	if err != nil {
		log.Fatalf("failed to marshal json for accounts %v", err)
	}

	ReturnJSON(w, a)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var categories = []Category{}

	db.Find(&categories)
	c, err := json.Marshal(&categories)
	if err != nil {
		log.Fatalf("failed to marshal json for categories %v", err)
	}

	ReturnJSON(w, c)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	pageSize := 100
	params := mux.Vars(r)
	page, _ := strconv.Atoi(params["page"])
	offset := (page - 1) * pageSize

	var transactions = []Transaction{}

	db.Offset(offset).Limit(pageSize).Joins("Account").Joins("Category").Joins("Payee").Find(&transactions)

	t, err := json.Marshal(&transactions)
	if err != nil {
		log.Fatalf("failed to marshal json for transactions %v", err)
	}

	ReturnJSON(w, t)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	var tt struct {
		Payee    string
		Amount   float32
		Category string
		Account  string
		Date     string
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
	p := Payee{Name: tt.Payee}
	c := Category{}
	a := Account{}

	// return an error if no account or category
	cres := db.Where("name = ?", tt.Category).First(&c)
	ares := db.Where("name = ?", tt.Account).First(&a)

	if cres.RowsAffected == 0 || ares.RowsAffected == 0 {
		log.Fatalf("Category or Account doesn't exist! These need to be defined in your configuration!")
	}

	// if we dont have a payee we can create it here
	db.Where("name = ?", tt.Payee).FirstOrCreate(&p)

	// format the date
	date, err := time.ParseInLocation("2006-01-02", tt.Date, time.Local)
	if err != nil {
		log.Fatalf("failed to parse date for transaction %v", err)
	}

	// save the transaction and return it
	var t = Transaction{
		Account:  a,
		Category: c,
		Payee:    p,
		Amount:   tt.Amount,
		Date:     date,
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
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	params := mux.Vars(r)
	t := Transaction{}
	db.Delete(&t, params["id"])
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // todo, make this more secure later
}

func main() {
	var b Budget
	b.GetBudget()

	InitDB(b)

	r := mux.NewRouter()

	r.HandleFunc("/accounts", GetAccounts).Methods("GET")
	r.HandleFunc("/categories", GetCategories).Methods("GET")
	r.HandleFunc("/transactions/{page}", GetTransactions).Methods("GET")
	r.HandleFunc("/transactions/create", CreateTransaction).Methods("POST")
	r.HandleFunc("/transactions/delete/{id}", DeleteTransaction).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Go Budget server running on port 8000")

	log.Fatal(srv.ListenAndServe())
}
