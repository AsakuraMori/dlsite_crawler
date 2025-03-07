package dlsite_crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func GetInfoFromSearch(keyWord, typeName string, client *http.Client) ([]map[string]interface{}, error) {
	sType := ""
	switch typeName {
	case "游戏":
		sType = "pro"
	case "同人":
		sType = "maniax"
	case "漫画":
		sType = "book"
	case "手机游戏":
		sType = "appx"
	default:
		errMsg := errors.New("未匹配到类型")
		return nil, errMsg
	}
	apiUrl := fmt.Sprintf("https://www.dlsite.com/%s/fsr/=/language/jp/sex_category[0]/male/keyword/%s", sType, keyWord)
	return getInfoFromSearch(client, apiUrl)
}

func getInfoFromSearch(client *http.Client, urlStr string) ([]map[string]interface{}, error) {
	// 创建HTTP客户端并设置User-Agent
	doc, err := doRequest(client, urlStr)
	if err != nil {
		return nil, err
	}
	// 存储结果的切片
	var items []map[string]interface{}

	// 遍历搜索结果条目
	doc.Find("ul.n_worklist").Each(func(i int, s *goquery.Selection) {
		maker := []string{}
		pimg := []string{}
		s.Find(".maker_name").Each(func(k int, u *goquery.Selection) {
			maker = append(maker, u.Text())
		})
		s.Find(".lazy").Each(func(k int, u *goquery.Selection) {
			img, ex := u.Attr("src")
			if ex {
				pimg = append(pimg, img)
			}
		})
		s.Find(".multiline_truncate").Each(func(j int, t *goquery.Selection) {
			titleLink := t.Find(".multiline_truncate a")
			title := strings.TrimSpace(titleLink.Text())
			href, _ := titleLink.Attr("href")
			index := strings.LastIndex(href, "/")
			Id := ""
			if index != -1 {
				Id = href[index+1 : strings.LastIndex(href, ".html")]
			} else {
				fmt.Println("没有找到匹配的子字符串")
			}
			// 构造数据项
			item := map[string]interface{}{
				"title": title,
				"url":   href,
				"Id":    Id,
				"image": pimg[j],
				"maker": maker[j],
			}
			items = append(items, item)
		})
	})
	return items, nil
}
