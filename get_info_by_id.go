package dlsite_crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)


func GetInfoByID(id, typeName string, client *http.Client) (map[string]interface{}, error) {
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
	apiUrl := fmt.Sprintf("https://www.dlsite.com/%s/work/=/product_id/%s", sType, id)
	return getInfoByID(client, apiUrl)
}

func getInfoByID(client *http.Client, urlStr string) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"success": false,
		"error":   "",
		"data":    nil,
	}
	doc, err := doRequest(client, urlStr)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	extractBasicInfo(doc, data, urlStr)
	extractDescription(doc, data)
	extractStaff(doc, data)

	result["success"] = true
	result["data"] = data
	return result, nil
}

func extractBasicInfo(doc *goquery.Document, data map[string]interface{}, urlStr string) {
	// 产品ID
	if u, err := url.Parse(urlStr); err == nil {
		parts := strings.Split(u.Path, "/")
		if len(parts) > 0 {
			data["product_id"] = strings.TrimSuffix(parts[len(parts)-1], ".html")
		}
	}

	// 标题
	doc.Find("h1#work_name").Each(func(i int, s *goquery.Selection) {
		data["title"] = strings.TrimSpace(s.Text())
	})

	// 制作商
	doc.Find(".maker_name").Each(func(i int, s *goquery.Selection) {
		data["circle"] = strings.TrimSpace(s.Text())
	})
	// img
	doc.Find(".product-slider-data").Each(func(i int, s *goquery.Selection) {
		//data["image"] = strings.TrimSpace(s.Text())
		img, ex := s.Find("div").Attr("data-src")
		if ex {
			data["image"] = img
		}
	})
}

func extractDescription(doc *goquery.Document, data map[string]interface{}) {
	details := ""
	var err error
	doc.Find(".work_parts_container").Each(func(i int, s *goquery.Selection) {
		details, err = fmtText(s.Text())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	})

	data["description"] = details
}

func extractStaff(doc *goquery.Document, data map[string]interface{}) {
	staff := ""
	var err error
	doc.Find("table#work_outline").Each(func(i int, s *goquery.Selection) {
		staff, err = fmtText(s.Text())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	})
	data["staff"] = staff
}
