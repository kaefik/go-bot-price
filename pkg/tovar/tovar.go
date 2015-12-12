// tovar
package tovar

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
	"go-bot-price/pkg/pick"
	"golang.org/x/net/html/charset"
)

// структура задания с информацией по товару
type TaskerTovar struct {
	Url string // ссылка на источник данных
	Tovar
	Tasker
}

//// структура книги
type Tovar struct {
	name          string // название товара
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
func TriggerisUslovie(tb []TaskerTovar) []TaskerTovar {
	for i := 0; i < len(tb); i++ {
		tb[i].isTrue(tb[i].Tovar)
	}
	return tb
}

//// -----------  функции для Tovar

// вывод  в файл лога и на экран информации о товаре
func (book0 *Tovar) Print() {
	LogFile.Println("Название товара: ", book0.name)
	LogFile.Println("Цена: ", book0.price)
	LogFile.Println("Цена со скидкой: ", book0.pricediscount)
	return
}

//сохранить данные Book в файл
func (db *Tovar) Savetocsvfile(namef string) error {
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
		stitle := "Дата выгрузки;Название товара;Цена;Цена со скидкой;Ссылка" + "\n"
		file.WriteString(stitle)
	}
	curdate := time.Now().String()
	str := curdate + ";" + db.name + ";" + strconv.Itoa(db.price) + ";" + strconv.Itoa(db.pricediscount) + "\n"
	//";"+db.url+
	file.WriteString(str)
	return err
}

//// ----------- END функции для Tovar

//// -----------  функции для Tasker

//проверка триггеров по массиву полученных данных по товарам
func (task *Tasker) isTrue(book0 Tovar) {
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

//// ----------- END  функции для Tasker

//// -----------  функции для TaskerTovar

//func (tb *TaskerBook) Print() {
//	tb.Book.Print()
//	LogFile.Println("Ссылка на книгу: ",tb.Url)
//	return
//}

//// если тригер сработал то возвращает строку сообщения, иначе пусто
func (task *TaskerTovar) Genmessage() string {
	var sprice, spricedisc string
	var smegtrigger, smegtrigger0, smsg string
	smsg = ""
	if task.result {
		b := task.Tovar
		sprice = strconv.Itoa(b.price)
		spricedisc = strconv.Itoa(b.pricediscount)
		smegtrigger = "Сбработал триггер по товару: \n" + "Название: " + b.name + "\n" + "Цена: " + sprice + "\n" + "Цена со скидкой: " + spricedisc + "\n" + "Ссылка: " + task.Url + "\n\n"
		sprice = strconv.Itoa(task.Tasker.price)
		smegtrigger0 = "Условие триггера: " + task.uslovie + "\n Цена триггера: " + sprice + "\n\n"
		smsg = smegtrigger + smegtrigger0
	}
	return smsg
}

//отправка сообщения если сработал триггер адресат toaddr
func (task *TaskerTovar) Sendmail(toaddr string) {
	smsg := task.Genmessage()
	if smsg != "" {
		sendmailyandex("сработал триггер", smsg, toaddr)
	}
	return
}

//// ----------- END  функции для TaskerTovar

////---------------- общие функции ---------------------

//удаление пробелов из строки s
func delspacefromstring(s string) string{
	r:=[]rune(s)
	rnew:=make([]rune,0)
	for i:=0;i<len(r);i++{
		if (r[i]!=160) && (r[i]!=32) { 
			rnew=append(rnew,r[i])
		}
	}
	return string(rnew)
}

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

//// чтение файла с именем namefи возвращение содержимое файла, иначе текст ошибки
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
func Readtaskercfg(namef string) []TaskerTovar {
	var res []TaskerTovar
	str := readfiletxt(namef)
	vv := strings.Split(str, "\n")
	for i := 0; i < len(vv); i++ {
		s := strings.Split(vv[i], ";")
		tt, _ := strconv.Atoi(s[2])
		dt := Tasker{uslovie: s[1], price: tt, result: false}
		t := TaskerTovar{Url: s[0]}
		t.Tasker = dt
		res = append(res, t)
	}
	return res
}

// предварительная обработка
func RunTovarPre(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//---- инициализация переменных
	//	var list_tasker []TaskerTovar

	namefurls := namestore + "-url.cfg"
	namelogfile := namestore + ".log"
	////---- END инициализация переменных

	LogFile = InitLogFile(namelogfile) // инициализация лог файла
	LogFile.Println("Starting programm")

	LogFile.Println("Имя магазина store: ", namestore)
	LogFile.Println("Э/почта для отправки уведомлений: ", toaddr)

	// получаем задания из файла
	list_tasker = Readtaskercfg(namefurls)	
	return list_tasker
}

// окончательная обработка
func RunTovarEnd(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//проверка на наличии срабатываний
	list_tasker = TriggerisUslovie(list_tasker)

	for i := 0; i < len(list_tasker); i++ {
		LogFile.Println(list_tasker[i].Genmessage())
		list_tasker[i].Sendmail(toaddr)
	}
	LogFile.Println("The end....!\n")
	return list_tasker
}

// вызов парсинга книжного магазина  - сюда добавляем вызов для новых магазинов для которых можем парсить
func RunTovar(namestore string, toaddr string) {
	//---- инициализация переменных
	var list_tasker []TaskerTovar

	list_tasker = RunTovarPre(list_tasker, namestore, toaddr)

	switch namestore {
	case "eldorado":
		list_tasker = RunTovarGetDataEldorado(list_tasker, namestore, toaddr)
	case "dns":
		list_tasker = RunTovarGetDataDns(list_tasker, namestore, toaddr)
	case "ulmart":
		list_tasker = RunTovarGetDataUlmart(list_tasker, namestore, toaddr)
	case "citilink":
		list_tasker = RunTovarGetDataCitilink(list_tasker, namestore, toaddr)		
	case "mvideo":
		list_tasker = RunTovarGetDataMvideo(list_tasker, namestore, toaddr)		
	case "aliexpress":
		list_tasker = RunTovarGetDataAliexpress(list_tasker, namestore, toaddr)						
	default:
		return
	}
	RunTovarEnd(list_tasker, namestore, toaddr)

}

// ---------------  парсинг магазинов ( состоит из двух функций RunTovarGetDataНазваниеМагазина  и GetdataTovarfromНазваниеМагазин

//получение данных товаров из магазина Эльдорадо
func RunTovarGetDataEldorado(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//получение данных товара
	namef := namestore + ".csv"
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].GetdataTovarfromEldorado(list_tasker[i].Url) // <-- тут меняем на нужную функцию парсинга
		list_tasker[i].Savetocsvfile(namef)
		list_tasker[i].Print()
	}
	return list_tasker
}

//получение данных товара из магазина Эльдорадо по урлу url
func (this *Tovar) GetdataTovarfromEldorado(url string) {
	var ss []string
	if url == "" {
		return
	}
	body := gethtmlpage(url)
	shtml := string(body)

	sname, _ := pick.PickText(&pick.Option{ // текст цены книги
		&shtml,
		"div",
		&pick.Attr{
			"class",
			"q-fixed-name no-mobile",
		},
	})

	for i := 0; i < len(sname); i++ {
		if strings.TrimSpace(sname[i]) != "" { // удаление пробелов
			ss = append(ss, sname[i])
		}
	}

	this.name = ss[0]

	sprice, _ := pick.PickText(&pick.Option{&shtml, "span", &pick.Attr{"itemprop", "price"}})

	ss = make([]string, 0)
	for i := 0; i < len(sprice); i++ {
		if strings.TrimSpace(sprice[i]) != "" { // удаление пробелов
			ss = append(ss, sprice[i])
		}
	}

	if len(ss) > 0 {
		this.price, _ = strconv.Atoi(ss[0])
	}

	return
}

//------ парсинг ДНС
//получение данных товаров из магазина ДНС
func RunTovarGetDataDns(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//получение данных товара
	namef := namestore + ".csv"
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].GetdataTovarfromDns(list_tasker[i].Url) // <-- тут меняем на нужную функцию парсинга
		if _, err := os.Stat(namestore); os.IsNotExist(err) {
			os.Mkdir(namestore,0666)
		}
		list_tasker[i].Savetocsvfile(namestore+"\\"+namef)	
		list_tasker[i].Print()
	}
	return list_tasker
}

//получение данных товара из магазина ДНС по урлу url
func (this *Tovar) GetdataTovarfromDns(url string) {
	var ss []string
	if url == "" {
		return
	}
	body := gethtmlpage(url)
	shtml := string(body)

	sname, _ := pick.PickText(&pick.Option{ // текст цены книги
		&shtml,
		"h1",
		&pick.Attr{
			"class",
			"page-title price-item-title",
		},
	})

	for i := 0; i < len(sname); i++ {
		if strings.TrimSpace(sname[i]) != "" { // удаление пробелов
			ss = append(ss, sname[i])
		}
	}
	this.name = ss[0]

   //<meta itemprop="price" content="3190.00" />

	sprice, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"itemprop", "price"}}, "content")
	
	if len(sprice)>0 {
		sprice1:=strings.Split(sprice[0],".")		
		this.price, _ = strconv.Atoi(sprice1[0])	
	}
	
	

	return
}

//------ парсинг Юлмарт
//получение данных товаров из магазина Юлмарт
func RunTovarGetDataUlmart(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//получение данных товара
	namef := namestore + ".csv"
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].GetdataTovarfromUlmart(list_tasker[i].Url) // <-- тут меняем на нужную функцию парсинга
		if _, err := os.Stat(namestore); os.IsNotExist(err) {
			os.Mkdir(namestore,0666)
		}
		list_tasker[i].Savetocsvfile(namestore+"\\"+namef)	
		list_tasker[i].Print()
	}
	return list_tasker
}

//получение данных товара из магазина Юлмарт по урлу url
func (this *Tovar) GetdataTovarfromUlmart(url string) {
	var ss []string
	if url == "" {
		return
	}
	body := gethtmlpage(url)
	shtml := string(body)
			
//	<meta name="keywords" content="Подгузники Mepsi L (9-16 кг), 38 шт, Mepsi, артикул 3421456"/>
	sname, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"name", "keywords"}}, "content")

	for i := 0; i < len(sname); i++ {
		if strings.TrimSpace(sname[i]) != "" { // удаление пробелов
			ss = append(ss, sname[i])
		}
	}
	this.name = ss[0]

	//    <meta itemprop="price" content="660">
	sprice, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"itemprop", "price"}}, "content")	
	if len(sprice)>0 {
		sprice1:=strings.Split(sprice[0],".")		
		this.price, _ = strconv.Atoi(sprice1[0])	
	}
	return
}

//------ парсинг Cитилинк
//получение данных товаров из магазина Cитилинк
func RunTovarGetDataCitilink(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//получение данных товара
	namef := namestore + ".csv"
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].GetdataTovarfromCitilink(list_tasker[i].Url) // <-- тут меняем на нужную функцию парсинга
		if _, err := os.Stat(namestore); os.IsNotExist(err) {
			os.Mkdir(namestore,0666)
		}
		list_tasker[i].Savetocsvfile(namestore+"\\"+namef)	
		list_tasker[i].Print()
	}
	return list_tasker
}

//получение данных товара из магазина Cитилинк по урлу url
func (this *Tovar) GetdataTovarfromCitilink(url string) {
	if url == "" {
		return
	}
	body := gethtmlpage(url)
	shtml := string(body)
			
//	<meta itemprop="name" content="Подгузники MERRIES Large" />
	sname, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"itemprop", "name"}}, "content")
	this.name = sname[0]

	//    <meta itemprop="price" content="1540.00" />
	sprice, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"itemprop", "price"}}, "content")	
	if len(sprice)>0 {
		sprice1:=strings.Split(sprice[0],".")		
		this.price, _ = strconv.Atoi(sprice1[0])	
	}

	return
}

//------ парсинг МВидео
//получение данных товаров из магазина МВидео
func RunTovarGetDataMvideo(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//получение данных товара
	namef := namestore + ".csv"
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].GetdataTovarfromMvideo(list_tasker[i].Url) // <-- тут меняем на нужную функцию парсинга
		if _, err := os.Stat(namestore); os.IsNotExist(err) {
			os.Mkdir(namestore,0666)
		}
		list_tasker[i].Savetocsvfile(namestore+"\\"+namef)		
		list_tasker[i].Print()
	}
	return list_tasker
}

//получение данных товара из магазина МВидео по урлу url
func (this *Tovar) GetdataTovarfromMvideo(url string) {
	if url == "" {
		return
	}
	body := gethtmlpage(url)
	shtml := string(body)
			
//	<meta property="og:title" content="Ультрабук ASUS Zenbook UX32LA-R3094H"/>
	sname, _ := pick.PickAttr(&pick.Option{&shtml, "meta", &pick.Attr{"property", "og:title"}}, "content")
	this.name = sname[0]

	//    <strong class="product-price-current">43990</strong>
//	sprice, _ := pick.PickAttr(&pick.Option{&shtml, "strong", &pick.Attr{"class", "product-price-current"}},)	
	
	sprice, _ := pick.PickText(&pick.Option{ 
		&shtml,
		"strong",
		&pick.Attr{
			"class",
			"product-price-current",
		},
	})
	
	if len(sprice)>0 {
		sprice1:=strings.Split(sprice[0],".")		
		this.price, _ = strconv.Atoi(sprice1[0])	
	}

	return
}


//------ парсинг aliexpress
//получение данных товаров из магазина Aliexpress
func RunTovarGetDataAliexpress(list_tasker []TaskerTovar, namestore string, toaddr string) []TaskerTovar {
	//получение данных товара
	namef := namestore + ".csv"
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].GetdataTovarfromAliexpress(list_tasker[i].Url) // <-- тут меняем на нужную функцию парсинга
        if _, err := os.Stat(namestore); os.IsNotExist(err) {
			os.Mkdir(namestore,0666)
		}
		list_tasker[i].Savetocsvfile(namestore+"\\"+namef)
		list_tasker[i].Print()
	}
	return list_tasker
}

//получение данных товара из магазина Aliexpress по урлу url
func (this *Tovar) GetdataTovarfromAliexpress(url string) {
	if url == "" {
		return
	}
	body := gethtmlpage(url)
	shtml := string(body)
			
//	<h1 class="product-name" itemprop="name">
	sname, _ := pick.PickText(&pick.Option{ 
		&shtml,
		"h1",
		&pick.Attr{
			"class",
			"product-name",
		},
	})
	this.name = sname[0]

	//    <span id="sku-price" itemprop="price">56.99</span>

	sprice, _ := pick.PickText(&pick.Option{ 
		&shtml,
		"span",
		&pick.Attr{
			"itemprop",
			"price",
		},
	})
	
	fmt.Println(sprice)
	
	if len(sprice)>0 {
		sprice1:=strings.Split(sprice[0],",")	
		fmt.Println(sprice1)
		fmt.Println(delspacefromstring(sprice1[0]))
		
			
		this.price, _ = strconv.Atoi(delspacefromstring(sprice1[0]))	
	}

	return
}


//// --------------- END  парсинг магазинов
