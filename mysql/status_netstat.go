package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"golang.org/x/crypto/ssh"
)

// Server represents a remote server with its IP address, username, password, and alias
type Server struct {
	IP       string
	Username string
	Password string
	Alias    string
}

func main() {
	// Open MySQL database
	db, err := sql.Open("mysql", "username:password@tcp(ip:3306)/db_name")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

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

	// List of servers with their IP addresses, usernames, passwords, and aliases
	servers := []Server{
		{"IP", "User", "Pass", "ID_unit"},
		// Add more servers as needed
	}

	var wg sync.WaitGroup
	wg.Add(len(servers))

	for _, server := range servers {
		go func(server Server) {
			defer wg.Done()
			connectToServer(db, server)
		}(server)
	}

	wg.Wait()
}

func connectToServer(db *sql.DB, server Server) {
	var maxRetries = 1
	var retryCount = 0

	for {
		fmt.Printf("Connecting to %s (%s)...\n", server.Alias, server.IP)

		// SSH connection configuration with password authentication
		config := &ssh.ClientConfig{
			User: server.Username,
			Auth: []ssh.AuthMethod{
				ssh.Password(server.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use only in testing environments
		}

		// Connect to the remote server
		client, err := ssh.Dial("tcp", server.IP+":22", config)
		if err != nil {
			log.Printf("Failed to dial to %s (%s): %v\n", server.Alias, server.IP, err)
			retryCount++
			if retryCount >= maxRetries {
				insertDataToDatabase(db, server, "", "Failed to Connect")
				return
			}
			time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
			continue
		}

		// Create a session
		session, err := client.NewSession()
		if err != nil {
			log.Printf("Failed to create session for %s (%s): %v\n", server.Alias, server.IP, err)
			err := client.Close()
			if err != nil {
				return
			}
			retryCount++
			if retryCount >= maxRetries {
				insertDataToDatabase(db, server, "", "Failed to Connect")
				return
			}
			time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
			continue
		}

		// Execute the netstat command
		output, err := session.CombinedOutput("netstat")
		if err != nil {
			log.Printf("Failed to execute command on %s (%s): %v\n", server.Alias, server.IP, err)
			err := client.Close()
			if err != nil {
				return
			}
			retryCount++
			if retryCount >= maxRetries {
				insertDataToDatabase(db, server, "", "")
				return
			}
			time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
			continue
		}

		// Close session and client
		err = session.Close()
		if err != nil {
			return
		}
		err = client.Close()
		if err != nil {
			return
		}

		// Process the output to find lines containing "master" in foreign address column
		var foreignAddress, statusOutput string
		lines := bytes.Split(output, []byte("\n"))
		for _, line := range lines {
			if strings.Contains(string(line), "master") {
				parts := strings.Split(string(line), " ")
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

		// Store data in SQLite database
		insertDataToDatabase(db, server, foreignAddress, statusOutput)

		return
	}
}

func insertDataToDatabase(db *sql.DB, server Server, foreignAddress, statusOutput string) {
	_, err := db.Exec("INSERT INTO display_status (id_unit, ip_unit, foreign_address, status) VALUES (?, ?, ?, ?)", server.Alias, server.IP, foreignAddress, statusOutput)
	if err != nil {
		log.Printf("Failed to insert data for %s (%s) into database: %v\n", server.Alias, server.IP, err)
	} else {
		log.Printf("Data inserted successfully for %s (%s) into database\n", server.Alias, server.IP)
	}
}
