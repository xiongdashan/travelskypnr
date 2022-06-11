中航信GDS订座PNR文本解析
===========

``` go
go get -u -v github.com/xiongdashan/travelskypnr/v2
```

## 更新
- 2022-05-09 因为有项目使用，回来了...
- 2018-11-14 支持解析UATP，支持婴儿
- 2018-11-22 多个航段且票号相同时，返回乘客对应票号去重



## 说明
对Blueskey的PNR文件解析，然后格式化输出

原文本

```
    1.ZHONG/XXXXXX H00000 
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
13.PEK000
```

解析后输出：

``` json
{
    "Code": "H00000",
    "Journey": [{
        "RPH": 1,
        "FlightNumber": "3U8085",
        "Combin": "H",
        "DepartDate": "2018-09-19T00:00:00Z",
        "DepartTime": "1145",
        "ArrDate": "2018-09-19T00:00:00Z",
        "ArrTime": "1755",
        "DepartCode": "CTU",
        "ArrCode": "NRT",
        "Terminal": "T11"
    }, {
        "RPH": 2,
        "FlightNumber": "3U8086",
        "Combin": "Q",
        "DepartDate": "2018-09-30T00:00:00Z",
        "DepartTime": "2020",
        "ArrDate": "2018-10-01T00:00:00Z",
        "ArrTime": "0030",
        "DepartCode": "NRT",
        "ArrCode": "CTU",
        "Terminal": "1"
    }],
    "Person": [{
        "RPH": 1,
        "Name": "ZHONG/XXXXXX",
        "Type": "成人",
        "IDType": "P",
        "IDNumber": "E80000000",
        "IDIssue": "CN",
        "Nationality": "CN",
        "Birthday": "12DEC53",
        "Gender": "M",
        "Expired": "08OCT26",
        "Mobile": "10000000000"
    }],
    "Price": null,
    "TicketNumer": null
}
```