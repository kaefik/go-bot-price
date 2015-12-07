package books

import (
//	"fmt"	
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"	
	"log"
	"io"
	"github.com/ddo/pick"
	"golang.org/x/net/html/charset"
)

//------------ Объявление типов и глобальных переменных

// структура задания с информацией по книге
type TaskerBook struct {
	Url string // ссылка на источник данных
	Book
	Tasker
}

//// структура книги
type Book struct {
	name          string // название книги
	autor         string // автор
	year          int    // год издания
	kolpages      int    // кол-во стрниц
	ves           int    // вес книги
	price         int    // цена для всех (обычная)
	pricediscount int    // цена со скидкой которая видна
}

// задание-триггер для срабатывания оповещения
type Tasker struct {
	uslovie string // условие < , > , =
	price   int    // цена триггера
	result  bool   // результат срабатывания триггера, если true , то триггер сработал
}

var LogFile *log.Logger

//------------ END Объявление типов и глобальных переменных

//проверка триггеров по массиву полученных данных по книгах

// проверки триггеров TaskerBook
func TriggerBookisUslovie(tb []TaskerBook) []TaskerBook {
	for i := 0; i < len(tb); i++ {
		tb[i].isTrue(tb[i].Book)
	}
	return tb
}

// ---------------  парсинг магазина Лабиринт
//получение данных книги из магазина лабиринт по урлу url
func (dbook *Book) Getlabirint(url string) {
	if url == "" {
		return
	}	
	body := gethtmlpage(url)
	shtml := string(body)
	scena, _ := pick.PickText(&pick.Option{ // текст цены книги
		&shtml,
		"span",
		&pick.Attr{
			"itemprop",
			"price",
		},
	})
	
	scenaskidka, _ := pick.PickText(&pick.Option{ // текст цены со скидкой книги
		&shtml,
		"span",
		&pick.Attr{
			"class",
			"buying-pricenew-val-number",
		},
	})

	sauthor, _ := pick.PickText(&pick.Option{ // текст описания книги
		&shtml,
		"span",
		&pick.Attr{
			"itemtype",
			"http://schema.org/ItemList",
		},
	})

	stitle, _ := pick.PickText(&pick.Option{&shtml, "span", &pick.Attr{"itemprop", "name"}})

	for i := 0; i < len(sauthor); i++ {
		switch sauthor[i] {
		case "Автор(ы)":
			dbook.autor = sauthor[i+1]
		case "Масса":
			dbook.ves, _ = strconv.Atoi(sauthor[i+1])
		case "Количество страниц":
			dbook.kolpages, _ = strconv.Atoi(sauthor[i+1])
		}
	}

	dbook.name = stitle[1]
	if len(scenaskidka) > 0 {
		dbook.pricediscount, _ = strconv.Atoi(scenaskidka[0])
	} else {
		dbook.pricediscount = 0
	}
	vv := strings.Split(scena[0], " ")
	dbook.price, _ = strconv.Atoi(vv[1])
	return
}

// --------------- END  парсинг магазина Лабиринт

// -----------  функции для Book

func (book0 *Book) Print() {
	LogFile.Println("Автор: ", book0.autor)
	LogFile.Println("Название книги: ", book0.name)
	LogFile.Println("Вес: ", book0.ves)
	LogFile.Println("Кол-во страниц: ", book0.kolpages)
	LogFile.Println("Цена: ", book0.price)
	LogFile.Println("Цена со скидкой: ", book0.pricediscount)
	return
}

//сохранить данные Book в файл
func (db *Book) Savetocsvfile(namef string) error {
	var fileflag bool = false
	if _, err := os.Stat(namef); os.IsNotExist(err) {
		// path/to/whatever does not exist
		fileflag = true
	}

	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// handle the error here
		return err
	}
	defer file.Close()
	if fileflag { // если не существует файл
		stitle := "Дата выгрузки;Автор;Название книги;Год издания;Кол-во стр.;Вес;Цена;Цена со скидкой;Ссылка" + "\n"
		file.WriteString(stitle)
	}
	curdate := time.Now().String()
	str := curdate + ";" + db.autor + ";" + db.name + ";" + strconv.Itoa(db.year) + ";" + strconv.Itoa(db.kolpages) + ";" + strconv.Itoa(db.ves) + ";" + strconv.Itoa(db.price) + ";" + strconv.Itoa(db.pricediscount) + "\n"
	//";"+db.url+
	file.WriteString(str)
	return err	
}

// ----------- END функции для Book



// -----------  функции для Tasker

//проверка триггеров по массиву полученных данных по книгах
func (task *Tasker) isTrue(book0 Book) {
	var res1, res2 bool
	res1 = false
	res2 = false
	switch task.uslovie {
	case ">":
		res1 = book0.price > task.price
	case "=":
		res1 = book0.price == task.price
	case "<":
		res1 = book0.price < task.price
	default:
		res1 = false
	}
	if book0.pricediscount > 0 { // если цена со скидкой больше нуля, то проверяем триггер на скидку
		switch task.uslovie {
		case ">":
			res2 = book0.pricediscount > task.price
		case "=":
			res2 = book0.pricediscount == task.price
		case "<":
			res2 = book0.pricediscount < task.price
		default:
			res2 = false
		}
	}
	task.result = res1 || res2
}

// ----------- END  функции для Tasker

// -----------  функции для TaskerBook

func (tb *TaskerBook) Print() {
	tb.Book.Print()
	LogFile.Println("Ссылка на книгу: ",tb.Url)
	return
}

// если тригер сработал то возвращает строку сообщения, иначе пусто
func (task *TaskerBook) Genmessage() string {
	var sprice, spricedisc string
	var smegtrigger, smegtrigger0, smsg string
	smsg = ""
	if task.result {
		b := task.Book
		sprice = strconv.Itoa(b.price)
		spricedisc = strconv.Itoa(b.pricediscount)
		smegtrigger = "Сбработал триггер по книге: \n\n" + "Автор: " + b.autor + "\n" + "Название: " + b.name + "\n" + "Цена: " + sprice + "\n" + "Цена со скидкой: " + spricedisc + "Ссылка: " + task.Url + "\n\n"
		sprice = strconv.Itoa(task.Tasker.price)
		smegtrigger0 = "Условие триггера: " + task.uslovie + "\n Цена триггера: " + sprice + "\n\n"
		smsg = smegtrigger + smegtrigger0
	}
	return smsg
}

//отправка сообщения если сработал триггер адресат toaddr
func (task *TaskerBook) Sendmail(toaddr string) {
	smsg := task.Genmessage()
	if smsg != "" {
		sendmailyandex("сработал триггер", smsg, toaddr)
	}
	return
}

// ----------- END  функции для TaskerBook

//---------------- общие функции ---------------------

//отправка почты через яндекс темой stema сообщение smsg адресату toaddr
func sendmailyandex(stema, smsg, toaddr string) bool {
	auth := smtp.PlainAuth("", "magazinebot@yandex.ru", "qwe123!!", "smtp.yandex.ru")
	to := []string{toaddr}
	msg := []byte("To: " + toaddr + "\r\n" +
		"Subject: " + stema + " \r\n" +
		"\r\n" +
		smsg + "\r\n")
	err := smtp.SendMail("smtp.yandex.ru:25", auth, "magazinebot@yandex.ru", to, msg)
	if err != nil {
		panic(err)
	}
	return true
}

//получение страницы из урла url
func gethtmlpage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		LogFile.Println("HTTP error:", err)
		panic("HTTP error")
	}
	defer resp.Body.Close()
	// вот здесь и начинается самое интересное
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		LogFile.Println("Encoding error:", err)
		panic("Encoding error")
	}
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		LogFile.Println("IO error:", err)
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

// инициализация файла логов
func InitLogFile(namef string) *log.Logger {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}
	multi := io.MultiWriter(file, os.Stdout)
	LFile := log.New(multi, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	return LFile
}

// чтение из текстового конфиг файла заданий и возращает массив заданий Tasker
func Readtaskerbookcfg(namef string) []TaskerBook {
	var res []TaskerBook
	str := readfiletxt(namef)
	vv := strings.Split(str, "\n")
	for i := 0; i < len(vv); i++ {
		s := strings.Split(vv[i], ";")
		tt, _ := strconv.Atoi(s[2])
		dt := Tasker{uslovie: s[1], price: tt, result: false}
		t := TaskerBook{Url: s[0]}
		t.Tasker = dt
		res = append(res, t)
	}
	return res
}


// вызов парсинга книжного магазина
func RunBooks(namestore string,toaddr string) {
	//---- инициализация переменных	
	var list_tasker []TaskerBook
	
	namefurls := namestore + "-url.cfg"
	namelogfile := namestore + ".log"
//---- END инициализация переменных		

	LogFile = InitLogFile(namelogfile) // инициализация лог файла
	LogFile.Println("Starting programm")
	
	LogFile.Println("Имя магазина store: ",namestore)
	LogFile.Println("Э/почта для отправки уведомлений: ",toaddr)	

	// получаем задания из файла
	list_tasker = Readtaskerbookcfg(namefurls)
	
	//получение данных книжек
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].Getlabirint(list_tasker[i].Url)
		namef := namestore + ".csv"
		list_tasker[i].Savetocsvfile(namef)
		list_tasker[i].Print()
	}

	//проверка на наличии срабатываний
	list_tasker = TriggerBookisUslovie(list_tasker)

	for i := 0; i < len(list_tasker); i++ {
		LogFile.Println(list_tasker[i].Genmessage())
		list_tasker[i].Sendmail(toaddr)
	}

	LogFile.Println("The end....!\n")
}