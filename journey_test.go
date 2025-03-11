package travelskypnr

import (
	"fmt"
	"testing"
)


func TestNewJourney(t *testing.T) {
	line := "UA505  C1  TH26DEC  SFOIAH HK3   2359 0549+1        E  3 C "
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	if j.Dep.StationName != "3" {
		t.Errorf("Dep.StationName is not 3")
	}
	if j.Arrival.StationName != "C" {
		t.Errorf("Arrival.StationName is not C")
	}
	fmt.Println(j)
}

func TestNewJourneyWithTerminal(t *testing.T) {
	line := "CZ3800 Z   FR18OCT  TAOCAN UN3   1705 2010          E --T2 S "
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	if j.Dep.StationName != "T2" {
		t.Errorf("Dep.StationName is not T2")
	}
	if j.Arrival.StationName != "" {
		t.Errorf("Arrival.StationName is not empty")
	}
}

func TestNewJourneyWithNoTerminal(t *testing.T) {
	line := "TK2126 Y1 TU10DEC ISTESB HK1 0800 0910 SEAME"
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	if j.FlightNumber != "TK2126" {
		t.Errorf("FlightNumber is not TK2126")
	}
}

func TestFormatDate(t *testing.T) {
	line := "AF111  B1  SA01FEB  PVGCDG HK1   2205 0550+1    SEAME  12E"
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	fmt.Println(j)
}

func TestFormatArrDepTime(t *testing.T) {
	line := "CA1230 L   MO10MAR25XIYPEK HK1   1930 2150          E T2T3"
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	fmt.Println(j)
}

