package main

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
)

type client struct {
	id           int
	Name         string
	pass         string
	Stocks       []stock
	Transactions []transaction
}

// Insert values into database. This inserts the password as plain-text. Do NOT
// do this in a production setting :). This is for debugging purposes.
func (c *client) insert() error {
	if c.Name == "" {
		return errors.New("Name is not set")
	} else if c.pass == "" {
		return errors.New("pass is not set")
	}

	_, err := db.Exec("INSERT INTO client (name, pass) VALUES ($1, $2)", c.Name, c.pass)
	if err != nil {
		return err
	}

	return nil
}

// Check if in database.
func (c *client) exist() bool {
	err := db.QueryRow("SELECT id FROM client WHERE name = $1", c.Name).Scan(&c.id)
	if err != nil {
		return false
	}

	return true
}

func (c *client) passCorrect() error {
	var pass string

	err := db.QueryRow("SELECT id, pass FROM client WHERE name = $1", c.Name).Scan(&c.id, &pass)
	if err != nil {
		return err
	}

	if c.pass != pass {
		return errors.New("invalid pass")
	}

	return nil
}

// If used in production, it might be preferably to create some sort of cookie
// session manager that manages cookies more securely. This function currently
// encrypts the user's ID and uses the runtime key for encryption with AES-256.
func (c *client) login(w http.ResponseWriter) error {
	if err := c.passCorrect(); err != nil {
		return err
	}

	cipher, err := encrypt([]byte(strconv.Itoa(c.id)))
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "key",
		Value:  hex.EncodeToString(cipher),
		MaxAge: 86400, // 24 hours from now
		Secure: true,
	})

	return nil
}

// Fetch name or ID
func (c *client) update() (err error) {
	if c.Name == "" && c.id == 0 {
		return errors.New("Name and id are not set")
	}

	if c.Name != "" {
		err = db.QueryRow("SELECT id FROM client WHERE name = $1", c.Name).Scan(&c.id)
	} else if c.id != 0 {
		err = db.QueryRow("SELECT name FROM client WHERE id = $1", c.id).Scan(&c.Name)
	}

	if err != nil {
		return err
	}

	return nil
}

// Reads id from the cookie set by us.
func (c *client) readCookie(w http.ResponseWriter, r *http.Request) error {
	crypt, err := r.Cookie("key")
	if err != nil {
		return err
	}

	cipher, err := hex.DecodeString(crypt.Value)
	if err != nil {
		return err
	}

	plain, err := decrypt(cipher)
	if err != nil {
		return err
	}

	c.id, err = strconv.Atoi(string(plain))
	if err != nil {
		return err
	}

	if err := c.update(); err != nil {
		return err
	}

	return nil
}

func (c *client) updateStocks() error {
	if c.id == 0 {
		return errors.New("client: id is not set")
	}

	rows, err := db.Query("SELECT symbol, quantity FROM stock WHERE client_id = $1", c.id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var s stock
		err := rows.Scan(&s.Symbol, &s.Quantity)
		if err != nil {
			return err
		}
		s.getPrice()
		c.Stocks = append(c.Stocks, s)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (c *client) updateTransactions() error {
	if c.id == 0 {
		return errors.New("client: id is not set")
	}

	rows, err := db.Query(`
		SELECT cat_id, amount, balance, description, time FROM transaction
		WHERE client_id = $1 ORDER BY time DESC LIMIT 20`, c.id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var t transaction
		err := rows.Scan(&t.Cat_id, &t.Amount, &t.Balance, &t.Description, &t.Time)
		if err != nil {
			return err
		}
		c.Transactions = append(c.Transactions, t)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
