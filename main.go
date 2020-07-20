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
	/*err := UpdateData()
	if err != nil {
		fmt.Println(err)
	}
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
	*/
	http.HandleFunc("/updatedata", UpdateHandler)
	http.HandleFunc("/", HomePageHandler)

	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/html/", http.StripPrefix("/html/", fs))

	fmt.Println("Start Host on port 65000")
	http.ListenAndServe(":65000", nil)
}

type dataResult struct {
	Result map[int]data.Data
	Init   string
	Day    string
	Level  string
	Time   string
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := UpdateData()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Update Complete"))
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	results := make(map[int]data.Data)
	p := dataResult{Init: "init"}
	if r.Method == http.MethodPost {
		p.Init = ""
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		cmd, err := data.GenerateSearchCmd(r.Form)
		if err == nil {
			fmt.Println(cmd)
			rows, err := db.Query(cmd)
			p.Day = r.Form["day"][0]
			p.Level = r.Form["lv"][0]
			p.Time = r.Form["traintime"][0]
			fmt.Println("finish Query")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer rows.Close()
			for rows.Next() {
				var recordData data.Data
				if err := rows.Scan(&recordData.UID,
					&recordData.Name,
					&recordData.Day,
					&recordData.StartTime,
					&recordData.EndTime,
					&recordData.CourtName,
					&recordData.Address,
					&recordData.FromLevel,
					&recordData.ToLevel,
					&recordData.CourtCount,
					&recordData.FeeM,
					&recordData.FeeF,
					&recordData.MinBallType,
					&recordData.Note); err != nil {
					// Check for a scan error.
					// Query rows will be closed with defer.
					log.Fatal(err)
				}
				recordData.Name = data.HexToString(recordData.Name)
				recordData.CourtName = data.HexToString(recordData.CourtName)
				recordData.Address = data.HexToString(recordData.Address)
				recordData.MinBallType = data.HexToString(recordData.MinBallType)
				recordData.Note = data.HexToString(recordData.Note)
				recordData.Daystr = data.Day2DayStr(recordData.Day)
				results[recordData.UID] = recordData

			}

		}
	}
	fmt.Println("returning")
	p.Result = results
	t, err := template.ParseFiles("./html/results.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.Execute(w, p)
}

//UpdateData download new data and refresh cache
func UpdateData() error {
	fmt.Println("Start UpdateData")

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
	fmt.Println("Finish UpdateData")
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
