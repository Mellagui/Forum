package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/GlobVar"
	"forum/Handlers"
	middleware "forum/Middleware"
	"forum/Migrations"

	_ "modernc.org/sqlite"
)

// handle min-length password in sign-up
// handle js message invalid userName
// handle userName-comments
// image handling
// avatar resize
// new-post textearia resize

func init() {
	var err error
	GlobVar.DB, err = sql.Open("sqlite", "../Database/database.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	Migrations.Migrate()
}

func main() {
	defer GlobVar.DB.Close()
	Handlers.HandleStatic()
	Handlers.HandleUploads()

	// Public routes
	http.HandleFunc("/", Handlers.HandleIndex)
	http.HandleFunc("/Sign_In", Handlers.HandleSignIn)
	http.HandleFunc("/Sign_Up", Handlers.HandleSignUp)
	http.HandleFunc("/api/auth/status", Handlers.HandleAuthStatus)

	// Protected routes
	http.HandleFunc("/Comment", middleware.ValidateSession(Handlers.HandleComment))
	http.HandleFunc("/IsLike", middleware.ValidateSession(Handlers.HandleLikeDislike))
	http.HandleFunc("/post/", middleware.ValidateSession(Handlers.HandlePostPage))
	http.HandleFunc("/Log_Out", middleware.ValidateSession(Handlers.HandleLogOut))
	http.HandleFunc("/Profile_Account", middleware.ValidateSession(Handlers.HandleProfileAccount))
	http.HandleFunc("/Update_Profile", middleware.ValidateSession(Handlers.HandleProfileUpdate))
	http.HandleFunc("/New_Post", middleware.ValidateSession(Handlers.HandleNewPost))

	log.Println("server start: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("server not listener: %v", err)
	}
}