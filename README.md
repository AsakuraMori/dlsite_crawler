# dlsite_crawler

简易的dlsite爬虫。

## 函数：

``` go
result, err := dlsite_crawler.GetInfoFromSearch([名字], [类型], client)
```

该函数用于搜索条目，类型只能是以下字符串：

```
游戏、同人、漫画、手机游戏
```

如果指定上述以外的字符串，则可能调用失败。

返回一个[]map[string]interface{}和一个error，如果err为nil，则没有错误。

=================================================================

``` go'
result, err := dlsite_crawler.GetInfoByID([ID], [类型], client)
```

该函数用于查看条目，类型只能是以下字符串：

```
游戏、同人、漫画、手机游戏
```

如果指定上述以外的字符串，则可能调用失败。

返回一个map[string]interface{}和一个error，如果err为nil，则没有错误。

## 关于client

指定模版如下：

``` go
proxyURL := "" // 替换为你的代理地址
proxy, err := url.Parse(proxyURL)
if err != nil {
	panic(err)
}
client := &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyURL(proxy),
	},
}
```

