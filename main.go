package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	DraganBallZ              = "dragonBallZ"
	DraganBallZTitleEncoding = "Dragon+Ball+Z+%28Dub%29"
)

var (
	start   = flag.Int("s", 0, "start (must be <=end)")
	end     = flag.Int("e", 0, "end (must be >=start)")
	title   = flag.String("t", DraganBallZ, "title (not nil)")
	numFile = flag.String("f", "", "num list file. if it exists, choose num files instead of start,end")
)

func main() {
	flag.Parse()

	if *title == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *numFile == "" && *start == 0 && *end == 0 {
		fmt.Printf("both num file and start/end is nil\n")
		flag.Usage()
		os.Exit(1)
	}

	var nums []int

	if *numFile == "" {
		nums = makeRange(*start, *end)
	} else {
		nums = readNums(*numFile)
	}

	crawlArray(*title, nums)
}

func readNums(file string) []int {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Printf("read file num file %s failed %v\n", file, err)
		os.Exit(1)
	}

	numstr := strings.Split(string(b), "\n")

	if len(numstr) == 0 {
		return nil
	}

	var r []int

	for _, s := range numstr {
		n, err := strconv.Atoi(s)

		if err != nil {
			fmt.Printf("find no num %s in file %s\n", s, file)
			os.Exit(1)
		}

		r = append(r, n)
	}

	return r
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

func makeRange(start, end int) []int {
	if start > end {
		return nil
	}

	var r []int
	for i := start; i <= end; i++ {
		r = append(r, i)
	}

	return r
}

func crawlArray(title string, nums []int) {
	for _, i := range nums {
		link, err := crawlEp(title, i)

		if err != nil {
			fmt.Printf("crawl %s failure ep%d\terror %v\n", title, i, err)
		} else {
			fmt.Printf("crawl %s success ep%d\t %s\n", title, i, link)
		}
	}
}
