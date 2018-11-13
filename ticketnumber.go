package travelskypnr

import (
	"strconv"
	"strings"
)

type TicketNumberLine struct {
	TicketNumberList []*TicketNumber
	isTn             bool
}

func NewTktLine() *TicketNumberLine {
	return &TicketNumberLine{}
}

func (t *TicketNumberLine) IsMatch(line string) bool {
	if strings.HasPrefix(line, "SSR TKNE") {
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

func (t *TicketNumberLine) Add(pos int, line string, pl *PersonLine) bool {
	if t.IsMatch(line) == false {
		return false
	}

	tkt := &TicketNumber{}
	//rph := ""

	// TN/000-000000000/P1
	if t.isTn {
		itemAry := strings.Split(line, "/")
		tkt.Number = itemAry[1]
		tkt.PersonRPH = t.rphToi(itemAry[2]) //strconv.Atoi(strings.Replace(itemAry[2], "P", "", -1))
	} else {
		tktItmes := strings.Fields(line)
		// SSR TKNE CA HK1 PEKMEL 165 W26SEP 9992876664435/1/P2
		if len(tktItmes) < 8 {
			return false
		}
		tkt.Airline = tktItmes[2]
		rph := strings.Split(tktItmes[7], "/")
		tkt.Number = rph[0]
		tkt.JourneyRPH, _ = strconv.Atoi(rph[1])
		tkt.PersonRPH = t.rphToi(rph[2]) //strconv.Atoi(strings.Replace(rph[2], "P", "", -1))
	}

	for _, v := range t.TicketNumberList {
		if v.Number == tkt.Number {
			return true
		}
	}

	pl.SetTktNumber(tkt.PersonRPH, tkt.Number)
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
}
