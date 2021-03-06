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
	"fmt"
	"go-bot-price/pkg/tovar"
	"os"
	//	"strings"
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

//возвращает список имен файлов в директории dirname
func Getlistfileindirectory(dirname string) []string {
	listfile := make([]string, 0)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			//fmt.Println(fi.Name(), fi.Size(), "bytes")
			listfile = append(listfile, fi.Name())
		}
	}
	return listfile
}

// получение массив имен магазинов которые получаются из имен файлов формата: названиемагазина-url.cfg
func GetNameStoreFromFilename(nf []string) []string {
	reslist := make([]string, 0)
	for _, v := range nf {
		if v[len(v)-3:] == "cfg" {
			reslist = append(reslist, v[:len(v)-8])
		}
	}
	return reslist
}

//---------------- END общие функции ---------------------

func main() {

	if !parse_args() {
		return
	}

	tovar.Homedirs = hd

	listff := Getlistfileindirectory(hd)
	liststore := GetNameStoreFromFilename(listff)

	///-------- для теста
	//		store="mvideo"
	///-------- END для теста

	for i := 0; i < len(liststore); i++ {
		store = liststore[i]
		tovar.RunTovar(store, toaddr)
	}

}
