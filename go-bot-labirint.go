// go-bot-labirint
// программа скачивает по ссылкам данные по книгам Лабиринт и проверяет условия по ссылкам
// Автор: Ильнур Сайфутдинов
// email: ilnursoft@gmail.com
// декабрь 2015

// рефакторинг кода: создание новых сущностей и объединение существующих
package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"  
  	"io/ioutil"	
	"io"
	"net/http"
	"time"
	"net/smtp"
//	"golang.org/x/net/html"
	"github.com/ddo/pick"
	"golang.org/x/net/html/charset"
	"log"
)

//------------ Объявление типов и глобальных переменных

// структура задания с информацией по книге
type TaskerBook struct {
	url string  // ссылка на источник данных
	Book
	Tasker
}


// структура книги
type Book struct {
	name  string // название книги
	autor string // автор
	year  int    // год издания
	kolpages int // кол-во стрниц
	ves  int   // вес книги
	price int // цена для всех (обычная)
	pricediscount int // цена со скидкой которая видна	
	url string  // ссылка на источник данных
}

// задание-триггер для срабатывания оповещения
type Tasker struct {
	url string  // ссылка на источник данных
	uslovie string // условие < , > , = 
	price int // цена триггера
	result bool // результат срабатывания триггера, если true , то триггер сработал 
}

var LogFile *log.Logger 

//------------ END Объявление типов и глобальных переменных



// проверки триггеров
func TriggersisUslovie(book0 []Book,task []Tasker) []Tasker {
	for i:=0;i<len(task);i++ {
		task[i].isUslovie(book0)
	}
	return task
}


// ---------------  парсинг магазина Лабиринт
//получение данных книги из магазина лабиринт по урлу url
func (dbook *Book) Getlabirint(url string) {
	body:=gethtmlpage(url)		
	shtml := string(body)	
//	fmt.Println(shtml)

	//book0:=parselabirintbook(shtml)
		scena, _ := pick.PickText(&pick.Option{   // текст цены книги
		&shtml,
		"span",
		&pick.Attr{
			"itemprop",
			"price",
		},
	})

	scenaskidka, _ := pick.PickText(&pick.Option{   // текст цены со скидкой книги
		&shtml,
		"span",
		&pick.Attr{
			"class",
			"buying-pricenew-val-number",
		},
	})		
	
	sauthor, _ := pick.PickText(&pick.Option{   // текст описания книги
		&shtml,
		"span",
		&pick.Attr{
			"itemtype",
			"http://schema.org/ItemList",
		},
	})

	stitle, _ := pick.PickText(&pick.Option{&shtml,"span",&pick.Attr{"itemprop","name"}})	
	//book0=parsedescribebook(sauthor)
	
	for i:=0;i<len(sauthor);i++ {
		switch sauthor[i] {
			case "Автор(ы)": dbook.autor=sauthor[i+1]
			case "Масса": dbook.ves,_=strconv.Atoi(sauthor[i+1])	
			case "Количество страниц": dbook.kolpages,_=strconv.Atoi(sauthor[i+1])		 
		}
	}
	
	dbook.name=stitle[1]
	if len(scenaskidka)>0 {
		dbook.pricediscount,_=strconv.Atoi(scenaskidka[0])
	}
	vv := strings.Split(scena[0], " ")
	dbook.price,_ =strconv.Atoi(vv[1])
	
	return
}


// --------------- END  парсинг магазина Лабиринт


// -----------  функции для Book

func (book0 *Book) print () {
	LogFile.Println("Автор: ",book0.autor)
	LogFile.Println("Название книги: ",book0.name)
	LogFile.Println("Вес: ",book0.ves)
	LogFile.Println("Кол-во страниц: ",book0.kolpages)
	LogFile.Println("Цена: ",book0.price)
	LogFile.Println("Цена со скидкой: ",book0.pricediscount)
	LogFile.Println("Ссылка на книгу: ",book0.url)
	return
}

func printbook (book0 Book) {
	LogFile.Println("Автор: ",book0.autor)
	LogFile.Println("Название книги: ",book0.name)
	LogFile.Println("Вес: ",book0.ves)
	LogFile.Println("Кол-во страниц: ",book0.kolpages)
	LogFile.Println("Цена: ",book0.price)
	LogFile.Println("Цена со скидкой: ",book0.pricediscount)
	LogFile.Println("Ссылка на книгу: ",book0.url)
	return
}

//сохранить данные Book в файл 
func (db *Book) savetocsvfile(namef string) error {
	var fileflag bool = false
	if _, err := os.Stat(namef); os.IsNotExist(err) {
 	 // path/to/whatever does not exist
		fileflag=true
	}
	
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
		// handle the error here
		return err
	}
	defer file.Close()
	if fileflag { // если не существует файл
		stitle:="Дата выгрузки;Автор;Название книги;Год издания;Кол-во стр.;Вес;Цена;Цена со скидкой;Ссылка"+"\n"
		file.WriteString(stitle)
	}
    curdate := time.Now().String()
	str:=curdate+";"+db.autor+";"+db.name+";"+strconv.Itoa(db.year)+";"+strconv.Itoa(db.kolpages)+";"+strconv.Itoa(db.ves)+";"+strconv.Itoa(db.price)+";"+strconv.Itoa(db.pricediscount)+";"+db.url+"\n"
	file.WriteString(str)
	return err
	return err
}



// ----------- END функции для Book


// -----------  функции для Tasker

// чтение из текстового конфиг файла заданий и возращает массив заданий Tasker
func readtaskercfg(namef string) []Tasker {
	var res []Tasker
	str := readfiletxt(namef)
	vv := strings.Split(str, "\n")	
	for i:=0;i<len(vv);i++ {
		s:=strings.Split(vv[i],";")		
		tt,_:= strconv.Atoi(s[2])
		res=append(res,Tasker{s[0],s[1],tt,false})		
	}
	return res
}

//проверка триггеров по массиву полученных данных по книгах
func (task *Tasker) isUslovie(book0 []Book) {
	var res1, res2 bool
	res1=false
	res2=false
	for j:=0;j<len(book0);j++ {
		if task.url == book0[j].url {
				switch task.uslovie {
					case ">": res1 = book0[j].price > task.price
					case "=": res1 = book0[j].price == task.price
					case "<": res1 = book0[j].price < task.price
					default: res1=false
				}
				if book0[j].pricediscount>0 { // если цена со скидкой больше нуля, то проверяем триггер на скидку
					switch task.uslovie {
						case ">": res2 = book0[j].pricediscount > task.price
						case "=": res2 = book0[j].pricediscount == task.price
						case "<": res2 = book0[j].pricediscount < task.price
						default: res2=false
					}
				}
				task.result=res1 || res2
				}  
	}
}

// если тригер сработал то возвращает строку сообщения, иначе пусто
func (task *Tasker) genmessage(book0 []Book) string {
	var sprice, spricedisc string
	var smegtrigger, smegtrigger0, smsg string
	smsg=""
	if task.result {
		for j:=0;j<len(book0);j++ {
			b:= book0[j]
			if task.url == b.url {
				sprice = strconv.Itoa(b.price)
				spricedisc = strconv.Itoa(b.pricediscount)
				smegtrigger="Сбработал триггер по книге: \n\n"+"Автор: "+b.autor+"\n"+"Название: "+b.name+"\n"+"Цена: "+sprice+"\n"+"Цена со скидкой: "+spricedisc+"\n"+"Ссылка: "+b.url+"\n\n"	
				sprice= strconv.Itoa(task.price)
				smegtrigger0="Условие триггера: "+task.uslovie + "\n Цена триггера: "+sprice+"\n\n"
				smsg=smegtrigger+smegtrigger0
			}
		}			
	}
 	return smsg
}

//отправка сообщения если сработал триггер
func (task *Tasker) sendmessage(book0 []Book, toaddr string){
	smsg:=task.genmessage(book0)
	if smsg!="" {
		sendmailyandex("сработал триггер",smsg, toaddr)
	}
	return
}

// ----------- END  функции для Tasker


//---------------- общие функции ---------------------

func InitLogFile(namef string) *log.Logger {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}
	multi:= io.MultiWriter(file, os.Stdout)
	LFile:= log.New(multi, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)	
	return LFile
}

//отправка почты через яндекс темой stema сообщение smsg адресату toaddr
func sendmailyandex(stema, smsg, toaddr string) bool {
	auth := smtp.PlainAuth("", "magazinebot@yandex.ru", "qwe123!!", "smtp.yandex.ru") 	
	to := []string{toaddr}
	msg := []byte("To: "+toaddr+"\r\n" +
    "Subject: "+stema+" \r\n" +
    "\r\n" +
    smsg+"\r\n")
	err := smtp.SendMail("smtp.yandex.ru:25", auth, "magazinebot@yandex.ru", to, msg)
	
    if err != nil { 
//        log.Fatal(err) 
		panic(err)
    } 
	return true
}

//получение страницы из урла url
func gethtmlpage(url string) []byte {
	resp, err := http.Get(url)
    if err != nil {
        fmt.Println("HTTP error:", err)
		panic("HTTP error")        
    }

    defer resp.Body.Close()
    // вот здесь и начинается самое интересное
    utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
    if err != nil {
        fmt.Println("Encoding error:", err)
        panic("Encoding error")
    }
    // оп-па-ча, готово
	
//	fmt.Println(namebook(utf8))
	
    body, err := ioutil.ReadAll(utf8)
    if err != nil {
        fmt.Println("IO error:", err)
		panic("IO error")
    }
	return body
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

//---------------- END общие функции ---------------------



func main() {
	var books []Book
	var book0 Book
	toaddr:="i.saifutdinov@kazan.2gis.ru"
	//sdir:="books"
	namestore:="labirint"	
	namefurls:=namestore+"-url.cfg"
	namelogfile:=namestore+".log"
	
	LogFile=InitLogFile(namelogfile)  // инициализация лог файла		
	LogFile.Println("Starting programm")	
	
	// получаем урлы из файлы
//    list_urls:=readcfgs(namefurls)
	
	// получаем задания из файла
	list_tasker:=readtaskercfg(namefurls)
	
	fmt.Println(readtaskercfg(namefurls))
	
	//получение данных книжек
	for i:=0;i<len(list_tasker);i++{
		
		//book0:=getbooklabirint(list_tasker[i].url)
		book0.Getlabirint(list_tasker[i].url)
		
		namef:=	namestore+".csv"
		book0.url=list_tasker[i].url
		book0.savetocsvfile(namef)
		books=append(books,book0)
		//printbook(book0)
	//	book0.print()
	}
	
	for i:=0;i<len(books);i++ {
		books[i].print()
	}
	
	//проверка на наличии срабатываний	
	TriggersisUslovie(books,list_tasker)
	
	for i:=0;i<len(list_tasker);i++{
		LogFile.Println(list_tasker[i].genmessage(books))
		list_tasker[i].sendmessage(books, toaddr)	
	}
	
	LogFile.Println("The end....!\n")
}
