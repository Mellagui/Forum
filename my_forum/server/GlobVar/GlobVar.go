package GlobVar

import (
	"database/sql"
	"sync"
	"time"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Name         string    `db:"user_name" json:"user_name"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	Image        string    `db:"user_image" json:"user_image"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	// UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Post struct {
	ID        string    `db:"id" json:"id"`
	UserId    string    `db:"user_id" json:"user_id"`
	Image     string    `db:"image_url" json:"image_url"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Category  string    `db:"category" json:"category"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Comment struct {
	ID        string    `db:"id" json:"id"`
	PostId    string    `db:"post_id" json:"post_id"`
	UserId    string    `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UserName string		`db:"UserName" json:"UserName"`
}

type Categories struct {
	ID              string `db:"id" json:"id"`
	CategoryName    string `db:"category_name" json:"category_name"`
	CreatedByUserId string `db:"created_by_user_id" json:"created_by_user_id"`
}

type PostCategory struct {
	PostId     string `db:"post_id" json:"post_id"`
	CategoryId string `db:"category_id" json:"category_id"`
}

type LikeDislike struct {
	ID     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
	PostId string `db:"post_id" json:"post_id"`
	IsLike bool   `db:"is_like" json:"is_like"`
}

type Home struct {
	ID string `json:"id"`
	// Post
	PostId        string    `json:"post_id"`        // Post.ID
	PostImage     string    `json:"post_image"`     // Post.Image
	PostTitle     string    `json:"post_title"`     // Post.Title
	PostContent   string    `json:"post_content"`   // Post.Content
	PostCreatedAt time.Time `json:"post_created_at"` // Post.CreatedAt
	// User
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`  // User.Name
	UserImage string `json:"user_image"` // User.Image
	// Comment
	NbrComment int `json:"nbr_comment"`
	// LikeDislike
	NbrLike    int `json:"nbr_like"`
	NbrDislike int `json:"nbr_dislike"`
	// Category
	CategoryName string `json:"category_name"`
}

var (
	UpdateLikeDislikeSuccess bool
	DB            *sql.DB
	Users         []User
	Posts         []Post
	Comments      []Comment
	LikesDislikes []LikeDislike

	UserId                   string
	AddAccountSucces         bool
	AddPostSucces            bool
	AddCommentSucces         bool
	AddCategorySucces        bool
	AddLikeDislikeSucces     bool
	DeleteLikeDislikeSuccess bool
	UserMutex sync.Mutex 
)

const (
	TemplatesPath = "../../client/templates/"
	StaticPath    = "../../client/static/"
	DefaultImage  = "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
)