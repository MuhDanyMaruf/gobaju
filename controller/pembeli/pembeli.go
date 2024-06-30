package pembeli

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"belajar/database"
	"belajar/model/pembeli"

	"github.com/gorilla/mux"
)

func GetPembeli(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM pembeli")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pembeliList []pembeli.Pembeli
	for rows.Next() {
		var p pembeli.Pembeli
		if err := rows.Scan(&p.PembeliId, &p.Nama, &p.Email, &p.JumlahPembelian); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pembeliList = append(pembeliList, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pembeliList)
}

func PostPembeli(w http.ResponseWriter, r *http.Request) {
	var p pembeli.Pembeli
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO pembeli (nama, email, jumlahpembelian) 
		VALUES (?, ?, ?)`

	res, err := database.DB.Exec(query, p.Nama, p.Email, p.JumlahPembelian)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
		http.Error(w, "Failed to insert pembeli: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pembeli added successfully",
		"id":      id,
	})
}

func PutPembeli(w http.ResponseWriter, r *http.Request) {
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

	var p pembeli.Pembeli
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE pembeli 
		SET nama=?, email=?, jumlahpembelian=?
		WHERE id_pembeli=?`

	result, err := database.DB.Exec(query, p.Nama, p.Email, p.JumlahPembelian, id)
	if err != nil {
		log.Printf("Error updating database: %v", err)
		http.Error(w, "Failed to update pembeli: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving affected rows: %v", err)
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pembeli updated successfully",
	})
}

func DeletePembeli(w http.ResponseWriter, r *http.Request) {
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
		DELETE FROM pembeli
		WHERE id_pembeli = ?`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting from database: %v", err)
		http.Error(w, "Failed to delete pembeli: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving affected rows: %v", err)
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pembeli deleted successfully",
	})
}
