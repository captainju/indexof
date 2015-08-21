package indexof

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"net/url"
	"strings"
)

var (
	nb_request = make(chan bool, 20)
)

func OldIndex(url_str string) (<-chan string) {
	visited_urls := make(map[string]string)
	visited_urls[url_str] = url_str
	result_chan := make(chan string, 10)
	newSearch(url_str, &visited_urls, result_chan)
	return result_chan
}

func newSearch(url_str string, visited_urls *map[string]string, result_chan chan string) {
	var err error
	var parsed_url *url.URL
	parsed_url, err = url.Parse(url_str)
	if err != nil {
		panic(err)
	}

	var resp *http.Response
	nb_request <- true
	resp, err = http.Get(parsed_url.String())

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	defer releaseConnection()

	doc, errgo := goquery.NewDocumentFromResponse(resp)
	if errgo != nil {
		fmt.Println(err)
		return
	}

	trs := doc.Find("table tr")
	trs.Each(func(index int, node *goquery.Selection) {
		tds := node.Find("td a")
		if tds.Size() > 0 {

			go extractData(tds, parsed_url, *visited_urls, result_chan)

		}
	})
	return
}

func releaseConnection() {
	<-nb_request
}

func extractData(tds *goquery.Selection, parsed_url *url.URL, visited_urls map[string]string, result_chan chan string) {

	val := tds.Get(0).Attr[0]

	new_path, err := url.Parse(val.Val)
	if err != nil {
		panic(err)
	}
	recomposed_url := parsed_url.ResolveReference(new_path)

	if _, ok := visited_urls[recomposed_url.String()]; !ok {

		var full_url = recomposed_url.String()

		if !strings.Contains(recomposed_url.Path, ".") {
			visited_urls[full_url] = full_url
			newSearch(full_url, &visited_urls, result_chan)
		} else {
			result_chan <- full_url
		}
	}

}
