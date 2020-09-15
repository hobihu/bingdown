package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	api  = "https://cn.bing.com/HPImageArchive.aspx?format=js&n=1&idx=0"
	host = "https://cn.bing.com"
)

type Bing struct {
	Images []struct {
		Enddate   string `json:"enddate"`
		Urlbase   string `json:"urlbase"`
		Copyright string `json:"copyright"`
	} `json:"images"`
}

var bing Bing
var imgURL, imgFileName, imgTXTFileName string
var dir string = "/home/bingdown/" //文件保存路径

func main() {
	getImgURL()
	getTXT()
}

func getImgURL() {
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println("图片地址获取失败！")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("图片地址获取失败！")
	}
	err = json.Unmarshal(body, &bing)
	if err != nil {
		fmt.Println("图片地址解析失败！")
	}
	imgType := [...]string{"_UHD.jpg", "_1920x1080.jpg", "_1280x720.jpg"}
	for i := 0; i < len(imgType); i++ {
		imgURL = host + bing.Images[0].Urlbase + imgType[i]
		imgFileName = bing.Images[0].Enddate + imgType[i]
		fmt.Println("下载文件：", imgFileName, "\n地址：", imgURL)
		getImg() //获取UHD图片
	}
}

func getImg() {
	resp, err := http.Get(imgURL)
	if err != nil {
		fmt.Println("图片获取失败！")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("图片获取失败！")
	}
	f, err := os.Create(dir + imgFileName)
	defer f.Close()
	if err != nil {
		fmt.Println("图片保存失败！")
	}
	size, err := f.Write(body)
	if err != nil {
		fmt.Println("图片保存失败！")
	}
	fmt.Printf("图片保存成功：%v bytes。\n", size)
}

func getTXT() {
	//获取Copyright
	binghome := dir + bing.Images[0].Enddate + ".txt"
	fmt.Println("下载Copyright文件：", binghome)
	err := ioutil.WriteFile(binghome, []byte(bing.Images[0].Copyright), 0755)
	if err != nil {
		fmt.Printf("ioutil.WriteFile failure, err=[%v]\n", err)
	}
	fmt.Println("下载完成")
}
