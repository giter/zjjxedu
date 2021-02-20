package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Login struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Ticket string `json:"ticket"`
	} `json:"data"`
}

var r = resty.New()

func punch(ids ...string) {

	for _, id := range ids {

		rlogin := Login{}

		fmt.Println("自动登录...")
		rep, err := r.NewRequest().SetHeader("User-Agent", "User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat").
			SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=utf-8").
			SetHeader("Accept-Encoding", "gzip, deflate, br").
			SetQueryParam("loginname", id).
			SetQueryParam("password", "hy123456").
			Post("https://jk.zjjxedu.gov.cn/sso/mobi/WxLogin2")

		if rep.StatusCode() != 200 {

			fmt.Println("登录失败:", rep.Status())
			continue
		}

		err = json.Unmarshal([]byte(rep.Body()), &rlogin)

		if err != nil {
			fmt.Println("登录失败:", err.Error())
			continue
		}

		if rlogin.Code != 0 {

			fmt.Println("登录失败:", rlogin.Msg)
			continue
		}

		fmt.Println("自动签到...")

		rep, err = r.NewRequest().
			SetHeader("User-Agent", "User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat").
			SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=utf-8").
			SetHeader("ticket", rlogin.Data.Ticket).
			SetHeader("Referer", "https://servicewechat.com/wxbee484c9ef018d43/23/page-frame.html").
			SetHeader("Accept-Encoding", "gzip, deflate, br").
			SetBody("temperature=36&is_whether=1&amorpm=0&symptom=&remark=&is_famwhether=1&famremark=").
			Post("https://jk.zjjxedu.gov.cn/health/mobiapi/savePunchclock")

		if err != nil {
			fmt.Println("签到失败", err)
			continue
		}

		fmt.Println(id, "签到成功...")
	}
}

func main() {
	punch(os.Args[1:]...)
}
