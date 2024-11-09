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
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/db_name")
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(db)

	http.HandleFunc("/data2", getData(db))
	log.Fatal(http.ListenAndServe(":port", nil))
}

func getData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /data2")
		// Optimized query to get the first syn_sent after the last established and the latest status for each id_unit
		query := `
			WITH LastEstablished AS (
				SELECT id_unit, MAX(date_time) AS last_established
				FROM display_status
				WHERE status = 'established'
				GROUP BY id_unit
			),
			RankedSynSent AS (
				SELECT ds.id, ds.date_time, ds.id_unit, ds.ip_unit, ds.foreign_address, ds.status,
					   ROW_NUMBER() OVER (PARTITION BY ds.id_unit ORDER BY ds.date_time) AS rn
				FROM display_status ds
				INNER JOIN LastEstablished le ON ds.id_unit = le.id_unit
				WHERE ds.date_time > le.last_established AND ds.status = 'syn_sent'
			),
			FirstSynSent AS (
				SELECT id, date_time, id_unit, ip_unit, foreign_address, status
				FROM RankedSynSent
				WHERE rn = 1
			),
			RankedLatestStatus AS (
				SELECT ds.id, ds.date_time, ds.id_unit, ds.ip_unit, ds.foreign_address, ds.status,
					   ROW_NUMBER() OVER (PARTITION BY ds.id_unit ORDER BY ds.date_time DESC) AS rn
				FROM display_status ds
			),
			LatestStatus AS (
				SELECT id, date_time, id_unit, ip_unit, foreign_address, status
				FROM RankedLatestStatus
				WHERE rn = 1
			)
			SELECT id, date_time, id_unit, ip_unit, foreign_address, status
			FROM FirstSynSent
			UNION
			SELECT id, date_time, id_unit, ip_unit, foreign_address, status
			FROM LatestStatus
			WHERE id_unit NOT IN (SELECT id_unit FROM FirstSynSent)
			ORDER BY id_unit, date_time;
		`

		log.Println("Executing query")
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error executing query: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Printf("Error closing rows: %v", err)
			}
		}(rows)

		var data []Data

		log.Println("Processing query results")
		for rows.Next() {
			var d Data
			err := rows.Scan(&d.ID, &d.DateTime, &d.IDUnit, &d.IPUnit, &d.ForeignAddr, &d.StatusID)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data = append(data, d)
		}

		err = rows.Err()
		if err != nil {
			log.Printf("Error iterating rows: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Marshaling JSON response")
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			log.Printf("Error writing response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Response successfully written")
	}
}
