package Cruds

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/GlobVar"
	cookies "forum/Cookies"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


// Insert Data
func InsertUser(name, image, email, password string) string {
    id := GenerateUUID()
    hashedPassword, err := HashPassword(password)
    if err != nil {
        log.Printf("error hashing password: %v", err)
        return ""
    }
    query := `INSERT INTO users (id, email, user_name, password_hash, user_image) VALUES (?, ?, ?, ?, ?)`
    _, err = GlobVar.DB.Exec(query, id, email, name, hashedPassword, image)
    if err != nil {
        log.Printf("error inserting user: %v", err)
        return ""
    }
    return id
}

func InsertPost(userId, image, title, content, category string) bool {
	id := GenerateUUID()
	query := `INSERT INTO posts (id, user_id, title, content, image_url, category) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, userId, title, content, image, category)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return false
	}
	return true
}

func InsertComment(postId, userId, content string) {
	id := GenerateUUID()
	query := `INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, postId, userId, content)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
	GlobVar.AddCommentSucces = true
}

func InsertCategory(byUserId, categoryName string) {
	id := GenerateUUID()
	query := `INSERT INTO categories (id, category_name, created_by_user_id) VALUES (?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, categoryName, byUserId)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
	GlobVar.AddCategorySucces = true
}

func InsertLikeDislike(userId, postId string, isLike bool) {
	id := GenerateUUID()
	query := `INSERT INTO likeDislike (id, user_id, post_id, is_like) VALUES (?, ?, ?, ?)`
	_, err := GlobVar.DB.Exec(query, id, userId, postId, isLike)
	if err != nil {
		log.Printf("error exec query: %v", err)
		return
	}
	GlobVar.AddLikeDislikeSucces = true
}

// this function deletes the like and dislike from the database
func DeleteLikeDislike(userId, postId string) {
	query := `DELETE FROM likeDislike WHERE user_id = ? AND post_id = ?`
	_, err := GlobVar.DB.Exec(query, userId, postId)
	if err != nil {
		log.Printf("error deleting like or dislike: %v", err)
		GlobVar.DeleteLikeDislikeSuccess = false
		return
	}
	GlobVar.DeleteLikeDislikeSuccess = true
}

// this function checks if the like or dislike already exist
func CheckLikeDislikeExists(userId, postId string) (bool, bool) {
	var isLike bool
	query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ?`
	err := GlobVar.DB.QueryRow(query, userId, postId).Scan(&isLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false
		}
		return false, false
	}
	return true, isLike
}

func GetPostByID(postID string) (string, *GlobVar.Post, error) {
    query := `SELECT id, user_id, image_url, title, content, category, created_at FROM posts WHERE id = ?`
    var post GlobVar.Post
    err := GlobVar.DB.QueryRow(query, postID).Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return "" ,nil, fmt.Errorf("post not found")
        }
        return "" ,nil, err
    }
    return "" ,&post, nil
}

// Update Data
func UpdateUser(email, name, image, password, userId string) {
	var err error
	var hashedPassword string
	if len(password) != 0 {
		hashedPassword, err = HashPassword(password)
		if err != nil {
			log.Printf("error hashing password: %v", err)
			return
		}
		query := `UPDATE users SET user_name = ?, user_image = ?, email = ?, password_hash = ? WHERE id = ?`
		_, err = GlobVar.DB.Exec(query, name, image, email, hashedPassword, userId)
		if err != nil {
			log.Printf("error exec query Update: %v", err)
		}
	} else {
		query := `UPDATE users SET user_name = ?, user_image = ?, email = ? WHERE id = ?`
		_, err = GlobVar.DB.Exec(query, name, image, email, userId)
		if err != nil {
			log.Printf("error exec query Update: %v", err)
		}
	}
}

// Get Data
func GetUserByAny(required string) *GlobVar.User {
    query := `SELECT id, email, user_name, password_hash, user_image, created_at FROM users WHERE id = ?`
    var user GlobVar.User
    err := GlobVar.DB.QueryRow(query, required).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("User not found:", required) // Debugging
            return nil
        }
        log.Printf("error getUserByName: %v", err)
        return nil
    }
    return &user
}



// Get All Data
func GetAllUsers() ([]GlobVar.User, error) {
	query := `SELECT id, email, user_name, password_hash, user_image, created_at FROM users`
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []GlobVar.User
	for rows.Next() {
		var user GlobVar.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllPosts() ([]GlobVar.Post, error) {
	query := `SELECT id, user_id, image_url, title, content, category, created_at FROM posts`
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []GlobVar.Post
	for rows.Next() {
		var post GlobVar.Post
		if err := rows.Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.Category, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetAllComments() ([]GlobVar.Comment, error) {
    query := `
        SELECT c.id, c.post_id, c.user_id, c.content, u.user_name 
        FROM comments c
        JOIN users u ON c.user_id = u.id
    `
    rows, err := GlobVar.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []GlobVar.Comment
    for rows.Next() {
        var cmt GlobVar.Comment
        if err := rows.Scan(&cmt.ID, &cmt.PostId, &cmt.UserId, &cmt.Content, &cmt.UserName); err != nil {
            return nil, err
        }
        comments = append(comments, cmt)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return comments, nil
}

func GetAllLikeDislike() ([]GlobVar.LikeDislike, error) {
	query := `SELECT id, user_id, post_id, is_like FROM likeDislike`
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likeDislike []GlobVar.LikeDislike
	for rows.Next() {
		var lk GlobVar.LikeDislike
		if err := rows.Scan(&lk.ID, &lk.UserId, &lk.PostId, &lk.IsLike); err != nil {
			return nil, err
		}
		likeDislike = append(likeDislike, lk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return likeDislike, nil
}

func ValidateSessionIDAndGetUserID(sessionID string) (string, bool) {
    var expiresAt time.Time
    var userID string
    query := `SELECT user_id, expires_at FROM Session WHERE id = ?`
    err := GlobVar.DB.QueryRow(query, sessionID).Scan(&userID, &expiresAt)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("Session not found:", sessionID) // Debugging
            return "", false
        }
        log.Printf("Error validating session ID: %v", err)
        return "", false
    }

    // Check if the session is expired
    if time.Now().After(expiresAt) {
        fmt.Println("Session expired:", sessionID)
        return "", false
    }


    return userID, true
}

func Set_Cookies_Handler(w http.ResponseWriter, r *http.Request, userID string) {
    var sessionID, token string
    var err error

    // Generate a unique session ID
    sessionID, err = cookies.Generate_Cookie_session()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        log.Printf("Error generating session ID: %v", err)
        return
    }

    // Generate a unique token for the session
    for {
        token, err = cookies.Generate_Cookie_session()
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            log.Printf("Error generating session token: %v", err)
            return
        }

        // Check if the token already exists in the database
        var exists bool
        query := `SELECT EXISTS(SELECT 1 FROM Session WHERE token = ?)`
        err = GlobVar.DB.QueryRow(query, token).Scan(&exists)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            log.Printf("Error checking token existence: %v", err)
            return
        }

        if !exists {
            break // Token is unique, exit the loop
        }
    }

    // Insert the session into the database
    expiresAt := time.Now().Add(7 * 24 * time.Hour) // Session expires in 7 days
    query := `INSERT INTO Session (id, user_id, token, expires_at) VALUES (?, ?, ?, ?)`
    _, err = GlobVar.DB.Exec(query, sessionID, userID, token, expiresAt)
    if err != nil {
        log.Printf("Error storing session in database: %v", err) // Debugging
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    fmt.Println("Session created for user:", userID) // Debugging

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

func Delete_Cookie_Handler(w http.ResponseWriter, r *http.Request) {
	// Get the session cookie
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