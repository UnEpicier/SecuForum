package forum

import (
	"database/sql"
	"fmt"
	f "forum"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./static/pages/user/login.html", "./static/layout/base.html"}
	tplt := template.Must(template.ParseFiles(files...))

	var page f.Page
	page.Logged = false

	cookie, _ := r.Cookie("user")
	if cookie != nil {
		page.Logged = true
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}

		email := r.FormValue("email")
		password := r.FormValue("passwd")
		keepAlive := r.FormValue("keep-alive")

		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}

		password, _ = f.HashPassword(password)

		row, err := db.Query("SELECT uuid FROM user WHERE email = "+ email +" AND password = "+ password + " LIMIT 1")

		if err != nil {
			log.Fatal(err)
		}
		var db_uuid string
		for row.Next() {
			err = row.Scan(&db_uuid)
			if err != nil {
				log.Fatal(err)
			}
		}
		row.Close()

		fmt.Println(db_uuid)

		cookie := http.Cookie{
			Name:       "user",
			Value:      db_uuid,
			Path:       "/",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   []string{},
		}
		if keepAlive == "on" {
			cookie.Expires = time.Now().AddDate(20, 0, 0)
		}
		http.SetCookie(w, &cookie)

		_, err = db.Exec("UPDATE user SET last_seen = ? WHERE email = ?", time.Now(), email)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
	}

	err := tplt.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
