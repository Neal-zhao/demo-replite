package main

import (
	"demo/database"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var url string = "https://v.qq.com/channel/cartoon"
var db database.DB

func main() {
	db.InitSqlx()
	//抓取网页内容
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get err：", err)
		return
	}

	//re := regexp.MustCompile(`<div class="video-banner-item"(.*?)><div class="video-banner-item`)
	//video-banner-item	video-card-wrap
	re := regexp.MustCompile(`<div class="video-banner-item"(.*?)></div></div></div>`)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(" err io.ReadAll：", err)
		return
	}

	divs := re.FindAllStringSubmatch(string(body), -1)
	for _, div := range divs {
		src := getImgSrc2(div[0])
		title, link := getLincInfo(div[0])
		if !strings.Contains(src, "http") {
			src = fmt.Sprintf("http:%s", src)
		}
		fmt.Println(src, title, link)

		//filename := GetFilenameFromUrl2(src, title)
		//DownloadFile(url, filename)
		xxxDown(url)

		r, err := database.SqlxDB.Exec("INSERT INTO `cartoon`.`hot`(`name`, `link`, `img`) VALUES (?, ?, ?) ", title, link, src)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(r.LastInsertId())
	}
}

// 下载图片
func DownloadImg() {
	//for url := range chanImageUrls {
	//	filename := GetFilenameFromUrl(url)
	//	ok := DownloadFile(url, filename)
	//	if ok {
	//		fmt.Printf("%s 下载成功\n", filename)
	//	} else {
	//		fmt.Printf("%s 下载失败\n", filename)
	//	}
	//}
}

// 下载图片，传入的是图片叫什么
func DownloadFile(url string, filename string) (ok bool) {
	resp, err := http.Get(url)
	HandleError(err, "http.get.url")
	defer resp.Body.Close()

	fmt.Println(resp.Header)
	//contentType := resp.Header.Get("Content-Type")
	//fileExt, err := getFileExt(contentType)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//filename = filename + fileExt
	bytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "resp.body")
	//filename = "./img/" + filename
	// 写出数据
	err = ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		return false
	} else {
		return true
	}
}
func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}
func getFileExt(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	default:
		return ".jpg", fmt.Errorf("unknown content type")
	}
}

// 截取url名字
func GetFilenameFromUrl2(url, name string) (filename string) {

	// 时间戳解决重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + name
	return
}

// 截取url名字
func GetFilenameFromUrl(url string) (filename string) {
	// 返回最后一个/的位置
	lastIndex := strings.LastIndex(url, "/")
	// 切出来
	filename = url[lastIndex+1:]
	// 时间戳解决重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + filename
	return
}
func getImgSrc(elem string) (src [][]string) {
	imgRe := regexp.MustCompile(`<img.*data-src="(.*?)"\s+class="banner-cover".*>`)
	matchs := imgRe.FindAllStringSubmatch(elem, -1)
	fmt.Println(len(matchs), matchs[0][1])
	for _, match := range matchs {
		src = append(src, []string{match[1]})
	}
	return
}
func getImgSrc2(elem string) (src string) {
	imgRe := regexp.MustCompile(`<img.*data-src="(.*?)"\s+class="banner-cover".*>`)
	matchs := imgRe.FindAllStringSubmatch(elem, -1)
	fmt.Println(len(matchs), matchs[0][1])
	for _, match := range matchs {
		src = match[1]
	}
	return
}

func getLincInfo(elem string) (title, link string) {
	//return title link
	//fmt.Println(elem)
	imgRe := regexp.MustCompile(`<a\s+href="(.*?)".*class="banner-title".*>(.*?)</a>`)
	matchs := imgRe.FindAllStringSubmatch(elem, -1)
	//fmt.Println(len(matchs), matchs[0][1], matchs[0][2])
	for _, match := range matchs {
		link = match[1]
		title = match[2]
	}
	return
}

func downImg(url string) {

}

func xxxDown(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 创建要保存文件的目录
	dir := "./img"
	_ = os.MkdirAll(dir, 0755)

	// 获取文件名
	fileName := filepath.Base(url)
	if !strings.Contains(fileName, ".") {
		//contentType := resp.Header.Get("Content-Type")
		//fileExt, err := getFileExt(contentType)
		//if err != nil {
		//	//fileExt = ".jpg"
		//}
		//fileName += fileExt
	}

	// 创建文件
	file, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 将响应体数据写入文件
	if _, err := io.Copy(file, resp.Body); err != nil {
		panic(err)
	}

	fmt.Printf("文件已保存：'%s'\n", fileName)
}
