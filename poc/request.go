package main
import (
	"net/http"
	"net/http/cookiejar"
	//"io/ioutil"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func main() {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{Jar : jar}
	resp, err := client.Get("http://example.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(body))

	doc, err := goquery.NewDocumentFromResponse(resp)

	fmt.Println(doc.Find("p").First().Text())
}
