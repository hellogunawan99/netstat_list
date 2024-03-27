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
	db, err := sql.Open("mysql", "gunawan:CANcer99@tcp(127.0.0.1:3306)/netstat_status")
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
		{"172.16.100.49", "jigsaw", "jigsaw", "D3-049"},
		{"172.16.100.64", "jigsaw", "jigsaw", "D3-064"},
		{"172.16.100.81", "jigsaw", "jigsaw", "D3-081"},
		{"172.16.100.84", "jigsaw", "jigsaw", "D3-084"},
		{"172.16.100.85", "jigsaw", "jigsaw", "D3-085"},
		{"172.16.100.86", "jigsaw", "jigsaw", "D3-086"},
		{"172.16.100.87", "jigsaw", "jigsaw", "D3-087"},
		{"172.16.100.88", "jigsaw", "jigsaw", "D3-088"},
		{"172.16.100.91", "jigsaw", "jigsaw", "D3-091"},
		{"172.16.100.101", "jigsaw", "jigsaw", "D3-101"},
		{"172.16.100.102", "jigsaw", "jigsaw", "D3-102"},
		{"172.16.100.108", "jigsaw", "jigsaw", "D3-108"},
		{"172.16.100.112", "jigsaw", "jigsaw", "D3-112"},
		{"172.16.100.113", "jigsaw", "jigsaw", "D3-113"},
		{"172.16.100.115", "jigsaw", "jigsaw", "D3-115"},
		{"172.16.100.124", "jigsaw", "jigsaw", "D3-124"},
		{"172.16.100.132", "jigsaw", "jigsaw", "D3-132"},
		{"172.16.100.134", "jigsaw", "jigsaw", "D3-134"},
		{"172.16.100.141", "jigsaw", "jigsaw", "D3-141"},
		{"172.16.100.142", "jigsaw", "jigsaw", "D3-142"},
		{"172.16.100.143", "jigsaw", "jigsaw", "D3-143"},
		{"172.16.100.146", "jigsaw", "jigsaw", "D3-146"},
		{"172.16.100.147", "jigsaw", "jigsaw", "D3-147"},
		{"172.16.100.148", "jigsaw", "jigsaw", "D3-148"},
		{"172.16.100.149", "jigsaw", "jigsaw", "D3-149"},
		{"172.16.100.150", "jigsaw", "jigsaw", "D3-150"},
		{"172.16.102.154", "jigsaw", "jigsaw", "D3-154"},
		{"172.16.100.155", "jigsaw", "jigsaw", "D3-155"},
		{"172.16.100.158", "jigsaw", "jigsaw", "D3-158"},
		{"172.16.100.161", "jigsaw", "jigsaw", "D3-161"},
		{"172.16.100.168", "jigsaw", "jigsaw", "D3-168"},
		{"172.16.100.171", "jigsaw", "jigsaw", "D3-171"},
		{"172.16.100.173", "jigsaw", "jigsaw", "D3-173"},
		{"172.16.100.174", "jigsaw", "jigsaw", "D3-174"},
		{"172.16.100.175", "jigsaw", "jigsaw", "D3-175"},
		{"172.16.100.176", "jigsaw", "jigsaw", "D3-176"},
		{"172.16.100.177", "jigsaw", "jigsaw", "D3-177"},
		{"172.16.100.179", "jigsaw", "jigsaw", "D3-179"},
		{"172.16.100.181", "jigsaw", "jigsaw", "D3-181"},
		{"172.16.100.183", "jigsaw", "jigsaw", "D3-183"},
		{"172.16.100.186", "jigsaw", "jigsaw", "D3-186"},
		{"172.16.100.196", "jigsaw", "jigsaw", "D3-196"},
		{"172.16.100.197", "jigsaw", "jigsaw", "D3-197"},
		{"172.16.100.198", "jigsaw", "jigsaw", "D3-198"},
		{"172.16.100.199", "jigsaw", "jigsaw", "D3-199"},
		{"172.16.100.200", "jigsaw", "jigsaw", "D3-200"},
		{"172.16.100.207", "jigsaw", "jigsaw", "D3-207"},
		{"172.16.100.208", "jigsaw", "jigsaw", "D3-208"},
		{"172.16.100.209", "jigsaw", "jigsaw", "D3-209"},
		{"172.16.100.210", "jigsaw", "jigsaw", "D3-210"},
		{"172.16.100.211", "jigsaw", "jigsaw", "D3-211"},
		{"172.16.100.212", "jigsaw", "jigsaw", "D3-212"},
		{"172.16.100.213", "jigsaw", "jigsaw", "D3-213"},
		{"172.16.100.214", "jigsaw", "jigsaw", "D3-214"},
		{"172.16.100.215", "jigsaw", "jigsaw", "D3-215"},
		{"172.16.100.219", "jigsaw", "jigsaw", "D3-219"},
		{"172.16.100.221", "jigsaw", "jigsaw", "D3-221"},
		{"172.16.100.222", "jigsaw", "jigsaw", "D3-222"},
		{"172.16.100.224", "jigsaw", "jigsaw", "D3-224"},
		{"172.16.100.229", "jigsaw", "jigsaw", "D3-229"},
		{"172.16.100.230", "jigsaw", "jigsaw", "D3-230"},
		{"172.16.100.231", "jigsaw", "jigsaw", "D3-231"},
		{"172.16.100.232", "jigsaw", "jigsaw", "D3-232"},
		{"172.16.100.233", "jigsaw", "jigsaw", "D3-233"},
		{"172.16.100.234", "jigsaw", "jigsaw", "D3-234"},
		{"172.16.101.165", "jigsaw", "jigsaw", "D3-245"},
		{"172.16.100.246", "jigsaw", "jigsaw", "D3-246"},
		{"172.16.100.247", "jigsaw", "jigsaw", "D3-247"},
		{"172.16.100.248", "jigsaw", "jigsaw", "D3-248"},
		{"172.16.100.250", "jigsaw", "jigsaw", "D3-250"},
		{"172.16.100.251", "jigsaw", "jigsaw", "D3-251"},
		{"172.16.100.252", "jigsaw", "jigsaw", "D3-252"},
		{"172.16.101.176", "jigsaw", "jigsaw", "D3-256"},
		{"172.16.101.177", "jigsaw", "jigsaw", "D3-257"},
		{"172.16.101.178", "jigsaw", "jigsaw", "D3-258"},
		{"172.16.101.180", "jigsaw", "jigsaw", "D3-260"},
		{"172.16.101.182", "jigsaw", "jigsaw", "D3-262"},
		{"172.16.101.184", "jigsaw", "jigsaw", "D3-264"},
		{"172.16.101.186", "jigsaw", "jigsaw", "D3-266"},
		{"172.16.101.189", "jigsaw", "jigsaw", "D3-269"},
		{"172.16.101.191", "jigsaw", "jigsaw", "D3-271"},
		{"172.16.101.192", "jigsaw", "jigsaw", "D3-272"},
		{"172.16.101.194", "jigsaw", "jigsaw", "D3-274"},
		{"172.16.101.197", "jigsaw", "jigsaw", "D3-277"},
		{"172.16.101.199", "jigsaw", "jigsaw", "D3-279"},
		{"172.16.101.200", "jigsaw", "jigsaw", "D3-280"},
		{"172.16.101.201", "jigsaw", "jigsaw", "D3-281"},
		{"172.16.101.206", "jigsaw", "jigsaw", "D3-286"},
		{"172.16.101.209", "jigsaw", "jigsaw", "D3-289"},
		{"172.16.101.210", "jigsaw", "jigsaw", "D3-290"},
		{"172.16.101.211", "jigsaw", "jigsaw", "D3-291"},
		{"172.16.101.212", "jigsaw", "jigsaw", "D3-292"},
		{"172.16.101.213", "jigsaw", "jigsaw", "D3-293"},
		{"172.16.101.214", "jigsaw", "jigsaw", "D3-294"},
		{"172.16.101.215", "jigsaw", "jigsaw", "D3-295"},
		{"172.16.101.224", "jigsaw", "jigsaw", "D3-304"},
		{"172.16.101.225", "jigsaw", "jigsaw", "D3-305"},
		{"172.16.101.228", "jigsaw", "jigsaw", "D3-308"},
		{"172.16.101.231", "jigsaw", "jigsaw", "D3-311"},
		{"172.16.101.240", "jigsaw", "jigsaw", "D3-320"},
		{"172.16.101.241", "jigsaw", "jigsaw", "D3-321"},
		{"172.16.101.242", "jigsaw", "jigsaw", "D3-322"},
		{"172.16.101.243", "jigsaw", "jigsaw", "D3-323"},
		{"172.16.101.244", "jigsaw", "jigsaw", "D3-324"},
		{"172.16.101.245", "jigsaw", "jigsaw", "D3-325"},
		{"172.16.101.246", "jigsaw", "jigsaw", "D3-326"},
		{"172.16.101.247", "jigsaw", "jigsaw", "D3-327"},
		{"172.16.101.248", "jigsaw", "jigsaw", "D3-328"},
		{"172.16.101.249", "jigsaw", "jigsaw", "D3-329"},
		{"172.16.101.250", "jigsaw", "jigsaw", "D3-330"},
		{"172.16.101.251", "jigsaw", "jigsaw", "D3-331"},
		{"172.16.101.252", "jigsaw", "jigsaw", "D3-332"},
		{"172.16.101.253", "jigsaw", "jigsaw", "D3-333"},
		{"172.16.101.254", "jigsaw", "jigsaw", "D3-334"},
		{"172.16.103.35", "jigsaw", "jigsaw", "D3-335"},
		{"172.16.103.37", "jigsaw", "jigsaw", "D3-337"},
		{"172.16.103.38", "jigsaw", "jigsaw", "D3-338"},
		{"172.16.103.39", "jigsaw", "jigsaw", "D3-339"},
		{"172.16.103.41", "jigsaw", "jigsaw", "D3-341"},
		{"172.16.103.48", "jigsaw", "jigsaw", "D3-348"},
		{"172.16.103.49", "jigsaw", "jigsaw", "D3-349"},
		{"172.16.103.50", "jigsaw", "jigsaw", "D3-350"},
		{"172.16.103.55", "jigsaw", "jigsaw", "D3-355"},
		{"172.16.103.56", "jigsaw", "jigsaw", "D3-356"},
		{"172.16.103.57", "jigsaw", "jigsaw", "D3-357"},
		{"172.16.103.58", "jigsaw", "jigsaw", "D3-358"},
		{"172.16.103.59", "jigsaw", "jigsaw", "D3-359"},
		{"172.16.103.60", "jigsaw", "jigsaw", "D3-360"},
		{"172.16.103.61", "jigsaw", "jigsaw", "D3-361"},
		{"172.16.103.63", "jigsaw", "jigsaw", "D3-363"},
		{"172.16.103.64", "jigsaw", "jigsaw", "D3-364"},
		{"172.16.103.70", "jigsaw", "jigsaw", "D3-370"},
		{"172.16.103.71", "jigsaw", "jigsaw", "D3-371"},
		{"172.16.103.72", "jigsaw", "jigsaw", "D3-372"},
		{"172.16.103.73", "jigsaw", "jigsaw", "D3-373"},
		{"172.16.103.74", "jigsaw", "jigsaw", "D3-374"},
		{"172.16.103.75", "jigsaw", "jigsaw", "D3-375"},
		{"172.16.103.76", "jigsaw", "jigsaw", "D3-376"},
		{"172.16.103.77", "jigsaw", "jigsaw", "D3-377"},
		{"172.16.103.78", "jigsaw", "jigsaw", "D3-378"},
		{"172.16.103.79", "jigsaw", "jigsaw", "D3-379"},
		{"172.16.103.80", "jigsaw", "jigsaw", "D3-380"},
		{"172.16.103.81", "jigsaw", "jigsaw", "D3-381"},
		{"172.16.103.82", "jigsaw", "jigsaw", "D3-382"},
		{"172.16.103.83", "jigsaw", "jigsaw", "D3-383"},
		{"172.16.103.84", "jigsaw", "jigsaw", "D3-384"},
		{"172.16.103.85", "jigsaw", "jigsaw", "D3-385"},
		{"172.16.103.86", "jigsaw", "jigsaw", "D3-386"},
		{"172.16.103.87", "jigsaw", "jigsaw", "D3-387"},
		{"172.16.103.88", "jigsaw", "jigsaw", "D3-388"},
		{"172.16.103.89", "jigsaw", "jigsaw", "D3-389"},
		{"172.16.103.90", "jigsaw", "jigsaw", "D3-390"},
		{"172.16.103.91", "jigsaw", "jigsaw", "D3-391"},
		{"172.16.103.92", "jigsaw", "jigsaw", "D3-392"},
		{"172.16.103.93", "jigsaw", "jigsaw", "D3-393"},
		{"172.16.103.94", "jigsaw", "jigsaw", "D3-394"},
		{"172.16.103.95", "jigsaw", "jigsaw", "D3-395"},
		{"172.16.103.96", "jigsaw", "jigsaw", "D3-396"},
		{"172.16.103.97", "jigsaw", "jigsaw", "D3-397"},
		{"172.16.103.98", "jigsaw", "jigsaw", "D3-398"},
		{"172.16.103.99", "jigsaw", "jigsaw", "D3-399"},
		{"172.16.103.100", "jigsaw", "jigsaw", "D3-400"},
		{"172.16.103.101", "jigsaw", "jigsaw", "D3-401"},
		{"172.16.103.102", "jigsaw", "jigsaw", "D3-402"},
		{"172.16.103.103", "jigsaw", "jigsaw", "D3-403"},
		{"172.16.103.104", "jigsaw", "jigsaw", "D3-404"},
		{"172.16.103.105", "jigsaw", "jigsaw", "D3-405"},
		{"172.16.103.106", "jigsaw", "jigsaw", "D3-406"},
		{"172.16.103.107", "jigsaw", "jigsaw", "D3-407"},
		{"172.16.103.108", "jigsaw", "jigsaw", "D3-408"},
		{"172.16.103.109", "jigsaw", "jigsaw", "D3-409"},
		{"172.16.103.110", "jigsaw", "jigsaw", "D3-410"},
		{"172.16.103.175", "jigsaw", "jigsaw", "D3-411"},
		{"172.16.103.115", "jigsaw", "jigsaw", "D3-415"},
		{"172.16.103.116", "jigsaw", "jigsaw", "D3-416"},
		{"172.16.103.117", "jigsaw", "jigsaw", "D3-417"},
		{"172.16.103.118", "jigsaw", "jigsaw", "D3-418"},
		{"172.16.103.119", "jigsaw", "jigsaw", "D3-419"},
		{"172.16.103.120", "jigsaw", "jigsaw", "D3-420"},
		{"172.16.103.192", "jigsaw", "jigsaw", "D3-421"},
		{"172.16.103.127", "jigsaw", "jigsaw", "D3-427"},
		{"172.16.103.128", "jigsaw", "jigsaw", "D3-428"},
		{"172.16.103.129", "jigsaw", "jigsaw", "D3-429"},
		{"172.16.103.130", "jigsaw", "jigsaw", "D3-430"},
		{"172.16.103.131", "jigsaw", "jigsaw", "D3-431"},
		{"172.16.103.164", "jigsaw", "jigsaw", "D3-464"},
		{"172.16.103.165", "jigsaw", "jigsaw", "D3-465"},
		{"172.16.103.166", "jigsaw", "jigsaw", "D3-466"},
		{"172.16.103.167", "jigsaw", "jigsaw", "D3-467"},
		{"172.16.103.178", "jigsaw", "jigsaw", "D3-478"},
		{"172.16.103.179", "jigsaw", "jigsaw", "D3-479"},
		{"172.16.103.180", "jigsaw", "jigsaw", "D3-480"},
		{"172.16.103.181", "jigsaw", "jigsaw", "D3-481"},
		{"172.16.103.182", "jigsaw", "jigsaw", "D3-482"},
		{"172.16.103.183", "jigsaw", "jigsaw", "D3-483"},
		{"172.16.103.184", "jigsaw", "jigsaw", "D3-484"},
		{"172.16.103.199", "jigsaw", "jigsaw", "D3-499"},
		{"172.16.103.200", "jigsaw", "jigsaw", "D3-500"},
		{"172.16.103.201", "jigsaw", "jigsaw", "D3-501"},
		{"172.16.103.207", "jigsaw", "jigsaw", "D3-507"},
		{"172.16.103.208", "jigsaw", "jigsaw", "D3-508"},
		{"172.16.103.209", "jigsaw", "jigsaw", "D3-509"},
		{"172.16.103.210", "jigsaw", "jigsaw", "D3-510"},
		{"172.16.103.211", "jigsaw", "jigsaw", "D3-511"},
		{"172.16.103.212", "jigsaw", "jigsaw", "D3-512"},
		{"172.16.103.215", "jigsaw", "jigsaw", "D3-515"},
		{"172.16.103.216", "jigsaw", "jigsaw", "D3-516"},
		{"172.16.103.219", "jigsaw", "jigsaw", "D3-519"},
		{"172.16.103.174", "jigsaw", "jigsaw", "D3-521"},
		{"172.16.103.223", "jigsaw", "jigsaw", "D3-523"},
		{"172.16.103.224", "jigsaw", "jigsaw", "D3-524"},
		{"172.16.103.225", "jigsaw", "jigsaw", "D3-525"},
		{"172.16.103.226", "jigsaw", "jigsaw", "D3-526"},
		{"172.16.103.229", "jigsaw", "jigsaw", "D3-529"},
		{"172.16.103.230", "jigsaw", "jigsaw", "D3-530"},
		{"172.16.103.231", "jigsaw", "jigsaw", "D3-531"},
		{"172.16.103.232", "jigsaw", "jigsaw", "D3-532"},
		{"172.16.103.239", "jigsaw", "jigsaw", "D3-539"},
		{"172.16.103.240", "jigsaw", "jigsaw", "D3-540"},
		{"172.16.103.193", "jigsaw", "jigsaw", "D3-541"},
		{"172.16.103.244", "jigsaw", "jigsaw", "D3-544"},
		{"172.16.103.245", "jigsaw", "jigsaw", "D3-545"},
		{"172.16.103.246", "jigsaw", "jigsaw", "D3-546"},
		{"172.16.103.247", "jigsaw", "jigsaw", "D3-547"},
		{"172.16.103.248", "jigsaw", "jigsaw", "D3-548"},
		{"172.16.103.253", "jigsaw", "jigsaw", "D3-553"},
		{"172.16.103.254", "jigsaw", "jigsaw", "D3-554"},
		{"172.16.102.220", "jigsaw", "jigsaw", "D3-555"},
		{"172.16.102.221", "jigsaw", "jigsaw", "D3-556"},
		{"172.16.102.57", "jigsaw", "jigsaw", "D3-557"},
		{"172.16.102.58", "jigsaw", "jigsaw", "D3-558"},
		{"172.16.102.59", "jigsaw", "jigsaw", "D3-559"},
		{"172.16.102.60", "jigsaw", "jigsaw", "D3-560"},
		{"172.16.103.145", "jigsaw", "jigsaw", "D3-565"},
		{"172.16.103.146", "jigsaw", "jigsaw", "D3-566"},
		{"172.16.103.147", "jigsaw", "jigsaw", "D3-568"},
		{"172.16.103.148", "jigsaw", "jigsaw", "D3-569"},
		{"172.16.103.149", "jigsaw", "jigsaw", "D3-570"},
		{"172.16.103.150", "jigsaw", "jigsaw", "D3-571"},
		{"172.16.103.151", "jigsaw", "jigsaw", "D3-572"},
		{"172.16.102.80", "jigsaw", "jigsaw", "D3-580"},
		{"172.16.102.195", "jigsaw", "jigsaw", "D3-581"},
		{"172.16.102.196", "jigsaw", "jigsaw", "D3-582"},
		{"172.16.102.197", "jigsaw", "jigsaw", "D3-583"},
		{"172.16.102.198", "jigsaw", "jigsaw", "D3-584"},
		{"172.16.102.199", "jigsaw", "jigsaw", "D3-587"},
		{"172.16.102.200", "jigsaw", "jigsaw", "D3-588"},
		{"172.16.102.201", "jigsaw", "jigsaw", "D3-590"},
		{"172.16.102.202", "jigsaw", "jigsaw", "D3-591"},
		{"172.16.102.203", "jigsaw", "jigsaw", "D3-592"},
		{"172.16.102.204", "jigsaw", "jigsaw", "D3-593"},
		{"172.16.102.205", "jigsaw", "jigsaw", "D3-594"},
		{"172.16.102.206", "jigsaw", "jigsaw", "D3-595"},
		{"172.16.102.156", "jigsaw", "jigsaw", "D3-596"},
		{"172.16.102.157", "jigsaw", "jigsaw", "D3-597"},
		{"172.16.102.158", "jigsaw", "jigsaw", "D3-598"},
		{"172.16.102.159", "jigsaw", "jigsaw", "D3-599"},
		{"172.16.102.160", "jigsaw", "jigsaw", "D3-600"},
		{"172.16.102.161", "jigsaw", "jigsaw", "D3-601"},
		{"172.16.102.162", "jigsaw", "jigsaw", "D3-602"},
		{"172.16.102.163", "jigsaw", "jigsaw", "D3-603"},
		{"172.16.102.164", "jigsaw", "jigsaw", "D3-604"},
		{"172.16.102.165", "jigsaw", "jigsaw", "D3-605"},
		{"172.16.102.166", "jigsaw", "jigsaw", "D3-606"},
		{"172.16.102.167", "jigsaw", "jigsaw", "D3-607"},
		{"172.16.102.168", "jigsaw", "jigsaw", "D3-608"},
		{"172.16.102.169", "jigsaw", "jigsaw", "D3-609"},
		{"172.16.102.170", "jigsaw", "jigsaw", "D3-610"},
		{"172.16.102.171", "jigsaw", "jigsaw", "D3-611"},
		{"172.16.102.172", "jigsaw", "jigsaw", "D3-612"},
		{"172.16.102.173", "jigsaw", "jigsaw", "D3-613"},
		{"172.16.102.174", "jigsaw", "jigsaw", "D3-614"},
		{"172.16.102.175", "jigsaw", "jigsaw", "D3-615"},
		{"172.16.102.176", "jigsaw", "jigsaw", "D3-616"},
		{"172.16.102.193", "jigsaw", "jigsaw", "D3-633"},
		{"172.16.102.194", "jigsaw", "jigsaw", "D3-634"},
		{"172.16.101.121", "jigsaw", "jigsaw", "D3-635"},
		{"172.16.101.122", "jigsaw", "jigsaw", "D3-636"},
		{"172.16.101.123", "jigsaw", "jigsaw", "D3-637"},
		{"172.16.101.124", "jigsaw", "jigsaw", "D3-638"},
		{"172.16.101.125", "jigsaw", "jigsaw", "D3-649"},
		{"172.16.101.126", "jigsaw", "jigsaw", "D3-650"},
		{"172.16.101.127", "jigsaw", "jigsaw", "D3-651"},
		{"172.16.101.128", "jigsaw", "jigsaw", "D3-652"},
		{"172.16.101.129", "jigsaw", "jigsaw", "D3-653"},
		{"172.16.101.130", "jigsaw", "jigsaw", "D3-654"},
		{"172.16.101.131", "jigsaw", "jigsaw", "D3-655"},
		{"172.16.101.132", "jigsaw", "jigsaw", "D3-656"},
		{"172.16.101.133", "jigsaw", "jigsaw", "D3-657"},
		{"172.16.101.134", "jigsaw", "jigsaw", "D3-658"},
		{"172.16.101.149", "jigsaw", "jigsaw", "D3-659"},
		{"172.16.101.150", "jigsaw", "jigsaw", "D3-660"},
		{"172.16.101.151", "jigsaw", "jigsaw", "D3-661"},
		{"172.16.100.162", "jigsaw", "jigsaw", "D3-662"},
		{"172.16.100.193", "jigsaw", "jigsaw", "D3-663"},
		{"172.16.100.164", "jigsaw", "jigsaw", "D3-664"},
		{"172.16.101.137", "jigsaw", "jigsaw", "D3-665"},
		{"172.16.101.138", "jigsaw", "jigsaw", "D3-666"},
		{"172.16.101.139", "jigsaw", "jigsaw", "D3-667"},
		{"172.16.101.140", "jigsaw", "jigsaw", "D3-668"},
		{"172.16.101.141", "jigsaw", "jigsaw", "D3-669"},
		{"172.16.101.142", "jigsaw", "jigsaw", "D3-670"},
		{"172.16.101.143", "jigsaw", "jigsaw", "D3-671"},
		{"172.16.101.144", "jigsaw", "jigsaw", "D3-672"},
		{"172.16.101.145", "jigsaw", "jigsaw", "D3-673"},
		{"172.16.102.254", "jigsaw", "jigsaw", "D3-674"},
		{"172.16.101.147", "jigsaw", "jigsaw", "D3-675"},
		{"172.16.103.176", "jigsaw", "jigsaw", "D3-676"},
		{"172.16.102.177", "jigsaw", "jigsaw", "D3-677"},
		{"172.16.102.179", "jigsaw", "jigsaw", "D3-679"},
		{"172.16.102.180", "jigsaw", "jigsaw", "D3-680"},
		{"172.16.102.183", "jigsaw", "jigsaw", "D3-683"},
		{"172.16.100.184", "jigsaw", "jigsaw", "D3-684"},
		{"172.16.100.185", "jigsaw", "jigsaw", "D3-685"},
		{"172.16.100.156", "jigsaw", "jigsaw", "D3-686"},
		{"172.16.100.160", "jigsaw", "jigsaw", "D3-700"},
		{"172.16.103.3", "jigsaw", "jigsaw", "D3-703"},
		{"172.16.103.4", "jigsaw", "jigsaw", "D3-704"},
		{"172.16.103.5", "jigsaw", "jigsaw", "D3-705"},
		{"172.16.103.6", "jigsaw", "jigsaw", "D3-706"},
		{"172.16.103.13", "jigsaw", "jigsaw", "D3-713"},
		{"172.16.103.14", "jigsaw", "jigsaw", "D3-714"},
		{"172.16.103.15", "jigsaw", "jigsaw", "D3-715"},
		{"172.16.103.19", "jigsaw", "jigsaw", "D3-719"},
		{"172.16.103.20", "jigsaw", "jigsaw", "D3-720"},
		{"172.16.103.21", "jigsaw", "jigsaw", "D3-721"},
		{"172.16.103.22", "jigsaw", "jigsaw", "D3-722"},
		{"172.16.103.23", "jigsaw", "jigsaw", "D3-723"},
		{"172.16.103.24", "jigsaw", "jigsaw", "D3-724"},
		{"172.16.103.25", "jigsaw", "jigsaw", "D3-725"},
		{"172.16.103.26", "jigsaw", "jigsaw", "D3-726"},
		{"172.16.103.27", "jigsaw", "jigsaw", "D3-727"},
		{"172.16.103.28", "jigsaw", "jigsaw", "D3-728"},
		{"172.16.103.29", "jigsaw", "jigsaw", "D3-729"},
		{"172.16.103.30", "jigsaw", "jigsaw", "D3-730"},
		{"172.16.103.31", "jigsaw", "jigsaw", "D3-731"},
		{"172.16.103.32", "jigsaw", "jigsaw", "D3-732"},
		{"172.16.103.33", "jigsaw", "jigsaw", "D3-733"},
		{"172.16.103.155", "jigsaw", "jigsaw", "D3-735"},
		{"172.16.103.156", "jigsaw", "jigsaw", "D3-736"},
		{"172.16.103.157", "jigsaw", "jigsaw", "D3-737"},
		{"172.16.103.138", "jigsaw", "jigsaw", "D3-738"},
		{"172.16.103.139", "jigsaw", "jigsaw", "D3-739"},
		{"172.16.103.140", "jigsaw", "jigsaw", "D3-740"},
		{"172.16.103.141", "jigsaw", "jigsaw", "D3-741"},
		{"172.16.100.202", "jigsaw", "jigsaw", "D3-742"},
		{"172.16.100.203", "jigsaw", "jigsaw", "D3-743"},
		{"172.16.100.204", "jigsaw", "jigsaw", "D3-744"},
		{"172.16.102.185", "jigsaw", "jigsaw", "D3-745"},
		{"172.16.102.186", "jigsaw", "jigsaw", "D3-746"},
		{"172.16.102.187", "jigsaw", "jigsaw", "D3-747"},
		{"172.16.102.188", "jigsaw", "jigsaw", "D3-748"},
		{"172.16.102.189", "jigsaw", "jigsaw", "D3-749"},
		{"172.16.102.190", "jigsaw", "jigsaw", "D3-750"},
		{"172.16.102.191", "jigsaw", "jigsaw", "D3-751"},
		{"172.16.102.192", "jigsaw", "jigsaw", "D3-752"},
		{"172.16.102.247", "jigsaw", "jigsaw", "D3-753"},
		{"172.16.102.248", "jigsaw", "jigsaw", "D3-754"},
		{"172.16.103.185", "jigsaw", "jigsaw", "D3-755"},
		{"172.16.103.186", "jigsaw", "jigsaw", "D3-756"},
		{"172.16.103.187", "jigsaw", "jigsaw", "D3-757"},
		{"172.16.103.188", "jigsaw", "jigsaw", "D3-758"},
		{"172.16.103.189", "jigsaw", "jigsaw", "D3-759"},
		{"172.16.103.190", "jigsaw", "jigsaw", "D3-760"},
		{"172.16.103.158", "jigsaw", "jigsaw", "D3-774"},
		{"172.16.103.159", "jigsaw", "jigsaw", "D3-775"},
		{"172.16.103.160", "jigsaw", "jigsaw", "D3-776"},
		{"172.16.103.161", "jigsaw", "jigsaw", "D3-777"},
		{"172.16.103.162", "jigsaw", "jigsaw", "D3-778"},
		{"172.16.103.163", "jigsaw", "jigsaw", "D3-779"},
		{"172.16.103.168", "jigsaw", "jigsaw", "D3-780"},
		{"172.16.103.169", "jigsaw", "jigsaw", "D3-781"},
		{"172.16.103.170", "jigsaw", "jigsaw", "D3-782"},
		{"172.16.103.171", "jigsaw", "jigsaw", "D3-783"},
		{"172.16.103.173", "jigsaw", "jigsaw", "D3-785"},
		{"172.16.101.1", "jigsaw", "jigsaw", "D4-001"},
		{"172.16.101.2", "jigsaw", "jigsaw", "D4-002"},
		{"172.16.101.3", "jigsaw", "jigsaw", "D4-003"},
		{"172.16.101.4", "jigsaw", "jigsaw", "D4-004"},
		{"172.16.101.5", "jigsaw", "jigsaw", "D4-005"},
		{"172.16.101.6", "jigsaw", "jigsaw", "D4-006"},
		{"172.16.101.7", "jigsaw", "jigsaw", "D4-007"},
		{"172.16.101.8", "jigsaw", "jigsaw", "D4-008"},
		{"172.16.101.9", "jigsaw", "jigsaw", "D4-009"},
		{"172.16.101.10", "jigsaw", "jigsaw", "D4-010"},
		{"172.16.101.11", "jigsaw", "jigsaw", "D4-011"},
		{"172.16.101.12", "jigsaw", "jigsaw", "D4-012"},
		{"172.16.101.13", "jigsaw", "jigsaw", "D4-013"},
		{"172.16.101.14", "jigsaw", "jigsaw", "D4-014"},
		{"172.16.101.15", "jigsaw", "jigsaw", "D4-015"},
		{"172.16.101.17", "jigsaw", "jigsaw", "D4-017"},
		{"172.16.101.19", "jigsaw", "jigsaw", "D4-019"},
		{"172.16.101.22", "jigsaw", "jigsaw", "D4-022"},
		{"172.16.101.23", "jigsaw", "jigsaw", "D4-023"},
		{"172.16.101.26", "jigsaw", "jigsaw", "D4-026"},
		{"172.16.101.28", "jigsaw", "jigsaw", "D4-028"},
		{"172.16.101.31", "jigsaw", "jigsaw", "D4-031"},
		{"172.16.101.32", "jigsaw", "jigsaw", "D4-032"},
		{"172.16.101.33", "jigsaw", "jigsaw", "D4-033"},
		{"172.16.101.34", "jigsaw", "jigsaw", "D4-034"},
		{"172.16.101.35", "jigsaw", "jigsaw", "D4-035"},
		{"172.16.101.37", "jigsaw", "jigsaw", "D4-037"},
		{"172.16.101.39", "jigsaw", "jigsaw", "D4-039"},
		{"172.16.101.41", "jigsaw", "jigsaw", "D4-041"},
		{"172.16.101.43", "jigsaw", "jigsaw", "D4-043"},
		{"172.16.101.44", "jigsaw", "jigsaw", "D4-044"},
		{"172.16.101.45", "jigsaw", "jigsaw", "D4-045"},
		{"172.16.101.46", "jigsaw", "jigsaw", "D4-046"},
		{"172.16.101.47", "jigsaw", "jigsaw", "D4-047"},
		{"172.16.101.48", "jigsaw", "jigsaw", "D4-048"},
		{"172.16.101.49", "jigsaw", "jigsaw", "D4-049"},
		{"172.16.101.50", "jigsaw", "jigsaw", "D4-050"},
		{"172.16.101.51", "jigsaw", "jigsaw", "D4-051"},
		{"172.16.101.52", "jigsaw", "jigsaw", "D4-052"},
		{"172.16.101.53", "jigsaw", "jigsaw", "D4-053"},
		{"172.16.101.54", "jigsaw", "jigsaw", "D4-054"},
		{"172.16.101.55", "jigsaw", "jigsaw", "D4-055"},
		{"172.16.101.56", "jigsaw", "jigsaw", "D4-056"},
		{"172.16.101.57", "jigsaw", "jigsaw", "D4-057"},
		{"172.16.101.58", "jigsaw", "jigsaw", "D4-058"},
		{"172.16.101.59", "jigsaw", "jigsaw", "D4-059"},
		{"172.16.101.60", "jigsaw", "jigsaw", "D4-060"},
		{"172.16.101.61", "jigsaw", "jigsaw", "D4-061"},
		{"172.16.101.62", "jigsaw", "jigsaw", "D4-062"},
		{"172.16.101.63", "jigsaw", "jigsaw", "D4-063"},
		{"172.16.101.64", "jigsaw", "jigsaw", "D4-064"},
		{"172.16.101.65", "jigsaw", "jigsaw", "D4-065"},
		{"172.16.101.66", "jigsaw", "jigsaw", "D4-066"},
		{"172.16.101.67", "jigsaw", "jigsaw", "D4-067"},
		{"172.16.101.68", "jigsaw", "jigsaw", "D4-068"},
		{"172.16.101.69", "jigsaw", "jigsaw", "D4-069"},
		{"172.16.101.70", "jigsaw", "jigsaw", "D4-070"},
		{"172.16.101.71", "jigsaw", "jigsaw", "D4-071"},
		{"172.16.101.72", "jigsaw", "jigsaw", "D4-072"},
		{"172.16.101.73", "jigsaw", "jigsaw", "D4-073"},
		{"172.16.101.74", "jigsaw", "jigsaw", "D4-074"},
		{"172.16.101.75", "jigsaw", "jigsaw", "D4-075"},
		{"172.16.101.76", "jigsaw", "jigsaw", "D4-076"},
		{"172.16.101.77", "jigsaw", "jigsaw", "D4-077"},
		{"172.16.101.78", "jigsaw", "jigsaw", "D4-078"},
		{"172.16.101.79", "jigsaw", "jigsaw", "D4-079"},
		{"172.16.101.80", "jigsaw", "jigsaw", "D4-080"},
		{"172.16.101.81", "jigsaw", "jigsaw", "D4-081"},
		{"172.16.101.82", "jigsaw", "jigsaw", "D4-082"},
		{"172.16.101.83", "jigsaw", "jigsaw", "D4-083"},
		{"172.16.101.84", "jigsaw", "jigsaw", "D4-084"},
		{"172.16.101.85", "jigsaw", "jigsaw", "D4-085"},
		{"172.16.101.86", "jigsaw", "jigsaw", "D4-086"},
		{"172.16.101.87", "jigsaw", "jigsaw", "D4-087"},
		{"172.16.101.88", "jigsaw", "jigsaw", "D4-088"},
		{"172.16.101.89", "jigsaw", "jigsaw", "D4-089"},
		{"172.16.101.90", "jigsaw", "jigsaw", "D4-090"},
		{"172.16.101.91", "jigsaw", "jigsaw", "D4-091"},
		{"172.16.101.92", "jigsaw", "jigsaw", "D4-092"},
		{"172.16.101.95", "jigsaw", "jigsaw", "D4-095"},
		{"172.16.101.96", "jigsaw", "jigsaw", "D4-096"},
		{"172.16.101.97", "jigsaw", "jigsaw", "D4-097"},
		{"172.16.101.98", "jigsaw", "jigsaw", "D4-098"},
		{"172.16.101.99", "jigsaw", "jigsaw", "D4-099"},
		{"172.16.101.101", "jigsaw", "jigsaw", "D4-101"},
		{"172.16.101.102", "jigsaw", "jigsaw", "D4-102"},
		{"172.16.101.103", "jigsaw", "jigsaw", "D4-103"},
		{"172.16.102.8", "jigsaw", "jigsaw", "D5-008"},
		{"172.16.102.11", "jigsaw", "jigsaw", "D5-011"},
		{"172.16.102.12", "jigsaw", "jigsaw", "D5-012"},
		{"172.16.102.13", "jigsaw", "jigsaw", "D5-013"},
		{"172.16.102.14", "jigsaw", "jigsaw", "D5-014"},
		{"172.16.102.15", "jigsaw", "jigsaw", "D5-015"},
		{"172.16.102.16", "jigsaw", "jigsaw", "D5-016"},
		{"172.16.102.17", "jigsaw", "jigsaw", "D5-017"},
		{"172.16.102.18", "jigsaw", "jigsaw", "D5-018"},
		{"172.16.102.19", "jigsaw", "jigsaw", "D5-019"},
		{"172.16.102.20", "jigsaw", "jigsaw", "D5-020"},
		{"172.16.102.21", "jigsaw", "jigsaw", "D5-021"},
		{"172.16.102.22", "jigsaw", "jigsaw", "D5-022"},
		{"172.16.102.23", "jigsaw", "jigsaw", "D5-023"},
		{"172.16.102.24", "jigsaw", "jigsaw", "D5-024"},
		{"172.16.102.25", "jigsaw", "jigsaw", "D5-025"},
		{"172.16.102.26", "jigsaw", "jigsaw", "D5-026"},
		{"172.16.102.27", "jigsaw", "jigsaw", "D5-027"},
		{"172.16.102.28", "jigsaw", "jigsaw", "D5-028"},
		{"172.16.102.113", "jigsaw", "jigsaw", "D5-029"},
		{"172.16.102.30", "jigsaw", "jigsaw", "D5-030"},
		{"172.16.102.31", "jigsaw", "jigsaw", "D5-031"},
		{"172.16.102.32", "jigsaw", "jigsaw", "D5-032"},
		{"172.16.102.33", "jigsaw", "jigsaw", "D5-033"},
		{"172.16.102.34", "jigsaw", "jigsaw", "D5-034"},
		{"172.16.102.35", "jigsaw", "jigsaw", "D5-035"},
		{"172.16.102.36", "jigsaw", "jigsaw", "D5-036"},
		{"172.16.102.37", "jigsaw", "jigsaw", "D5-037"},
		{"172.16.102.38", "jigsaw", "jigsaw", "D5-038"},
		{"172.16.102.39", "jigsaw", "jigsaw", "D5-039"},
		{"172.16.102.52", "jigsaw", "jigsaw", "G4-013"},
		{"172.16.100.73", "jigsaw", "jigsaw", "G4-024"},
		{"172.16.100.72", "jigsaw", "jigsaw", "G4-025"},
		{"172.16.102.47", "jigsaw", "jigsaw", "G4-026"},
		{"172.16.102.42", "jigsaw", "jigsaw", "G4-030"},
		{"172.16.102.50", "jigsaw", "jigsaw", "G4-046"},
		{"172.16.103.68", "jigsaw", "jigsaw", "G4-047"},
		{"172.16.100.75", "jigsaw", "jigsaw", "G4-048"},
		{"172.16.103.43", "jigsaw", "jigsaw", "G4-049"},
		{"172.16.101.166", "jigsaw", "jigsaw", "G4-050"},
		{"172.16.102.90", "jigsaw", "jigsaw", "G4-063"},
		{"172.16.102.91", "jigsaw", "jigsaw", "G4-065"},
		{"172.16.103.69", "jigsaw", "jigsaw", "G4-068"},
		{"172.16.103.133", "jigsaw", "jigsaw", "G4-069"},
		{"172.16.103.135", "jigsaw", "jigsaw", "G4-073"},
		{"172.16.103.134", "jigsaw", "jigsaw", "G4-074"},
		{"172.16.103.136", "jigsaw", "jigsaw", "G4-075"},
		{"172.16.102.51", "jigsaw", "jigsaw", "G4-076"},
		{"172.16.102.53", "jigsaw", "jigsaw", "G4-077"},
		{"172.16.102.54", "jigsaw", "jigsaw", "G4-081"},
		{"172.16.102.55", "jigsaw", "jigsaw", "G4-082"},
		{"172.16.102.87", "jigsaw", "jigsaw", "G4-083"},
		{"172.16.102.88", "jigsaw", "jigsaw", "G4-084"},
		{"172.16.102.89", "jigsaw", "jigsaw", "G4-085"},
		{"172.16.102.92", "jigsaw", "jigsaw", "G4-087"},
		{"172.16.102.93", "jigsaw", "jigsaw", "G4-088"},
		{"172.16.102.94", "jigsaw", "jigsaw", "G4-089"},
		{"172.16.102.96", "jigsaw", "jigsaw", "G4-090"},
		{"172.16.102.97", "jigsaw", "jigsaw", "G4-091"},
		{"172.16.102.98", "jigsaw", "jigsaw", "G4-092"},
		{"172.16.101.232", "jigsaw", "jigsaw", "G4-093"},
		{"172.16.101.233", "jigsaw", "jigsaw", "G4-094"},
		{"172.16.100.104", "jigsaw", "jigsaw", "G4-096"},
		{"172.16.100.105", "jigsaw", "jigsaw", "G4-098"},
		{"172.16.100.106", "jigsaw", "jigsaw", "G4-099"},
		{"172.16.102.48", "jigsaw", "jigsaw", "G5-002"},
		{"172.16.102.49", "jigsaw", "jigsaw", "G5-003"},
		{"172.16.101.167", "jigsaw", "jigsaw", "G5-004"},
		{"172.16.100.89", "jigsaw", "jigsaw", "G5-005"},
		{"172.16.100.60", "jigsaw", "jigsaw", "G5-006"},
		{"172.16.102.101", "jigsaw", "jigsaw", "X4-001"},
		{"172.16.102.126", "jigsaw", "jigsaw", "X4-026"},
		{"172.16.102.142", "jigsaw", "jigsaw", "X4-042"},
		{"172.16.102.145", "jigsaw", "jigsaw", "X4-045"},
		{"172.16.102.147", "jigsaw", "jigsaw", "X4-047"},
		{"172.16.102.148", "jigsaw", "jigsaw", "X4-048"},
		{"172.16.103.203", "jigsaw", "jigsaw", "X4-050"},
		{"172.16.101.174", "jigsaw", "jigsaw", "X4-074"},
		{"172.16.101.175", "jigsaw", "jigsaw", "X4-075"},
		{"172.16.101.114", "jigsaw", "jigsaw", "X4-077"},
		{"172.16.101.115", "jigsaw", "jigsaw", "X4-078"},
		{"172.16.101.136", "jigsaw", "jigsaw", "X5-036"},
		{"172.16.101.146", "jigsaw", "jigsaw", "X5-046"},
		{"172.16.101.154", "jigsaw", "jigsaw", "X5-054"},
		{"172.16.101.155", "jigsaw", "jigsaw", "X5-055"},
		{"172.16.101.156", "jigsaw", "jigsaw", "X5-056"},
		{"172.16.101.157", "jigsaw", "jigsaw", "X5-057"},
		{"172.16.101.158", "jigsaw", "jigsaw", "X5-058"},
		{"172.16.101.159", "jigsaw", "jigsaw", "X5-059"},
		{"172.16.101.160", "jigsaw", "jigsaw", "X5-060"},
		{"172.16.101.161", "jigsaw", "jigsaw", "X5-061"},
		{"172.16.101.164", "jigsaw", "jigsaw", "X5-064"},
		{"172.16.102.222", "jigsaw", "jigsaw", "X5-065"},
		{"172.16.102.223", "jigsaw", "jigsaw", "X5-066"},
		{"172.16.102.224", "jigsaw", "jigsaw", "X5-067"},
		{"172.16.102.225", "jigsaw", "jigsaw", "X5-068"},
		{"172.16.101.116", "jigsaw", "jigsaw", "X5-069"},
		{"172.16.101.117", "jigsaw", "jigsaw", "X5-070"},
		{"172.16.101.118", "jigsaw", "jigsaw", "X5-071"},
		{"172.16.101.119", "jigsaw", "jigsaw", "X5-072"},
		{"172.16.101.120", "jigsaw", "jigsaw", "X5-073"},
		{"172.16.102.61", "jigsaw", "jigsaw", "X7-001"},
		{"172.16.102.68", "jigsaw", "jigsaw", "X7-008"},
		{"172.16.102.70", "jigsaw", "jigsaw", "X7-010"},
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
