package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hongping1224/BadTeam/data"
)

const (
	combineOutputPath = "./Combine.csv"
	savePath          = "./data.csv"
	locationPath      = "./location.csv"
)

var db *sql.DB

func main() {
	/*fmt.Println("Start UpdateData")
	err := UpdateData()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Finish UpdateData")
	user := os.Getenv("SQLUSER")
	pass := os.Getenv("SQLPASS")
	loginstr := fmt.Sprintf("%s:%s@tcp(badteam.ccz3kc9rn8lq.ap-southeast-1.rds.amazonaws.com:3306)/BADMINTON", user, pass)
	db, err = sql.Open("mysql", loginstr)
	if err != nil {
		fmt.Printf(" sql.Open Error: %v\n", err)
	}
	defer db.Close()
	if err != nil {
		fmt.Printf("DropTable Error: %v\n", err)
	}
	data.DropTable(db)
	data.CreateTable(db)
	fmt.Println("Start Upload")
	data.UploadDataToDatabase(db, combineOutputPath)
	fmt.Println("Finish Upload")
	/*
		fmt.Println("Start Test")
		test(db)
		fmt.Println("Finish Test")
	*/
	http.HandleFunc("/", HomePageHandler)

	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/html/", http.StripPrefix("/html/", fs))

	fmt.Println("Start Host on port 65000")
	http.ListenAndServe(":65000", nil)
}

func test(db *sql.DB) {
	rows, err := db.Query("SELECT address FROM TeamData")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		fmt.Println(data.HexToString(name))
	}
}

type dataResult struct {
	Result map[int]data.Data
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {

	results := make(map[int]data.Data)

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		cmd, err := data.GenerateSearchCmd(r.Form)
		if err == nil {
			fmt.Println(cmd)
			rows, err := db.Query(cmd)
			fmt.Println("finish Query")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer rows.Close()
			fmt.Println(rows)
			for rows.Next() {
				var name data.Data
				if err := rows.Scan(&name.UID,
					&name.Name,
					&name.Day,
					&name.StartTime,
					&name.EndTime,
					&name.CourtName,
					&name.Address,
					&name.FromLevel,
					&name.ToLevel,
					&name.CourtCount,
					&name.FeeM,
					&name.FeeF,
					&name.MinBallType,
					&name.Note); err != nil {
					// Check for a scan error.
					// Query rows will be closed with defer.
					log.Fatal(err)
				}
				name.Name = data.HexToString(name.Name)
				name.CourtName = data.HexToString(name.CourtName)
				name.Address = data.HexToString(name.Address)
				name.MinBallType = data.HexToString(name.MinBallType)
				name.Note = data.HexToString(name.Note)
				results[name.UID] = name
			}
		}
	}
	fmt.Println("returning")
	p := dataResult{Result: results}
	t, err := template.ParseFiles("./html/results.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.Execute(w, p)
}

//UpdateData download new data and refresh cache
func UpdateData() error {
	dataLink := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRGDH-jfWULsmOHH5jDTgDZDZPxdgMmnrM6TOrF8FzV6FJEYtbSTcRhONDNG21hfKge04nZ96oKA78I/pub?gid=0&single=true&output=csv"

	err := DownloadFile(savePath, dataLink)
	if err != nil {
		return err
	}
	locationLink := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRGDH-jfWULsmOHH5jDTgDZDZPxdgMmnrM6TOrF8FzV6FJEYtbSTcRhONDNG21hfKge04nZ96oKA78I/pub?gid=1619494999&single=true&output=csv"

	err = DownloadFile(locationPath, locationLink)
	if err != nil {
		return err
	}
	locations, err := CreateLocationMap(locationPath)
	if err != nil {
		return err
	}
	err = MapLocationToData(savePath, combineOutputPath, locations)
	if err != nil {
		return err
	}
	return nil
}

func MapLocationToData(filepath, outpath string, locations map[string]data.Location) error {
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

func CreateLocationMap(filePath string) (map[string]data.Location, error) {
	locations := make(map[string]data.Location)
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
		locations[record[0]] = data.Location{Lon: lon, Lat: lat}
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
