package books

//import (
//	"fmt"
//)

// структура книги
type Book struct {
	name          string // название книги
	autor         string // автор
	year          int    // год издания
	kolpages      int    // кол-во стрниц
	ves           int    // вес книги
	price         int    // цена для всех (обычная)
	pricediscount int    // цена со скидкой которая видна
}