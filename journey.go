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
}

func NewJourneyLine() *JourneyLine {
	j := &JourneyLine{}
	j.Regex = regexp.MustCompile(`(\w+)\s+([A-Z0-9]{1,2})\s+([A-Z]{2})(\d{2})([A-Z]{3})\s+([A-Z]{6})\s+([A-Z0-9]{2,3})\s+(\d{4})\s+(\d{4})(\+(\d{1}))?\s+([A-Z]{1})`)
	return j
}

func (j *JourneyLine) Data() []*Journey {
	return j.JourneyList
}

func (j *JourneyLine) IsMatch(line string) bool {
	return j.Regex.MatchString(strings.TrimSpace(line))
}

func (j *JourneyLine) Add(pos int, line string) bool {
	if j.IsMatch(line) == false {
		return false
	}
	line = strings.TrimSpace(line)

	var jny *Journey

	//地面段
	if strings.HasPrefix(line, "ARNK") {
		jny = &Journey{
			FlightNumber: "ARNK",
		}
	} else {
		jny = newJourney(line)
	}

	jny.RPH = len(j.JourneyList) + 1
	j.JourneyList = append(j.JourneyList, jny)
	//fmt.Println(jny.FlightNumber)
	return true
}

type Journey struct {
	RPH          int
	FlightNumber string    `json:"flightNumber"`
	Combin       string    `json:"combin"`
	DepartDate   time.Time `json:"departDate"`
	DepartTime   string    `json:"departTime"`
	ArrDate      time.Time `json:"arrDate"`
	ArrTime      string    `json:"arrTime"`
	DepartCode   string    `json:"departCode"`
	ArrCode      string    `json:"arrCode"`
	Terminal     string    `json:"terminal"`
}

func newJourney(line string) *Journey {
	line = strings.TrimSpace(line)
	regex := regexp.MustCompile(`\s+`)
	itemAry := regex.Split(line, -1)
	j := &Journey{}
	j.FlightNumber = itemAry[0]
	j.Combin = itemAry[1]
	j.DepartDate = j.formatDate(itemAry[2])
	j.DepartCode = itemAry[3][:3]
	j.ArrCode = itemAry[3][3:]
	j.DepartTime = itemAry[5]
	j.ArrTime = j.formatTime(itemAry[6])
	if len(itemAry) >= 9 {
		j.Terminal = itemAry[8]
	}
	return j
}

func (j *Journey) formatDate(input string) time.Time {
	val := fmt.Sprintf("%s%d", input[2:], time.Now().Year())
	t, _ := time.Parse("02Jan2006", val)
	//t = t.AddDate(1, 0, 0)

	return t
}

func (j *Journey) formatTime(input string) string {
	regex := regexp.MustCompile(`(\d{4})\+(\d+)`)
	if regex.MatchString(input) == false {
		j.ArrDate = j.DepartDate
		return input
	}
	match := regex.FindAllStringSubmatch(input, -1)[0]
	val, _ := strconv.Atoi(match[2])
	j.ArrDate = j.DepartDate.AddDate(0, 0, val)
	return match[1]
}
