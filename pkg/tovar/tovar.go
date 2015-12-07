// tovar
package tovar


// структура задания с информацией по книге
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

// вызов парсинга книжного магазина
func RunTovar(namestore string,toaddr string) {
	//---- инициализация переменных	
//	var list_tasker []TaskerTovar
	
	namefurls := namestore + "-url.cfg"
	namelogfile := namestore + ".log"
////---- END инициализация переменных		

//	LogFile = InitLogFile(namelogfile) // инициализация лог файла
//	LogFile.Println("Starting programm")
	
//	LogFile.Println("Имя магазина store: ",namestore)
//	LogFile.Println("Э/почта для отправки уведомлений: ",toaddr)	

//	// получаем задания из файла
//	list_tasker = Readtaskerbookcfg(namefurls)
	
//	//получение данных книжек
//	for i := 0; i < len(list_tasker); i++ {
//		list_tasker[i].Getlabirint(list_tasker[i].Url)
//		namef := namestore + ".csv"
//		list_tasker[i].Savetocsvfile(namef)
//		list_tasker[i].Print()
//	}

//	//проверка на наличии срабатываний
//	list_tasker = TriggerBookisUslovie(list_tasker)

//	for i := 0; i < len(list_tasker); i++ {
//		LogFile.Println(list_tasker[i].Genmessage())
//		list_tasker[i].Sendmail(toaddr)
//	}

//	LogFile.Println("The end....!\n")
}