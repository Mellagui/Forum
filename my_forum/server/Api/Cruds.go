package Cruds

import (
	"database/sql"
	"fmt"
	"log"

	"forum/GlobVar"
)

// Insert Data
func InsertUser(name, image, email, password string) {
	var lastId int
	query := `insert into users (id, email, user_name, password_hash, user_image) values (?,?,?,?,?)`
	maxId := "select COALESCE(MAX(id), 0) from users"
	err := GlobVar.DB.QueryRow(maxId).Scan(&lastId)
	if err != nil {
		log.Printf("error queryrow maxid: %v", err)
		return
	}
	_, err = GlobVar.DB.Exec(query, lastId+1, email, name, password, image)
	if err != nil {
		log.Printf("error exec queryyy: %v", err)
		return
	}
	GlobVar.AddAccountSucces = true
}

func InsertPost(userId int, image, title, content, category string) {
	var lastId int
	query := `insert into posts (id, user_id, title, content, image_url, category) values (?,?,?,?,?,?)`
	maxId := "select COALESCE(MAX(id), 0) from posts"
	err := GlobVar.DB.QueryRow(maxId).Scan(&lastId)
	if err != nil {
		log.Printf("error queryrow maxid: %v", err)
		return
	}
	_, err = GlobVar.DB.Exec(query, lastId+1, userId, title, content, image, category)
	if err != nil {
		log.Printf("error exec queryyy: %v", err)
		return
	}
	GlobVar.AddPostSucces = true
}

func InsertComment(postId, userId int, content string) {
	var lastId int
	query := `insert into comments (id, post_id, user_id, content) values (?,?,?,?)`
	maxId := "select COALESCE(MAX(id), 0) from comments"
	err := GlobVar.DB.QueryRow(maxId).Scan(&lastId)
	if err != nil {
		log.Printf("error queryrow maxid: %v", err)
		return
	}
	_, err = GlobVar.DB.Exec(query, lastId+1, postId, userId, content)
	if err != nil {
		log.Printf("error exec queryyy: %v", err)
		return
	}
	GlobVar.AddCommentSucces = true
}

func InsertCategory(byUserId int, categoryName string) {
	var lastId int
	query := `insert into categories (id, category_name, created_by_user_id) values (?,?,?)`
	maxId := "select COALESCE(MAX(id), 0) from categories"
	err := GlobVar.DB.QueryRow(maxId).Scan(&lastId)
	if err != nil {
		log.Printf("error queryrow maxid: %v", err)
		return
	}
	_, err = GlobVar.DB.Exec(query, lastId+1, categoryName, byUserId)
	if err != nil {
		log.Printf("error exec queryyy: %v", err)
		return
	}
	GlobVar.AddCategorySucces = true
}

func InsertLikeDislike(userId, postId int, isLike bool) {
	var lastId int
	query := `insert into likeDislike (id, user_id, post_id, is_like) values (?,?,?,?)`
	maxId := "select COALESCE(MAX(id), 0) from likeDislike"
	err := GlobVar.DB.QueryRow(maxId).Scan(&lastId)
	if err != nil {
		log.Printf("error queryrow maxid: %v", err)
		return
	}
	_, err = GlobVar.DB.Exec(query, lastId+1, userId, postId, isLike)
	if err != nil {
		log.Printf("error exec queryyy: %v", err)
		return
	}
	GlobVar.AddLikeDislikeSucces = true
}

// this function deletes the like and dislike from the database
func DeleteLikeDislike(userId, postId int) {
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
func CheckLikeDislikeExists(userId, postId int) (bool, bool) {
    var isLike bool
    query := `SELECT is_like FROM likeDislike WHERE user_id = ? AND post_id = ?`

    err := GlobVar.DB.QueryRow(query, userId, postId).Scan(&isLike)
    if err != nil {
        if err == sql.ErrNoRows {
			// makayn la like la dislike
            return false, false
        }
        return false, false
    }
    // Return true o isLike
    return true, isLike
}



// Update Data
func UpdateUser(email, name, image, password, userEmail string) {
	query := `UPDATE users SET user_name = ?, user_image = ?, email = ?, password_hash = ? WHERE email = ?`
	_, err := GlobVar.DB.Exec(query, name, image, email, password, userEmail)
	if err != nil {
		log.Printf("error exec query Update: %v", err)
	}
}

// Get Data
func GetUserByAny(required string) *GlobVar.User {
	// iyner join
	query := `SELECT id, email, user_name, password_hash, user_image, created_at FROM users WHERE email = ?`
	var user GlobVar.User
	err := GlobVar.DB.QueryRow(query, required).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt)
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

// func GetUserById(id int) *User {
// 	// iyner join
// 	query := "SELECT id, email, user_name, password_hash, user_image, created_at FROM users WHERE id = ?"
// 	var user User
// 	err := GlobVar.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("sql.ErrRows")
// 			return nil // No user found
// 		}
// 		log.Fatalf("error getUserByName: %v", err)
// 		return nil
// 	}
// 	return &user
// }

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
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Image, &user.CreatedAt); err != nil { //&user.Password,
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
		if err := rows.Scan(&post.ID, &post.UserId, &post.Image, &post.Title, &post.Content, &post.Category, &post.CreatedAt); err != nil { //&post.Password,
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
	query := `SELECT id, post_id, user_id, content FROM comments`
	rows, err := GlobVar.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []GlobVar.Comment
	for rows.Next() {
		var cmt GlobVar.Comment
		if err := rows.Scan(&cmt.ID, &cmt.PostId, &cmt.UserId, &cmt.Content); err != nil {
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
