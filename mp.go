package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hb0730/go-request"
	"strconv"
	"time"
)

type Mp struct {
	Query   string
	Token   string
	Cookies string
	Count   int
}

func NewMp(query, token, cookies string) *Mp {
	return &Mp{
		Query:   query,
		Token:   token,
		Cookies: cookies,
	}
}

// GetFirstMp 获取第一页第一个数据
func (m *Mp) GetFirstMp() (MpList, error) {
	list, err := m.GetMp(0)
	if err != nil {
		return MpList{}, err
	}
	return list.List[0], nil
}

// GetAllMp 获取所有的MP
func (m *Mp) GetAllMp() ([]MpList, error) {
	r, err := m.GetMp(0)
	if err != nil {
		return nil, err
	}
	totalPageNum := (r.Total + m.Count - 1) / m.Count
	list := make([]MpList, 0)
	for i := 0; i < totalPageNum; i++ {
		time.Sleep(3 * time.Second)
		m, err := m.GetMp(i)
		if err != nil {
			break
		}
		list = append(list, m.List...)
	}
	return list, nil
}

// GetMp 获取第一页的数据
func (m *Mp) GetMp(begin int) (MpResult, error) {
	var result MpResult
	params := map[string]string{
		"action": "search_biz",
		"begin":  strconv.Itoa(begin),
		"count":  strconv.Itoa(m.Count),
		"query":  m.Query,
		"token":  m.Token,
		"lang":   "zh_CN",
		"f":      "json",
		"ajax":   "1",
	}
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36",
	}
	req, err := request.CreateRequest("GET", "https://mp.weixin.qq.com/cgi-bin/searchbiz", "")
	if err != nil {
		return result, err
	}
	request.SetGetParams(req, params)
	req.SetHeaders(headers)
	req.SetCookies(m.Cookies)
	err = req.Do()
	if err != nil {
		return result, err
	}
	bt, err := req.GetBody()
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return result, err
	}
	if result.BaseResp.ErrMsg != "ok" {
		return result, errors.New(fmt.Sprintf("查找失败,Code:%d Message: %s", result.BaseResp.Ret, result.BaseResp.ErrMsg))
	}
	return result, nil
}

type MpResult struct {
	BaseResp BaseResp `json:"base_resp"`
	List     []MpList `json:"list"`
	Total    int      `json:"total"`
}

type MpList struct {
	Fakeid       string `json:"fakeid"`
	Nickname     string `json:"nickname"`
	Alias        string `json:"alias"`
	SoundHeadImg string `json:"round_head_img"`
	ServiceType  int    `json:"service_type"`
	Signature    string `json:"signature"`
}
