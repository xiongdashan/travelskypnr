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
		t.Log(len(outer.Person))
		t.Error("Person is not 2")
		return
	}
}

func TestPsgCHD(t *testing.T) {
	txt := readPnrfile("psg_chd")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}

	if len(outer.Price) != 2 {
		t.Error("Price is not 2")
		return
	}
}

func TestNoPrice(t *testing.T) {
	txt := readPnrfile("no_price")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}
	if len(outer.Price) == 0 {
		t.Error("Price is not 0")
	}
}

func TestSimple(t *testing.T) {
	txt := readPnrfile("simple")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}
	if len(outer.Journey) != 1 {
		t.Error("Journey is not 2")
		return
	}
}

func TestPsgNoEqu(t *testing.T) {
	txt := readPnrfile("psg_no")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}
	if len(outer.Person) != 7 {
		t.Error("Person is not 7")
		return
	}
}

func TestBigPnr(t *testing.T) {
	txt := readPnrfile("bigpnr")
	al := NewAnalysis(txt)
	outer := al.Output()
	if outer == nil {
		t.Error("Output is nil")
		return
	}
	if len(outer.Journey) != 1 {
		t.Error("Journey is not 1")
		return
	}
	if outer.BigCode != "NBMLJF" {
		t.Error("BigCode is not CA/NBMLJF")
		return
	}
}