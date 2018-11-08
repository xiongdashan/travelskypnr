package pnrhelper

import (
	"regexp"
	"strings"
)

type Analysis struct {
	PnrTxt string
}

func NewAnalysis(txt string) *Analysis {
	return &Analysis{
		PnrTxt: txt,
	}
}

type PNRInfo struct {
	Code        string          `json:"code"`
	Journey     []*Journey      `json:"journey"`
	Person      []*Person       `json:"person"`
	Price       []*Price        `json:"price"`
	TicketNumer []*TicketNumber `json:"tktNumber"`
}

func (a *Analysis) Output() *PNRInfo {
	pos := strings.Index(a.PnrTxt, "1.")
	if pos == -1 {
		return nil
	}
	pnr := strings.TrimSpace(a.PnrTxt[pos:])
	//fmt.Println(pnr)
	regex := regexp.MustCompile(`\b(\d+)\.`)
	lines := regex.Split(pnr, -1)
	pl := NewPersonLine()
	j := NewJourneyLine()
	tl := NewTktLine()
	priceLn := NewPriceLine()

	for i, l := range lines {

		if l == "" {
			continue
		}

		l = strings.Replace(l, "\n", " ", -1)

		if j.Add(i, l) {
			pl.End()
		}
		// 扫描区姓名
		if pl.Add(i, l) {
			continue
		}

		// 票号
		if tl.Add(i, l, pl) {
			continue
		}

		priceLn.Add(i, l)
	}

	rev := &PNRInfo{
		Code:        pl.PnrCode,
		Journey:     j.Data(),
		Person:      pl.Data(),
		Price:       priceLn.Data(),
		TicketNumer: tl.Data(),
	}

	return rev
}
