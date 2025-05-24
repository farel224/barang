package login

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	Username string
	Email    string
}

var db *sql.DB

// Fungsi ini akan dipanggil dari main.go untuk setup semua handler login
func CekLogin() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic("Gagal koneksi ke database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("Gagal ping database: " + err.Error())
	}

	// Daftarkan semua handler login
	http.Handle("/foto/", http.StripPrefix("/foto/", http.FileServer(http.Dir("foto"))))
	http.HandleFunc("/", renderLoginPage)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/dashboard", renderDashboard)
	http.HandleFunc("/barang", renderBarangPage)
	http.HandleFunc("/logout", handleLogout)

	fmt.Println("Login handler berhasil disiapkan")
}

// halaman login
func renderLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("username")
	if err == nil && cookie.Value != "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("login/template/login.html")
	if err != nil {
		http.Error(w, "Gagal membuka login.html: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Gagal render template login: "+err.Error(), http.StatusInternalServerError)
	}
}

// proses login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	identity := r.FormValue("identity")
	password := r.FormValue("password")

	var username, email, storedPassword string
	err := db.QueryRow("SELECT username, email, password FROM users WHERE username = ? OR email = ?", identity, identity).
		Scan(&username, &email, &storedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Login gagal: username/email tidak ditemukan", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error query database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if password != storedPassword {
		http.Error(w, "Login gagal: password salah", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// render dashboard
func renderDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := cookie.Value
	var email string
	err = db.QueryRow("SELECT email FROM users WHERE username = ?", username).Scan(&email)
	if err != nil {
		http.Error(w, "Gagal mengambil data user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := UserInfo{
		Username: username,
		Email:    email,
	}

	tmpl, err := template.ParseFiles("login/template/dashboard.html")
	if err != nil {
		http.Error(w, "Gagal membuka dashboard.html: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Gagal render template dashboard: "+err.Error(), http.StatusInternalServerError)
	}
}

// render halaman barang
func renderBarangPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("login/template/barang.html")
	if err != nil {
		http.Error(w, "Gagal membuka barang.html: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Gagal render template barang: "+err.Error(), http.StatusInternalServerError)
	}
}

// logout
func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
