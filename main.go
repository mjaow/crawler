package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	DraganBallZ              = "dragonBallZ"
	DraganBallZTitleEncoding = "Dragon+Ball+Z+%28Dub%29"
)

var (
	start = flag.Int("s", 0, "start (must be <=end)")
	end   = flag.Int("e", 0, "end (must be >=start)")
	title = flag.String("t", DraganBallZ, "title")
)

func main() {
	flag.Parse()

	if start == nil || end == nil || *start == 0 || *end == 0 || title == nil || *title == "" {
		fmt.Println("start/end/title is nil")
		flag.Usage()
		os.Exit(1)
	}

	if *start > *end {
		flag.Usage()
		os.Exit(1)
	}

	crawlRange(*title, *start, *end)
}

func int2Base64(d int) string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", d)))
}

func parseQuery(title string, episode int) (string, string, error) {
	if title == DraganBallZ {
		if episode >= 1 && episode <= 151 {
			return int2Base64(episode + 76107), DraganBallZTitleEncoding, nil
		} else if episode >= 152 && episode <= 207 {
			return int2Base64(episode + 76507), DraganBallZTitleEncoding, nil
		} else {
			return int2Base64(episode + 76513), DraganBallZTitleEncoding, nil
		}
	} else {
		return "", "", fmt.Errorf("title %s not supported", title)
	}
}

func crawlEp(title string, episode int) (string, error) {
	ep, t, err := parseQuery(title, episode)

	if err != nil {
		return "", fmt.Errorf("parseQuery err %v", err)
	}

	res, err := http.Get("https://gogo-stream.com/download?id=" + ep +
		"&title=" + t +
		"&typesub=SUB&" +
		"sub=W10=" +
		"&cover=Y292ZXIvZHJhZ29uLWJhbGwtei1kdWIuanBn")

	if err != nil {
		return "", fmt.Errorf("http get err %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("read doc error: %v", err)
	}

	link := doc.Find(".mirror_link .dowload a[href*=\"360p\"]").First().AttrOr("href", "")

	if link == "" {
		return "", fmt.Errorf("no 360p")
	}

	if !strings.Contains(link, fmt.Sprintf("EP.%d.", episode)) {
		return "", fmt.Errorf("find link %s, but not match", link)
	}

	return link, nil
}

func crawlRange(title string, start, end int) {
	for i := start; i <= end; i++ {
		link, err := crawlEp(title, i)

		if err != nil {
			fmt.Printf("crawl %s failure ep%d\terror %v\n", title, i, err)
		} else {
			fmt.Printf("crawl %s success ep%d\t %s\n", title, i, link)
		}
	}
}
