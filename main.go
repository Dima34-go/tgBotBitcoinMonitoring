package main


import (
"database/sql"
"fmt"
"github.com/go-telegram-bot-api/telegram-bot-api"
"log"
"strconv"
"time"
)
//инлайн кнопки для выбора  режима трекинга
var setRegime= tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1⃣","trackingRegime1"),
		tgbotapi.NewInlineKeyboardButtonData("2⃣","trackingRegime2"),
		tgbotapi.NewInlineKeyboardButtonData("3⃣","trackingRegime3"),
	),
)
// Вступительное сообщение
var helloMessage  =`Я слежу за курсом Bitcoin'a, а также последними новостями и свежей аналитикой динамики роста.
/news - последние новости
/analytics - аналитика изменения стоимости криптовалюты
/rate - текущий курс Bitcoin'a

Также вы можете настроить систему мониторинга стоимость биткоина, а именно уведомления о цене и её изменении за последние 24 часах, 3-х типов:
1) уведомления раз в час
2) уведомление при достижение валютой указанной вами цены
3) уведомление при росте валюты за 24 часа большем, чем вы указали

/tracking - настроить оповещения
/stop_tracking - отменить все оповещения`

func main() {
	//Запуск бота и настройка обновлений
	botToken:="1940781427:AAHHOpliPiUV0SQcKrLRQ6R9Ytkv1JZf5YE"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	//Запуск базы данных
	db, err := sql.Open("mysql", "root:password@/userschat")
	if err != nil {
		log.Panic(err)
	}
	defer func(){
		err=db.Close()
		if err!=nil{
			log.Println("Ошибка закрытия БД: ",err)
		}
	}()
	//массив с данными о состоянии пользователя
	status := make(map[int64]string,0)
	bitcoinNow,err:=parsAllInfoAboutBitcoin()
	if err!=nil{
		log.Println("Проблема с парсингом информации о биткойне: ",err)
	}
	allNews,err:=parsNews()
	if err!=nil{
		fmt.Println(err)
	}
	allAnalysis,err:=parsAnalysis()
	if err!=nil{
		fmt.Println(err)
	}
	//--------------------------------получение update в бесконечном цикле, вынесенное в отдельную горутину
	go func() {
		for update := range updates {
			if update.Message == nil && update.CallbackQuery == nil{//пустое обновление
				continue
			}else if update.CallbackQuery != nil{//обработка callback ответа
				switch update.CallbackQuery.Data {
				case "trackingRegime1":
					msg :=tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,"✔ Режим отслеживания включен")
					err=changeInformation(int(update.CallbackQuery.Message.Chat.ID),true,db)
					errorsWorkDB(InfoBitcoinDB,changeInfo,err)
					_,err =bot.Send(msg)
					errorsMessage(placeCallbackQuery,err)
					continue
				case "trackingRegime2":
					msg :=tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,"Укажите стоимость Bitcoin`a, о которой нужно сообщить "+
						"для этого отправьте число в формате: '123.456'")
					//status
					status[update.CallbackQuery.Message.Chat.ID]="setCost"
					_,err =bot.Send(msg)
					errorsMessage(placeCallbackQuery,err)
					continue
				case "trackingRegime3":
					msg :=tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,"Укажите величину изменения стоимости Bitcoin`a, о которой нужно сообщить "+
						"для этого отправьте число в формате: '123.456'")
					//status
					status[update.CallbackQuery.Message.Chat.ID]="setChangeCost"
					_,err =bot.Send(msg)
					errorsMessage(placeCallbackQuery,err)
					continue
				}
			}else if update.Message.IsCommand(){//обработка команд
				cmdText := update.Message.Command()
				switch cmdText {
				case "help" :
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
						helloMessage)
					_,err =bot.Send(msg)
					errorsMessage(placeMessageCommand,err)
				case "start":
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
						helloMessage)
					err=addNewUser(int(update.Message.Chat.ID),true,db)
					errorsWorkDB(InfoBitcoinDB,addInfo,err)
					_,err =bot.Send(msg)
					errorsMessage(placeMessageCommand,err)
				case "tracking":
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
						"Здесь вы можете настроить оповещения.\n" +
							"На выбор предоставляются три режима:\n" +
							"1⃣. Ежечасное оповещение о цене\n" +
							"2⃣. Оповещение при достижении определенной цены\n" +
							"3⃣. Оповещение при резком подъеме криптовалюты\n" +
							"Для включения/настройки одного из режимов нажмите на кнопку с соответствующим номером")
					//добавляние кнопок
					msg.ReplyMarkup=setRegime
					_,err =bot.Send(msg)
					errorsMessage(placeMessageCommand,err)
				case "rate":
					var msgText string
					if bitcoinNow.isIncrease{
						msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подорожал на %v💲 (+%v%%)",
							bitcoinNow.Cost,bitcoinNow.changeCostUSD,bitcoinNow.changeCostPr)
					}else {
						msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подешевел на %v💲 (%v%%)",
							bitcoinNow.Cost,bitcoinNow.changeCostUSD*(-1),bitcoinNow.changeCostPr)
					}
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,msgText)
					_,err =bot.Send(msg)
					errorsMessage(placeMessageCommand,err)
				case "news":
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"🔊 Новости о Bitcoin`e:\n")
					_,err =bot.Send(msg)
					errorsMessage(placeMessageCommand,err)
					for _,news:= range allNews{
						msg =tgbotapi.NewMessage(update.Message.Chat.ID,news)
						_,err =bot.Send(msg)
						errorsMessage(placeMessageCommand,err)
					}
				case "analytics":
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,"🎓 Аналитика динамики Bitcoin`a:\n")
					_,err = bot.Send(msg)
					errorsMessage(placeMessageCommand,err)
					for _,analysis:= range allAnalysis{
						msg :=tgbotapi.NewMessage(update.Message.Chat.ID,analysis)
						_,err =bot.Send(msg)
						errorsMessage(placeMessageCommand,err)
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
					errorsMessage(placeMessageCommand,err)
				}
			} else if !update.Message.IsCommand() {//обработка сообщений
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
						msg =tgbotapi.NewMessage(update.Message.Chat.ID,"Некоректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
					}
					_,err =bot.Send(msg)
					errorsMessage(placeMessageNotCommand,err)
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
						msg =tgbotapi.NewMessage(update.Message.Chat.ID,"Некоректный формат ввода, пожалуйста, отправьте число в формате: '123.456' ")
					}
					_,err =bot.Send(msg)
					errorsMessage(placeMessageNotCommand,err)
				default:
					msg :=tgbotapi.NewMessage(update.Message.Chat.ID,
						helloMessage)
					_,err =bot.Send(msg)
					errorsMessage(placeMessageNotCommand,err)
				}
			}

		}
	}()
	//Часть отвечающая за сообщения связанные с достижением фиксированной цены
	//проверим цену на биткоин и сравним ее с теми, что есть на запросах в таблице
	go func(){
		for;;{
			time.Sleep(60*time.Second)
			//из бд получим юзеров, которым нужно отправить оповещение
			users,err:=allChatIdCostDB(bitcoinNow.Cost,db)
			errorsWorkDB(ChatIdCostDB,giveInfo,err)
			//отправим юзерам оповещение
			for _,user:=range users{
				msg :=tgbotapi.NewMessage(int64(user.ChatId),"Bitcoin достиг стоимости в "+
					strconv.FormatFloat(user.cost,'f',2,64)+" USD\n"+
					"Сейчас цена составляет "+strconv.FormatFloat(bitcoinNow.Cost,'f',2,64)+" USD.")
				_,err =bot.Send(msg)
				errorsMessage(placeMessageNotCommand,err)
			}
			////////////////////////////////////////////////////////////////////////////////////
		}
	}()
	//Часть отвечающая за сообщения связанные с фиксированным ростом цены
	//проверим изменение цены на биткоин и сравним его с теми, что есть  в таблице
	go func(){
		for;;{
			time.Sleep(60*time.Second)
			//из бд получим юзеров, которым нужно отправить оповещение
			users,err:=allChatIdChangeCostDB(bitcoinNow.changeCostUSD,db)
			errorsWorkDB(ChatIdChangeCostDB,giveInfo,err)
			//отправим юзерам оповещение
			for _,user:=range users{
				msg :=tgbotapi.NewMessage(int64(user.ChatId),"Изменение цены Bitcoin`a за последние 24 часа превысило "+
					strconv.FormatFloat(user.cost,'f',2,64)+" USD\n"+
					"И сейчас составляет "+strconv.FormatFloat(bitcoinNow.changeCostUSD,'f',2,64)+" USD.")
				_,err =bot.Send(msg)
				errorsMessage(placeSendMessageAboutBitcoin,err)
			}
			////////////////////////////////////////////////////////////////////////////////////
		}
	}()
	//обновление информации о биткойне
	go func(){
		for ;;{
			bitcoinNow,err=parsAllInfoAboutBitcoin()
			if err!=nil{
				log.Println("Проблема с парсингом информации о биткойне: ",err)
			}
			fmt.Println("Обновление информации о биткоине")
			time.Sleep(30*time.Second)
		}
	}()
	go func(){
		for ;;{
			time.Sleep(30*time.Second)
			allAnalysis,err=parsAnalysis()
			if err!=nil{
				fmt.Println(err)
			}
			allNews,err=parsNews()
			if err!=nil{
				fmt.Println(err)
			}
		}
	}()
	//Отправка информации о цене биткойна
	for ; ; {
		var msgText string
		if bitcoinNow.isIncrease{
			msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подорожал на %v💲 (+%v%%)",
				bitcoinNow.Cost,bitcoinNow.changeCostUSD,bitcoinNow.changeCostPr)
		}else {
			msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подешевел на %v💲 (%v%%)",
				bitcoinNow.Cost,bitcoinNow.changeCostUSD*(-1),bitcoinNow.changeCostPr)
		}
		time.Sleep(30*time.Second)
		allChatId,err:=allChatIdInfoBitcoinDB(db)
		errorsWorkDB(InfoBitcoinDB,giveInfo,err)
		for _,chatId:=range allChatId{
			msg :=tgbotapi.NewMessage(chatId,msgText)
			_,err =bot.Send(msg)
			errorsMessage(placeSendMessageAboutBitcoin,err)
		}
		time.Sleep(30*time.Second)
	}
}