package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	path := "./demo.sqlite"
	_ = os.Remove(path)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	must(db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT,
			age INTEGER,
			active BOOLEAN,
			balance REAL,
			created_at TEXT
		);
	`))
	must(db.Exec(`
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			product TEXT,
			amount REAL,
			status TEXT,
			placed_at TEXT
		);
	`))
	must(db.Exec(`
		CREATE TABLE products (
			sku TEXT PRIMARY KEY,
			name TEXT,
			price REAL,
			stock INTEGER
		);
	`))

	rand.Seed(1)

	names := []string{"alice", "bob", "carol", "dave", "eve", "frank", "grace", "heidi", "ivan", "judy", "mallory", "oscar", "peggy", "trent", "victor", "walter"}
	statuses := []string{"pending", "paid", "shipped", "refunded", "cancelled"}
	products := []string{"keyboard", "mouse", "monitor", "headphones", "webcam", "desk", "chair", "lamp", "mug", "notebook"}

	for i, n := range names {
		email := sql.NullString{}
		if i%4 != 0 {
			email = sql.NullString{String: n + "@example.com", Valid: true}
		}
		_, err := db.Exec(
			`INSERT INTO users (name, email, age, active, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			n, email, 18+rand.Intn(50), i%3 != 0, float64(rand.Intn(100000))/100,
			time.Now().Add(-time.Duration(i*24)*time.Hour).Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 200; i++ {
		_, err := db.Exec(
			`INSERT INTO orders (user_id, product, amount, status, placed_at) VALUES (?, ?, ?, ?, ?)`,
			1+rand.Intn(len(names)),
			products[rand.Intn(len(products))],
			float64(rand.Intn(50000))/100,
			statuses[rand.Intn(len(statuses))],
			time.Now().Add(-time.Duration(rand.Intn(720))*time.Hour).Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			panic(err)
		}
	}

	for _, p := range products {
		_, err := db.Exec(
			`INSERT INTO products (sku, name, price, stock) VALUES (?, ?, ?, ?)`,
			fmt.Sprintf("SKU-%s", p[:3]), p, float64(rand.Intn(20000))/100, rand.Intn(500),
		)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("seeded ./demo.sqlite")
}

func must(_ sql.Result, err error) {
	if err != nil {
		panic(err)
	}
}
