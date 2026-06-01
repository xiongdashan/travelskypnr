package travelskypnr

import (
	"fmt"
	"testing"
	"time"
)

// withFrozenNow 替换包级 nowFn，让 formatDate 的窗口推断在测试里可重复。
func withFrozenNow(t *testing.T, now time.Time) {
	t.Helper()
	orig := nowFn
	nowFn = func() time.Time { return now }
	t.Cleanup(func() { nowFn = orig })
}

func TestNewJourney(t *testing.T) {
	line := "UA505  C1  TH26DEC  SFOIAH HK3   2359 0549+1        E  3 C "
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	if j.Dep.StationName != "3" {
		t.Errorf("Dep.StationName is not 3")
	}
	if j.Arrival.StationName != "C" {
		t.Errorf("Arrival.StationName is not C")
	}
	fmt.Println(j)
}

func TestNewJourneyWithTerminal(t *testing.T) {
	line := "CZ3800 Z   FR18OCT  TAOCAN UN3   1705 2010          E --T2 S "
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	if j.Dep.StationName != "T2" {
		t.Errorf("Dep.StationName is not T2")
	}
	if j.Arrival.StationName != "" {
		t.Errorf("Arrival.StationName is not empty")
	}
}

func TestNewJourneyWithNoTerminal(t *testing.T) {
	line := "TK2126 Y1 TU10DEC ISTESB HK1 0800 0910 SEAME"
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	if j.FlightNumber != "TK2126" {
		t.Errorf("FlightNumber is not TK2126")
	}
}

func TestFormatDate(t *testing.T) {
	line := "AF111  B1  SA01FEB  PVGCDG HK1   2205 0550+1    SEAME  12E"
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	fmt.Println(j)
}

func TestFormatArrDepTime(t *testing.T) {
	line := "CA1230 L   MO10MAR25XIYPEK HK1   1930 2150          E T2T3"
	jl := &JourneyLine{}
	j := jl.newJourney(line)
	fmt.Println(j)
}

// TestFormatDate_YearInference 锁住"7 天窗口"语义。所有断言都依赖 nowFn 被冻结。
func TestFormatDate_YearInference(t *testing.T) {
	type tc struct {
		name     string
		now      time.Time
		input    string // fields[2] 部分
		wantDate string // formatted as 2006-01-02
	}
	cases := []tc{
		{
			name:     "未来 31 天保持当前年",
			now:      time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
			input:    "TU02JUL",
			wantDate: "2026-07-02",
		},
		{
			name:     "今天保持当前年",
			now:      time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
			input:    "MO01JUN",
			wantDate: "2026-06-01",
		},
		{
			name:     "过去 6 天仍在窗口内 -> 保持当前年",
			now:      time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
			input:    "MO08JUN",
			wantDate: "2026-06-08",
		},
		{
			name:     "过去 7 天恰在窗口边界 -> 保持当前年",
			now:      time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
			input:    "SU07JUN",
			wantDate: "2026-06-07",
		},
		{
			name:     "过去 8 天超过窗口 -> 滚到下一年",
			now:      time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
			input:    "SA06JUN",
			wantDate: "2027-06-06",
		},
		{
			name:     "跨月过去 30 天 -> 滚到下一年",
			now:      time.Date(2026, 6, 14, 12, 0, 0, 0, time.UTC),
			input:    "FR15MAY",
			wantDate: "2027-05-15",
		},
		{
			name:     "年末解析下个月（跨年） -> 滚到下一年",
			now:      time.Date(2026, 12, 30, 12, 0, 0, 0, time.UTC),
			input:    "FR15JAN",
			wantDate: "2027-01-15",
		},
		{
			name:     "年末解析下周（跨年但 8 天） -> 滚到下一年",
			now:      time.Date(2026, 12, 30, 12, 0, 0, 0, time.UTC),
			input:    "TH07JAN",
			wantDate: "2027-01-07",
		},
		{
			name:     "长格式带年份 25 -> 直接用，不猜",
			now:      time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
			input:    "MO10MAR25XIYPEK",
			wantDate: "2025-03-10",
		},
		{
			name:     "长格式带年份 27 -> 直接用，不猜",
			now:      time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
			input:    "TU01JUL27XIYPEK",
			wantDate: "2027-07-01",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withFrozenNow(t, c.now)
			j := &Journey{Dep: &ArrDep{}, Arrival: &ArrDep{}}
			got := j.formatDate(c.input)
			if got.Format("2006-01-02") != c.wantDate {
				t.Errorf("got %s, want %s", got.Format("2006-01-02"), c.wantDate)
			}
		})
	}
}

