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
	ArrTime      string  `json:"arrTime"`
	DepDate      string  `json:"depDate"`
	DepTime      string  `json:"depTime"`
	offset       int
	innerYear    int // 长格式 PNR 行里嵌的两位年份（2025 -> 25），0 表示未提供
}

// nowFn 让测试可替换"今天"。生产代码不应 SET，单测里用。
var nowFn = time.Now

// pastTolerance 是"刚刚出发"窗口：解析出来的日期在过去 ≤ 这个窗口内时，
// 仍按当前年算（视为同年起飞、补录场景）；超出窗口才滚到下一年。
const pastTolerance = 7 * 24 * time.Hour

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

	line = strings.ReplaceAll(line, "\x1c", " ")
	line = strings.ReplaceAll(line, "\x1d", " ")

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
	j.Dep.AircaftScheduledDateTime = j.FormatArrDepTime(fields[2], fields[5-j.offset])
	j.Arrival.AircaftScheduledDateTime = j.FormatArrDepTime(fields[2], fields[6-j.offset])
	j.DepDate = fields[2]
	j.DepTime = fields[5-j.offset]
	j.ArrTime = fields[6-j.offset]
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

// formatDate 把 PNR 里的日期字段 ("01FEB", "SA01FEB", "MO10MAR25XIYPEK") 解出 time.Time。
//
// 年份推断：
//  1. 长格式行（"MO10MAR25XIYPEK"）里嵌有两位年份的，直接用，不再猜。
//  2. 其它情况按"7 天窗口"规则推断：先按当前年解析；只有当结果比今天早超过 7 天，
//     才认为这是一笔为来年的预订，加 1 年。窗口内的过去日期（同月已起飞 / 跨月补录）
//     保持在当前年——避免把刚起飞或补录的航班错推到明年。
func (j *Journey) formatDate(input string) time.Time {

	if len(input) > 12 {
		input = j.formatDateWithWeek(input)
	}

	now := nowFn()
	year := now.Year()
	if j.innerYear > 0 {
		year = j.innerYear
	}

	val := fmt.Sprintf("%s%d", input[2:], year)
	t, err := time.Parse("02Jan2006", val)
	if err != nil {
		return time.Time{}
	}

	// 长格式里带了年，照单收下。
	if j.innerYear > 0 {
		return t
	}

	// 7 天窗口：只有"显著在过去"才滚明年。
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, t.Location())
	if today.Sub(t) > pastTolerance {
		t = t.AddDate(1, 0, 0)
	}
	return t
}

// 处理带星期的日期。
// 长格式形如 "MO10MAR25XIYPEK"：前 2 位星期，3-4 日，5-7 月，8-9 两位年（可选），
// 后面接 6 位 dep+arr 三字码。
func (j *Journey) formatDateWithWeek(input string) string {

	j.offset = 1
	if len(input) >= 9 && isTwoDigit(input[7:9]) {
		if yy, err := strconv.Atoi(input[7:9]); err == nil {
			j.innerYear = 2000 + yy
		}
	}
	str := input[:7]
	j.Dep.IATA_LocationCode = input[9:12]
	j.Arrival.IATA_LocationCode = input[12:]
	return str
}

func isTwoDigit(s string) bool {
	if len(s) != 2 {
		return false
	}
	return s[0] >= '0' && s[0] <= '9' && s[1] >= '0' && s[1] <= '9'
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
