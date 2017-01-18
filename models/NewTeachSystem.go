package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	simplejson "github.com/bitly/go-simplejson"
)

type cxScore struct {
	Semestre   string `json:"semestre"`
	CXType     string `json:"cxType"`
	Name       string `json:"name"`
	Score      string `json:"score"`
	InsertTime string `json:"insertTime"`
}

func login(userName string, password string) string {
	LoginURL := beego.AppConfig.String("SYSTEM_LOGIN")

	v := url.Values{"j_username": {userName + ",undergraduate"}, "j_password": {password}}
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, LoginURL, body)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	// data, _ := ioutil.ReadAll(resp.Body)

	res := fmt.Sprintf("%s", resp.Request.URL)
	tagCookies := strings.Split(strings.Split(res, ";")[1], "=")[1]

	return tagCookies
}

func GetCourseFromLogin(userName string, password string, semestre string) map[string]interface{} {
	cookie := login(userName, password)

	return getTrueCourse(cookie, userName, semestre)
}

func getTrueCourse(thatCookie string, userName string, semestre string) map[string]interface{} {
	getTrueCourseURL := beego.AppConfig.String("COURSE_TABLE")

	v := url.Values{"listXnxq": {semestre}, "uid": {userName}}
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, getTrueCourseURL, body)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "JSESSIONID="+thatCookie)

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	json, err := simplejson.NewJson(data)
	var nodes = make(map[string]interface{})

	nodes, _ = json.Map()

	return nodes
}

func GetCXScoreFromLogin(userName string, password string) map[string]interface{} {
	cookie := login(userName, password)

	return getTrueCXScore(cookie, userName)
}

func getTrueCXScore(thatCookie string, userName string) map[string]interface{} {
	getTrueCourseURL := beego.AppConfig.String("INNVOATION_SCORE")

	v := url.Values{"uid": {userName}}
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, getTrueCourseURL, body)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "JSESSIONID="+thatCookie)

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	str := strings.NewReader(string(data))

	doc, _ := goquery.NewDocumentFromReader(str)

	var finalCXScore []cxScore

	doc.Find(".gridtable tbody tr").Each(func(i int, a *goquery.Selection) {
		a.Find("td").Each(func(j int, m *goquery.Selection) {
			finalCXScore[i].Name = "s"
			finalCXScore[i].CXType = "s"
			finalCXScore[i].InsertTime = "s"
			finalCXScore[i].Score = "s"
			finalCXScore[i].Semestre = "s"
			fmt.Println(m.Text())
		})
	})
	a, err := json.Marshal(finalCXScore)

	jsons, err := simplejson.NewJson(a)
	// var nodes = make(map[string]interface{})

	nodes, err := jsons.Map()

	return nodes
}
