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
	Code         string          `json:"code"`
	IsUATP       bool            `jons:"isUATP"`
	Journey      []*Journey      `json:"journey"`
	Person       []*Person       `json:"person"`
	Price        []*Price        `json:"price"`
	TicketNumer  []*TicketNumber `json:"tktNumber"`
	OfficeNumber string          `json:"officeNumber"`
	Expired      string          `json:"expired"`
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

		if _, ok := j.Add(i, l); ok {
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
		rev.Person = append(rev.Person, p)

	}

	rev.Price = priceLn.PriceList
	rev.Journey = j.JourneyList
	rev.OfficeNumber = a.getOfficeNum(lines)
	rev.TicketNumer = tl.TicketNumberList
	rev.Price = a.setPrice(priceLn, rev)

	return rev
}



func (a *Analysis) getOfficeNum(lines []string) string {
	tail := strings.TrimSpace(lines[len(lines)-1])
	if match, _ := regexp.MatchString(`$[A-Z0-9]{6}^`, tail); match {
		return tail
	}
	return ""
}



//重新设置返回价格，因为FN是多人价格，返回时只需要返回每种乘客类型的价格即可

func (a *Analysis) setPrice(priceLn *PriceLine, rev *PNRInfo) []*Price {

	priceMap := make(map[string]*Price)

	for _, p := range rev.Person {
		for _, pr := range priceLn.PriceList {
			if pr.include(p.RPH) {
				if _, ok := priceMap[p.PTC]; !ok {
					priceMap[p.PTC] = pr
					priceMap[p.PTC].NumberOfPeople = 1;
					priceMap[p.PTC].PTC = p.PTC
				} else {
					priceMap[p.PTC].NumberOfPeople++
				}
			}
		}
	}

	rev.Price = make([]*Price, 0)
	for _, v := range priceMap {
		rev.Price = append(rev.Price, v)
	}

	return rev.Price
}
