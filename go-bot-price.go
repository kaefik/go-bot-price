// go-bot-labirint
// программа скачивает по ссылкам данные по книгам Лабиринт и проверяет условия по ссылкам
// Автор: Ильнур Сайфутдинов
// email: ilnursoft@gmail.com
// декабрь 2015

// отдельный пакет pkg/books
// изм 07.12.2015

package main

import (
	"fmt"
	"flag"
	"go-bot-price/pkg"
	"go-bot-price/pkg/tovar"
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
	
	var tv tovar.TaskerTovar
	
	fmt.Println(tv)
	
	if !parse_args() {
	   return
 	}	
	
	///-------- для теста
	
	store="ulmart"
			
	///-------- END для теста
	
	
	
	switch store {	
		case "labirint": books.RunBooks(store,toaddr) // вызов парсинга книжного магазина
		default: tovar.RunTovar(store,toaddr) // вызов парсинга магазина электроники
	}

}
