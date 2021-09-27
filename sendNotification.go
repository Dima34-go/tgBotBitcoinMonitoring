package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"time"
)

func sendMessageBitcoin(bitcoinNow infoBitcoin,db *sql.DB,bot *tgbotapi.BotAPI,updateRate time.Duration,errInfoBitcoinPars *error){
	for ; ; {
		time.Sleep(updateRate)
		if *errInfoBitcoinPars==nil{
			var msgText string
			if bitcoinNow.isIncrease{
				msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подорожал на %v💲 (+%v%%)",
					bitcoinNow.Cost,bitcoinNow.changeCostUSD,bitcoinNow.changeCostPr)
			}else {
				msgText=fmt.Sprintf("Курс на данный момент:1₿  = %v💲 \nЗа последние 24 часа Bitcoin подешевел на %v💲 (%v%%)",
					bitcoinNow.Cost,bitcoinNow.changeCostUSD*(-1),bitcoinNow.changeCostPr)
			}
			allChatId,err:=allChatIdInfoBitcoinDB(db)
			errorsWorkDB(InfoBitcoinDB,giveInfo,err)
			for _,chatId:=range allChatId{
				msg :=tgbotapi.NewMessage(chatId,msgText)
				_,err =bot.Send(msg)
				if err!=nil{
					errorsMessage(placeSendMessageAboutBitcoin,err,msg,db)
				}
			}
		}
	}
}
//Часть отвечающая за сообщения связанные с достижением фиксированной цены
//проверим цену на биткоин и сравним ее с теми, что есть на запросах в таблице
func sendMessageAboutCostBitcoin(bitcoinNow infoBitcoin,db *sql.DB,bot *tgbotapi.BotAPI,updateRate time.Duration,errInfoBitcoinPars *error){
	for;;{
		time.Sleep(updateRate)
		if *errInfoBitcoinPars==nil{
			//из бд получим юзеров, которым нужно отправить оповещение
			users,err:=allChatIdCostDB(bitcoinNow.Cost,db)
			errorsWorkDB(ChatIdCostDB,giveInfo,err)
			//отправим юзерам оповещение
			for _,user:=range users{
				msg :=tgbotapi.NewMessage(int64(user.ChatId),"Bitcoin достиг стоимости в "+
					strconv.FormatFloat(user.cost,'f',2,64)+" USD\n"+
					"Сейчас цена составляет "+strconv.FormatFloat(bitcoinNow.Cost,'f',2,64)+" USD.")
				_,err =bot.Send(msg)
				if err!=nil{
					errorsMessage(placeMessageNotCommand,err,msg,db)
				}
			}
		}
	}
}
//Часть отвечающая за сообщения связанные с фиксированным ростом цены
//проверим изменение цены на биткоин и сравним его с теми, что есть  в таблице
func sendMessageAboutChangeCostBitcoin(bitcoinNow infoBitcoin,db *sql.DB,bot *tgbotapi.BotAPI,updateRate time.Duration,errInfoBitcoinPars *error){
	for;;{
		time.Sleep(updateRate)
		if *errInfoBitcoinPars==nil{
			//из бд получим юзеров, которым нужно отправить оповещение
			users,err:=allChatIdChangeCostDB(bitcoinNow.changeCostUSD,db)
			errorsWorkDB(ChatIdChangeCostDB,giveInfo,err)
			//отправим юзерам оповещение
			for _,user:=range users{
				msg :=tgbotapi.NewMessage(int64(user.ChatId),"Изменение цены Bitcoin`a за последние 24 часа превысило "+
					strconv.FormatFloat(user.cost,'f',2,64)+" USD\n"+
					"И сейчас составляет "+strconv.FormatFloat(bitcoinNow.changeCostUSD,'f',2,64)+" USD.")
				_,err =bot.Send(msg)
				if err!=nil{
					errorsMessage(placeSendMessageAboutBitcoin,err,msg,db)
				}
			}
		}
	}
}