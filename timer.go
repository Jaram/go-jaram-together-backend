package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func request(msgType string, podName string, group string, uid string) {
	resp, err := http.PostForm("http://localhost:8080/msgRequest", url.Values{"msgType": {msgType}, "podName": {podName}, "group": {group}, "uid":{uid}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}
}

func scheduler(startHour string, startMinute string, endHour string, endMinute string, podName string, group string, uid string) {
	year, month, day := time.Now().Date()
	startTimeStr := fmt.Sprint(year, "-", int(month), "-", day, " ", startHour, ":", startMinute)
	endTimeStr := fmt.Sprint(year, "-", int(month), "-", day, " ", endHour, ":", endMinute)
	loc, _ := time.LoadLocation("Asia/Seoul")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04", startTimeStr, loc)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04", endTimeStr, loc)
	startDuration := startTime.Sub(time.Now())
	endDuration := endTime.Sub(time.Now())

	time.AfterFunc(startDuration-10*time.Minute, func() {
		request("readyReminder", podName, group, uid)
	})
	time.AfterFunc(startDuration, func() {
		request("startPod", podName, group, uid)
	})
	time.AfterFunc(endDuration, func() {
		request("endPod", podName, group, uid)
	})
	time.AfterFunc(endDuration+10*time.Minute, func() {
		request("deletePod", podName, group, uid)
	})
}
