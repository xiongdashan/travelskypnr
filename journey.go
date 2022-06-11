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
	j.Regex = regexp.MustCompile(`(\w+)\s+([A-Z0-9]{1,2})\s+(([A-Z]{2})(\d{2})([A-Z]{3}))(\s+|\d{2})([A-Z]{6})\s?([A-Z0-9]{2,3})\s+(\d{4})\s+((\d{4})(\+\d{1})?)\s+([A-Z]{1})`)
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

func (j *JourneyLine) Add(pos int, line string) bool {

	line = strings.TrimSpace(line)

	if !j.IsMatch(line) {
		return false
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
	return true
}

type Journey struct {
	RPH          int
	FlightNumber string `json:"flightNumber"`
	CabinClass   string `json:"cabinClass"`
	Terminal     string `json:"terminal"`
	innerDptDate time.Time
	Arrival      *ArrDep `json:"arrival"`
	Dep          *ArrDep `json:"dep"`
}

type ArrDep struct {
	AircaftScheduledDateTime string `json:"aircraftScheduledDateTime"`
	BoardingGateID           string `json:"boardingGateID"`
	IATA_LocationCode        string `json:"iataLocationCode"`
	StationName              string `json:"stationName"`
	TerminalName             string `json:"terminalName"`
}

func (jl *JourneyLine) newJourney(line string) *Journey {
	line = strings.TrimSpace(line)

	matche := jl.Regex.FindAllStringSubmatch(line, -1)[0]

	j := &Journey{
		Arrival: &ArrDep{},
		Dep:     &ArrDep{},
	}
	j.FlightNumber = matche[1]
	j.CabinClass = matche[2]
	j.innerDptDate = j.formatDate(matche[3])
	j.Terminal = matche[14]
	j.Dep.IATA_LocationCode = matche[8][:3]
	j.Dep.AircaftScheduledDateTime = j.FormatArrDepTime(matche[3], matche[10])
	j.Arrival.AircaftScheduledDateTime = j.FormatArrDepTime(matche[3], matche[11])
	j.Arrival.TerminalName = matche[14]
	j.Arrival.IATA_LocationCode = matche[8][3:]

	return j
}

func (j *Journey) formatDate(input string) time.Time {
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

func (j *Journey) FormatArrDepTime(date, timeVal string) string {
	formatedDate := j.formatDate(date)
	splitedTime := strings.Split(timeVal, "+")
	houre, _ := strconv.Atoi(splitedTime[0][:2])
	minute, _ := strconv.Atoi(splitedTime[0][2:])
	if len(splitedTime) >= 2 {
		day, _ := strconv.Atoi(splitedTime[1])
		formatedDate = formatedDate.AddDate(0, 0, day)
	}
	return time.Date(formatedDate.Year(), formatedDate.Month(), formatedDate.Day(), houre, minute, 0, 0, formatedDate.Location()).Format("2006-01-02 15:04:05")
}
