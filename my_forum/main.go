package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID           int
	Name         string
	Image        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
type Usss struct {
	IDS    int    `json:"id"`
	NameS  string `json:"name"`
	ImageS string `json:"image"`
	EmailS string `json:"email"`
}

var (
	usss             []Usss
	db               *sql.DB
	users            []User
	addAccountSucces bool
	guest            bool
	userEmail        string
)

const defaultImage = "https://groupietrackers.herokuapp.com/api/images/queen.jpeg"

func init() {
	var err error
	db, err = sql.Open("sqlite3", "database/database.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	Migrate()
	getData()
}

func Migrate() {
	query, err := os.ReadFile("modules.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database migrated successfully!")
}

func getData() {
	URL := "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("...........", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("cant get data status code !200")
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&usss)
	if err != nil {
		fmt.Println("............", err)
	}
}

func main() {
	defer db.Close()
	guest = true

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/Sign_In", HandleSignIn)
	http.HandleFunc("/Sign_Up", HandleSignUp)
	http.HandleFunc("/Home", HandleHome)
	http.HandleFunc("/Profile_Account", HandleProfileAccount)
	http.HandleFunc("/Update_Profile", HandleProfileUpdate)
	// http.HandleFunc("/New_Post", HandlePost)
	// http.HandleFunc("/Update_Post", HandlePostUpdate)
	// http.HandleFunc("/Filter_Page", Handle5)
	http.HandleFunc("/Chat_Rooms", HandleChatRoom)
	log.Println("server start: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("server not listener: %v", err)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_In" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		userEmail = r.FormValue("email")
		password := r.FormValue("password")

		// GetUserByEmail from Table users
		data := GetUserByAny(userEmail)

		// Is Valid Acconut Redirect to /Home
		if data.Email == userEmail && data.PasswordHash == password {
			guest = false
			http.Redirect(w, r, "/Home", 303) // 301: Moved Permanently // 302: Found // 303: See Other
			return
		}

		tmpl, err := template.ParseFiles("sign-in-page.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, r)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("sign-in-page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	return
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_Up" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// insertUser into Table users
		InsertUser(name, email, password)

		// if addAccountSucces Redirect to /Home
		if addAccountSucces {
			guest = false
			addAccountSucces = false
			http.Redirect(w, r, "/Home", 303)
			return
		}

		tmpl, err := template.ParseFiles("sign-up-page.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, r)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("sign-up-page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Home" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("home-page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, usss)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleProfileAccount(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Profile_Account" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	// if user hase no account riderect to /Sign_In
	if guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "page - not found", 404)
		return
	}

	data := GetUserByAny(userEmail)
	fmt.Println(data)

	tmpl, err := template.ParseFiles("account-page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Update_Profile" {
		http.Error(w, "page - not found", 404)
		return
	}
	// Redirect guests to sign in
	if guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}

	data := GetUserByAny(userEmail)

	if r.Method == http.MethodPost {
		// Handle text fields
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		if len(name) == 0 {
			name = data.Name
		}
		if len(email) == 0 {
			email = data.Email
		}

		// Handle file upload
		file, handler, err := r.FormFile("image")
		if err != nil && err != http.ErrMissingFile { // Allow updates without file
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}
		defer func() {
			if file != nil {
				file.Close()
			}
		}()

		imagePath := data.Image // Keep the existing image if no new file is uploaded
		if file != nil {
			log.Printf("Saving file: %s", handler.Filename)
			dst, err := os.Create("./uploads/" + handler.Filename)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Error saving the file", http.StatusInternalServerError)
				return
			}
			imagePath = "/uploads/" + handler.Filename
			log.Printf("File saved successfully: %s", imagePath)
		}

		// Update user in the database
		UpdateUser(email, name, imagePath, password, userEmail)
		userEmail = email
		http.Redirect(w, r, "/Profile_Account", 303)
		return
	}

	// Render the update page
	tmpl, err := template.ParseFiles("update-account-page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// func HandleProfileUpdate(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/Update_Profile" {
// 		http.Error(w, "page - not found", 404)
// 		return
// 	}
// 	// if user hase no account riderect to /Sign_In
// 	if guest {
// 		http.Redirect(w, r, "/Sign_In", 303)
// 		return
// 	}

// 	data := GetUserByAny(userEmail)

// 	if r.Method == http.MethodPost {
// 		// fmt.Println("222222222")
// 		name := r.FormValue("name")
// 		email := r.FormValue("email")
// 		if len(name) == 0 {
// 			name = data.Name
// 		}
// 		if len(email) == 0 {
// 			email = data.Email
// 		}
// 		// 	// fetch data ?

// 		//  // update data ?
// 		UpdateUser(email, name, userEmail)
// 		userEmail = email
// 		http.Redirect(w, r, "/Profile_Account", 303)
// 		return
// 	}

// 	// fmt.Println("1111111111111")
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "page - not found", 404)
// 		return
// 	}

// 	tmpl, err := template.ParseFiles("update-account-page.html")
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	err = tmpl.Execute(w, data)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 	}
// }

func HandleChatRoom(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Chat_Rooms" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	// if user hase no account riderect to /Sign_In
	if guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}

	tmpl, err := template.ParseFiles("messages-page.html")
	if err != nil {
		fmt.Println("111111111")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	data := GetUserByAny(userEmail)
	fmt.Println(data)

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("22222222222")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

func InsertUser(name, email, password string) {
	var lastId int
	query := `insert into users (id, user_name, user_image, email, password_hash) values (?,?,?,?,?)`
	maxId := "select COALESCE(MAX(id), 0) from users"
	err := db.QueryRow(maxId).Scan(&lastId)
	if err != nil {
		log.Printf("error queryrow maxid: %v", err)
		return
	}
	_, err = db.Exec(query, lastId+1, name, defaultImage, email, password)
	if err != nil {
		log.Printf("error exec queryyy: %v", err)
		return
	}
	addAccountSucces = true
}

func GetUserByAny(required string) *User {
	// iyner join
	query := "SELECT id, user_name, user_image, email, password_hash, created_at FROM users WHERE email = ?"
	var user User
	err := db.QueryRow(query, required).Scan(&user.ID, &user.Name, &user.Image, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("sql.ErrRows")
			return nil // No user found
		}
		log.Fatalf("error getUserByName: %v", err)
		return nil
	}
	return &user
}

func UpdateUser(email, name, image, password, userEmail string) {
    log.Printf("Updating user: name=%s, image=%s, email=%s, password=%s, oldEmail=%s", name, image, email, password, userEmail)
    query := `UPDATE users SET user_name = ?, user_image = ?, email = ?, password_hash = ? WHERE email = ?`
    _, err := db.Exec(query, name, image, email, password, userEmail)
    if err != nil {
        log.Printf("error exec query Update: %v", err)
        return
    }
    log.Println("User updated successfully")
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Post" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	// if user hase no account riderect to /Sign_In
	if guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}

	tmpl, err := template.ParseFiles("post-page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// func HandlePostUpdate(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/Post_Update" {
// 		http.Error(w, "404 not found", http.StatusNotFound)
// 		return
// 	}

// 	// if user hase no account riderect to /Sign_In
// 	if guest {
// 		http.Redirect(w, r, "/Sign_In", 303)
// 		return
// 	}

// 	tmpl, err := template.ParseFiles("post-update-page.html")
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	err = tmpl.Execute(w, r)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 	}
// }
