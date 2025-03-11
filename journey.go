package travelskypnr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type JourneyLine struct {
	Regex       *regexp.Regexp
	JourneyList []*Journey
	isARNK      bool
}

var arnk = "ARNK"

func NewJourneyLine() *JourneyLine {
	j := &JourneyLine{}
	j.Regex = regexp.MustCompile(`(\w+)\s+([A-Z0-9]{1,2})\s+[A-Z]{2}\d{2}[A-Z]{3}`)
	return j
}

func (j *JourneyLine) Data() []*Journey {
	return j.JourneyList
}

func (j *JourneyLine) IsMatch(line string) bool {

	if strings.HasPrefix(line, arnk) {
		j.isARNK = true
		return true
	}
	return j.Regex.MatchString(strings.TrimSpace(line))
}

func (j *JourneyLine) Add(pos int, line string) (*Journey, bool) {

	line = strings.TrimSpace(line)

	if !j.IsMatch(line) {
		return nil, false
	}

	var jny *Journey

	//地面段
	if j.isARNK {
		jny = &Journey{
			FlightNumber: "ARNK",
		}
		j.isARNK = false
	} else {
		jny = j.newJourney(line)
	}

	jny.RPH = len(j.JourneyList) + 1
	j.JourneyList = append(j.JourneyList, jny)
	return jny, true
}

type Journey struct {
	RPH          int
	FlightNumber string `json:"flightNumber"`
	CabinClass   string `json:"cabinClass"`
	innerDptDate time.Time
	Arrival      *ArrDep `json:"arrival"`
	Dep          *ArrDep `json:"dep"`
	ArrTime      string   `json:"arrTime"`
	DepDate      string   `json:"depDate"`
	DepTime      string   `json:"depTime"`
	offset       int
}

type ArrDep struct {
	AircaftScheduledDateTime string `json:"aircraftScheduledDateTime"`
	BoardingGateID           string `json:"boardingGateID"`
	IATA_LocationCode        string `json:"iataLocationCode"`
	StationName              string `json:"stationName"`
}



// 新版解析

/***************


[0] =
"CZ8233"
[1] =
"Z"
[2] =
"WE14AUG"
[3] =
"CANTFU"
[4] =
"RR4"
[5] =
"1420"
[6] =
"1640"
[7] =
"E"
[8] =
"T2T2"
[9] =
"-CA-NT0ER3"
*****************/

func (jl *JourneyLine) newJourney(line string) *Journey {
	line = strings.TrimSpace(line)

	fields := strings.Fields(line)

	j := &Journey{
		Arrival: &ArrDep{},
		Dep:     &ArrDep{},
	}
	j.FlightNumber = fields[0]
	j.CabinClass = fields[1]
	j.innerDptDate = j.formatDate(fields[2])
	if j.Dep.IATA_LocationCode == "" {
		j.Dep.IATA_LocationCode = fields[3][:3]
	}
	if j.Arrival.IATA_LocationCode == "" {
		j.Arrival.IATA_LocationCode = fields[3][3:]
	}
	j.Dep.AircaftScheduledDateTime = j.FormatArrDepTime(fields[2], fields[5 - j.offset])
	j.Arrival.AircaftScheduledDateTime = j.FormatArrDepTime(fields[2], fields[6 - j.offset])
	j.DepDate = fields[2]
	j.DepTime = fields[5 - j.offset]
	j.ArrTime = fields[6 - j.offset]
	jl.formatTerminal(fields, j)
	return j
}



func (j *JourneyLine) formatTerminal(fields []string, jny *Journey) {
	if len(fields) < 9 {
		return
	}
	t := fields[8]
	if strings.HasPrefix(t, "--") {
		jny.Dep.StationName = strings.TrimPrefix(t, "--")
		return
	}
	if strings.HasSuffix(t, "--") {
		jny.Arrival.StationName = strings.TrimSuffix(t, "--")
		return
	}
	reg := regexp.MustCompile(`T(\d{1,2})T(\d{1,2})`)
	if reg.MatchString(t) {
		matche := reg.FindStringSubmatch(t)
		jny.Dep.StationName = fmt.Sprintf("T%s", matche[1])
		jny.Arrival.BoardingGateID = fmt.Sprintf("T%s", matche[2])
		return
	}
	if len(fields) >= 10 {
		jny.Dep.StationName = fields[8]
		jny.Arrival.StationName = fields[9]
	}
}


func (j *Journey) formatDate(input string) time.Time {
   
	if len(input) > 12 {
		input = j.formatDateWithWeek(input)
	}

	val := fmt.Sprintf("%s%d", input[2:], time.Now().Year())
	t, err := time.Parse("02Jan2006", val)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.Format("2006-01-02"))
	if t.Month() < time.Now().Month() {
		t = t.AddDate(1, 0, 0)
	}

	return t
}

// 处理带星期的日期
func (j *Journey) formatDateWithWeek(input string) string {
	
  //MO10MAR25XIYPEK
  j.offset = 1
  str := input[:7]	
  j.Dep.IATA_LocationCode = input[9:12]
  j.Arrival.IATA_LocationCode = input[12:]
  return str
}



func (j *Journey) FormatArrDepTime(date, timeVal string) string {
	formatedDate := j.innerDptDate
	splitedTime := strings.Split(timeVal, "+")
	houre, _ := strconv.Atoi(splitedTime[0][:2])
	minute, _ := strconv.Atoi(splitedTime[0][2:])
	if len(splitedTime) >= 2 {
		day, _ := strconv.Atoi(splitedTime[1])
		formatedDate = formatedDate.AddDate(0, 0, day)
	}
	return time.Date(formatedDate.Year(), formatedDate.Month(), formatedDate.Day(), houre, minute, 0, 0, formatedDate.Location()).Format("2006-01-02 15:04")
}
