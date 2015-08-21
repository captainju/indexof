package indexof

import (
	"testing"
	"fmt"
)


func TestIndex(t *testing.T) {

	urlTest := "https://jaimeleschips.fr/private/library/"
	/*
	http://ledarkangel.com/stream/
	https://jaimeleschips.fr/private/library/
	http://12.228.222.153/
	http://rahsia.info/temp/Kids%20Movies/

	 */

	results := make(chan string, 10)
	Index(urlTest, &results)

	counter := 0
	for res := range results {
		fmt.Println(res)
		counter++
	}
	fmt.Println(counter)

}
