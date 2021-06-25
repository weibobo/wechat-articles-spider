package main

import "log"

var (
	Cookie = ""
	Token  = ""
)

func main() {
	mp := NewMp("知一码园", Token, Cookie)
	m, err := mp.GetFirstMp()
	if err != nil {
		panic(err)
	}
	a := NewArticles(m.Fakeid, Token, Cookie, 5)
	all, err := a.GetAllArticles()
	if err != nil {
		panic(err)
	}
	for _, v := range all {
		log.Println(v.Aid)
	}
}
