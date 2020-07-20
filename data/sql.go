package data

import (
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func GenerateSearchCmd(form map[string][]string) (string, error) {
	var trainTime int16
	if val, ok := form["traintime"]; ok {
		if len(val) > 0 {
			trainTime = parseRequestTime(val[0])
		}
	}
	day := parseRequestAllDay(form)
	dayquery := getDayQuery(day)
	level := parseRequestAllLevel(form)
	levelquery := getLevelQuery(level)
	cmd := fmt.Sprintf("SELECT * FROM TeamData WHERE (startTime<=%d AND endTime>=%d) %s %s;", level, level, trainTime, trainTime, dayquery, levelquery)
	return cmd, nil
}

func getDayQuery(b []bool) string {
	s := ""
	for i, bo := range b {
		if bo == true {
			if s == "" {
				s = fmt.Sprintf("AND ((fromLevel<=%d AND toLevel>=%d)", i, i)
			} else {
				s += fmt.Sprintf(" OR (fromLevel<=%d AND toLevel>=%d)", i, i)
			}
		}
	}
	if s != "" {
		s += ")"
	}
	return s
}

func getLevelQuery(b []bool) string {
	s := ""
	for i, bo := range b {
		if bo == true {
			if s == "" {
				s = fmt.Sprintf("AND (day=%d", i+1)
			} else {
				s += fmt.Sprintf(" OR day=%d", i+1)
			}
		}
	}
	if s != "" {
		s += ")"
	}
	return s
}

func parseRequestLevel(s string) int8 {
	if val, err := strconv.Atoi(s); err == nil {
		return int8(val)
	}
	return 0
}

func parseRequestAllLevel(form map[string][]string) []bool {
	levels := []string{"lv0", "lv1", "lv2", "lv3", "lv4", "lv5"}
	b := make([]bool, 6)
	print(form)
	for i, level := range levels {
		if val, ok := form[level]; ok {
			if len(val) > 0 {
				b[i] = true
			}
		}
	}
	return b
}

func parseRequestAllDay(form map[string][]string) []bool {
	days := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	b := make([]bool, 7)
	print(form)
	for i, day := range days {
		if val, ok := form[day]; ok {
			if len(val) > 0 {
				b[i] = true
			}
		}
	}
	return b
}

func parseRequestDay(s string) int8 {
	switch s {
	case "mon":
		return 1
	case "tue":
		return 2
	case "wed":
		return 3
	case "thu":
		return 4
	case "fri":
		return 5
	case "sat":
		return 6
	case "sun":
		return 7
	}

	return 0
}

func parseRequestTime(s string) int16 {
	sp := strings.Split(s, ":")
	if len(sp) < 2 {
		return 0
	}
	h, err := strconv.Atoi(sp[0])
	if err != nil {
		return 0
	}
	m, err := strconv.Atoi(sp[1])
	if err != nil {
		return 0
	}
	return int16((h * 100) + m)
}

func ToStoreSQLCmd(data Data) string {
	return fmt.Sprintf(
		`INSERT INTO TeamData 
	(name, day, startTime, endTime, courtName, address,fromLevel,toLevel, courtCount,feeM,feeF,minBallType,note) VALUES  ("%s",%d,%d,%d,"%s","%s",%d,%d, %d,%d,%d,"%s","%s");`,
		StringToHex(data.Name), data.Day, data.StartTime, data.EndTime, StringToHex(data.CourtName), StringToHex(data.Address), data.FromLevel, data.ToLevel, data.CourtCount, data.FeeM, data.FeeF, StringToHex(data.MinBallType), StringToHex(data.Note))
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
		minBallType VARCHAR(50),
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
