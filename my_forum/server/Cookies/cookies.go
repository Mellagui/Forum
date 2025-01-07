package cookies

import (
	"crypto/rand"
	"encoding/base64"
	//"forum/GlobVar"
	"log"
	"net/http"
	"time"
)
// this function generate a crypted random cookie ID
func Generate_Cookie_session() string {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		log.Fatal(err)
	}
	return base64.RawStdEncoding.EncodeToString(id)
}

// This Function set cookies
func Set_Cookies_Handler(w http.ResponseWriter, r *http.Request) {
	session_id := Generate_Cookie_session()
	cookies := &http.Cookie{
		Name: "Session_ID", // cookie name
		Value: session_id,
		Path: "/",
		Secure: true,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}
	http.SetCookie(w, cookies)
}

// This function delete the cookie
func Delete_Cookie_Handler(w http.ResponseWriter, r *http.Request) {
	cookies := &http.Cookie{
		Name: "Session_ID", // cookie name
		Value: "",
		Path: "/",
		Secure: true,
		Expires: time.Now().Add(-1 * time.Hour),	
	}
	http.SetCookie(w, cookies)
}