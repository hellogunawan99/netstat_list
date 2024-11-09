package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	ID          int    `json:"id"`
	DateTime    string `json:"date_time"`
	IDUnit      string `json:"id_unit"`
	IPUnit      string `json:"ip_unit"`
	ForeignAddr string `json:"foreign_address"`
	StatusID    string `json:"status"`
}

func main() {
	// Set up the database connection
	db, err := sql.Open("mysql", "username:password@tcp(ip:3306)/db_name")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	http.HandleFunc("/data2", getData(db))
	log.Fatal(http.ListenAndServe(":port", nil))
}

func getData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the latest data for each different id_unit
		query := `
			SELECT id, date_time, id_unit, ip_unit, foreign_address, status
			FROM display_status
			WHERE status IN ('SYN_SENT', 'ESTABLISHED', 'Failed to Connect', '')
			ORDER BY date_time DESC;
		`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		var data []Data
		seen := make(map[string]bool)

		for rows.Next() {
			var d Data
			err := rows.Scan(&d.ID, &d.DateTime, &d.IDUnit, &d.IPUnit, &d.ForeignAddr, &d.StatusID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !seen[d.IDUnit] && (d.StatusID == "SYN_SENT" || d.StatusID == "ESTABLISHED" || d.StatusID == "Failed to Connect" || d.StatusID == "") {
				if d.StatusID == "" {
					d.StatusID = "Netstat not detect Master"
				}
				data = append(data, d)
				seen[d.IDUnit] = true
			}
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			return
		}
	}
}
