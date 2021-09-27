package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)
//Кнопки для выбора режима
var setRegime= tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1⃣","trackingRegime1"),
		tgbotapi.NewInlineKeyboardButtonData("2⃣","trackingRegime2"),
		tgbotapi.NewInlineKeyboardButtonData("3⃣","trackingRegime3"),
	),
)
//Вступительное сообщение
var (
	helloMessage  = `Я слежу за курсом Bitcoin'a, а также последними новостями и свежей аналитикой динамики роста.
/news - последние новости
/analytics - аналитика изменения стоимости криптовалюты
/rate - текущий курс Bitcoin'a

Также вы можете настроить систему мониторинга стоимость биткоина, а именно уведомления о цене и её изменении за последние 24 часах, 3-х типов:
1) уведомления раз в час
2) уведомление при достижение валютой указанной вами цены
3) уведомление при росте валюты за 24 часа большем, чем вы указали

/tracking - настроить оповещения
/stop_tracking - отменить все оповещения`
	trackingMessage = `Здесь вы можете настроить оповещения.
		На выбор предоставляются три режима:
		1⃣. Ежечасное оповещение о цене
		2⃣. Оповещение при достижении определенной цены
		3⃣. Оповещение при резком подъеме криптовалюты
		Для включения/настройки одного из режимов нажмите на кнопку с соответствующим номером`
)

//___________________________________
//Обработка Команд
//___________________________________
func isCommandCase(update *tgbotapi.Update,bot *tgbotapi.BotAPI,db *sql.DB,bitcoinNow infoBitcoin,allNews , allAnalysis []string,
errBitcoinInfoPars ,errNewsPars,errAnalysisPars *error){
	var err error
	cmdText := update.Message.Command()
	switch cmdText {
	case "help" :
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
			helloMessage)
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
	case "start":
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
			helloMessage)
		err=addNewUser(int(update.Message.Chat.ID),true,db)
		errorsWorkDB(InfoBitcoinDB,addInfo,err)
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
	case "tracking":
		msgText:=trackingMessage
		if *errBitcoinInfoPars!=nil{
			msgText+="\nВ данный момент имеются проблемы с доступом к сайту биржи, вы начнете получать оповещения, как только проблема будет решена."
		}
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
			msgText)
		//добавление кнопок
		msg.ReplyMarkup=setRegime
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
	case "rate":
		var msgText string
		if *errBitcoinInfoPars==nil{
			if bitcoinNow.isIncrease{
				msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подорожал на %v💲 (+%v%%)",
					bitcoinNow.Cost,bitcoinNow.changeCostUSD,bitcoinNow.changeCostPr)
			}else {
				msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подешевел на %v💲 (%v%%)",
					bitcoinNow.Cost,bitcoinNow.changeCostUSD*(-1),bitcoinNow.changeCostPr)
			}
		}else{
			msgText="В данный момент имеются проблемы с доступом к сайту биржи, " +
				"вы сможете ознакомиться с курсом криптовалюты, как только проблема будет решена." +
				"Приносим извинения за доставленные неудобства."
		}
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,msgText)
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
	case "news":
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"🔊 Новости о Bitcoin`e:\n")
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
		if *errNewsPars==nil{
			for _,news:= range allNews{
				msg =tgbotapi.NewMessage(update.Message.Chat.ID,news)
				_,err =bot.Send(msg)
				if err!=nil{
					errorsMessage(placeMessageCommand,err,msg,db)
				}
			}
		}else{
			msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"В данный момент имеются проблемы с доступом к сайту биржи, " +
				"вы сможете ознакомиться с последними новостями, как только проблема будет решена." +
				"Приносим извинения за доставленные неудобства.")
			_,err =bot.Send(msg)
			if err!=nil{
				errorsMessage(placeMessageCommand,err,msg,db)
			}
		}
	case "analytics":
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"🎓 Аналитика динамики Bitcoin`a:\n")
		_,err = bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
		if *errAnalysisPars==nil{
			for _,analysis:= range allAnalysis{
				msg :=tgbotapi.NewMessage(update.Message.Chat.ID,analysis)
				_,err =bot.Send(msg)
				if err!=nil{
					errorsMessage(placeMessageCommand,err,msg,db)
				}
			}
		}else{
			msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"В данный момент имеются проблемы с доступом к сайту биржи, " +
				"вы сможете ознакомиться с последними прогнозами на динамику криптовалюты, как только проблема будет решена." +
				"Приносим извинения за доставленные неудобства.")
			_,err = bot.Send(msg)
			if err!=nil{
				errorsMessage(placeMessageCommand,err,msg,db)
			}
		}
	case "stop_tracking":
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"❌ Оповещения отключены ")
		err=changeInformation(int(update.Message.Chat.ID),false,db)
		errorsWorkDB(InfoBitcoinDB,changeInfo,err)
		err=deleteUserChatIdCostDB(int(update.Message.Chat.ID),db)
		errorsWorkDB(ChatIdCostDB,deleteInfo,err)
		err=deleteUserChatIdChangeCostDB(int(update.Message.Chat.ID),db)
		errorsWorkDB(ChatIdChangeCostDB,deleteInfo,err)
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageCommand,err,msg,db)
		}
	}
}
//___________________________________
//Обработка Callback ответов
//___________________________________
func isCallbackQuery(update *tgbotapi.Update,bot *tgbotapi.BotAPI,db *sql.DB,status map[int64]string){
	var err error
	switch update.CallbackQuery.Data {
	case "trackingRegime1":
		msg :=tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,"✔ Режим отслеживания включен")
		err=changeInformation(int(update.CallbackQuery.Message.Chat.ID),true,db)
		errorsWorkDB(InfoBitcoinDB,changeInfo,err)
		_,err :=bot.Send(msg)
		if err!=nil{
			errorsMessage(placeCallbackQuery,err,msg,db)
		}
	case "trackingRegime2":
		msg :=tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,"Укажите стоимость Bitcoin`a, о которой нужно сообщить "+
			"для этого отправьте число в формате: '123.456'")
		//status
		status[update.CallbackQuery.Message.Chat.ID]="setCost"
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeCallbackQuery,err,msg,db)
		}
	case "trackingRegime3":
		msg :=tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,"Укажите величину изменения стоимости Bitcoin`a, о которой нужно сообщить "+
			"для этого отправьте число в формате: '123.456'")
		//status
		status[update.CallbackQuery.Message.Chat.ID]="setChangeCost"
		_,err :=bot.Send(msg)
		if err!=nil{
			errorsMessage(placeCallbackQuery,err,msg,db)
		}
	}
}
//____________________________________
//Обработка обычных сообщений
//___________________________________
func isUsualMessage(update *tgbotapi.Update,bot *tgbotapi.BotAPI,db *sql.DB,status map[int64]string){
	var err error
	switch status[update.Message.Chat.ID] {
	case "setCost":// установление цены для оповещений
		cost,ok:=validAndPrepare(update.Message.Text)
		var msg tgbotapi.MessageConfig
		if ok{
			msg =tgbotapi.NewMessage(update.Message.Chat.ID,"Когда стоимость Bitcoin`a достигнет "+strconv.FormatFloat(cost,'f',2,64)+" USD, "+
				"вы получите уведомление.")
			//добавление в базу данных
			err=addUserCostDB(int(update.Message.Chat.ID),cost,db)
			errorsWorkDB(ChatIdCostDB,addInfo,err)
			delete(status,update.Message.Chat.ID)
		}else{
			msg =tgbotapi.NewMessage(update.Message.Chat.ID,"Некорректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
		}
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageNotCommand,err,msg,db)
		}
	case "setChangeCost": // установление изменения цены для оповещений
		changeCost,ok:=validAndPrepare(update.Message.Text)
		var msg tgbotapi.MessageConfig
		if ok{
			msg =tgbotapi.NewMessage(update.Message.Chat.ID,"Когда скачок цены Bitcoin`a окажется больше чем "+strconv.FormatFloat(changeCost,'f',2,64)+" USD, "+
				"вы получите уведомление.")
			err=addUserChangeCostDB(int(update.Message.Chat.ID),changeCost,db)
			errorsWorkDB(ChatIdChangeCostDB,addInfo,err)
			delete(status,update.Message.Chat.ID)
		}else{
			msg =tgbotapi.NewMessage(update.Message.Chat.ID,"Некорректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
		}
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageNotCommand,err,msg,db)
		}
	default:
		msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
			helloMessage)
		_,err =bot.Send(msg)
		if err!=nil{
			errorsMessage(placeMessageNotCommand,err,msg,db)
		}
	}
}
