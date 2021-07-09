package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	getinfo(&godinfo)
	saveinfo()
	fmt.Printf("%v;%v;%v", godinfo.Name, godinfo.Jiage, godinfo.Timess)

}
func getinfo(godinfo *godinfo) {

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

func saveinfo() {}
