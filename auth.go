package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var key []byte

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	p, err := rand.Prime(rand.Reader, 256)
	if err != nil {
		return key, err
	}

	p.FillBytes(key)
	return key, nil
}

func encrypt(payload []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, payload, nil)
	fmt.Printf("%x\n", ciphertext)
	return ciphertext, nil
}

func decrypt(payload []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(payload) < nonceSize {
		return errors.New("payload is less than nonce size")
	}

	nonce, payload := payload[:nonceSize], payload[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return err
	}
	fmt.Printf("%q\n", plaintext)
	return nil
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var c client

	if c.readCookie(w, r) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Name = r.FormValue("name")
		c.pass = r.FormValue("pass")
		action := r.FormValue("action")
		fmt.Println(c.Name)

		if action == "login" {
			if err := c.login(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if action == "register" {
			if c.exist() {
				http.Error(w, "username taken", http.StatusInternalServerError)
				return
			}
			if err := c.insert(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/auth", http.StatusFound)
	}

	if err := templates.ExecuteTemplate(w, "auth.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
