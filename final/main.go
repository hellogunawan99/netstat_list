package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"

	_ "github.com/mattn/go-sqlite3"
)

// Global variables for SSH credentials
var (
	username = "your_ssh_username"
	password = "your_ssh_password"
)

// Server represents a remote server with its IP address, and alias
type Server struct {
	IP    string
	Alias string
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./server_data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create server data table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS server_data (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						alias TEXT,
						ip TEXT,
						foreign_address TEXT,
						status_output TEXT
					);`)
	if err != nil {
		log.Fatal(err)
	}

	// List of servers with their IP addresses and aliases
	servers := []Server{
		{"server1_ip", "server1_alias"},
		{"server2_ip", "server2_alias"},
		// Add more servers as needed
	}

	// Loop over each server
	for _, server := range servers {
		fmt.Printf("Connecting to %s (%s)...\n", server.Alias, server.IP)

		// SSH connection configuration with password authentication
		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use only in testing environments
		}

		// Connect to the remote server
		client, err := ssh.Dial("tcp", server.IP+":22", config)
		if err != nil {
			log.Printf("Failed to dial to %s (%s): %v\n", server.Alias, server.IP, err)
			continue
		}

		// Create a session
		session, err := client.NewSession()
		if err != nil {
			log.Printf("Failed to create session for %s (%s): %v\n", server.Alias, server.IP, err)
			client.Close()
			continue
		}

		// Execute the netstat command
		output, err := session.CombinedOutput("netstat")
		if err != nil {
			log.Printf("Failed to execute command on %s (%s): %v\n", server.Alias, server.IP, err)
			client.Close()
			continue
		}

		// Close session and client
		session.Close()
		client.Close()

		// Process the output to find lines containing "bootps" in foreign address column
		var foreignAddress, statusOutput string
		lines := bytes.Split(output, []byte("\n"))
		for _, line := range lines {
			if strings.Contains(string(line), "bootps") {
				foreignAddress = string(line)
				break
			}
		}
		statusOutput = string(output)

		// Store data in SQLite database
		_, err = db.Exec("INSERT INTO server_data (alias, ip, foreign_address, status_output) VALUES (?, ?, ?, ?)",
			server.Alias, server.IP, foreignAddress, statusOutput)
		if err != nil {
			log.Printf("Failed to insert data for %s (%s) into database: %v\n", server.Alias, server.IP, err)
		} else {
			log.Printf("Data inserted successfully for %s (%s) into database\n", server.Alias, server.IP)
		}
	}
}
