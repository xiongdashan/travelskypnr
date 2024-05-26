package travelskypnr

import (
	"regexp"
	"strconv"
	"strings"
)

type TicketNumberLine struct {
	TicketNumberList []*TicketNumber
	isTn             bool
	ssrError         bool
	ssr              int
}

func NewTktLine() *TicketNumberLine {
	return &TicketNumberLine{}
}

func (t *TicketNumberLine) IsMatch(line string) bool {
	if strings.HasPrefix(line, "SSR TKNE") {
		t.ssr++
		return true
	}
	if strings.HasPrefix(line, "TN/") {
		t.isTn = true
		return true
	}
	return false
}

func (t *TicketNumberLine) Data() []*TicketNumber {
	return t.TicketNumberList
}

const combinDatePattern = `([A-Z]{1}\d{2}[A-Z]{3})`

func (t *TicketNumberLine) Add(pos int, line string) bool {
	if !t.IsMatch(line) {
		return false
	}

	tkt := &TicketNumber{}

	// TN/000-000000000/P1
	if t.isTn && (t.ssrError || t.ssr == 0) {
		regex := regexp.MustCompile(`TN(\/IN)?\/([0-9\-]+)\/P(\d+)`)
		if !regex.MatchString(line) {
			return true
		}
		match := regex.FindAllStringSubmatch(line, -1)[0]
		tkt.Number = match[2]
		tkt.PersonRPH = t.rphToi(match[3])

		//婴儿票
		if match[1] != "" {
			tkt.Type = Infant
		}

		//return true

	} else {

		reg := regexp.MustCompile(combinDatePattern)

		idx := reg.FindAllStringIndex(line, -1)

		if len(idx) == 0 {
			t.ssrError = true
			return true
		}

		prefix := line[:idx[0][0]]

		tktItmes := strings.Fields(prefix)

		//tktItmes := strings.Fields(line)
		// SSR TKNE CA HK1 PEKMEL 165 W26SEP (INF)? 9992876664435/1/P2

		if len(tktItmes) <= 2 {
			// 非自有配制出票，SSR信息解析异常，从TN项获取票号
			t.ssrError = true
			return true
		}

		tkt.Airline = tktItmes[2]

		suffix := line[idx[0][1]:]

		rph := strings.Split(suffix, "/")

		tkt.Number = strings.TrimSpace(rph[0])

		//如果是婴儿票号
		if strings.HasPrefix(tkt.Number, Infant) {
			tkt.Number = tkt.Number[3:]
			tkt.Type = Infant
		}
		tkt.JourneyRPH, _ = strconv.Atoi(rph[1])
		tkt.PersonRPH = t.rphToi(rph[2])
	}

	tkt.Number = strings.Replace(tkt.Number, "-", "", -1)
	tkt.Number = strings.TrimSpace(tkt.Number)

	for _, v := range t.TicketNumberList {
		if v.Number == tkt.Number {
			return true
		}
	}

	t.TicketNumberList = append(t.TicketNumberList, tkt)
	//fmt.Println(tkt.Number)
	return true
}

func (t *TicketNumberLine) rphToi(rph string) int {
	rph = strings.TrimSpace(rph)
	rev, _ := strconv.Atoi(strings.Replace(rph, "P", "", -1))
	return rev
}

type TicketNumber struct {
	Airline    string `json:"airLine"`
	Number     string `json:"number"`
	JourneyRPH int    `json:"journeyRPH"`
	PersonRPH  int    `json:"personRPH"`
	Type       string `json:"type"`
}
