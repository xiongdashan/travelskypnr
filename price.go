package travelskypnr

import (
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

type PriceLine struct {
	PriceList []*Price
	IsUATP    bool
	isINF     bool
	matches   map[string]string
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

const uatpMatch = `RMK(\s+)TP[0-9X]+`

func (p *PriceLine) Add(pos int, line string) bool {

	if ok, _ := regexp.MatchString(uatpMatch, line); ok {
		p.IsUATP = true
	}

	if !p.IsMatch(line) {
		return false
	}
	pItems := strings.Split(line, "/")

	//婴儿运价
	if pItems[1] == "IN" {
		p.isINF = true
	}

	p.matcheLine(line)

	var price *Price

	if p.IsUATP {
		price = p.uatpPrice(pItems)
	} else {
		price = p.bspPrice(pItems)
	}

	// 目前只从扫描区里获取姓名，不包含婴，如果价格中含P，为儿童价或婴儿价
	if price.PersonRPH > 0 {
		if p.isINF {
			price.PTC = Infant
		} else {
			price.PTC = Child
		}

	} else {
		price.PTC = Adult
	}
	p.PriceList = append(p.PriceList, price)
	return true
}

// BSP 支付运价
func (p *PriceLine) bspPrice(priceItem []string) *Price {

	price := &Price{}
	price.BaseAmount = cast.ToFloat64(p.matches["SCNY"])
	price.Tax = cast.ToFloat64(p.matches["XCNY"])
	price.AgencyFee = cast.ToFloat64(p.matches["C"])
	price.YQ = cast.ToFloat64(p.matches["TCNY"])
	price.ToRefPsg(p.matches["P"])
	return price
}

// 各种正则表达式
var expressions = map[string]*regexp.Regexp{
	"FCNY": regexp.MustCompile(`FCNY(\d+\.\d+)`),
	"RCNY": regexp.MustCompile(`RCNY(\d+\.\d+)`),
	"SCNY": regexp.MustCompile(`SCNY(\d+\.\d+)`),
	"XCNY": regexp.MustCompile(`XCNY(\d+\.\d+)`),
	"C":    regexp.MustCompile(`C(\d+\.\d+)`),
	"TCNY": regexp.MustCompile(`TCNY(\d+\.\d+)(?:CN|YQ)?`),
	"ACNY": regexp.MustCompile(`ACNY(\d+\.\d+)`),
	"P":    regexp.MustCompile(`P(\d+(?:/\d+)*)`),
}

func (p *PriceLine) matcheLine(l string) {
	p.matches = make(map[string]string)
	for k, v := range expressions {
		match := v.FindStringSubmatch(l)
		if len(match) > 1 {
			p.matches[k] = match[1]
		}
	}
}

var UatpExpressions = map[string]*regexp.Regexp{
	"RCNY": regexp.MustCompile(`RCNY(\d+\.\d+)`),
	"SCNY": regexp.MustCompile(`SCNY(\d+\.\d+)`),
	"BCNY": regexp.MustCompile(`XCNY(\d+\.\d+)`),
	"C":    regexp.MustCompile(`C(\d+\.\d+)`),
	"TCNY": regexp.MustCompile(`TCNY(\d+\.\d+)(?:CN|YQ)?`),
	"ACNY": regexp.MustCompile(`ACNY(\d+\.\d+)`),
	"P":    regexp.MustCompile(`P(\d+(?:/\d+)*)`),
}

// UATP支付价格计算
func (p *PriceLine) uatpPrice(priceItem []string) *Price {
	price := &Price{}
	price.BaseAmount = cast.ToFloat64(p.matches["RCNY"])
	price.Tax = cast.ToFloat64(p.matches["BCNY"])
	price.AgencyFee = cast.ToFloat64(p.matches["C"])
	price.YQ = cast.ToFloat64(p.matches["TCNY"])
	price.ToRefPsg(p.matches["P"])
	return price
}

type Price struct {
	PersonRPH      int
	BaseAmount     float64 `json:"baseAmount"`
	Tax            float64 `json:"tax"`
	YQ             float64 `json:"yq"`
	AgencyFee      float64 `json:"agencyFee"`
	NumberOfPeople int     `json:"numberOfPeople"`
	PTC            string  `json:"ptc"`
	RefPsg         []int   `json:"refPsg"`
}

func (p *Price) ToRefPsg(matchedP string) {

	ptc := strings.Split(matchedP, "/")

	for _, v := range ptc {
		if v == "" {
			continue
		}
		p.NumberOfPeople++
		p.RefPsg = append(p.RefPsg, cast.ToInt(v))
	}
}

func (p *Price) include(psgPtc int) bool {
	for _, v := range p.RefPsg {
		if v == psgPtc {
			return true
		}
	}
	return false
}
