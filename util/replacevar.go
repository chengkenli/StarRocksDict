package util

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewReplaceVar(s string) string {
	old := s
	re := regexp.MustCompile(`{(yyyy)?(-?mm)?(-?dd)?(,-?\d+)?}`)
	//fmt.Printf("%+v", re.FindAllString(old, -1))
	allMathString := re.FindAllString(old, -1)
	//[{yyyy-mm,1} {yyyymmdd,-2} {yyyymm,3} {yyyy-mm-dd} {dd,-6} {mm,1} {mm} {yyyy,1}]
	for _, m := range allMathString {
		var (
			goDateFmt string
			dateFmt   string
			interval  string
			step      int
		)
		goDateFmt = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(m, "yyyy", "2006"), "mm", "01"), "dd", "02")
		goDateFmt = strings.ReplaceAll(goDateFmt, "{", "")
		goDateFmt = strings.ReplaceAll(goDateFmt, "}", "")
		lst := strings.Split(goDateFmt, ",")

		if len(lst) > 1 {
			dateFmt, interval = lst[0], lst[1]
			step, _ = strconv.Atoi(interval)
		} else {
			dateFmt, step = lst[0], 0
		}
		//fmt.Println(lst, dateFmt)
		dateFmtM := strings.ReplaceAll(dateFmt, "-", "")

		switch dateFmtM {
		case "20060102", "02":
			old = strings.ReplaceAll(old, m, time.Now().AddDate(0, 0, step).Format(dateFmt))
		case "200601", "01":
			old = strings.ReplaceAll(old, m, time.Now().AddDate(0, step, 0).Format(dateFmt))
		case "2006":
			old = strings.ReplaceAll(old, m, time.Now().AddDate(step, 0, 0).Format(dateFmt))
		}
	}
	return old
}
