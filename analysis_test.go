package travelskypnr

import (
	"os"
	"testing"
)

func readPnrfile(filename string) string {
	fullPath := "./data/" + filename + ".txt"
	buf, _ := os.ReadFile(fullPath)
	return string(buf)
}

func TestOutput(t *testing.T) {
	pnrText := readPnrfile("JWQ2G4")

	analyzer := NewAnalysis(pnrText)
	outer := analyzer.Output()

	if outer == nil {
		t.Error("Output is nil")
		return

	}

	j := outer.Journey[0]
	if j == nil || j.Arrival == nil {
		t.Fail()
		return
	}
	deptime := outer.Journey[0].Dep.AircaftScheduledDateTime
	t.Log(deptime)

}

func TestPsgZero(t *testing.T) {
	txt := readPnrfile("psg_zero_err")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}
	if len(outer.Price) == 0 {
		t.Error("Price is nil")
		return
	}
	for _, p := range outer.Price {
		if p.NumberOfPeople == 0 {
			t.Error("psg is zero")
		}
		t.Log(p.NumberOfPeople)
	}
}

func TestJourneyLost(t *testing.T) {
	txt := readPnrfile("journey_lost")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}
	if len(outer.Journey) != 4 {
		t.Log(len(outer.Journey))
		t.Error("Journey is not 4")
		return
	}
	if len(outer.Person) != 2 {
		t.Log(len(outer.Price))
		t.Error("Price is not 2")
		return
	}
}
