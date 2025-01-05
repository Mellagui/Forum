package GlobVar

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Email        string
	Name         string
	PasswordHash string
	Image        string
	CreatedAt    time.Time
	// UpdatedAt    time.Time
}

type Post struct {
	ID        int
	UserId    int
	Image     string
	Title     string
	Content   string
	Category  string
	CreatedAt time.Time
	// UpdatedAt time.Time
}

type Comment struct {
	ID        int
	PostId    int
	UserId    int
	Content   string
	CreatedAt time.Time
	UpdatedAT time.Time
}

type Categories struct {
	ID              int
	CategoryName    string
	CreatedByUserId int
}

type PostCategory struct {
	PostId     int
	CategoryId int
}

type LikeDislike struct {
	ID     int
	UserId int
	PostId int
	IsLike bool
}

type Home struct {
	ID int
	//post
	PostId        int       //Post.ID
	PostImage     string    //Post.Image
	PostTitle     string    //Post.Title
	PostContent   string    //Post.Content
	PostCreatedAt time.Time //Post.CreatedAt
	//user
	UserId    int
	UserName  string //User.Name
	UserImage string //User.Image
	//comment
	NbrComment int
	// PostComments map[string]string
	// PostComments map[UserName]Comment
	//likedislike
	NbrLike    int
	NbrDislike int
	//category
	CategoryName string
}

var (
	DB            *sql.DB
	Users         []User
	Posts         []Post
	Comments      []Comment
	LikesDislikes []LikeDislike

	UserEmail            string
	Guest                bool
	AddAccountSucces     bool
	AddPostSucces        bool
	AddCommentSucces     bool
	AddCategorySucces    bool
	AddLikeDislikeSucces bool
)

const TemplatesPath = "../../client/templates/"
const StaticPath = "../../client/static/"
const DefaultImage = "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
