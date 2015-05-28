package qlog

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slene/iploc"
	"io"
	"io/ioutil"
	"net/http"
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

type TaobaoIpRet struct {
	Code     int      `json:"code"`
	TaobaoIp TaobaoIp `json:"data"`
}
type TaobaoIp struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Isp       string `json:"isp"`
	IspId     string `json:"isp_id"`
	Ip        string `json:"ip"`
}

const INVALID_TAOBAO_IP = `{"code":1,"data":"invaild ip."}`

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

func sha1Hash(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//template functions
func Empty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func NotEmpty(str string) bool {
	return strings.TrimSpace(str) != ""
}

func UrlFor(path string) string {
	serverRoot := "http://" + GlbConf.ListenHost
	if GlbConf.ListenPort != 80 {
		serverRoot += fmt.Sprintf(":%d", GlbConf.ListenPort)
	}
	return fmt.Sprintf("%s%s", serverRoot, path)
}

func GetIpInfo(ip string) (code, country, region, city, isp, note string, err error) {
	ipInfo, ipErr := iploc.GetIpInfo(ip)
	if ipErr != nil {
		err = errors.New("get ip location info failed due to," + ipErr.Error())
		return
	}
	switch ipInfo.Flag {
	case iploc.FLAG_INUSE:
		if ipInfo.Code == "CN" {
			code = ipInfo.Code
			country = ipInfo.Country
			region = ipInfo.Region
			city = ipInfo.City
			isp = ipInfo.Isp
		} else {
			code = ipInfo.Code
			country = ipInfo.Country
		}
	case iploc.FLAG_RESERVED:
		note = ipInfo.Note
	case iploc.FLAG_NOTUSE:
		note = ipInfo.Note
	}
	return
}

func GetTaobaoIpInfo(ip string) (code, country, region, city, isp, note string, err error) {
	url := fmt.Sprintf("http://ip.taobao.com/service/getIpInfo.php?ip=%s", ip)
	resp, respErr := http.Get(url)
	if respErr != nil {
		err = errors.New(fmt.Sprintf("get taobao ip error, %s", respErr))
		return
	}
	defer resp.Body.Close()
	respData, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		err = errors.New(fmt.Sprintf("read taobao ip data error, %s", readErr))
		return
	}
	if string(respData) == INVALID_TAOBAO_IP {
		err = errors.New(fmt.Sprintf("invalid ip %s", ip))
		return
	}
	taobaoIpRet := TaobaoIpRet{}
	pErr := json.Unmarshal(respData, &taobaoIpRet)
	if pErr != nil {
		err = errors.New(fmt.Sprintf("parse taobao ip data error, %s", pErr))
		return
	}
	code = taobaoIpRet.TaobaoIp.CountryId
	country = taobaoIpRet.TaobaoIp.Country
	region = taobaoIpRet.TaobaoIp.Region
	city = taobaoIpRet.TaobaoIp.City
	isp = taobaoIpRet.TaobaoIp.Isp
	return
}
