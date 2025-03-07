package dlsite_crawler

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

func doRequest(client *http.Client, urlStr string) (*goquery.Document, error) {
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	//req.Header.Set("Accept-Language", "ja,en-US;q=0.9")

	var resp *http.Response
	maxRetries := 3
	var err error
	for i := 0; i < maxRetries; i++ {
		resp, err = client.Do(req) // 注意这里使用=而不是:=
		if err == nil && resp != nil && resp.StatusCode == 200 {
			break
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	// 错误处理优先
	if err != nil {
		errMsg := errors.New("请求失败: " + err.Error())
		return nil, errMsg
	}
	if resp == nil {
		errMsg := errors.New("所有重试均失败")
		return nil, errMsg
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != 200 {
		errMsg := errors.New(fmt.Sprintf("状态码异常: %d", resp.StatusCode))
		return nil, errMsg
	}

	// 安全读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMsg := errors.New(fmt.Sprintf("读取响应失败: %v", err))
		return nil, errMsg
	}

	// 安全解析HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		errMsg := errors.New(fmt.Sprintf("HTML解析失败: %v", err))
		return nil, errMsg
	}

	return doc, nil

}

func fmtText(text string) (string, error) {
	ret := ""
	scanner := bufio.NewScanner(strings.NewReader(text))

	// 设置分隔符为换行符，这样Scanner会在每行结束时停止
	scanner.Split(bufio.ScanLines)

	// 逐行读取并打印
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的内容
		//fmt.Println(line)      // 打印当前行的内容
		line = strings.Replace(line, "  ", "", -1)
		if len(line) > 0 {
			//fmt.Println(line)
			ret += line + "\n"
		}
	}

	// 检查是否发生错误
	if err := scanner.Err(); err != nil {
		errMsg := errors.New("读取时发生错误:" + err.Error())
		return "", errMsg
	}
	return ret, nil
}
