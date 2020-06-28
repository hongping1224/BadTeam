package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//UpdateData()
	s, _ := test("create")
	fmt.Println(s)
}

func CreateTable(db *sql.DB, tablename, colume string) error {
	q := fmt.Sprintf("CREATE TABLE %s(%s);", tablename, colume)
	rows, err := db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
func test(name string) (string, error) {
	db, err := sql.Open("mysql", "admin:NanaDatabasePassword@tcp(badteam.ccz3kc9rn8lq.ap-southeast-1.rds.amazonaws.com:3306)/badteam")
	if err != nil {
		return fmt.Sprintf(" sql.Open Error: %v", err), nil
	}
	defer db.Close()
	if name == "create" {
		err = CreateTable(db, "test", "id INT NOT NULL AUTO_INCREMENT,PRIMARY KEY(id), data INT, quo VARCHAR(30)")
		if err != nil {
			return fmt.Sprintf(" Create Table Error: %v", err), nil
		}
		return "Create Success", nil
	}
	return "", nil
}

func SetupCache(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	for {
		// Read each record from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
	}
	return nil
}

//UpdateData download new data and refresh cache
func UpdateData() error {
	datalink := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRGDH-jfWULsmOHH5jDTgDZDZPxdgMmnrM6TOrF8FzV6FJEYtbSTcRhONDNG21hfKge04nZ96oKA78I/pub?gid=0&single=true&output=csv"
	savepath := "./data.csv"
	err := DownloadFile(savepath, datalink)
	if err != nil {
		return err
	}

	return SetupCache(savepath)
}

//DownloadFile from url and save to filepath
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
