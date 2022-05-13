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

//BSP 支付运价
func (p *PriceLine) bspPrice(priceItem []string) *Price {

	price := &Price{}
	for _, v := range priceItem {
		v = strings.TrimSpace(v)
		//票面
		scny := "SCNY"
		if strings.HasPrefix(v, scny) {
			price.BaseAmount = cast.ToFloat64(v[4:])
			continue
		}
		// 代理费
		c := "C"
		if len(v) > 2 && strings.HasPrefix(v, c) {
			price.AgencyFee = cast.ToFloat64(v[1:])
			continue
		}
		// 税总和
		xcny := "XCNY"
		if strings.HasPrefix(v, xcny) {
			price.Tax = cast.ToFloat64(v[4:])
			continue
		}
		// 乘客序号
		p := "P"
		if len(v) > 1 && strings.HasPrefix(v, p) {
			price.PersonRPH = cast.ToInt(v[1:])
		}
	}
	return price
}

// UATP支付价格计算
func (p *PriceLine) uatpPrice(priceItem []string) *Price {
	price := &Price{}
	for _, v := range priceItem {
		v = strings.TrimSpace(v)
		//票面
		scny := "RCNY"
		if strings.HasPrefix(v, scny) {
			price.BaseAmount = cast.ToFloat64(v[4:])
			continue
		}
		// 代理费
		c := "C"
		if len(v) > 2 && strings.HasPrefix(v, c) {
			price.AgencyFee = cast.ToFloat64(v[1:])
			continue
		}
		// 税总和
		xcny := "BCNY"
		if strings.HasPrefix(v, xcny) {
			price.Tax = cast.ToFloat64(v[4:])
			continue
		}
		// 乘客序号
		p := "P"
		if len(v) > 1 && strings.HasPrefix(v, p) {
			price.PersonRPH = cast.ToInt(v[1:])
		}
	}
	return price
}

type Price struct {
	PersonRPH      int
	BaseAmount     float64 `json:"baseAmount"`
	Tax            float64 `json:"tax"`
	AgencyFee      float64 `json:"agencyFee"`
	NumberOfPeople int     `json:"numberOfPeople"`
	PTC            string  `json:"ptc"`
}
