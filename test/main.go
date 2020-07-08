package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	UpdateData()
}

//SetupCache by reading csv file from path
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
	dataLink := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRGDH-jfWULsmOHH5jDTgDZDZPxdgMmnrM6TOrF8FzV6FJEYtbSTcRhONDNG21hfKge04nZ96oKA78I/pub?gid=0&single=true&output=csv"
	savePath := "./data.csv"
	err := DownloadFile(savePath, dataLink)
	if err != nil {
		return err
	}
	locationLink := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRGDH-jfWULsmOHH5jDTgDZDZPxdgMmnrM6TOrF8FzV6FJEYtbSTcRhONDNG21hfKge04nZ96oKA78I/pub?gid=1619494999&single=true&output=csv"
	locationPath = "./location.csv"
	err = DownloadFile(locationPath, locationLink)
	if err != nil {
		return err
	}
	locations ,err :=CreateLocationMap(locationPath)
	if err != nil {
		return err
	}

	return nil
}

type Location struct {
	Lon float64
	Lat float64
}

func CreateLocationMap(filePath string) map[string]Location ,error {
	locations := make(map[string]Location)
	csvfile, err := os.Open(filePath)
	if err != nil {
		return location, err
	}
	r := csv.NewReader(csvfile)
	_, _ = r.Read()
	_, _ = r.Read()
	lonColume := 1
	latColume := 2
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return location, err
		}
		lon, err := strconv.ParseFloat(record[lonColume], 64)
		if err != nil {
			return location, err
		}
		lat, err := strconv.ParseFloat(record[latColume], 64)
		if err != nil {
			return location, err
		}
		locations[record[0]] = Location{Lon: lon, Lat: lat}
	}
	return locations , nil
}

func UploadDataToDatabase(client *sql.DB,filePath string , locations map[string]Location) error{
	csvfile, err := os.Open(filePath)
	if err != nil {
		return location, err
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
			return location, err
		}
		if val, ok := locations[data[addressColume]]; ok {
			UploadToSQL(client)
		}else{
			continue
		}
	}
	return nil
}

func UploadToSQL(client *sql.DB){
	//球隊名稱,星期,時間,球館,地址,強度,場地數,收費(男),收費(女),用球,隊長,備註
`
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
	captain VARCHAR(30),
	note VARCHAR(60));
`
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
