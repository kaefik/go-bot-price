// go-bot-labirint
package main

import (
	"fmt"
	"github.com/headzoo/surf"
)

// возвращает содержимое html страницы по урлу suri и ошибку 
func gethtmlfromurl(suri string) (string,error) {
	// Create a new browser and open url.
    bow := surf.NewBrowser()
    err := bow.Open(suri)
   // if err != nil {
 //       panic(err)
 //   }
	shtml:=bow.Body()
	
	return shtml,err
}



func main() {
	fmt.Println("Hello World!")
	gethtmlfromurl()
}
