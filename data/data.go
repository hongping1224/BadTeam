package data

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Location struct {
	Lon float64
	Lat float64
}

type Data struct {
	UID         int
	Name        string
	Day         int8
	StartTime   int16
	EndTime     int16
	CourtName   string
	Address     string
	FromLevel   int8
	ToLevel     int8
	CourtCount  int8
	FeeM        int16
	FeeF        int16
	MinBallType int8
	Note        string
	Location    Location
}

func NewData(record []string) Data {
	d := Data{Name: record[0], CourtName: record[3], Address: record[4]}
	day, _ := strconv.Atoi(record[1])
	d.Day = int8(day)
	timestr := strings.Split(record[2], "-")
	startTime, _ := strconv.Atoi(timestr[0])
	d.StartTime = int16(startTime)
	endtime, _ := strconv.Atoi(timestr[1])
	d.EndTime = int16(endtime)
	return d
}

func ToStoreSQLCmd(data Data) string {
	return fmt.Sprintf(
		`INSERT INTO TeamData 
	(name, day, startTime, endTime, courtName, address) VALUES  (%v,%d,%d,%d,%v,%v);`,
		[]byte(data.Name), data.Day, data.StartTime, data.EndTime, []byte(data.CourtName), []byte(data.Address))
}

func CreateTable(client *sql.DB) error {
	/*DROP TABLE IF EXISTS TeamData; */
	q := `
	CREATE TABLE IF NOT EXISTS TeamData (
		uid INT AUTO_INCREMENT NOT NULL UNIQUE KEY,
		PRIMARY KEY(uid),
		name VARBINARY(100),
		day TINYINT,
		startTime SMALLINT,
		endTime SMALLINT,
		courtName VARBINARY(60),
		address VARBINARY(100),
		fromLevel TINYINT,
		toLevel TINYINT,
		courtCount TINYINT,
		feeM SMALLINT,
		feeF SMALLINT,
		minBallType TINYINT,
		note VARBINARY(100));
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
		fmt.Println(string([]byte(name)))
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

func UploadDataToDatabase(client *sql.DB, filePath string) error {
	csvfile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r := csv.NewReader(csvfile)
	_, _ = r.Read()
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		data := NewData(record)
		err = UploadToSQL(client, data)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func UploadToSQL(client *sql.DB, data Data) error {
	cmd := ToStoreSQLCmd(data)
	fmt.Println(cmd)
	rows, err := client.Query(cmd)
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
