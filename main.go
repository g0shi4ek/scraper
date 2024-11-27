package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
)

func getLink(t html.Token) (bool, string) {
	for _, el := range t.Attr {
		if el.Key == "href" {
			return true, el.Val
		}
	}
	return false, ""
}

func crawl(url string, chFinish chan bool, chLinks chan string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	defer func() {
		chFinish <- true // при завершении
	}()

	tokenizer := html.NewTokenizer(resp.Body) // для получения токенов

	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			return // конец файла
		case tt == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				_, link := getLink(token)
				hasProto := strings.Index(link, "http") == 0 //проверяет начало с http
				if hasProto {
					chLinks <- link
				}
			}
		}
	}
}

func main() {
	found := make(map[string]bool)
	urlLinks := os.Args[1:] // из командной строки 1 тк название файла
	c := 0
	// каналы
	chLinks := make(chan string)
	chFinish := make(chan bool)

	for _, url := range urlLinks {
		go crawl(url, chFinish, chLinks)
	}

	for c < len(urlLinks) {
		select { // чтобы избежать дедлоков
		case url := <-chLinks:
			found[url] = true
		case <-chFinish:
			c++
		}
	}

	fmt.Println("\nFounded links:\n")
	for link := range found {
		fmt.Println(link)
	}

	close(chLinks)
	close(chFinish)
}
