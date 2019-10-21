package helper

import (
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	CpuPerc float64
	MemTotal int
	MemUsed int
	MemCached int
	MemFree int
	RxBytes int
	TxBytes int
	DiskUsed int
	DiskFree int
	DiskWrite int
	DiskRead int
	Uptime string
	Time string
	Tables_in_stats string
)

func StatsFromHostname(hostname string) {
	db := connect()
	row, err := db.Query("SELECT * FROM stats." + hostname + " ORDER BY Time DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	for row.Next() {
		err = row.Scan(&CpuPerc, &MemTotal, &MemUsed, &MemCached, &MemFree, &RxBytes, &TxBytes, &DiskUsed, &DiskFree, &DiskRead, &DiskWrite, &Uptime, &Time)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = row.Err()
	if err != nil {
		log.Fatal()
	}
}

func ListOfTables() []string {
	var tables []string

	db := connect()
	rows, err := db.Query("SHOW TABLES FROM stats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Tables_in_stats)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, Tables_in_stats)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal()
	}

	return tables
}

func connect() (*sql.DB) {
	db, err := sql.Open("mysql", DBSource())
	if err != nil {
		panic(err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err)
	} else {
		//commented for debugging
		//fmt.Println("DB connection working")
	}
	 return db
}

func DBSource() string {
	file, err := os.Open(".login.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	var user string
	var pass string
	var ip string
	for scanner.Scan() {
		if i == 0 {
			user = scanner.Text()
		} else if i == 1 {
			pass = scanner.Text()
		} else if i == 2 {
			ip = scanner.Text()
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	output := user + ":" + pass + "@tcp(" + ip + ")/stats"

	return output
}