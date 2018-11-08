package pnrhelper

import "strings"

import "github.com/otwdev/galaxylib"

type PriceLine struct {
	PriceList []*Price
}

func NewPriceLine() *PriceLine {
	return &PriceLine{}
}

func (p *PriceLine) Data() []*Price {
	return p.PriceList
}

func (p *PriceLine) IsMatch(line string) bool {
	return strings.HasPrefix(line, "FN/")
}

func (p *PriceLine) Add(pos int, line string) bool {
	if p.IsMatch(line) == false {
		return false
	}
	pItems := strings.Split(line, "/")
	price := &Price{}
	for _, v := range pItems {
		scny := "SCNY"
		if strings.HasPrefix(v, scny) {
			price.ActualPrice = galaxylib.DefaultGalaxyConverter.MustFloat(v[4:])
			continue
		}
		c := "C"
		if len(v) > 2 && strings.HasPrefix(v, c) {
			price.AgencyFees = galaxylib.DefaultGalaxyConverter.MustFloat(v[1:])
			continue
		}
		xcny := "XCNY"
		if strings.HasPrefix(v, xcny) {
			price.Tax = galaxylib.DefaultGalaxyConverter.MustFloat(v[4:])
			continue
		}
		p := "P"
		if len(v) > 1 && strings.HasPrefix(v, p) {
			price.PersonRPH = galaxylib.DefaultGalaxyConverter.MustInt(v[1:])
		}
	}
	p.PriceList = append(p.PriceList, price)
	return true
}

type Price struct {
	PersonRPH      int
	ActualPrice    float64 `json:"amount"`
	Tax            float64 `json:"tax"`
	AgencyFees     float64 `json:"agencyFees"`
	NumberOfPeople int     `json:"numberOfPeople"`
}
