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

const (
	nameCol        = 0 //d
	dayCol         = 1 //d
	timeCol        = 2 //d
	courtCol       = 3 //d
	addressCol     = 4 //d
	levelCol       = 5
	courtCountCol  = 6
	feeMCol        = 7
	feeFCol        = 8
	minBallTypeCol = 9
	lastUpdateCol  = 10
	noteCol        = 11
	lonCol         = 12 //d
	latCol         = 13 //d
)

func NewData(record []string) Data {
	d := Data{Name: record[nameCol], CourtName: record[courtCol], Address: record[addressCol]}
	d.Day = parseDay(record[dayCol])
	d.StartTime, d.EndTime = parseTime(record[timeCol])
	d.Location = parseLocation(record[lonCol], record[latCol])
	d.FromLevel, d.ToLevel = parseLevel(record[levelCol])
	return d
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
