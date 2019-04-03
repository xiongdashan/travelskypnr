package travelskypnr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PersonLine struct {
	PnrCode string
	Dict    map[string]*Person
	isPass  bool
	isSSR   bool
}

func NewPersonLine() *PersonLine {
	p := &PersonLine{}
	p.Dict = make(map[string]*Person)
	return p
}

func (p *PersonLine) Data() (rev []*Person) {
	for _, v := range p.Dict {
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

// XN/IN/Com/ILILLYROSE (AUG16)/P1
const patternINF = `XN\/IN\/(.*[^\/]+)\/(.*[^\()]+)\(([A-Z0-9]{5})\)\/P(\d+)`

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

		// 匹配婴儿...
		if match, _ := regexp.MatchString(patternINF, line); match {
			p.setInft(line)
			return true
		}

		return false
	}

	person := &Person{}
	person.RPH = pos
	pnrNum := person.splitName(line)
	key := fmt.Sprintf("P%d", pos)
	p.Dict[key] = person
	if pnrNum != "" {
		p.PnrCode = pnrNum
	}

	fmt.Println(person.Name)

	return true

}

// 婴儿信息
func (p *PersonLine) setInft(line string) {
	regex := regexp.MustCompile(patternINF)
	match := regex.FindAllStringSubmatch(line, -1)[0]
	person := &Person{
		Name:     fmt.Sprintf("%s/%s", match[1], match[2]),
		Birthday: match[3],
		Type:     Infant,
	}
	person.RPH, _ = strconv.Atoi(match[4])

	key := fmt.Sprintf("P%sINF", match[4])
	p.Dict[key] = person
}

// 设置票号
func (p *PersonLine) SetTktNumber(rph int, num string, tType string) {
	for _, v := range p.Dict {
		if v.RPH == rph {
			// 如果是婴儿，找相同类型的乘客
			if tType == Infant && v.Type != Infant {
				continue
			}
			v.TicketNumber = append(v.TicketNumber, num)
		}
	}
}

// 统计人数类型
func (p *PersonLine) TypeCount(ty string) (rev int) {
	for _, v := range p.Dict {
		if v.Type == ty {
			rev++
		}
	}
	return
}

// ..AddSSR 证件信息
func (p *PersonLine) AddSSR(line string) bool {

	if strings.HasPrefix(line, "SSR DOCS") {
		p.ssr(line)
		return true
	}
	if strings.HasPrefix(line, "SSR FOID") {
		p.foid(line)
		return true
	}
	return false
}

// foid SSR FOID CA HK1 NI220182198906185118/P1
func (p *PersonLine) foid(line string) {
	aryItem := strings.Split(line, "/")
	if len(aryItem) < 2 {
		return
	}
	idAry := strings.Fields(aryItem[0])
	idInfo := idAry[len(idAry)-1]
	key := strings.TrimSpace(aryItem[len(aryItem)-1])
	psn, ok := p.Dict[key]
	if !ok {
		fmt.Print("无乘客信e")
		return
	}
	psn.IDType = "NI"
	psn.IDNumber = idInfo[2:]
}

// 证件类型/发证国家/证件号码/国籍/出生日期/性别/证件有效期限/SURNAME(姓)/FIRST-NAME(名)/MID-NAME(中间名)/持有人标识H/P1
// SSR DOC AM HK1 0P/      1 CN/   2 E30028197/3 CN/4 24AUG79/5 F/ 6   12SEP23/   LU/       FANGFANG/ P4
func (p *PersonLine) ssr(line string) {
	aryItem := strings.Split(line, "/")
	if len(aryItem) < 5 {
		return
	}

	//证件类型 -- SSR DOC AM HK1 P
	idAry := strings.Fields(aryItem[0])

	//idInfostr := aryItem[4]
	//idItem := strings.Split(idInfostr, "/")
	key := strings.TrimSpace(aryItem[len(aryItem)-1])

	psn, ok := p.Dict[key]
	if !ok {
		fmt.Printf("无乘客信息--%s\n", key)
		return
	}
	psn.IDType = idAry[len(idAry)-1] //idItem[0]

	psn.IDIssue = aryItem[1]
	psn.IDNumber = aryItem[2]
	psn.Nationality = aryItem[3]
	psn.Birthday = aryItem[4]
	psn.Gender = aryItem[5]
	psn.Expired = aryItem[6]
	fmt.Println(p.Dict[key].IDNumber)
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
		p.Dict[v].Mobile = mobile

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
	if strings.HasSuffix(name, Child) {
		//fmt.Println(name[:3])
		p.Name = strings.TrimSpace(name[:len(name)-3])
		p.Type = Child
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
		p.Gender = "M"
	} else if strings.HasSuffix(name, "MS") {
		p.Name = strings.TrimSpace(strings.TrimRight(name, "MS"))
		p.Gender = "F"
	} else {
		p.Name = name
	}

	//p.Name = name
	p.Type = Adult
	return ""
}
