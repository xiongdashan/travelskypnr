package travelskypnr

import (
	"testing"
)

func TestOutput(t *testing.T) {
	pnrText := ` **ELECTRONIC TICKET PNR**
 1.侯XX KQ7X31
 2.  3U8948 N   MO14MAY  TNAKMG RR1   0050 1355+1          E      N1 
     -CA-PHLF63
 3.SHA/T SHA/T021-64000000/XXX XXX BUSINESS INFORMATION CONSULTING      CO.,LTD/BJ ABCDEFG 
 4.Tb
 5.SSR FOID 3U HK1 NI370103198409270000/P1 
 6.SSR FQTV 3U HK1 TNAKMG 8948 N14MAY 3U618320570/C/P1 
 7.SSR ADTK 1E BY SHA10MAY18/1314 OR CXL 3U8948 N14MAY 
 8.SSR TKNE 3U HK1 TNAKMG 8948 N14MAY 8762090171140/1/P1
 9.OSI 3U CTCM13791040000/P1                                                   
10.RMK CMS/A/**                                                                 11.RMK MP 13791040000   

12.RMK TJ BJS000
13.RMK CA/PHLF63
14.RMK AU REQ0000000
15.RMK AUTOMATIC FARE QUOTE
16.RMK  TP1876XXXXXXX8820 0123 000690.00CNY 8344
17.FN/A/FCNY640.00/SCNY0.00/RCNY640.00/BCNY50.00/C0.00/XCNY50.00/TCNY50.00CN/       TEXEMPTYQ/ACNY690.00
18.TN/876-2090170000/P1
19.FP/CC/TP1876XXXXXXX8820 0123 000690.00CNY 8344   
20.BJS001`

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
