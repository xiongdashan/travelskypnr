package travelskypnr

import (
	"strconv"
	"strings"
)

type TicketNumberLine struct {
	TicketNumberList []*TicketNumber
}

func NewTktLine() *TicketNumberLine {
	return &TicketNumberLine{}
}

func (t *TicketNumberLine) IsMatch(line string) bool {
	return strings.HasPrefix(line, "SSR TKNE")
}

func (t *TicketNumberLine) Data() []*TicketNumber {
	return t.TicketNumberList
}

func (t *TicketNumberLine) Add(pos int, line string, pl *PersonLine) bool {
	if t.IsMatch(line) == false {
		return false
	}
	tktItmes := strings.Fields(line)
	tkt := &TicketNumber{}
	// SSR TKNE CA HK1 PEKMEL 165 W26SEP 9992876664435/1/P2
	tkt.Airline = tktItmes[2]
	rph := strings.Split(tktItmes[7], "/")
	tkt.Number = rph[0]

	for _, v := range t.TicketNumberList {
		if v.Number == tkt.Number {
			return true
		}
	}

	tkt.JourneyRPH, _ = strconv.Atoi(rph[1])
	tkt.PersonRPH, _ = strconv.Atoi(strings.Replace(rph[2], "P", "", -1))
	pl.SetTktNumber(tkt.PersonRPH, tkt.Number)
	t.TicketNumberList = append(t.TicketNumberList, tkt)
	//fmt.Println(tkt.Number)
	return true
}

type TicketNumber struct {
	Airline    string `json:"airLine"`
	Number     string `json:"number"`
	JourneyRPH int    `json:"journeyRPH"`
	PersonRPH  int    `json:"personRPH"`
}
