// go-bot-labirint
package main

import (
	"fmt"
	"os"
	"strings"
//	"strconv"  
  	"io/ioutil"	
	"io"
	"net/http"
	"golang.org/x/net/html"
)

// структура книги
type dataBook struct {
	name  string // название книги
	autor string // автор
	year  int    // год издания
	kolpages int // кол-во стрниц
	ves  int   // вес книги
	price float32 // цена для всех (обычная)
	pricediscount float32 // цена со скидкой которая видна

}

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

//----------------
// возвращает содержимое html страницы по урлу suri и ошибку 
func gethtmlfromurl0(suri string) (string,error){
	resp, err := http.Get(suri)
  	body, err := ioutil.ReadAll(resp.Body)   
	return string(body),err	
}

// сохраняет содержимое html страницы по урлу suri в файл с именем namef и ошибку 	
func savehtmlfromurl0(suri string,namef string) (string,error) {
	// Create a new browser and open url.
    shtml,err:=gethtmlfromurl0(suri)
    if err != nil {
        panic(err)
    }else{
			err=savestrtofile(namef,shtml)	
		}
	return shtml,err
}  
//----------------

func gethtmlfromurl(suri string) (io.ReadCloser,error){
	resp, err := http.Get(suri)
  	body:=resp.Body 
	//defer b.Close() // close Body when the function returns 
	return body,err	
}


//функция парсинга страницы shtml
func parsehtmlbookean(shtml io.Reader) dataBook {
	var resdata dataBook
	z := html.NewTokenizer(shtml)
	fmt.Println(z)
	return resdata
}

func main() {
	namestore:="bookean"	
	namefurls:=namestore+"-url.cfg"
//	namefhtml:=namestore+"-page.html"
	
	fmt.Println("Start programm....!")
	// сохраняем урлы в файлы
	list_urls:=readcfgs(namefurls)
	for i:=0;i<len(list_urls)-1;i++{
		//s, _:=savehtmlfromurl(list_urls[i],strconv.Itoa(i)+namefhtml)		
		s, _:=gethtmlfromurl(list_urls[i])
		fmt.Println(parsehtmlbookean(s))
	}
	fmt.Println("The end....!")
}
