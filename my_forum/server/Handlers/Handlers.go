package Handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	Cruds "forum/Api"
	cookies "forum/Cookies"
	"forum/GlobVar"
)

func HandleStatic() {
	fs := http.FileServer(http.Dir(GlobVar.StaticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
func HandleUploads() {
	http.Handle("/Uploads/", http.StripPrefix("/Uploads/", http.FileServer(http.Dir("../Uploads"))))
}

func HandlePostPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	// Extract the post ID from the URL
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Fetch the post details
	_, post, err := Cruds.GetPostByID(postID)

	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Fetch the user who created the post
	user := Cruds.GetUserByAny(post.UserId)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch comments for the post
	comments, err := Cruds.GetAllComments()
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	// Filter comments for this post
	var postComments []GlobVar.Comment
	for _, comment := range comments {
		if comment.PostId == postID {
			postComments = append(postComments, comment)
		}
	}

	// Fetch likes and dislikes for the post
	likesDislikes, err := Cruds.GetAllLikeDislike()
	if err != nil {
		http.Error(w, "Failed to fetch likes/dislikes", http.StatusInternalServerError)
		return
	}

	var likes, dislikes int
	for _, ld := range likesDislikes {
		if ld.PostId == postID {
			if ld.IsLike {
				likes++
			} else {
				dislikes++
			}
		}
	}

	// Prepare the data to be passed to the template
	data := struct {
		Post     *GlobVar.Post
		User     *GlobVar.User
		Comments []GlobVar.Comment
		Likes    int
		Dislikes int
	}{
		Post:     post,
		User:     user,
		Comments: postComments,
		Likes:    likes,
		Dislikes: dislikes,
	}

	// Render the post page
	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "post_page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func HandleComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Comment" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodPost {
		comment := r.FormValue("content")
		postId := r.FormValue("postId")
		userId := r.FormValue("userId")

		Cruds.InsertComment(postId, userId, comment)
		http.Redirect(w, r, "/post/?id="+postId, http.StatusSeeOther)
		return
	}

	http.Error(w, "404", http.StatusNotFound)
}

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/IsLike" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		postId := r.FormValue("postId")
		userId := r.FormValue("userId")
		isLike := r.FormValue("isLike") == "true"

		exists, currentIsLike := Cruds.CheckLikeDislikeExists(userId, postId)

		if exists {
			if isLike == currentIsLike {
				Cruds.DeleteLikeDislike(userId, postId)
			} else {
				Cruds.DeleteLikeDislike(userId, postId)
				Cruds.InsertLikeDislike(userId, postId, isLike)
			}
		} else {
			Cruds.InsertLikeDislike(userId, postId, isLike)
		}

		http.Redirect(w, r, "/post/?id="+postId, http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		http.Error(w, "404 - Page Not Found", 404)
	}
}

func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Log_Out" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		// Delete the session cookie and session from the database
		Delete_Cookie_Handler(w, r)

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "404 - Page Not Found", http.StatusNotFound)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_In" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Fetch the user
		GlobVar.UserMutex.Lock()
		user := Cruds.GetUserByAny(GlobVar.UserId)
		GlobVar.UserMutex.Unlock()
		if user == nil {
			fmt.Println("User not found for email:", email) // Debugging
			http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
			return
		}

		// Compare the password
		if !Cruds.CheckPasswordHash(password, user.PasswordHash) {
			fmt.Println("Password mismatch for user:", email) // Debugging
			http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
			return
		}

		// Set the session cookie
		Set_Cookies_Handler(w, r, user.ID)

		GlobVar.UserMutex.Lock()
		GlobVar.UserId = user.ID
		GlobVar.UserMutex.Unlock()

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the sign-in page
	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-in-page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_Up" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		image := GlobVar.DefaultImage

		// Check if the email or username already exists
		GlobVar.UserMutex.Lock()
		users, err := Cruds.GetAllUsers()
		if err != nil {
			log.Printf("Error fetching users: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			GlobVar.UserMutex.Unlock()
			return
		}

		for _, user := range users {
			if user.Email == email || user.Name == name {
				http.Redirect(w, r, "/Sign_Up", http.StatusSeeOther)
				GlobVar.UserMutex.Unlock()
				return
			}
		}

		// Insert the new user into the database
		userID := Cruds.InsertUser(name, image, email, password)
		if userID == "" {
			http.Redirect(w, r, "/Sign_Up", http.StatusSeeOther)
			GlobVar.UserMutex.Unlock()
			return
		}

		// Set the session cookie for the newly signed-up user
		Set_Cookies_Handler(w, r, userID)

		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		GlobVar.UserMutex.Unlock()
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Render the sign-up page
	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-up-page.html"))
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405", http.StatusMethodNotAllowed)
		return
	}

	posts, err := Cruds.GetAllPosts()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	var data []GlobVar.Home
	var home GlobVar.Home
	if len(posts) > 0 {
		users, err := Cruds.GetAllUsers()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		comments, err := Cruds.GetAllComments()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		likesDislikes, err := Cruds.GetAllLikeDislike()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		for _, p := range posts {
			home.ID = p.ID
			//post
			home.PostId = p.ID
			home.PostImage = p.Image
			home.PostTitle = p.Title
			home.PostContent = p.Content
			home.PostCreatedAt = p.CreatedAt
			//user
			for _, u := range users {
				if p.UserId == u.ID {
					home.UserId = u.ID
					home.UserName = u.Name
					home.UserImage = u.Image
					break
				}
			}
			//comment
			home.NbrComment = 0
			if len(comments) > 0 {
				for _, c := range comments {
					if c.PostId == p.ID {
						home.NbrComment++
					}
				}
			}
			//likedislike
			home.NbrDislike = 0
			home.NbrLike = 0
			if len(likesDislikes) > 0 {
				for _, d := range likesDislikes {
					if d.PostId == p.ID {
						if !d.IsLike {
							home.NbrDislike++
						} else {
							home.NbrLike++
						}
					}
				}
			}
			//category
			home.CategoryName = p.Category
			data = append(data, home)
			home = GlobVar.Home{}
		}
	}

	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "index.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleProfileAccount(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Profile_Account" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Query the user using the user ID from the context
	data := Cruds.GetUserByAny(userID)
	if data == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "account-page.html"))
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

	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println("rani hna")
	data := Cruds.GetUserByAny(userID)
	fmt.Println("rani hna o jebt l user data:", data)

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println("rani hna o jebt l inputs valus:", email, " ", name, " ", password, ".")
		if len(name) == 0 {
			name = data.Name
		}
		if len(email) == 0 {
			email = data.Email
		}
		if len(password) == 0 {
			password = ""
		}

		// Handle file upload
		// To be Impelented !!!!!!!!!!!!!

		// Default to existing image
		imagePath := data.Image
		// Update user in the database
		Cruds.UpdateUser(email, name, imagePath, password, userID)
		http.Redirect(w, r, "/Profile_Account", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "update-account-page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HandleNewPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/New_Post" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := Cruds.GetUserByAny(userID)
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		// category := r.FormValue("category")
		content := r.FormValue("content")
		fmt.Println("test chi le3ba : ", title, " ", content)

		if Cruds.InsertPost(data.ID, GlobVar.DefaultImage, title, content, "test") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/New_Post", http.StatusSeeOther)
		}
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "new-post-page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func Set_Cookies_Handler(w http.ResponseWriter, r *http.Request, userID string) {
	sessionID, err := cookies.Generate_Cookie_session()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error generating session ID: %v", err)
		return
	}

	// Insert the session into the database
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // Session expires in 7 days
	query := `INSERT INTO Session (id, user_id, token, expires_at) VALUES (?, ?, ?, ?)`
	_, err = GlobVar.DB.Exec(query, sessionID, userID, sessionID, expiresAt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error storing session in database: %v", err)
		return
	}

	// Set the session cookie
	cookie := &http.Cookie{
		Name:     "Session_ID",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  expiresAt,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// Delete_Cookie_Handler deletes the session cookie.
func Delete_Cookie_Handler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Session_ID")
	if err != nil {
		// No session cookie found
		http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
		return
	}

	// Delete the session from the database
	query := `DELETE FROM Session WHERE id = ?`
	_, err = GlobVar.DB.Exec(query, cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error deleting session from database: %v", err)
		return
	}

	// Clear the session cookie
	cookie = &http.Cookie{
		Name:     "Session_ID",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func HandleAuthStatus(w http.ResponseWriter, r *http.Request) {
	var isAuthenticated bool
	cookie, err := r.Cookie("Session_ID")
	if err == nil {
		// Validate the session ID
		_, isAuthenticated = Cruds.ValidateSessionIDAndGetUserID(cookie.Value)
	}

	// Return the authentication status as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"isAuthenticated": isAuthenticated,
	})
}
