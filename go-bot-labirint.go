// go-bot-labirint
package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"net/http"    
  	"io/ioutil"
	"github.com/headzoo/surf"
)

// чтение файла с именем namefи возвращение содержимое файла, иначе текст ошибки
func readfiletxt(namef string) string {
	file, err := os.Open(namef)
	if err != nil {
		return "handle the error here"
	}
	defer file.Close()
	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return "error here"
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return "error here"
	}
	return string(bs)
}

// чтение из текстового конфиг файла и возращает массив строк
func readcfgs(namef string) []string {
	str := readfiletxt(namef)
	vv := strings.Split(str, "\n")		
	return vv
}

//сохранение строки str в файл с именем namef
func savestrtofile(namef string, str string) error {
	file, err := os.Create(namef)
	if err != nil {
		// handle the error here
		return err
	}
	defer file.Close()

	file.WriteString(str)
	return err
}

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

// сохраняет содержимое html страницы по урлу suri в файл с именем namef и ошибку 
func savehtmlfromurl(suri string,namef string) (string,error) {
	// Create a new browser and open url.
    bow := surf.NewBrowser()
    err := bow.Open(suri)
	shtml:=""
    if err != nil {
        panic(err)
    }else{
		shtml=bow.Body()
		err=savestrtofile(namef,shtml)	
		}
	return shtml,err
}

//----------------
func gethtmlfromurl0(suri string) (string,error)
	resp, err := http.Get(suri)    
  	//fmt.Println("http transport error is:", err)
  	body, err := ioutil.ReadAll(resp.Body)                                           
  	fmt.Println("read error is:", err)

	return string(body),err	
	  
//----------------


func main() {
	namestore:="bookean"	
	namefurls:=namestore+"-url.cfg"
	namefhtml:=namestore+"-page.html"
	
	fmt.Println("Hello World!")
	
	list_urls:=readcfgs(namefurls)
	for i:=0;i<len(list_urls)-1;i++{
		savehtmlfromurl(list_urls[i],strconv.Itoa(i)+namefhtml)		
	}
	
}
