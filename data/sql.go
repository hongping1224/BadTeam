package data

import (
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func ToStoreSQLCmd(data Data) string {
	return fmt.Sprintf(
		`INSERT INTO TeamData 
	(name, day, startTime, endTime, courtName, address) VALUES  ("%s",%d,%d,%d,"%s","%s");`,
		StringToHex(data.Name), data.Day, data.StartTime, data.EndTime, StringToHex(data.CourtName), StringToHex(data.Address))
}

func StringToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

func HexToString(s string) string {
	b, _ := hex.DecodeString(s)
	return string(b)
}

func CreateTable(client *sql.DB) error {
	/*DROP TABLE IF EXISTS TeamData; */
	q := `
	CREATE TABLE IF NOT EXISTS TeamData (
		uid INT AUTO_INCREMENT NOT NULL UNIQUE KEY,
		PRIMARY KEY(uid),
		name VARCHAR(100),
		day TINYINT,
		startTime SMALLINT,
		endTime SMALLINT,
		courtName VARCHAR(60),
		address VARCHAR(100),
		fromLevel TINYINT,
		toLevel TINYINT,
		courtCount TINYINT,
		feeM SMALLINT,
		feeF SMALLINT,
		minBallType TINYINT,
		note VARCHAR(100));
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
	//fmt.Println(cmd)
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
