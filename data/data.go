package data

import (
	"strconv"
	"strings"
)

type Location struct {
	Lon float64
	Lat float64
}

type Data struct {
	UID           int
	Name          string
	Day           int8
	Daystr        string
	StartTime     int16
	EndTime       int16
	CourtName     string
	Address       string
	FromLevel     int8
	ToLevel       int8
	CourtCount    int8
	CourtCountStr string
	FeeM          int16
	FeeMStr       string
	FeeF          int16
	FeeFStr       string
	MinBallType   string
	Note          string
	Location      Location
}

const (
	nameCol        = 0 //d
	dayCol         = 1 //d
	timeCol        = 2 //d
	courtCol       = 3 //d
	addressCol     = 4 //d
	levelCol       = 5 //d
	courtCountCol  = 6 //d
	feeMCol        = 7 //d
	feeFCol        = 8 //d
	minBallTypeCol = 9 //d
	lastUpdateCol  = 10
	noteCol        = 11 //d
	lonCol         = 12 //d
	latCol         = 13 //d
)

func NewData(record []string) Data {
	d := Data{Name: record[nameCol], CourtName: record[courtCol], Address: record[addressCol], MinBallType: record[minBallTypeCol], Note: record[noteCol]}
	d.Day = parseDay(record[dayCol])
	d.StartTime, d.EndTime = parseTime(record[timeCol])
	d.Location = parseLocation(record[lonCol], record[latCol])
	d.FromLevel, d.ToLevel = parseLevel(record[levelCol])
	d.CourtCount = parseCourtCount(record[courtCountCol])
	d.FeeM = parseFee(record[feeMCol])
	d.FeeF = parseFee(record[feeFCol])
	return d
}

func parseCourtCount(s string) int8 {
	if val, err := strconv.Atoi(s); err == nil {
		return int8(val)
	}
	return -1
}

func parseFee(s string) int16 {
	if val, err := strconv.Atoi(s); err == nil {
		return int16(val)
	}
	return -1
}
func parseLevel(s string) (from, to int8) {
	if len(s) == 1 {
		if val, err := strconv.Atoi(s); err == nil {
			return int8(val), int8(val)
		}
	} else if strings.Contains(s, "-") {
		ss := strings.Split(s, "-")
		if val, err := strconv.Atoi(ss[0]); err == nil {
			from = int8(val)
		} else {
			return
		}
		if val, err := strconv.Atoi(ss[1]); err == nil {
			to = int8(val)
		} else {
			return 0, 0
		}
	}
	return
}

func Day2DayStr(d int8) string {
	switch d {
	case 1:
		return "一"
	case 2:
		return "二"
	case 3:
		return "三"
	case 4:
		return "四"
	case 5:
		return "五"
	case 6:
		return "六"
	case 7:
		return "日"
	}
	return ""
}

func parseDay(s string) int8 {
	switch strings.TrimSpace(s) {
	case "一":
		return 1
	case "二":
		return 2
	case "三":
		return 3
	case "四":
		return 4
	case "五":
		return 5
	case "六":
		return 6
	case "日":
		return 7
	}
	return 0
}

func parseLocation(lon, lat string) Location {
	lontitude, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return Location{}
	}
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return Location{}
	}
	return Location{Lon: lontitude, Lat: latitude}
}

func parseTime(s string) (startTime, endtime int16) {
	startTime, endtime = int16(0), int16(0)
	timestr := strings.Split(s, "-")
	if len(timestr) == 2 {
		st, err := strconv.Atoi(timestr[0])
		if err != nil {
			return
		}
		et, err := strconv.Atoi(timestr[1])
		if err != nil {
			return
		}
		startTime = int16(st)
		endtime = int16(et)
	}
	return startTime, endtime
}
