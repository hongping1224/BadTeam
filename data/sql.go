package data

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type Location struct {
	Lon float64
	Lat float64
}

type Data struct {
	UID         int
	Name        string
	Day         int8
	StartTime   int8
	EndTime     int8
	CourtName   string
	Address     string
	FromLevel   int8
	ToLevel     int8
	CourtCount  int8
	FeeM        int16
	FeeF        int16
	MinBallType int8
	Note        string
}

func NewData(record []string) {

}

func CreateTable(client *sql.DB) error {
	/*DROP TABLE IF EXISTS TeamData; */
	q := `
	CREATE TABLE IF NOT EXISTS TeamData (
		uid INT AUTO_INCREMENT NOT NULL UNIQUE KEY,
		PRIMARY KEY(uid),
		name VARCHAR(30),
		day TINYINT,
		startTime TINYINT,
		endTime TINYINT,
		courtName VARCHAR(30),
		address VARCHAR(60),
		fromLevel TINYINT,
		toLevel TINYINT,
		courtCount TINYINT,
		feeM SMALLINT,
		feeF SMALLINT,
		minBallType TINYINT,
		note VARCHAR(60));
	`
	rows, err := client.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		fmt.Println(name)
	}
	return nil
}

func DropTable(client *sql.DB) error {
	/*DROP TABLE IF EXISTS TeamData; */
	q := "DROP TABLE IF EXISTS TeamData;"
	rows, err := client.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		fmt.Println(name)
	}
	return nil
}

func UploadDataToDatabase(client *sql.DB, filePath string, locations map[string]Location) error {
	csvfile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r := csv.NewReader(csvfile)
	_, _ = r.Read()
	addressColume := 4
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if val, ok := locations[record[addressColume]]; ok {
			UploadToSQL(client)
			fmt.Println(val)
		} else {
			continue
		}
	}
	return nil
}

func UploadToSQL(client *sql.DB) {
	//球隊名稱,星期,時間,球館,地址,強度,場地數,收費(男),收費(女),用球,隊長,備註
	/*`
	CREATE TABLE IF NOT EXISTS TeamData (
		uid INT AUTO_INCREMENT NOT NULL UNIQUE KEY,
		PRIMARY KEY(uid),
		name VARCHAR(30),
		day TINYINT,
		startTime TINYINT,
		endTime TINYINT,
		courtName VARCHAR(30),
		address VARCHAR(60),
		fromLevel TINYINT,
		toLevel TINYINT,
		courtCount TINYINT,
		feeM SMALLINT,
		feeF SMALLINT,
		minBallType TINYINT,
		note VARCHAR(60));
	`*/
}
