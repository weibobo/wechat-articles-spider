package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_get(t *testing.T) {
	totalPageNum := (423 + 5 - 1) / 5
	////85
	log.Println(totalPageNum)
	totleThread := 5
	////29
	pageData := (totalPageNum + totleThread - 1) / totleThread //每个线程处理的数据页数

	for i := 1; i <= totleThread; i++ {
		k := i
		b := totalPageNum
		p := pageData
		for j := (k - 1) * p; j < p*k; j++ {
			//if j == 0 {
			//	j = 1
			//}
			if j > b {
				return
			}
			fmt.Println(j)
		}
	}
}
func TestArticlesTest(t *testing.T) {
	totalPageNum := (63 + 5 - 1) / 5
	////85
	log.Println(totalPageNum)
	totalThread := 5
	////29
	pageData := (totalPageNum + totalThread - 1) / totalThread //每个线程处理的数据页数
	for i := 1; i <= totalThread; i++ {
		c := make(chan []AppMsg, totalPageNum*i)
		go func(k, p, t int) {
			list := make([]AppMsg, 0)
			for j := (k - 1) * p; j < p*k; j++ {
				if j > t {
					return
				}
				applist := []AppMsg{
					{Aid: ""},
				}
				list = append(list, applist...)
				c <- list
				fmt.Printf("%d : %d \n", j*5, 5)
				time.Sleep(5 * time.Second)
			}
			close(c)
		}(i, totalPageNum, pageData)

		for ch := range c {
			println(ch)
		}
	}
}

func TestChan(t *testing.T) {
	for j := 0; j < 10; j++ {
		c := make(chan int)
		go func() {
			for i := 0; i < 3; i++ {
				c <- i
				fmt.Printf("send %d\n", i)
				time.Sleep(time.Second)
			}
			fmt.Println("ready close channel")
			close(c)
		}()

		for i := range c {
			fmt.Printf("receive %d\n", i)
		}
		fmt.Println("quit for loop")
	}
}

func TestArticles_run(t *testing.T) {
	//
	a := NewArticles("", "", "", 5)
	a.run(13, 5)

}
