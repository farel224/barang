package pinjam

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Peminjaman struct {
	ID                  string
	NamaPeminjam        string
	NamaBarang          string
	Jumlah              string
	TanggalPengembalian string
	TanggalPeminjaman   string
}

type Response struct {
	Status bool
	Pesan  string
	Data   []Peminjaman
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

func tampil(pesan string) Response {
	db, err := koneksi()
	if err != nil {
		return Response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM peminjaman")
	if err != nil {
		return Response{false, "Gagal query: " + err.Error(), nil}
	}
	defer rows.Close()

	var hasil []Peminjaman
	for rows.Next() {
		var p Peminjaman
		if err := rows.Scan(&p.ID, &p.NamaPeminjam, &p.NamaBarang, &p.Jumlah, &p.TanggalPengembalian, &p.TanggalPeminjaman); err != nil {
			return Response{false, "Gagal baca: " + err.Error(), nil}
		}
		hasil = append(hasil, p)
	}

	return Response{true, pesan, hasil}
}

func getP(id string) Response {
	db, err := koneksi()
	if err != nil {
		return Response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM peminjaman WHERE id = ?", id)

	var p Peminjaman
	if err := row.Scan(&p.ID, &p.NamaPeminjam, &p.NamaBarang, &p.Jumlah, &p.TanggalPengembalian, &p.TanggalPeminjaman); err != nil {
		return Response{false, "Data tidak ditemukan: " + err.Error(), nil}
	}

	return Response{true, "Berhasil", []Peminjaman{p}}
}

func tambah(id, namapeminjam, namabarang, jumlah, tanggalpengembalian, tanggalpeminjaman string) Response {
	db, err := koneksi()
	if err != nil {
		return Response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO peminjaman VALUES (?, ?, ?, ?, ?, ?)", id, namapeminjam, namabarang, jumlah, tanggalpengembalian, tanggalpeminjaman)
	if err != nil {
		return Response{false, "Gagal insert: " + err.Error(), nil}
	}

	return Response{true, "Berhasil tambah", nil}
}

func ubah(id, namapeminjam, namabarang, jumlah, tanggalpengembalian, tanggalpeminjaman string) Response {
	db, err := koneksi()
	if err != nil {
		return Response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	_, err = db.Exec("UPDATE peminjaman SET namapeminjam=?, namabarang=?, jumlah=?, tanggalpengembalian=?, tanggalpeminjaman=? WHERE id=?",
		namapeminjam, namabarang, jumlah, tanggalpengembalian, tanggalpeminjaman, id)
	if err != nil {
		return Response{false, "Gagal update: " + err.Error(), nil}
	}

	return Response{true, "Berhasil ubah", nil}
}

func hapus(id string) Response {
	db, err := koneksi()
	if err != nil {
		return Response{false, "Gagal koneksi: " + err.Error(), nil}
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM peminjaman WHERE id=?", id)
	if err != nil {
		return Response{false, "Gagal delete: " + err.Error(), nil}
	}

	return Response{true, "Berhasil hapus", nil}
}

// Fungsi utama sebagai handler HTTP
func Kontroler(w http.ResponseWriter, r *http.Request) {
	tmplIndex := template.Must(template.ParseFiles("pinjam/template/index.html"))
	tmplTambah := template.Must(template.ParseFiles("pinjam/template/tambah.html"))
	tmplEdit := template.Must(template.ParseFiles("pinjam/template/edit.html"))
	tmplHapus := template.Must(template.ParseFiles("pinjam/template/hapus.html"))

	aksi := r.URL.Query().Get("aksi")

	switch r.Method {
	case http.MethodGet:
		switch aksi {
		case "tambah":
			tmplTambah.Execute(w, nil)
		case "edit":
			id := r.URL.Query().Get("id")
			tmplEdit.Execute(w, getP(id))
		case "hapus":
			id := r.URL.Query().Get("id")
			tmplHapus.Execute(w, getP(id))
		default:
			tmplIndex.Execute(w, tampil("Berhasil tampil"))
		}

	case http.MethodPost:
		r.ParseForm()
		id := r.FormValue("id")
		namapeminjam := r.FormValue("namapeminjam")
		namabarang := r.FormValue("namabarang")
		jumlah := r.FormValue("jumlah")
		tanggalpengembalian := r.FormValue("tanggalpengembalian")
		tanggalpeminjaman := r.FormValue("tanggalpeminjaman")

		var hasil Response
		switch aksi {
		case "tambah":
			hasil = tambah(id, namapeminjam, namabarang, jumlah, tanggalpengembalian, tanggalpeminjaman)
		case "edit":
			hasil = ubah(id, namapeminjam, namabarang, jumlah, tanggalpengembalian, tanggalpeminjaman)
		case "hapus":
			hasil = hapus(id)
		default:
			hasil = tampil("Aksi tidak dikenali")
		}

		tmplIndex.Execute(w, tampil(hasil.Pesan))
	default:
		fmt.Fprint(w, "Hanya mendukung metode GET dan POST")
	}
}
