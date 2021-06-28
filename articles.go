package main

import (
	"encoding/json"
	"github.com/hb0730/go-request"
	"strconv"
	"time"
)

//Articles 文章
// 通过公众号的方式
type Articles struct {
	// Fakeid 作者
	Fakeid string
	//Token 公众号token
	Token string
	//Cookie 公众号Cookie
	Cookie string
	//Count 一页多少条
	Count int
}

func NewArticles(fakeid, token, cookie string, count int) *Articles {
	return &Articles{
		Fakeid: fakeid,
		Token:  token,
		Cookie: Cookie,
		Count:  count,
	}
}

func (a *Articles) GetArticles(begin int) (ArticlesResult, error) {
	var result ArticlesResult
	params := map[string]string{
		"action": "list_ex",
		"begin":  strconv.Itoa(begin),
		"count":  strconv.Itoa(a.Count),
		"fakeid": a.Fakeid,
		"type":   "9",
		"query":  "",
		"token":  a.Token,
		"lang":   "zh_CN",
		"f":      "json",
		"ajax":   "1",
	}
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36",
	}
	req, err := request.CreateRequest("GET", "https://mp.weixin.qq.com/cgi-bin/appmsg", "")
	if err != nil {
		return result, nil
	}
	request.SetGetParams(req, params)
	req.SetHeaders(headers)
	req.SetCookies(a.Cookie)
	err = req.Do()
	if err != nil {
		return result, nil
	}
	bt, err := req.GetBody()
	if err != nil {
		return result, nil
	}
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (a Articles) GetAllArticles() ([]AppMsg, error) {
	r, err := a.GetArticles(0)
	if err != nil {
		return nil, err
	}
	//总页数
	totalPageNum := (r.AppMsgCnt + a.Count - 1) / a.Count
	dataList := make([]AppMsg, 0)
	for i := 0; i < totalPageNum; i++ {
		r, err := a.GetArticles(i * a.Count)
		if err != nil {
			break
		}
		dataList = append(dataList, r.AppMsgList...)
		time.Sleep(5 * time.Second)
	}
	return dataList, nil
}

func (a Articles) GetAllArticlesByThread(totalThread int) ([]AppMsg, error) {
	r, err := a.GetArticles(0)
	if err != nil {
		return nil, err
	}
	//总页数
	totalPageNum := (r.AppMsgCnt + a.Count - 1) / a.Count
	return a.run(totalPageNum, totalThread), nil

}

func (a Articles) run(totalPageNum, totalThread int) []AppMsg {
	pageData := (totalPageNum + totalThread - 1) / totalThread //每个线程处理的数据页数
	list := make([]AppMsg, 0)
	//ch := make(chan int, totalThread)
	dataChan := make(chan []AppMsg, totalThread)
	func() {
		for i := 1; i <= totalThread; i++ {
			go func(k, p, t int) {
				//c := 0

				dataList := make([]AppMsg, 0)
				for j := (k - 1) * p; j < p*k; j++ {
					if j > t {
						break
					}
					r, err := a.GetArticles(j * a.Count)
					if err != nil {
						break
					}
					dataList = append(dataList, r.AppMsgList...)

					//fmt.Printf("%d : %d \n", j*a.Count, a.Count)
					//c += j
					time.Sleep(5 * time.Second)
				}
				//ch <- c
				dataChan <- dataList

			}(i, pageData, totalPageNum)
		}
	}()
	//for i := 0; i < totalThread; i++ {
	//	x := <-ch
	//	fmt.Printf(" <-ch :%d \n", x)
	//}
	for i := 0; i < totalThread; i++ {
		data := <-dataChan
		list = append(list, data...)
	}
	return list
}

type ChanList struct {
	List []AppMsg
}

type ArticlesResult struct {
	AppMsgCnt  int      `json:"app_msg_cnt"`
	AppMsgList []AppMsg `json:"app_msg_list"`
	BaseResp   BaseResp `json:"base_resp"`
}

type AppMsg struct {
	Aid              string `json:"aid"`
	AlbumId          string `json:"album_id"`
	AppmsgAlbumInfos []struct {
		AlbumId          int64         `json:"album_id"`
		AppmsgAlbumInfos []interface{} `json:"appmsg_album_infos"`
		Id               string        `json:"id"`
		Title            string        `json:"title"`
	} `json:"appmsg_album_infos"`
	Appmsgid              int64  `json:"appmsgid"`
	Checking              int    `json:"checking"`
	CopyrightType         int    `json:"copyright_type"`
	Cover                 string `json:"cover"`
	CreateTime            int    `json:"create_time"`
	Digest                string `json:"digest"`
	HasRedPacketCover     int    `json:"has_red_packet_cover"`
	IsPaySubscribe        int    `json:"is_pay_subscribe"`
	ItemShowType          int    `json:"item_show_type"`
	Itemidx               int    `json:"itemidx"`
	Link                  string `json:"link"`
	MediaDuration         string `json:"media_duration"`
	MediaapiPublishStatus int    `json:"mediaapi_publish_status"`
	PayAlbumInfo          struct {
		AppmsgAlbumInfos []interface{} `json:"appmsg_album_infos"`
	} `json:"pay_album_info"`
	Tagid      []interface{} `json:"tagid"`
	Title      string        `json:"title"`
	UpdateTime int           `json:"update_time"`
}
