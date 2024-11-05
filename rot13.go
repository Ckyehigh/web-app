package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/process", processHandler)

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
	// Check if the request method is POST
	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		key, err := strconv.Atoi(r.FormValue("key"))
		if err != nil {
			http.Error(w, "Invalid key", http.StatusBadRequest)
			return
		}
		action := r.FormValue("action")

		var result string
		if action == "E" || action == "e" {
			result = Encrypt(text, key)
		} else if action == "D" || action == "d" {
			result = Decrypt(text, key)
		}

		// Render the index template with results
		tpl.Execute(w, struct {
			Text   string
			Key    int
			Action string
			Result string
		}{
			Text:   text,
			Key:    key,
			Action: action,
			Result: result,
		})
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Encrypt(text string, key int) string {
	var result strings.Builder
	for _, letter := range text {
		if letter >= 'A' && letter <= 'Z' { // Uppercase letters
			newLetter := (int(letter)-'A'+key)%26 + 'A'
			result.WriteRune(rune(newLetter))
		} else if letter >= 'a' && letter <= 'z' { // Lowercase letters
			newLetter := (int(letter)-'a'+key)%26 + 'a'
			result.WriteRune(rune(newLetter))
		} else {
			result.WriteRune(letter) // Keep non-alphabetic characters unchanged
		}
	}
	return result.String()
}

func Decrypt(encrypted string, key int) string {
	var result strings.Builder
	for _, letter := range encrypted {
		if letter >= 'A' && letter <= 'Z' { // Uppercase letters
			newLetter := (int(letter)-'A'-key+26)%26 + 'A'
			result.WriteRune(rune(newLetter))
		} else if letter >= 'a' && letter <= 'z' { // Lowercase letters
			newLetter := (int(letter)-'a'-key+26)%26 + 'a'
			result.WriteRune(rune(newLetter))
		} else {
			result.WriteRune(letter) // Keep non-alphabetic characters unchanged
		}
	}
	return result.String()
}
