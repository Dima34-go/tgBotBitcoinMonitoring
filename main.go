package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)
func main() {
	//Запуск бота и настройка обновлений
	botToken:="1940781427:AAHHOpliPiUV0SQcKrLRQ6R9Ytkv1JZf5YE"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	var errNewsPars,errAnalysisPars,errInfoBitcoinPars error
	bot.Debug = true
	fmt.Printf("Authorized on account %s", bot.Self.UserName)
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
			fmt.Println("Ошибка закрытия БД: ",err)
		}
	}()
	//массив с данными о статусе пользователя
	status := make(map[int64]string,0)
	//Получение информации с сайта, до запуска бота
	bitcoinNow,err:=parsAllInfoAboutBitcoin()
	if err!=nil{
		log.Panic("Проблема с парсингом информации о биткойне: ",err)
	}
	allNews,err:=parsNews()
	if err!=nil{
		log.Panic("Проблема с парсингом новостей",err)
	}
	allAnalysis,err:=parsAnalysis()
	if err!=nil{
		log.Panic("Проблема с парсингом аналитики",err)
	}
	//--------------------------------получение update в бесконечном цикле
	go func() {
		for update := range updates {
			if update.Message == nil && update.CallbackQuery == nil{//пустое обновление
				continue
			}else if update.CallbackQuery != nil{//обработка callback ответа
				isCallbackQuery(&update,bot,db,status)
			}else if update.Message.IsCommand(){//обработка команд
				isCommandCase(&update,bot,db,bitcoinNow,allNews,allAnalysis,&errInfoBitcoinPars,&errNewsPars,&errAnalysisPars)
			} else if !update.Message.IsCommand() {//обработка сообщений
				isUsualMessage(&update,bot,db,status)
			}

		}
	}()
	//оповещение пользователей о проблемах работы бота
	go sendMessageUserAboutError(db,bot,&errInfoBitcoinPars)
	//обновление информации о bitcoin`e и изменении его цены
	go updateInfoAboutBitcoin(30*time.Second,&bitcoinNow,&errInfoBitcoinPars)
	go updateNewsAnalysis(30*time.Second,allAnalysis,allNews,&errNewsPars,&errAnalysisPars)
	//Отправка информации о цене bitcoin`a и ее изменении
	go sendMessageAboutCostBitcoin(bitcoinNow,db,bot,60*time.Second,&errInfoBitcoinPars)
	go sendMessageAboutChangeCostBitcoin(bitcoinNow,db,bot,60*time.Second,&errInfoBitcoinPars)
	sendMessageBitcoin(bitcoinNow,db,bot,60*time.Minute,&errInfoBitcoinPars)
}
