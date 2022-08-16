package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var filename string = "log_new.csv"

type Select struct {
	FMID string `json:"FMID"`
	FPID string `json:"FPID"`
	FOID string `json:"FOID"`
}
type SQLBuilderItem struct {
	SQLBuilderID string `json:"SQLBuilderID"`
	Select       Select `json:"Select"`
}
type info struct {
	Context struct {
	} `json:"Context"`
	SQLBuilderItem []SQLBuilderItem `json:"SQLBuilderItem"`
}
type result struct {
	JSONResult  string `json:"JsonResult"`
	JSONMessage struct {
		MessageIndex string `json:"MessageIndex"`
		Remark       string `json:"Remark"`
		MessageInfo  string `json:"MessageInfo"`
	} `json:"JsonMessage"`
	JSONData []struct {
		SQLBuilderID string `json:"SQLBuilderID"`
		FIELD        []struct {
			Attrname  string `json:"attrname"`
			Fieldtype string `json:"fieldtype"`
			WIDTH     string `json:"WIDTH,omitempty"`
		} `json:"FIELD"`
		ROW []struct {
			FKINDNAME  string `json:"FKIND_NAME"`
			FPRICEBASE string `json:"FPRICE_BASE"`
			FNEWTIME   string `json:"FNEWTIME"`
			FTOPREMARK string `json:"FTOP_REMARK"`
			FREMARK    string `json:"FREMARK"`
		} `json:"ROW"`
	} `json:"JsonData"`
}

type godinfo struct {
	Name   string
	Jiage  string
	Timess string
}

func main() {
	var godinfo godinfo
	Getinfo(&godinfo)
	Saveinfo(&godinfo)

}

//获取金价
func Getinfo(godinfo *godinfo) {

	var info info
	var Select Select
	Select.FMID = "{005A5001-B9AD-41CB-8409-8F7675D19143}"
	Select.FMID = "{A97855D9-C5D1-B8E1-7B86-30D964BC9666}"
	Select.FPID = "{A97855D9-C5D1-B8E1-7B86-30D964BC9666}"
	Select.FOID = "{7D77D027-9824-4156-A25E-12FC59527DDE}"
	var SQLBuilderItem SQLBuilderItem
	SQLBuilderItem.Select = Select
	SQLBuilderItem.SQLBuilderID = "{005A5001-B9AD-41CB-8409-8F7675D19143}"
	info.SQLBuilderItem = append(info.SQLBuilderItem, SQLBuilderItem)

	b, err := json.Marshal(info)
	if err != nil {
		fmt.Println("转换出错", err)
	}
	resp, err := http.Post("http://111.198.86.222/BAP/OpenApi", "application/json",
		bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("请求出错", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result result
	json.Unmarshal([]byte(body), &result)
	//fmt.Println(result.JSONData[0].ROW)
	for _, v := range result.JSONData[0].ROW {
		if v.FKINDNAME == "菜百投资基础金价" {
			//fmt.Println(v.FKINDNAME, v.FPRICEBASE, v.FNEWTIME)
			godinfo.Name = v.FKINDNAME
			godinfo.Jiage = v.FPRICEBASE
			godinfo.Timess = v.FNEWTIME
		}
	}
}

func Saveinfo(godinfo *godinfo) {

	t, err := os.Open(filename)
	defer t.Close()
	if err != nil && os.IsNotExist(err) {
		flname, err1 := os.Create(filename)
		defer flname.Close()
		if err1 != nil {
			log.Print("创建文件出错,信息:", err1)
		}
		t.WriteString(("\xEF\xBB\xBF"))
		wr := csv.NewWriter(flname)
		wr.Write([]string{"时间", "价格(元)"})
		wr.Flush()

	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("打开已有文件出错,信息:", err)
	}
	writer := csv.NewWriter(file)
	filetxt := []string{godinfo.Timess, strings.Split(godinfo.Jiage, " ")[0]}
	writer.Write(filetxt)
	writer.Flush()
}
