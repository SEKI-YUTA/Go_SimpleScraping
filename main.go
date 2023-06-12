package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Flags に net.FlagUp があれば接続されてる
	connected := checkNetworkConnect()

	if !connected{
		fmt.Println("ネットに接続されていません")
		return
	}

	// res, _ := http.Get("https://yahoo.co.jp")
	client := &http.Client{}
	header := http.Header{}
	// デフォルトだとスマホ版のページが帰ってくる 2023-0612現在
	// そのため動作確認のしやすさなどからパソコン版のページを使用するようにするためにパソコン版のユーザーエージェントを設定する
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36");
	req, err := http.NewRequest("GET", "https://google.com/search?q=坂井泉水", nil)
	if err != nil {
		fmt.Println("failed to request page")
		os.Exit(1)
	}
	req.Header = header
	res, _ := client.Do(req)
	// res, _ := http.Get()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	
	if err != nil {
		fmt.Println("failed to html convert to document")
		os.Exit(1)
	}

	// titles := doc.Find(".zBAuLc.l97dzf") // devToolのiPhoneSEで動作確認
	titles := doc.Find(".LC20lb.MBeuO.DKV0Md")

	fmt.Println(titles.Size())
	titles.Each(
		func(i int, s *goquery.Selection) {
			fmt.Println(s.Text())
	})

}

func checkNetworkConnect() bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("failed to get network status")
	}

	for _, v := range interfaces {
		if v.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}

		if strings.ToLower(v.Name) == "wi-fi" {
			if v.Flags&net.FlagUp == net.FlagUp {
				fmt.Println("Wi-Fi is up")
				return true
			}
		}
	}
	fmt.Println("Wi-Fi is not connected")
	return false
}