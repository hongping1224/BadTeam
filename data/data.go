package data

import (
	"fmt"
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
	LevelStr      string
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
	UpdatedCol     = 0
	nameCol        = 1  //d
	dayCol         = 2  //d
	timeCol        = 3  //d
	courtCol       = 4  //d
	AddressCol     = 5  //d
	levelCol       = 6  //d
	courtCountCol  = -1 //d
	feeMCol        = 7  //d
	feeFCol        = 8  //d
	minBallTypeCol = -1 //d
	lastUpdateCol  = -1
	noteCol        = 9  //d
	lonCol         = 10 //d
	latCol         = 11 //d
)

func NewData(record []string) Data {
	d := Data{Name: record[nameCol], CourtName: record[courtCol], Address: record[AddressCol], Note: record[noteCol]}
	d.Day = parseDay(record[dayCol])
	d.StartTime, d.EndTime = parseTime(record[timeCol])
	d.Location = parseLocation(record[lonCol], record[latCol])
	d.FromLevel, d.ToLevel = parseLevel(record[levelCol])
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

func LevelToStr(FromLevel, ToLevel int8) (LevelStr string) {
	label := []string{"無標示", "歡樂", "中下", "中等", "中上", "強"}
	if FromLevel == 0 && ToLevel == 0 {
		LevelStr = "-"
	} else if FromLevel == ToLevel {
		LevelStr = fmt.Sprintf("%s", label[FromLevel])
	} else {
		LevelStr = fmt.Sprintf("%s-%s", label[FromLevel], label[ToLevel])
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
