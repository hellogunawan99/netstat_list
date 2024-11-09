package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

// Data represents the structure of the data to be returned as JSON
type Data struct {
	ID          int    `json:"id"`
	DateTime    string `json:"date_time"`
	IDUnit      string `json:"id_unit"`
	IPUnit      string `json:"ip_unit"`
	ForeignAddr string `json:"foreign_address"`
	StatusID    string `json:"status"`
}

var (
	// ctx is a context used for Redis operations
	ctx = context.Background()
)

func main() {
	// Set up the database connection
	db, err := sql.Open("mysql", "username:password@tcp(ip:3306)/db_name")
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Set up the Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Define the HTTP handler for /data2 endpoint
	http.HandleFunc("/redisapi", getData(db, rdb))
	log.Fatal(http.ListenAndServe(":port", nil))
}

// getData is the HTTP handler function that fetches data from the database or cache
func getData(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /data2")

		// Check cache first
		cachedData, err := rdb.Get(ctx, "data2_cache").Result()
		if err == redis.Nil {
			// Cache miss, fetching data from database
			log.Println("Cache miss, fetching data from database")

			// SQL query to fetch data from the database
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
			defer rows.Close()

			var data []Data

			log.Println("Processing query results")
			for rows.Next() {
				var d Data
				// Scan the result into the Data struct
				err := rows.Scan(&d.ID, &d.DateTime, &d.IDUnit, &d.IPUnit, &d.ForeignAddr, &d.StatusID)
				if err != nil {
					log.Printf("Error scanning row: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data = append(data, d)
			}

			// Check for any error encountered during iteration
			err = rows.Err()
			if err != nil {
				log.Printf("Error iterating rows: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("Marshaling JSON response")
			// Convert the data to JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling JSON: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Cache the JSON data with a 10-minute expiration
			err = rdb.Set(ctx, "data2_cache", jsonData, 10*time.Minute).Err()
			if err != nil {
				log.Printf("Error caching data: %v", err)
			}

			// Set the response header and write the JSON data to the response
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(jsonData)
			if err != nil {
				log.Printf("Error writing response: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("Response successfully written from database")
		} else if err != nil {
			// Error accessing Redis
			log.Printf("Error getting cache: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// Cache hit, returning data from cache
			log.Println("Cache hit, returning data from cache")

			// Set the response header and write the cached data to the response
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write([]byte(cachedData))
			if err != nil {
				log.Printf("Error writing cached response: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("Response successfully written from cache")
		}
	}
}
