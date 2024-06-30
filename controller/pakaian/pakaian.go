package pakaian

import (
	"encoding/json"
	"net/http"
	"strconv"

	"belajar/database"
	"github.com/gorilla/mux"
	"belajar/model/pakaian"
)

func GetPakaian(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM pakaian")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pakaianList []pakaian.Pakaian
	for rows.Next() {
		var p pakaian.Pakaian
		if err := rows.Scan(&p.PakaianId, &p.NamaPakaian, &p.Deskripsi, &p.UkuranBaju, &p.StokBaju, &p.HargaPakaian); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pakaianList = append(pakaianList, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pakaianList)
}

func PostPakaian(w http.ResponseWriter, r *http.Request) {
	var p pakaian.Pakaian
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO pakaian (namapakaian, deskripsi, ukuran, stok, hargapakaian) 
		VALUES (?, ?, ?, ?, ?)`

	res, err := database.DB.Exec(query, p.NamaPakaian, p.Deskripsi, p.UkuranBaju, p.StokBaju, p.HargaPakaian)
	if err != nil {
		http.Error(w, "Failed to insert pakaian: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pakaian added successfully",
		"id":      id,
	})
}

func PutPakaian(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p pakaian.Pakaian
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE pakaian 
		SET namapakaian=?, deskripsi=?, ukuran=?, stok=?, hargapakaian=?
		WHERE id_pakaian=?`

	result, err := database.DB.Exec(query, p.NamaPakaian, p.Deskripsi, p.UkuranBaju, p.StokBaju, p.HargaPakaian, id)
	if err != nil {
		http.Error(w, "Failed to update pakaian: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pakaian updated successfully",
	})
}

func DeletePakaian(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM pakaian
		WHERE id_pakaian = ?`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete pakaian: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pakaian deleted successfully",
	})
}
