package main

import (
	"PERCOBAAN/login"
	"PERCOBAAN/pinjam"
	"PERCOBAAN/siswa"
	"net/http"
)

func main() {
	login.CekLogin()
	http.HandleFunc("/pinjam", pinjam.Kontroler)
	http.HandleFunc("/siswa", siswa.Kontroler)

	http.ListenAndServe(":8000", nil)
}
