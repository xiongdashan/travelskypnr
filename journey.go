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
	j.Regex = regexp.MustCompile(`(\w+)\s+([A-Z0-9]{1,2})\s+(([A-Z]{2})(\d{2})([A-Z]{3}))(\s+|\d{2})([A-Z]{6})\s+([A-Z0-9]{2,3})\s+(\d{4})\s+((\d{4})(\+\d{1})?)\s+([A-Z]{1})`)
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

	if j.IsMatch(line) == false {
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

func (jl *JourneyLine) newJourney(line string) *Journey {
	line = strings.TrimSpace(line)

	matche := jl.Regex.FindAllStringSubmatch(line, -1)[0]

	//regex := regexp.MustCompile(`\s+`)
	//itemAry := regex.Split(line, -1)
	j := &Journey{}
	j.FlightNumber = matche[1] //itemAry[0]
	j.Combin = matche[2]       //itemAry[1]
	j.DepartDate = j.formatDate(matche[3] /*itemAry[2]*/)
	j.DepartCode = matche[8][:3]         //itemAry[3][:3]
	j.ArrCode = matche[8][3:]            //itemAry[3][3:]
	j.DepartTime = matche[10]            //itemAry[5]
	j.ArrTime = j.formatTime(matche[11]) //j.formatTime(itemAry[6])
	j.Terminal = matche[14]
	// if len(itemAry) >= 9 {
	// 	j.Terminal = itemAry[8]
	// }
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
