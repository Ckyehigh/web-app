package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/process", processHandler)
	http.HandleFunc("/generate-key", generateKeyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	fmt.Println("Listening on port:", port)
	http.ListenAndServe(":"+port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	text := r.FormValue("text")
	key, err := hex.DecodeString(r.FormValue("key"))
	if err != nil {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}
	action := r.FormValue("action")

	var result string
	if action == "E" || action == "e" {
		result, err = Encrypt(text, key)
		if err != nil {
			http.Error(w, "Encryption error", http.StatusInternalServerError)
			return
		}
	} else if action == "D" || action == "d" {
		result, err = Decrypt(text, key)
		if err != nil {
			http.Error(w, "Decryption error", http.StatusInternalServerError)
			return
		}
	}

	tpl.Execute(w, struct {
		Text   string
		Key    string
		Action string
		Result string
	}{
		Text:   text,
		Key:    hex.EncodeToString(key),
		Action: action,
		Result: result,
	})
}

func generateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key, err := generateRandomKey()
	if err != nil {
		http.Error(w, "Key generation error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, key)
}

func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func Encrypt(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainText))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string, key []byte) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
