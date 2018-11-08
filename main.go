package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/otwdev/blueskypnr/models"
)

const pnrstr = `1.ZHONG/XXXXXX H00000 
2.  3U8085 H   WE19SEP  CTUNRT HK1   1145 1755          E T11  
	-CA-M99999 
3.  3U8086 Q   SU30SEP  NRTCTU HK1   2020 0030+1        E 1 T1 
	-CA-M99999 
4.PEK/T PEK/T010-00000000/xxxxxxxx (BEIJING) xxxxxxxxxx TRAVEL SERVICE CO.,   
   LTD ABCDEFG 
5.TL/0945/19SEP/PEK000 
6.SSR ADTK 1E BY PEK08SEP18/2147 OR CXL 3U8085 H19SEP  
7.SSR DOCS 3U HK1 P/CN/E80000000/CN/12DEC53/M/08OCT26/ZHONG/xxxxxx/P1  
8.OSI 3U CTCT10000000000   
9.OSI 3U CTCM10000000000/P1                                                   
10.RMK NAME xiiiiiixii                                                             

12.RMK CA/M99999
13.PEK000`

//FCNY票面价/SCNY是实收票款/XCNY是所有税的总和/TCNY机场建设费cn/TCNY燃油费YQ/ACNY票面+税的和如果是含儿童票的，后面跟有乘客序号/Pn例：

func main() {

	pos := strings.Index(pnrstr, "1.")
	pnr := strings.TrimSpace(pnrstr[pos:])
	//fmt.Println(pnr)
	regex := regexp.MustCompile(`\b(\d+)\.`)
	lines := regex.Split(pnr, -1)
	pl := models.NewPersonLine()
	j := models.NewJourneyLine()
	tl := models.NewTktLine()
	priceLn := models.NewPriceLine()

	for i, l := range lines {
		if l == "" {
			continue
		}

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

		priceLn.Add(i, l)
	}

	output := &struct {
		Code        string
		Journey     interface{}
		Person      interface{}
		Price       interface{}
		TicketNumer interface{}
	}{
		Code:        pl.PnrCode,
		Journey:     j.Data(),
		Person:      pl.Data(),
		Price:       priceLn.Data(),
		TicketNumer: tl.Data(),
	}

	buf, _ := json.Marshal(output)

	fmt.Println(string(buf))
	// t, _ := time.Parse("02Jan", "WE26SEP"[2:])
	// fmt.Println(t)

}
