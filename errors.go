package main


import "log"

//места возможного возникновения ошибок в коде с парсингом сайта
var (
	placeCallbackQuery ="CallbackQuery message"
	placeMessageCommand = "MessageCommand message"
	placeMessageNotCommand = "MessageNotCommand message"
	placeSendMessageAboutBitcoin = "SendMessageAboutBitcoin message"
)
func errorsMessage(place string,err error){
		log.Println(place,": ",err)
}
//места возможного возникновения ошибок в коде с базой данных
var (
	InfoBitcoinDB ="problem in work InfoBitcoinDB "
	ChatIdCostDB ="problem in work ChatIdCostDB "
	ChatIdChangeCostDB ="problem in work ChatIdChangeCostDB "
)
//операции с базой данных
var (
	changeInfo = "changeInfo function not work"
	giveInfo = "giveInfo function not work"
	addInfo = "addInfo function not work"
	deleteInfo= "deleteInfo function not work"
)
func errorsWorkDB(place string,operation string ,err error){
	log.Println(place,operation,": ",err)
}