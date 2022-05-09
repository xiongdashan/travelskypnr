package travelskypnr

import (
	"fmt"
	"regexp"
	"strings"
)

type Analysis struct {
	PnrTxt string
}

func NewAnalysis(txt string) *Analysis {

	txt = strings.Replace(txt, "�b", "", -1)

	return &Analysis{
		PnrTxt: txt,
	}
}

type PNRInfo struct {
	Code        string          `json:"code"`
	IsUATP      bool            `jons:"isUATP"`
	Journey     []*Journey      `json:"journey"`
	Person      []*Person       `json:"person"`
	Price       []*Price        `json:"price"`
	TicketNumer []*TicketNumber `json:"tktNumber"`
}

const (
	Adult  = "ADT"
	Child  = "CHD"
	Infant = "INF"
)

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
		if tl.Add(i, l) {
			continue
		}

		//价格
		priceLn.Add(i, l)
	}

	// rev := &PNRInfo{
	// 	Code:        pl.PnrCode,
	// 	Journey:     j.Data(),
	// 	Person:      pl.Data(),
	// 	Price:       priceLn.Data(),
	// 	TicketNumer: tl.Data(),
	// }

	rev := &PNRInfo{
		Code:   pl.PnrCode,
		IsUATP: priceLn.IsUATP,
	}

	for k, p := range pl.Dict {
		for _, t := range tl.TicketNumberList {
			key := fmt.Sprintf("P%d", t.PersonRPH)
			if t.Type == Infant {
				key = fmt.Sprintf("P%dINF", t.PersonRPH)
			}
			if k == key {
				// 判断是否已经存在
				has := false
				for _, pt := range p.TktAry {
					if pt == t.Number {
						has = true
						continue
					}
				}
				if !has {
					p.TktAry = append(p.TktAry, t.Number)
				}
			}
		}

		p.TktStr()

		for _, pr := range priceLn.PriceList {
			if p.PTC == pr.PTC {
				pr.NumberOfPeople++
			}
		}
		rev.Person = append(rev.Person, p)
	}

	rev.Price = priceLn.PriceList
	rev.Journey = j.JourneyList
	return rev
}
