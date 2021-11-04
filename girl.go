package main

/**
 * auther: caiqm
 * date: 2021-11-03
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

var (
	// 请求接口
	picUrl = "https://pic.sogou.com/napi/pc/searchList?mode=1&start=%d&xml_len=%d&query=%s"
	// 定义协程
	waitGroup = new(sync.WaitGroup)
)

/**
 * @param start 表示从第几张图片开始检索
 * @param len 表示从第几张往后获取48张图片
 * @param query 搜索关键词
 */
func requestPic(start, len int, query string) {
	// 中文urlencode转义
	query = url.QueryEscape(query)
	requestUrl := fmt.Sprintf(picUrl, start, len, query)
	// 请求接口
	rsp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	// 定义返回格式
	var girlJson map[string]interface{}
	json.Unmarshal([]byte(body), &girlJson)
	// 格式转换
	dataJson := girlJson["data"].(map[string]interface{})
	items := dataJson["items"].([]interface{})
	// 读取信息
	waitGroup.Add(len)
	for _, val := range items {
		item := val.(map[string]interface{})
		// 启动线程
		go func(pic map[string]interface{}) {
			// 图片路径
			picUrl := pic["picUrl"].(string)
			// 图片名称，沿用原本名称
			picName := pic["name"].(string)
			// 下载图片
			downloadPic(picUrl, picName)
			defer waitGroup.Done()
		}(item)
	}
	// 等待所有协程操作完成
	waitGroup.Wait()
}

// 判断目录是否存在
func fileExist(path, fileName string) string {
	// 路径拼接
	filePath := filepath.Join(path, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		os.Create(filePath)
	}
	return filePath
}

// 下载图片
func downloadPic(pic, name string) {
	// 请求图片链接
	rsp, err := http.Get(pic)
	if err != nil {
		fmt.Println("request fail, file name is : ", name)
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	// 当前链接
	path, _ := os.Getwd()
	// 判断文件夹是否存在
	filePath := fileExist(path, "pic")
	// 拼接文件
	fileName := filepath.Join(filePath, name)
	// 下载文件
	err = ioutil.WriteFile(fileName, []byte(body), 0777)
	if err != nil {
		fmt.Println("download fail, file name is : ", name)
		fmt.Println(err)
		return
	}
	// 输出结果
	fmt.Println("download success, file name is : ", name)
}

func main() {
	requestPic(1, 48, "黑丝,美女")
}
