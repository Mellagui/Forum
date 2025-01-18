package LogIn

import (
	"fmt"
	"forum/Cookies"
	g "forum/GlobVar"
	"forum/Api"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// this function hash user password
func PasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

// this function insert user to database (for sign up)
func insertUser(name, email, password string) error {
	imagePath := "../Uploads/image.jpg"
	query := `insert into users (email, user_name, password_hash, user_image) values (?,?,?,?)`
	_, err := g.DB.Exec(query, email, name, password, imagePath)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

// this function handle user sign up
func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_Up" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		name:= r.FormValue("name")
		email := r.FormValue("email")
        password := r.FormValue("password")
		password = PasswordHash(password)
		err := insertUser(name, email, password)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			return
		}
		if g.AddAccountSucces {
			g.Guest = false
			g.AddAccountSucces = false
			cookies.Set_Cookies_Handler(w, r)
			http.Redirect(w, r, "/", 303)
			return
		}
		http.Redirect(w, r, "/Sign_Up", 303)
		return
		
	} else if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("../../client/templates/sign-up-page.html" , "Sign_Up.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
    } else {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_In" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	users, err := Cruds.GetAllUsers()
	if err != nil {
		log.Fatal("failed to get users: ", err)
	}
	if len(users) == 0 {
		http.Redirect(w, r, "/Sign_Up", 303)
		return
	}
	if r.Method == http.MethodPost {
		mail := r.FormValue("email")
		password := r.FormValue("password")
		password = PasswordHash(password)

		for _, user := range users {
			if user.Email == mail && user.PasswordHash == password {
				g.UserEmail = mail
				g.Guest = false
				cookies.Set_Cookies_Handler(w, r)
				http.Redirect(w, r, "/", 303)
				return
			}
		}
		http.Redirect(w, r, "/Sign_In", 303)
		return
	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("../../client/templates/sign-in-page.html" , "Sign_In.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

}
