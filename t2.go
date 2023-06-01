package main

import (
	"demo/database"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var url1 string = "https://v.qq.com/channel/cartoon"
var db1 database.DB
var wg sync.WaitGroup

func main() {
	db1.InitSqlx()
	//抓取网页内容
	resp, err := http.Get(url1)
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

	for i := 0; i < 15; i++ {
		go func() {
			for img := range downImgChan {
				//time.Sleep(time.Second * 1)
				fmt.Println(img)

				src := img[0]
				title := img[1]
				dImg2(src, title)
				wg.Done()
			}

			fmt.Println("go down 结束", time.Now())
		}()
	}

	divs := re.FindAllStringSubmatch(string(body), -1)
	for _, div := range divs {
		src := getImgSrc3(div[0])
		title, link := getLincInfo1(div[0])
		if !strings.Contains(src, "http") {
			src = fmt.Sprintf("http:%s", src)
		}
		//fmt.Println(src, title, link)
		//dImg2(src, title)

		downImgChan <- []string{src, title}
		wg.Add(1)
		r, err := database.SqlxDB.Exec("INSERT INTO `cartoon`.`hot`(`name`, `link`, `img`) VALUES (?, ?, ?) ", title, link, src)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(r.LastInsertId())
		//break
	}
	wg.Wait()
	close(downImgChan)
}

var downImgChan = make(chan []string, 100)

func getImgSrc3(elem string) (src string) {
	imgRe := regexp.MustCompile(`<img.*data-src="(.*?)"\s+class="banner-cover".*>`)
	matchs := imgRe.FindAllStringSubmatch(elem, -1)
	for _, match := range matchs {
		src = match[1]
	}
	return
}
func getLincInfo1(elem string) (title, link string) {
	imgRe := regexp.MustCompile(`<a\s+href="(.*?)".*class="banner-title".*>(.*?)</a>`)
	matchs := imgRe.FindAllStringSubmatch(elem, -1)
	for _, match := range matchs {
		link = match[1]
		title = match[2]
	}
	return
}
func dImg2(src, title string) {
	//请求
	resp := request(src)
	//新建目录
	//title = filepath.Clean(title)
	title = strings.ReplaceAll(title, "?", "")
	dir := getSaveDir("")
	//文件名 扩展
	filename := getFileName(dir, title, resp.Header.Get("Content-Type"))

	// 创建本地文件
	saveFile(filename, resp.Body)
}
func request(src string) (resp *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", src, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	return
}
func saveFile(filename string, Body io.Reader) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 将响应Body中的数据写入本地文件
	size, err := io.Copy(file, Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloaded '%s' with '%v' bytes\n", filename, size)
}
func getFileName(dir, title, ContentType string) (filename string) {
	filename = fmt.Sprintf("%s_%s", strconv.FormatInt(time.Now().UnixNano(), 10), title)
	if !strings.Contains(filename, ".") {
		ext, err := getFileExt2(ContentType)
		if err != nil {
			fmt.Println(err)
		}
		filename = fmt.Sprintf("%s%s", filename, ext)
	}
	// 创建本地文件
	//fmt.Println(url, filename, dir)
	filename = filepath.Join(dir, filename)
	return
}
func getSaveDir(custom string) (dir string) {
	dir = "./img"
	dir += fmt.Sprintf("/%s/", time.Now().Format("20060102"))
	os.MkdirAll(dir, 0777)
	return
}
func getFileExt2(contentType string) (string, error) {
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
