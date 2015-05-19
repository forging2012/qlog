package qlog

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var shortMonthNames = map[string]time.Month{
	"Jan": time.January,
	"Feb": time.February,
	"Mar": time.March,
	"Apr": time.April,
	"May": time.May,
	"Jun": time.June,
	"Jul": time.July,
	"Aug": time.August,
	"Sep": time.September,
	"Oct": time.October,
	"Nov": time.November,
	"Dec": time.December,
}

func ParseDateTime(str string) (t time.Time, err error) {
	pattern := `^\d{2}/(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)/\d{4}:\d{2}:\d{2}:\d{2}\s\+0800$`
	match, mErr := regexp.MatchString(pattern, str)
	if mErr != nil {
		err = mErr
		return
	}
	if !match {
		err = errors.New("invalid log date time")
		return
	}
	loc, locErr := time.LoadLocation("Asia/Shanghai")
	if locErr != nil {
		err = locErr
		return
	}
	items := strings.Split(str, " ")
	dtime := items[0]
	dtimeItems := strings.SplitN(dtime, ":", 2)
	datePart := dtimeItems[0]
	timePart := dtimeItems[1]
	dateItems := strings.Split(datePart, "/")
	day, _ := strconv.Atoi(dateItems[0])
	month := shortMonthNames[dateItems[1]]
	year, _ := strconv.Atoi(dateItems[2])
	timeItems := strings.Split(timePart, ":")
	hour, _ := strconv.Atoi(timeItems[0])
	minute, _ := strconv.Atoi(timeItems[1])
	second, _ := strconv.Atoi(timeItems[2])
	t = time.Date(year, month, day, hour, minute, second, 0, loc)
	return
}

func Trim(str string, prefix string, suffix string) string {
	tstr := strings.TrimPrefix(str, prefix)
	tstr = strings.TrimSuffix(tstr, suffix)
	return tstr
}
