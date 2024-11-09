package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"golang.org/x/crypto/ssh"
)

// Server represents a remote server with its IP address and alias
type Server struct {
	IP    IPField `json:"ip"`
	Alias string  `json:"id"`
}

type IPField struct {
	String string `json:"String"`
	Valid  bool   `json:"Valid"`
}

// Configurations
const (
	maxConcurrentConnections = 100              // Max number of concurrent SSH connections
	maxRetries               = 2                // Max number of retries for SSH connection
	batchSize                = 50               // Number of records to insert in a single batch
	sshTimeout               = 10 * time.Second // SSH connection timeout
)

var (
	defaultUsername = "username"
	defaultPassword = "password"
)

func main() {
	// Open MySQL database
	db, err := sql.Open("mysql", "username:password@tcp(ip:port)/db_name")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create server data table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS display_status (
        id INT AUTO_INCREMENT PRIMARY KEY,
		date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        id_unit VARCHAR(255),
        ip_unit VARCHAR(255),
        foreign_address VARCHAR(255),
        status VARCHAR(255)
    );`)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch server list from API
	servers, err := fetchServerList("http://ip:port/ipunit")
	if err != nil {
		log.Fatalf("Failed to fetch server list: %v", err)
	}

	// Create a buffered channel to limit concurrent connections
	concurrencyLimiter := make(chan struct{}, maxConcurrentConnections)

	var wg sync.WaitGroup

	for _, server := range servers {
		wg.Add(1)
		concurrencyLimiter <- struct{}{} // Acquire a token
		go func(server Server) {
			defer wg.Done()
			connectToServer(db, server, defaultUsername, defaultPassword)
			<-concurrencyLimiter // Release the token
		}(server)
	}

	wg.Wait()
}

func fetchServerList(apiURL string) ([]Server, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server list: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var servers []Server
	err = json.Unmarshal(body, &servers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return servers, nil
}

func connectToServer(db *sql.DB, server Server, username, password string) {
	if !server.IP.Valid {
		log.Printf("Invalid IP for %s (%s)", server.Alias, server.IP.String)
		insertDataToDatabase(db, server, "", "Invalid IP")
		return
	}

	retryCount := 0

	for {
		fmt.Printf("Connecting to %s (%s)...\n", server.Alias, server.IP.String)

		// SSH connection configuration with password authentication
		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use only in testing environments
			Timeout:         sshTimeout,
		}

		// Connect to the remote server
		client, err := ssh.Dial("tcp", server.IP.String+":22", config)
		if err != nil {
			log.Printf("Failed to dial to %s (%s): %v\n", server.Alias, server.IP.String, err)
			retryCount++
			if retryCount >= maxRetries {
				insertDataToDatabase(db, server, "", "Failed to Connect")
				return
			}
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		// Create a session
		session, err := client.NewSession()
		if err != nil {
			log.Printf("Failed to create session for %s (%s): %v\n", server.Alias, server.IP.String, err)
			client.Close()
			retryCount++
			if retryCount >= maxRetries {
				insertDataToDatabase(db, server, "", "Failed to Connect")
				return
			}
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		// Execute the netstat command
		output, err := session.CombinedOutput("netstat")
		if err != nil {
			log.Printf("Failed to execute command on %s (%s): %v\n", server.Alias, server.IP.String, err)
			session.Close()
			client.Close()
			retryCount++
			if retryCount >= maxRetries {
				insertDataToDatabase(db, server, "", "Failed to Execute Command")
				return
			}
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		// Close session and client
		session.Close()
		client.Close()

		// Process the output to find lines containing "master" in foreign address column
		var foreignAddress, statusOutput string
		lines := bytes.Split(output, []byte("\n"))
		for _, line := range lines {
			if strings.Contains(string(line), "master") {
				parts := strings.Fields(string(line)) // Split by any whitespace
				for _, part := range parts {
					if strings.Contains(part, "master") {
						foreignAddress = part
					} else if part == "ESTABLISHED" || part == "SYN_SENT" {
						statusOutput = part
					}
				}
				break
			}
		}

		// Store data in the database
		insertDataToDatabase(db, server, foreignAddress, statusOutput)

		return
	}
}

func insertDataToDatabase(db *sql.DB, server Server, foreignAddress, statusOutput string) {
	_, err := db.Exec("INSERT INTO display_status (id_unit, ip_unit, foreign_address, status) VALUES (?, ?, ?, ?)", server.Alias, server.IP.String, foreignAddress, statusOutput)
	if err != nil {
		log.Printf("Failed to insert data for %s (%s) into database: %v\n", server.Alias, server.IP.String, err)
	} else {
		log.Printf("Data inserted successfully for %s (%s) into database\n", server.Alias, server.IP.String)
	}
}
