package pnrhelper

import (
	"fmt"
	"regexp"
	"strings"
)

type PersonLine struct {
	PnrCode string
	dict    map[string]*Person
	isPass  bool
	isSSR   bool
}

func NewPersonLine() *PersonLine {
	p := &PersonLine{}
	p.dict = make(map[string]*Person)
	return p
}

func (p *PersonLine) Data() (rev []*Person) {
	for _, v := range p.dict {
		rev = append(rev, v)
	}
	return
}

func (p *PersonLine) IsMatch(line string) bool {

	return p.isPass == false
}

// End 扫描区姓名结束...
func (p *PersonLine) End() {
	p.isPass = true
}

// 扫描区姓名
func (p *PersonLine) Add(pos int, line string) bool {

	if p.IsMatch(line) == false {
		// 如果扫描区已过，尝试匹配SSR信息
		if p.AddSSR(line) {
			return true
		}
		// 尝试CTCM
		if p.ctcm(line) {
			return true
		}
		return false
	}

	person := &Person{}
	person.RPH = pos
	pnrNum := person.splitName(line)
	key := fmt.Sprintf("P%d", pos)
	p.dict[key] = person
	if pnrNum != "" {
		p.PnrCode = pnrNum
	}

	fmt.Println(person.Name)

	return true

}

func (p *PersonLine) SetTktNumber(rph int, num string) {
	for _, v := range p.dict {
		if v.RPH == rph {
			v.TicketNumber = append(v.TicketNumber, num)
		}
	}
}

// ..AddSSR 证件信息
func (p *PersonLine) AddSSR(line string) bool {

	if strings.HasPrefix(line, "SSR DOCS") {

		p.ssr(line)
		return true
	}
	return false
}

// 证件类型/发证国家/证件号码/国籍/出生日期/性别/证件有效期限/SURNAME(姓)/FIRST-NAME(名)/MID-NAME(中间名)/持有人标识H/P1
// 0P/      1 CN/   2 E30028197/3 CN/4 24AUG79/5 F/ 6   12SEP23/   LU/       FANGFANG/ P4
func (p *PersonLine) ssr(line string) {
	aryItem := strings.Fields(line)
	if len(aryItem) < 5 {
		return
	}
	idInfostr := aryItem[4]
	idItem := strings.Split(idInfostr, "/")
	key := idItem[len(idItem)-1]
	psn := p.dict[key]
	psn.IDType = idItem[0]
	psn.IDIssue = idItem[1]
	psn.IDNumber = idItem[2]
	psn.Nationality = idItem[3]
	psn.Birthday = idItem[4]
	psn.Gender = idItem[5]
	psn.Expired = idItem[6]
	fmt.Println(p.dict[key].IDNumber)
}

func (p *PersonLine) ctcm(line string) bool {
	regex := regexp.MustCompile(`OSI\s+[A-Z0-9]{2}\s+CTCM(\d{11})\/([P0-9\/]+)`)
	if regex.MatchString(line) == false {
		return false
	}
	ctcmItems := regex.FindAllStringSubmatch(line, -1)[0]
	mobile := ctcmItems[1]
	pNumber := strings.Split(ctcmItems[2], "/")
	for _, v := range pNumber {
		if v == "" {
			continue
		}
		if strings.HasPrefix(v, "P") == false {
			v = fmt.Sprintf("P%s", v)
		}
		p.dict[v].Mobile = mobile

	}
	return true
}

type Person struct {
	RPH          int      `json:"rph"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	IDType       string   `json:"idType"`
	IDNumber     string   `json:"idNumber"`
	IDIssue      string   `json:"idIssue"`
	Nationality  string   `json:"nationality"`
	Birthday     string   `json:"birthday"`
	Gender       string   `json:"gender"`
	Expired      string   `json:"expired"`
	Mobile       string   `json:"mobile"`
	TicketNumber []string `json:"tktNumber"`
}

func (p *Person) splitName(name string) string {

	name = strings.ToUpper(strings.TrimSpace(name))
	if strings.HasSuffix(name, "CHD") {
		//fmt.Println(name[:3])
		p.Name = strings.TrimSpace(name[:len(name)-3])
		p.Type = "CHD"
		return ""
	}

	regex := regexp.MustCompile(`^(.+)\s+(\w{6})$`)
	if regex.MatchString(name) {
		matches := regex.FindAllStringSubmatch(name, -1)[0]
		p.splitName(matches[1])
		return matches[2]
	}

	if strings.HasSuffix(name, "MR") {
		p.Name = strings.TrimSpace(strings.TrimRight(name, "MR"))
		p.Gender = "M"
	} else if strings.HasSuffix(name, "MS") {
		p.Name = strings.TrimSpace(strings.TrimRight(name, "MS"))
		p.Gender = "F"
	} else {
		p.Name = name
	}

	//p.Name = name
	p.Type = "ADU"
	return ""
}
