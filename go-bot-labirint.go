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
	//surl:="http://labirint.ru"
	surl:="https://git-scm.com/book/ru/v1/%D0%9E%D1%81%D0%BD%D0%BE%D0%B2%D1%8B-Git-%D0%97%D0%B0%D0%BF%D0%B8%D1%81%D1%8C-%D0%B8%D0%B7%D0%BC%D0%B5%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9-%D0%B2-%D1%80%D0%B5%D0%BF%D0%BE%D0%B7%D0%B8%D1%82%D0%BE%D1%80%D0%B8%D0%B9"
	shtml,serr:=gethtmlfromurl(surl)
	fmt.Println(shtml)
	fmt.Println(serr)
}
