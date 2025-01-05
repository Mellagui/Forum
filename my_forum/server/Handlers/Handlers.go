package Handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	Cruds "forum/Api"
	"forum/GlobVar"
)

func HandleStatic() {
	fs := http.FileServer(http.Dir(GlobVar.StaticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
func HandleUploads() {
	http.Handle("/Uploads/", http.StripPrefix("/Uploads/", http.FileServer(http.Dir("../Uploads"))))
}

func HandleComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Comment" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	if GlobVar.Guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}
	if r.Method == http.MethodPost {
		comment := r.FormValue("content")
		postId := r.FormValue("postId")
		userId := r.FormValue("userId")
		p, _ := strconv.Atoi(postId)
		u, _ := strconv.Atoi(userId)

		Cruds.InsertComment(p, u, comment)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "404", http.StatusNotFound)
	return
}

func HandleLikeDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/IsLike" {
		http.Error(w, "404", http.StatusNotFound)
		return
	}
	if GlobVar.Guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}

	if r.Method == http.MethodPost {
		if r.FormValue("isLike") == "false" {
			postId, _ := strconv.Atoi(r.FormValue("postId"))
			userId, _ := strconv.Atoi(r.FormValue("userId"))
			Cruds.InsertLikeDislike(userId, postId, false)
		} else if r.FormValue("isLike") == "true" {
			postId, _ := strconv.Atoi(r.FormValue("postId"))
			userId, _ := strconv.Atoi(r.FormValue("userId"))
			Cruds.InsertLikeDislike(userId, postId, true)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		if r.FormValue("email") == GlobVar.UserEmail {
			GlobVar.UserEmail = ""
			GlobVar.Guest = true
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "404 - Page Not Found", 404)
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Sign_In" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	users, err := Cruds.GetAllUsers()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(users) == 0 {
		http.Redirect(w, r, "/Sign_Up", 303)
		return
	}

	if r.Method == http.MethodPost {
		GlobVar.UserEmail = r.FormValue("email")
		password := r.FormValue("password")

		// Is Valid Acconut Redirect to /
		for _, user := range users {
			if user.Email == GlobVar.UserEmail && user.PasswordHash == password {
				GlobVar.Guest = false
				http.Redirect(w, r, "/", 303) // 301: Moved Permanently // 302: Found // 303: See Other
				return
			}
		}
		http.Redirect(w, r, "/Sign_In", 303)
		return
		// for render error message to client 

		//     if r.Method == http.MethodPost {
		//         email := r.FormValue("email")
		//         password := r.FormValue("password")
		//         data := Cruds.GetUserByAny(email)
		//         if data == nil || data.Email != email || data.PasswordHash != password {
		//             // Render the sign-in page with an error message
		//             tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-in-page.html"))
		//             if err != nil {
		//                 http.Error(w, "Internal server error", http.StatusInternalServerError)
		//                 return
		//             }
		//             // Add error message to the data passed to the template
		//             tmplData := struct {
		//                 ErrorMessage string
		//             }{
		//                 ErrorMessage: "Incorrect email or password. Please try again.",
		//             }
		//             err = tmpl.Execute(w, tmplData)
		//             if err != nil {
		//                 http.Error(w, "Internal server error", http.StatusInternalServerError)
		//             }
		//             return
		//         }
		//         // Successful login
		//         GlobVar.Guest = false
		//         http.Redirect(w, r, "/Home", http.StatusSeeOther) // 303: See Other
		//         return
		//     }
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-in-page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// err = tmpl.Execute(w, nil)
	err = tmpl.Execute(w, r)
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
		GlobVar.UserEmail = r.FormValue("email")
		password := r.FormValue("password")
		image := GlobVar.DefaultImage

		users, _ := Cruds.GetAllUsers()
		for _, user := range users {
			if user.Email == GlobVar.UserEmail || user.Name == name {
				http.Redirect(w, r, "/Sign_Up", 303) // 301: Moved Permanently // 302: Found // 303: See Other
				return
			}
		}

		Cruds.InsertUser(name, image, GlobVar.UserEmail, password)
		if GlobVar.AddAccountSucces {
			GlobVar.Guest = false
			GlobVar.AddAccountSucces = false
			http.Redirect(w, r, "/", 303)
			return
		}
		http.Redirect(w, r, "/Sign_Up", 303)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "sign-up-page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, r)
	if err != nil {
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
						if d.IsLike == false {
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
	if GlobVar.Guest {
		http.Redirect(w, r, "/Sign_In", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := Cruds.GetUserByAny(GlobVar.UserEmail)

	tmpl, err := template.ParseFiles(filepath.Join(GlobVar.TemplatesPath, "account-page.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	return
}

func HandleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Update_Profile" {
		http.Error(w, "page - not found", 404)
		return
	}
	if GlobVar.Guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}

	data := Cruds.GetUserByAny(GlobVar.UserEmail)

	if r.Method == http.MethodPost {
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
		if file != nil {
			defer file.Close()
		}

		imagePath := data.Image // Keep the existing image if no new file is uploaded
		if file != nil {
			// log.Printf("Saving file: %s", handler.Filename)

			dst, err := os.Create("../Uploads/" + handler.Filename)
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
			imagePath = "/Uploads/" + handler.Filename
			// log.Printf("File saved successfully: %s", imagePath)
		}

		// Update user in the database
		Cruds.UpdateUser(email, name, imagePath, password, GlobVar.UserEmail)
		GlobVar.UserEmail = email
		http.Redirect(w, r, "/Profile_Account", 303)
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
	if GlobVar.Guest {
		http.Redirect(w, r, "/Sign_In", 303)
		return
	}

	data := Cruds.GetUserByAny(GlobVar.UserEmail)
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		Cruds.InsertPost(data.ID, GlobVar.DefaultImage, title, content, category)
		if GlobVar.AddPostSucces {
			GlobVar.AddPostSucces = false
			http.Redirect(w, r, "/", 303)
			return
		} else {
			http.Redirect(w, r, "/New_Post", 303)
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
