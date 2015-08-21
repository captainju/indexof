package indexof

import (
	"net/url"
	"fmt"
	"sync"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

var (
	nbWorker = make(chan bool, 20)
)

type visitedUrlsStruct struct {
	visitedUrlsAccessToken chan bool
	visitedUrls            map[string]int
	results                *chan string
}

func (v *visitedUrlsStruct) hasKey(key string) bool {
	v.visitedUrlsAccessToken <- true
	defer func() { <-v.visitedUrlsAccessToken }()
	_, hasKey := v.visitedUrls[key];
	return hasKey
}

func (v *visitedUrlsStruct) addKey(key string) {
	v.visitedUrlsAccessToken <- true
	defer func() { <-v.visitedUrlsAccessToken }()
	val, _ := v.visitedUrls[key]
	v.visitedUrls[key] = val +1
}

func Index(urlToVisit string, results *chan string) {
	go func() {
		visitedUrls := visitedUrlsStruct{
			visitedUrls:make(map[string]int),
			visitedUrlsAccessToken:make(chan bool, 1),
			results : results }

		parsedUrlToVisit, err := url.Parse(urlToVisit)
		if err != nil {
			panic(err)
		}

		wg := &sync.WaitGroup{}
		visitUrl(parsedUrlToVisit, wg, &visitedUrls)
		wg.Wait()
		close(*results)
	}()
}

func visitUrl(urlToVisit *url.URL, wg *sync.WaitGroup, visitedUrls *visitedUrlsStruct) {
	defer wg.Done()
	if visitedUrls.hasKey(urlToVisit.String()) {
		return
	}

	visitedUrls.addKey(urlToVisit.String())

	//fetch
	nbWorker <- true
	defer func() { <-nbWorker }()
	resp, err := http.Get(urlToVisit.String())

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()


	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	trs := doc.Find("table tr")
	trs.Each(func(index int, node *goquery.Selection) {
		tds := node.Find("td a")
		if tds.Size() > 0 {

			val := tds.Get(0).Attr[0]

			new_path, err := url.Parse(val.Val)
			if err != nil {
				panic(err)
			}
			recomposed_url := urlToVisit.ResolveReference(new_path)

			if strings.HasSuffix(recomposed_url.Path, "/") {

				//craw children
				wg.Add(1)
				go visitUrl(recomposed_url, wg, visitedUrls)

			} else {
				*(visitedUrls.results) <- recomposed_url.String()
			}
		}
	})
	return
}
