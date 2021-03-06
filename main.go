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
	DragonBallZ              = "dragonBallZ"
	DragonBallZTitleEncoding = "Dragon+Ball+Z+%28Dub%29"
	DragonBallZCoverEncoding = "Y292ZXIvZHJhZ29uLWJhbGwtei1kdWIuanBn"

	BlackLagoon1              = "blackLagoon1"
	BlackLagoon1TitleEncoding = "Black+Lagoon+%28Dub%29"
	BlackLagoon1CoverEncoding = "Y292ZXIvYmxhY2stbGFnb29uLWR1Yi5wbmc="

	BlackLagoon2              = "blackLagoon2"
	BlackLagoon2TitleEncoding = "Black+Lagoon%3A+The+Second+Barrage+%28Dub%29"
	BlackLagoon2CoverEncoding = "Y292ZXIvYmxhY2stbGFnb29uLXRoZS1zZWNvbmQtYmFycmFnZS1kdWIucG5n"

	GTO              = "GTO"
	GTOTitleEncoding = "Great+Teacher+Onizuka+%28Dub%29"
	GTOCoverEncoding = "Y292ZXIvZ3JlYXQtdGVhY2hlci1vbml6dWthLWR1Yi5wbmc="
)

var (
	start   = flag.Int("s", 0, "start (must be <=end)")
	end     = flag.Int("e", 0, "end (must be >=start)")
	title   = flag.String("t", DragonBallZ, "title (dragonBallZ,blackLagoon1,blackLagoon2)")
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

	crawl(*title, nums)
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

func parseQuery(title string, episode int) (string, string, string, error) {
	if title == DragonBallZ {
		if episode >= 1 && episode <= 151 {
			return int2Base64(episode + 76107), DragonBallZTitleEncoding, DragonBallZCoverEncoding, nil
		} else if episode >= 152 && episode <= 207 {
			return int2Base64(episode + 76507), DragonBallZTitleEncoding, DragonBallZCoverEncoding, nil
		} else {
			return int2Base64(episode + 76513), DragonBallZTitleEncoding, DragonBallZCoverEncoding, nil
		}
	} else if title == BlackLagoon1 {
		return int2Base64(episode + 76921), BlackLagoon1TitleEncoding, BlackLagoon1CoverEncoding, nil
	} else if title == BlackLagoon2 {
		return int2Base64(episode + 148295), BlackLagoon2TitleEncoding, BlackLagoon2CoverEncoding, nil
	} else if title == GTO {
		return int2Base64(episode + 82592), GTOTitleEncoding, GTOCoverEncoding, nil
	} else {
		return "", "", "", fmt.Errorf("title %s not supported", title)
	}
}

func crawlEp(title string, episode int) (string, error) {
	ep, t, c, err := parseQuery(title, episode)

	if err != nil {
		return "", fmt.Errorf("parseQuery err %v", err)
	}

	res, err := http.Get("https://gogo-stream.com/download?id=" + ep +
		"&title=" + t +
		"&typesub=SUB&" +
		"sub=W10=" +
		"&cover=" + c)

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

func crawl(title string, nums []int) {
	for _, i := range nums {
		link, err := crawlEp(title, i)

		if err != nil {
			fmt.Printf("crawl %s failure ep%d\terror %v\n", title, i, err)
		} else {
			fmt.Printf("crawl %s success ep%d\t %s\n", title, i, link)
		}
	}
}
