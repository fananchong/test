package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/alertmanager/template"
)

// 使用 webhook ，自定义扩展报警功能的例子

func alert(w http.ResponseWriter, r *http.Request) {
	// 这里演示， alertmanager 把警告投递给本进程后，本进程丢给企业微信
	data := &template.Data{}
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(jsonData, data); err != nil {
		panic(err)
	}
	fmt.Println(data)
	curl("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=18539bfc-6076-4871-bd05-0a82e9ffc7c2", data.Alerts[0].Annotations["description"])
}

func main() {
	http.HandleFunc("/alert", alert)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

type wechatTxt struct {
	Content string `json:"content"`
}
type wechatMsg struct {
	MsgType string     `json:"msgtype"`
	Text    *wechatTxt `json:"markdown"`
}

func curl(url string, txt string) {
	msg := &wechatMsg{
		MsgType: "markdown",
		Text: &wechatTxt{
			Content: txt,
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("resp.Status:", resp.Status)
	fmt.Println("resp.Headers:", resp.Header)
}
