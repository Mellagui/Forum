package main

import (
	"database/sql"
	"log"
	"net/http"
	
	"forum/Migrations"
	"forum/GlobVar"
	"forum/Handlers"

	_ "github.com/mattn/go-sqlite3"
)

// Category
// Comment
// Js
// Filter

func init() {
	var err error
	GlobVar.DB, err = sql.Open("sqlite3", "../Database/database.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	Migrations.Migrate()
}

func main() {
	defer GlobVar.DB.Close()
	GlobVar.Guest = true
	Handlers.HandleStatic()
	Handlers.HandleUploads()

	http.HandleFunc("/Comment", Handlers.HandleComment)
	http.HandleFunc("/IsLike", Handlers.HandleLikeDislike)
	http.HandleFunc("/Log_Out", Handlers.HandleLogOut)
	http.HandleFunc("/", Handlers.HandleIndex)
	http.HandleFunc("/Sign_In", Handlers.HandleSignIn)
	http.HandleFunc("/Sign_Up", Handlers.HandleSignUp)
	http.HandleFunc("/Profile_Account", Handlers.HandleProfileAccount)
	http.HandleFunc("/Update_Profile", Handlers.HandleProfileUpdate)
	http.HandleFunc("/New_Post", Handlers.HandleNewPost)
	// http.HandleFunc("/Filter_Page", HandleFilter)
	log.Println("server start: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("server not listener: %v", err)
	}
}
