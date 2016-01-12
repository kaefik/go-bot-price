// go-bot-labirint
// программа скачивает по ссылкам данные по книгам Лабиринт и проверяет условия по ссылкам
// Автор: Ильнур Сайфутдинов
// email: ilnursoft@gmail.com
// декабрь 2015

// отдельный пакет pkg/books
// изм 07.12.2015

package main

import (
	//	"fmt"
	"flag"
	"go-bot-price/pkg/tovar"
)

////------------ Объявление типов и глобальных переменных

var (
	store  string
	toaddr string
	hd     string
)

//------------ END Объявление типов и глобальных переменных

// функция парсинга аргументов программы
func parse_args() bool {
	flag.StringVar(&store, "store", "", "Название магазина по которому будут мониторить цены.")
	flag.StringVar(&toaddr, "toaddr", "", "Э/почта для отправки сообщений срабатываний триггера.")
	flag.StringVar(&hd, "hd", "", "Рабочая папка где нах-ся триггеры и будут выводится результаты работы.")
	flag.Parse()
	if store == "" {
		store = "labirint"
	}
	if toaddr == "" {
		toaddr = "ilnursoft@gmail.com"
	}
	if hd == "" {
		hd = "oilnur"
	}
	return true
}

//---------------- END общие функции ---------------------

func main() {

	if !parse_args() {
		return
	}

	tovar.Homedirs = hd

	///-------- для теста
	//		store="mvideo"
	///-------- END для теста

	tovar.RunTovar(store, toaddr)

}
