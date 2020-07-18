package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	sq "github.com/hongping1224/BadTeam/sql"
)

func main() {
	err := UpdateData()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("mysql", "admin:NanaDatabasePassword@tcp(badteam.ccz3kc9rn8lq.ap-southeast-1.rds.amazonaws.com:3306)/BADMINTON")
	if err != nil {
		fmt.Printf(" sql.Open Error: %v\n", err)
	}
	defer db.Close()
	err = CreateTable(db)
	//err = DropTable(db)
	if err != nil {
		fmt.Printf("DropTable Error: %v\n", err)
	}
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
	locationPath := "./location.csv"
	err = DownloadFile(locationPath, locationLink)
	if err != nil {
		return err
	}
	locations, err := CreateLocationMap(locationPath)
	if err != nil {
		return err
	}
	outputPath := "./Combine.csv"
	err = MapLocationToData(savePath, outputPath, locations)
	if err != nil {
		return err
	}
	return nil
}

func MapLocationToData(filepath, outpath string, locations map[string]sq.Location) error {
	csvfile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	outfile, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer outfile.Close()
	w := csv.NewWriter(outfile)

	header, _ := r.Read()
	header = append(header, "lon", "lat")
	w.Write(header)

	addressColume := 4
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if record[addressColume] == "" {
			continue
		}

		if p, ok := locations[record[addressColume]]; ok {
			record = append(record, fmt.Sprintf("%.8f", p.Lon), fmt.Sprintf("%.8f", p.Lat))
			w.Write(record)
			//fmt.Println(p)
		} else {
			fmt.Println("record not found", record[addressColume])
		}
	}
	w.Flush()
	return nil
}

func CreateLocationMap(filePath string) (map[string]sq.Location, error) {
	locations := make(map[string]Location)
	csvfile, err := os.Open(filePath)
	if err != nil {
		return locations, err
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
			return locations, err
		}
		lon, err := strconv.ParseFloat(record[lonColume], 64)
		if err != nil {
			return locations, err
		}
		lat, err := strconv.ParseFloat(record[latColume], 64)
		if err != nil {
			return locations, err
		}
		locations[record[0]] = Location{Lon: lon, Lat: lat}
	}
	return locations, nil
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
