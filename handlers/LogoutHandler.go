package forum

import (
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Cookie("user")
	if user != nil {
		cookie := http.Cookie{
			Name:    "user",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, &cookie)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
