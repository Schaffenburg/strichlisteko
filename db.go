package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"fmt"
	"log"
	"sync"
	"time"
)

var DB *sql.DB
var _DBonce sync.Once

func OpenDB() {
	_DBonce.Do(openDB)
}

func init() {
	OpenDB()
}

func openDB() {
	conf := GetConfig()

	var err error
	DB, err = sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatalf("Failed to open DB: %s", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS products (
	    id INT AUTO_INCREMENT PRIMARY KEY,
	    name TEXT NOT NULL,
	    stock INT NOT NULL,
	    EAN VARCHAR(255) NOT NULL,
	    price INT NOT NULL,
	    box_size INT NOT NULL,
	    amount VARCHAR(10) NOT NULL,
	    image UUID NOT NULL,
	    note TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to create table products: %s", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
	    userid INT AUTO_INCREMENT PRIMARY KEY,
	    username VARCHAR(255) UNIQUE KEY,
	    image UUID NOT NULL,
	    balance INT NOT NULL,
	    active BOOL NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to create table users: %s", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS transactions (
	    transaction_id INT AUTO_INCREMENT PRIMARY KEY,
	    value INT NOT NULL,
	    product VARCHAR(255) NOT NULL,
	    user_id INT NOT NULL,
	    time TIMESTAMP NOT NULL,
	    undone BOOL NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to create table transactions: %s", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS images (
	    image_id UUID PRIMARY KEY,
	    image MEDIUMBLOB NOT NULL,
	    mime VARCHAR(25) NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to create table images: %s", err)
	}
}

func getImage(id uuid.UUID) (image []byte, mime string, err error) {
	OpenDB()

	r, err := DB.Query("SELECT image, mime FROM images WHERE image_id = ?", id)
	if err != nil {
		log.Printf("failed to read image %s: %s", id, err)
		return nil, "", err
	}

	defer r.Close()

	if !r.Next() {
		return nil, "", sql.ErrNoRows
	}

	return image, mime, r.Scan(&image, &mime)
}

func getUser(id int) (u User, err error) {
	OpenDB() // ensure DB is opened

	r, err := DB.Query("SELECT userid, username, image, balance, active FROM users WHERE userid = ?", id)
	if err != nil {
		log.Printf("Error getting user %d: %s", id, err)
		return u, err
	}

	defer r.Close()

	if !r.Next() {
		return u, sql.ErrNoRows
	}

	return u, r.Scan(&u.ID, &u.Username, &u.Image, &u.Balance, &u.Active)
}

func setUser(u User) error {
	OpenDB()

	_, err := DB.Exec("UPDATE users SET username = ?, image = ?, balance = ?, active = ? WHERE userid = ?",
		u.Username, u.Image, u.Balance, u.Active, u.ID,
	)

	return err
}

func addProduct(prod Product) (int64, error) {
	OpenDB()

	r, err := DB.Exec("INSERT INTO products (name, stock, EAN, price, box_size, image, amount, note) VALUES (?, ?, ?, ? ,?, ?, ?, ?)",
		prod.Name, prod.Stock, prod.EAN, prod.Price, prod.BoxSize, prod.Image, prod.Amount, prod.Note,
	)
	if err != nil {
		return 0, err
	}

	return r.LastInsertId()
}

func setProduct(prod Product) error {
	OpenDB()

	_, err := DB.Exec("UPDATE products SET name = ?, stock = ?, EAN = ?, price = ?, box_size = ?, image = ? , amount = ?, note = ? WHERE id = ?",
		prod.Name, prod.Stock, prod.EAN, prod.Price, prod.BoxSize, prod.Image, prod.Amount, prod.Note, prod.ID,
	)

	return err
}

func deleteUser(id int) error {
	OpenDB()

	_, err := DB.Exec("DELETE FROM users WHERE userid = ?", id)
	if err != nil {
		return err
	}

	// delete all transactions
	_, err = DB.Exec("DELETE FROM transactions WHERE user_id = ?", id)

	return err
}

// returns id of created user
func addUser(username string, image uuid.UUID, balance int, active bool) (int, error) {
	OpenDB()

	r, err := DB.Exec(`INSERT INTO users (username, image, balance, active) VALUES (?, ?, ?, ?);`,
		username, image, balance, active,
	)
	if err != nil {
		return 0, err
	}

	i, err := r.LastInsertId()
	return int(i), err
}

func updateUsername(id int, username string) error {
	OpenDB()

	_, err := DB.Exec("UPDATE INTO users username = ? WHERE id = ?", username, id)

	return err
}

func updateBalance(id int, bal int) error {
	OpenDB()

	_, err := DB.Exec("UPDATE INTO users balance = ? WHERE id = ?", bal, id)

	return err
}

func delProduct(p int) error {
	OpenDB()

	_, err := DB.Exec("DELETE FROM products WHERE id = ?", p)

	return err
}

// limits to active users
func getUsers(active bool) (u []User, err error) {
	OpenDB() // ensure DB is opened

	r, err := DB.Query("SELECT userid, username, image, balance, active FROM users WHERE active = ?", active)
	if err != nil {
		log.Printf("Error getting products: %s", err)
		return
	}

	defer r.Close()

	var buf User
	for r.Next() {
		err = r.Scan(&buf.ID, &buf.Username, &buf.Image, &buf.Balance, &buf.Active)
		if err != nil {
			return nil, err
		}

		u = append(u, buf)
		buf = User{}
	}

	return u, nil

}

func getProducts() (p []Product, err error) {
	OpenDB() // ensure DB is opened

	r, err := DB.Query("SELECT id, name, stock, EAN, price, box_size, amount, image, note FROM products")
	if err != nil {
		log.Printf("Error getting products: %s", err)
	}

	defer r.Close()

	var buf Product
	for r.Next() {
		err = r.Scan(&buf.ID, &buf.Name, &buf.Stock, &buf.EAN, &buf.Price, &buf.BoxSize, &buf.Amount, &buf.Image, &buf.Note)
		if err != nil {
			return nil, err
		}

		p = append(p, buf)
		buf = Product{}
	}

	return p, nil
}

func getProduct(id int) (p Product, err error) {
	OpenDB() // ensure DB is opened

	r, err := DB.Query("SELECT id, name, stock, EAN, price, box_size, amount, image, note FROM products WHERE id = ?", id)
	if err != nil {
		log.Printf("Error getting product %d: %s", id, err)
	}

	defer r.Close()

	if !r.Next() {
		return p, sql.ErrNoRows
	}

	err = r.Scan(&p.ID, &p.Name, &p.Stock, &p.EAN, &p.Price, &p.BoxSize, &p.Amount, &p.Image, &p.Note)
	if err != nil {
		return p, err
	}

	return p, nil
}

func setImage(u uuid.UUID, mime string, i []byte) error {
	OpenDB()

	_, err := DB.Exec("INSERT INTO images (image_id, mime, image) VALUES (?, ?, ?);", u, mime, i)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Image    uuid.UUID `json:"image"`
	Balance  int       `json:"balance"`
	Active   bool      `json:"active"`
}

func (u User) BalanceNeg5() bool {
	return u.Balance < -500
}

type Product struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Stock   int       `json:"stock"`
	EAN     string    `json:"EAN"`
	Price   int       `json:"price"`
	BoxSize int       `json:"box_size"`
	Amount  string    `json:"amount"`
	Image   uuid.UUID `json:"image"`
	Note    string    `json:"note"`
}

func (p Product) PriceString() string {
	return fmt.Sprintf("%d.%02d €", p.Price/100, Abs(p.Price)%100)
}

func (u User) BalanceString() string {
	return fmt.Sprintf("%d.%02d €", u.Balance/100, Abs(u.Balance)%100)
}

func (u User) BalColor() string {
	if u.Balance < 0 {
		return "red"
	}

	return "green"
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func addTransaction(price int, name string, userid int) error {
	_, err := DB.Exec("INSERT INTO transactions (value, product, user_id, time, undone) VALUES (?, ?, ?, NOW(), false)",
		price, name, userid,
	)

	return err
}

type Transaction struct {
	ID      int
	Value   int
	Product string
	UserID  int
	Time    time.Time
	Undone  bool
}

func (t Transaction) ValueString() string {
	return fmt.Sprintf("%d.%02d €", t.Value/100, Abs(t.Value)%100)
}

func (t Transaction) TimeString() string {
	return t.Time.Format("Mon Jan _2 15:04:05 MST 2006")
}

func getTransaction(id int) (t Transaction, err error) {
	r, err := DB.Query("SELECT transaction_id, value, product, user_id, UNIX_TIMESTAMP(time), undone FROM transactions WHERE transaction_id = ?",
		id,
	)
	if err != nil {
		return
	}

	defer r.Close()

	if !r.Next() {
		err = sql.ErrNoRows

		return
	}

	var tbuf int64
	err = r.Scan(&t.ID, &t.Value, &t.Product, &t.UserID, &tbuf, &t.Undone)
	if err != nil {
		return
	}

	t.Time = time.Unix(tbuf, 0)

	return
}

func getTransactions(userid int) (t []Transaction, err error) {
	r, err := DB.Query("SELECT transaction_id, value, product, user_id, UNIX_TIMESTAMP(time), undone FROM transactions WHERE user_id = ? ORDER BY transaction_id DESC LIMIT 20",
		userid,
	)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	var buf Transaction
	var tbuf int64
	for r.Next() {
		err = r.Scan(&buf.ID, &buf.Value, &buf.Product, &buf.UserID, &tbuf, &buf.Undone)
		if err != nil {
			return nil, err
		}

		buf.Time = time.Unix(tbuf, 0)

		t = append(t, buf)
		buf = Transaction{}
	}

	return
}

func setTransactionUndo(id int, undone bool) error {
	_, err := DB.Exec("UPDATE transactions SET undone = ? WHERE transaction_id = ?",
		undone, id,
	)

	return err
}

func getImages() (s []uuid.UUID, err error) {
	r, err := DB.Query("SELECT image_id FROM images")
	if err != nil {
		return
	}

	defer r.Close()

	var buf uuid.UUID
	for r.Next() {
		err = r.Scan(&buf)
		if err != nil {
			return
		}

		s = append(s, buf)
		buf = uuid.UUID{}
	}

	return
}

func getProductByName(name string) (p Product, err error) {
	OpenDB() // ensure DB is opened

	r, err := DB.Query("SELECT id, name, stock, EAN, price, box_size, amount, image, note FROM products WHERE name = ?", name)
	if err != nil {
		return p, err
	}

	defer r.Close()

	if !r.Next() {
		return p, sql.ErrNoRows
	}

	err = r.Scan(&p.ID, &p.Name, &p.Stock, &p.EAN, &p.Price, &p.BoxSize, &p.Amount, &p.Image, &p.Note)
	if err != nil {
		return p, err
	}

	return p, nil
}
