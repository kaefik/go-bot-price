// go-bot-labirint
// программа скачивает по ссылкам данные по книгам Лабиринт и проверяет условия по ссылкам
// Автор: Ильнур Сайфутдинов
// email: ilnursoft@gmail.com
// декабрь 2015

// отдельный пакет pkg/books
// изм 07.12.2015

package main

import (
	"flag"
	"go-bot-price/pkg"
)

////------------ Объявление типов и глобальных переменных

var (
	store string
	toaddr string
)

//------------ END Объявление типов и глобальных переменных

// функция парсинга аргументов программы
func parse_args() bool {
	flag.StringVar(&store, "store", "", "Название магазина по которому будут мониторить цены.")
	flag.StringVar(&toaddr, "toaddr", "", "Э/почта для отправки сообщений срабатываний триггера.")
	
	flag.Parse()
	if store == "" {		
		store="labirint"
	}	
	if toaddr == "" {
		toaddr="ilnursoft@gmail.com"
	} 
	return true
}

//---------------- END общие функции ---------------------

func main() {
	var list_tasker []books.TaskerBook
	
	if !parse_args() {
	   return
 	}
	
//---- инициализация переменных
	namestore := store
	namefurls := namestore + "-url.cfg"
	namelogfile := namestore + ".log"
//---- END инициализация переменных		

	books.LogFile = books.InitLogFile(namelogfile) // инициализация лог файла
	books.LogFile.Println("Starting programm")
	
	books.LogFile.Println("Имя магазина store: ",store)
	books.LogFile.Println("Э/почта для отправки уведомлений: ",toaddr)	

	// получаем задания из файла
	list_tasker = books.Readtaskerbookcfg(namefurls)
	
	//получение данных книжек
	for i := 0; i < len(list_tasker); i++ {
		list_tasker[i].Getlabirint(list_tasker[i].Url)
		namef := namestore + ".csv"
		list_tasker[i].Savetocsvfile(namef)
		list_tasker[i].Print()
	}

	//проверка на наличии срабатываний
	list_tasker = books.TriggerBookisUslovie(list_tasker)

	for i := 0; i < len(list_tasker); i++ {
		books.LogFile.Println(list_tasker[i].Genmessage())
		list_tasker[i].Sendmail(toaddr)
	}

	books.LogFile.Println("The end....!\n")
}
