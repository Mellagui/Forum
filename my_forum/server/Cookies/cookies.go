package cookies

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

// Generate_Cookie_session generates a cryptographically secure random session ID.
func Generate_Cookie_session() (string, error) {
    id := make([]byte, 32)
    _, err := rand.Read(id)
    if err != nil {
        return "", err // Return the error instead of logging and exiting
    }
    return base64.RawStdEncoding.EncodeToString(id), nil
}

// This Function set cookies
func Set_Cookies_Handler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("Session_ID")
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name: "Session_ID", // cookie name
			Path: "/",
			Expires: time.Now().Add(-1 * time.Hour),
		})
	}

	session_id, err := Generate_Cookie_session()
	if err != nil {
		log.Println(err)
		return
	}
	cookies := &http.Cookie{
		Name: "Session_ID", // cookie name
		Value: session_id,
		Path: "/",
		HttpOnly: true, // protect from mallitious js injections
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
