package travelskypnr

import (
	"io/ioutil"
	"testing"
)

func readPnrfile(filename string) string {
	fullPath := "./data/" + filename + ".txt"
	buf, _ := ioutil.ReadFile(fullPath)
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
