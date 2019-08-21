package helper

import (
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	CpuUser float64
	CpuSystem float64
	CpuIdle float64
	MemTotal int
	MemUsed int
	MemCached int
	MemFree int
	RxBytes int
	TxBytes int
	Uptime string
	Time string
)

func StatsFromHostname(hostname string) {
	db := connect()
	row, err := db.Query("SELECT * FROM stats." + hostname + " ORDER BY Time DESC limit 1")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	for row.Next() {
		err = row.Scan(&CpuUser, &CpuSystem, &CpuIdle, &MemTotal, &MemUsed, &MemCached, &MemFree, &RxBytes, &TxBytes, &Uptime, &Time)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = row.Err()
	if err != nil {
		log.Fatal()
	}
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