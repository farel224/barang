package siswa

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type mahasiswa struct {
	Nim    string
	Nama   string
	Progdi string
	Smt    string
}

type response struct {
	Status bool
	Pesan  string
	Data   []mahasiswa
}

func koneksi() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func tampil(pesan string) response {
	db, err := koneksi()
	if err != nil {
		return response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM mahasiswa")
	if err != nil {
		return response{false, "Gagal Query: " + err.Error(), nil}
	}
	defer rows.Close()

	var hasil []mahasiswa
	for rows.Next() {
		var mhs mahasiswa
		if err := rows.Scan(&mhs.Nim, &mhs.Nama, &mhs.Progdi, &mhs.Smt); err != nil {
			return response{false, "Gagal Baca: " + err.Error(), nil}
		}
		hasil = append(hasil, mhs)
	}

	if err = rows.Err(); err != nil {
		return response{false, "Kesalahan: " + err.Error(), nil}
	}
	return response{true, pesan, hasil}
}

func getMhs(nim string) response {
	db, err := koneksi()
	if err != nil {
		return response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM mahasiswa WHERE nim = ?", nim)
	var mhs mahasiswa
	if err := row.Scan(&mhs.Nim, &mhs.Nama, &mhs.Progdi, &mhs.Smt); err != nil {
		return response{false, "Data tidak ditemukan: " + err.Error(), nil}
	}
	return response{true, "Berhasil Tampil", []mahasiswa{mhs}}
}

func tambah(nim, nama, progdi, smt string) response {
	db, err := koneksi()
	if err != nil {
		return response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO mahasiswa VALUES (?, ?, ?, ?)", nim, nama, progdi, smt)
	if err != nil {
		return response{false, "Gagal Insert: " + err.Error(), nil}
	}
	return response{true, "Berhasil Tambah", nil}
}

func ubah(nim, nama, progdi, smt string) response {
	db, err := koneksi()
	if err != nil {
		return response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	_, err = db.Exec("UPDATE mahasiswa SET nama=?, progdi=?, smt=? WHERE nim=?", nama, progdi, smt, nim)
	if err != nil {
		return response{false, "Gagal Update: " + err.Error(), nil}
	}
	return response{true, "Berhasil Ubah", nil}
}

func hapus(nim string) response {
	db, err := koneksi()
	if err != nil {
		return response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM mahasiswa WHERE nim=?", nim)
	if err != nil {
		return response{false, "Gagal Delete: " + err.Error(), nil}
	}
	return response{true, "Berhasil Hapus", nil}
}

func Kontroler(w http.ResponseWriter, r *http.Request) {
	tmplIndex, err := template.ParseFiles("siswa/template/index.html")
	if err != nil {
		http.Error(w, "Template index error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpltamb, err := template.ParseFiles("siswa/template/tamb.html")
	if err != nil {
		http.Error(w, "Template tambah error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpledit, err := template.ParseFiles("siswa/template/edit.html")
	if err != nil {
		http.Error(w, "Template edit error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmplhp, err := template.ParseFiles("siswa/template/hp.html") // Perbaikan di sini
	if err != nil {
		http.Error(w, "Template hapus error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	aksi := r.URL.Query().Get("aksi")

	switch r.Method {
	case http.MethodGet:
		switch aksi {
		case "tamb":
			tmpltamb.Execute(w, nil)
		case "edit":
			nim := r.URL.Query().Get("nim")
			tmpledit.Execute(w, getMhs(nim))
		case "hp":
			nim := r.URL.Query().Get("nim")
			tmplhp.Execute(w, getMhs(nim))
		default:
			tmplIndex.Execute(w, tampil("Berhasil Tampil"))
		}

	case http.MethodPost:
		r.ParseForm()
		nim := r.FormValue("nim")
		nama := r.FormValue("nama")
		progdi := r.FormValue("progdi")
		smt := r.FormValue("smt")

		var hasil response
		switch aksi {
		case "tamb":
			hasil = tambah(nim, nama, progdi, smt)
		case "edit":
			hasil = ubah(nim, nama, progdi, smt)
		case "hp":
			hasil = hapus(nim)
		default:
			hasil = tampil("Aksi tidak dikenali")
		}
		tmplIndex.Execute(w, tampil(hasil.Pesan))

	default:
		fmt.Fprint(w, "Hanya mendukung metode GET dan POST")
	}
}
