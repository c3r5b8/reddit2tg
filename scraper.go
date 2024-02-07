package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

var (
	BaseUrl      string
	SubredditUrl string
	NextUrl      string
)

func initScarper() {
	BaseUrl = os.Getenv("REDLIB_INSTANCE_URL")
	SubredditUrl = BaseUrl + os.Getenv("SUBREDDIT")
	NextUrl = SubredditUrl
}

func scroblePage(url string, chatId int) error {
	c := colly.NewCollector()
	var i int
	c.OnHTML("video", func(e *colly.HTMLElement) {
		postId := strings.Split(e.Attr("src"), "/")[2]
		post, _ := queries.GetPost(context.Background(), postId)
		if post.Shown {
			return
		}
		i++
		sendVideo(BaseUrl+e.Attr("src"), postId, chatId)
	})
	c.OnHTML("image", func(e *colly.HTMLElement) {
		if !strings.HasPrefix(e.Attr("href"), "/thumb/") {
			postId := strings.Split(e.Attr("href"), "/")[2]
			post, _ := queries.GetPost(context.Background(), postId)
			if post.Shown {
				return
			}
			i++
			sendPhoto(BaseUrl+e.Attr("href"), postId, chatId)
		}
	})
	c.OnHTML(".post_thumbnail", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c := colly.NewCollector()
		c.OnHTML("img", func(e *colly.HTMLElement) {
			if strings.HasPrefix(e.Attr("src"), "/static") ||
				strings.HasPrefix(e.Attr("src"), "/preview/external-pre/") {
				return
			}
			i++
			postId := strings.Split(strings.Split(e.Attr("src"), "/")[3], "?")[0]
			post, _ := queries.GetPost(context.Background(), postId)
			if post.Shown {
				return
			}

			sendPhoto(BaseUrl+e.Attr("src"), postId, chatId)
		})
		err := c.Visit(BaseUrl + link)
		if err != nil {
			fmt.Println(err)
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("accesskey") == "N" {
			NextUrl = SubredditUrl + e.Attr("href")
		}
	})

	err := c.Visit(url)
	if err != nil {
		return err
	}
	return nil
}
