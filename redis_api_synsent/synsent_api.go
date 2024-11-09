package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
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

type ExternalAPIResponse struct {
	ID string `json:"id"`
	IP struct {
		String string `json:"String"`
		Valid  bool   `json:"Valid"`
	} `json:"ip"`
}

var ctx = context.Background()

func main() {
	// Set up the database connection
	db, err := sql.Open("mysql", "username:password@tcp(IP:3306)/db_name")
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Set up Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust this if your Redis is on a different host/port
	})

	http.HandleFunc("/data2", getData(db, rdb))
	log.Fatal(http.ListenAndServe(":port", nil))
}

func getData(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for /data2")

		// Try to get data from Redis cache
		cachedData, err := rdb.Get(ctx, "data2_cache").Result()
		if err == redis.Nil {
			// Cache miss, fetch data from database and external API
			log.Println("Cache miss, fetching data from database and external API")
			data, err := fetchDataFromDBAndAPI(db)
			if err != nil {
				log.Printf("Error fetching data: %v", err)
				http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
				return
			}

			// Marshal data to JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling JSON: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Cache the data in Redis for 10 minutes
			err = rdb.Set(ctx, "data2_cache", jsonData, 10*time.Minute).Err()
			if err != nil {
				log.Printf("Error caching data in Redis: %v", err)
			}

			sendJSONResponse(w, jsonData)
		} else if err != nil {
			log.Printf("Error retrieving data from Redis: %v", err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
			return
		} else {
			// Cache hit, return cached data
			log.Println("Cache hit, returning data from Redis")
			sendJSONResponse(w, []byte(cachedData))
		}
	}
}

func fetchDataFromDBAndAPI(db *sql.DB) ([]Data, error) {
	// Fetch data from external API
	externalData, err := fetchExternalAPIData("http://IP:5010/ipunit")
	if err != nil {
		return nil, err
	}

	// Create a map of valid IDs from external API
	validIDUnits := make(map[string]bool)
	for _, entry := range externalData {
		validIDUnits[entry.ID] = true
	}

	// Fetch data from database
	rows, err := db.Query(getQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []Data

	for rows.Next() {
		var d Data
		err := rows.Scan(&d.ID, &d.DateTime, &d.IDUnit, &d.IPUnit, &d.ForeignAddr, &d.StatusID)
		if err != nil {
			return nil, err
		}
		// Filter data based on validIDUnits
		if validIDUnits[d.IDUnit] {
			data = append(data, d)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func getQuery() string {
	return `
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
}

func fetchExternalAPIData(url string) ([]ExternalAPIResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []ExternalAPIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func sendJSONResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
	log.Println("Response successfully written")
}
