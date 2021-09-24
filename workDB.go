package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct{
	UsersId int
	ChatId int
	UseBotFunc bool
}
type userCost struct {
	ChatId int
	cost float64
}
//база данных для обновлений с периодичностью 5 минут
func addNewUser(chatId int, useBotFunc bool,db *sql.DB) error {
	//проверка на присутствие в юзера в БД
	rows, err := db.Query("select * from userschat where chatId = ?",chatId)
	if err!=nil{
		return err
	}
	var us User
	for rows.Next(){
		err = rows.Scan(&us.UsersId, &us.ChatId, &us.UseBotFunc)
		if err != nil{
			return err
		}
	}
	//если пользователя нет - добавляем его в БД
	if us.UsersId==0 && us.ChatId==0 {
		_, err = db.Exec("insert into userschat (chatId, useBotFunc) values (?,?)",
			chatId,useBotFunc)
		if err != nil{
			return err
		}
	}
	return nil
}
func changeInformation(chatId int, useBotFunc bool,db *sql.DB) error{
	_ ,err := db.Exec("update userschat set useBotFunc = ? where chatId = ?",useBotFunc,chatId)
	return err
}
func allChatIdInfoBitcoinDB(db *sql.DB) ([]int64,error) {
	chatId:=make([]int64,0)
	rows, err := db.Query("select chatId from userschat where UseBotFunc=true")
	if err != nil {
		return chatId,err
	}
	//defer rows.Close()
	var newChatId int64
	for rows.Next(){
		err = rows.Scan(&newChatId)
		chatId=append(chatId,newChatId)
		if err != nil{
			return chatId,err
		}
	}
	err=rows.Close()
	return chatId,err
}
//база данных для уведомлений о дохождении биктойна до цены
func addUserCostDB(chatId int,cost float64,db *sql.DB) error{
	_, err := db.Exec("insert into notifCost (chatId, Cost) values (?,?)",
		chatId,cost)
	return err
}
func allChatIdCostDB(cost float64,db *sql.DB) ([]userCost,error) {
	us:=make([]userCost,0)
	rows, err := db.Query("select ChatId , Cost from notifCost where Cost < ?",cost)
	if err != nil {
		return us,err
	}
	i:=0
	for rows.Next(){
		us=append(us,userCost{0,0})
		err = rows.Scan(&us[i].ChatId, &us[i].cost)
		i++
		if err != nil{
			return us,err
		}
	}
	err=rows.Close()
	if err!=nil{
		return us,err
	}
	//удаление всех кого касался запрос
	_, err = db.Exec("delete from notifCost where Cost < ?",cost)
	return us,err
}
func deleteUserChatIdCostDB(chatId int,db *sql.DB) error {
	_, err := db.Exec("delete from notifCost where chatId = ?",chatId)
	return err
}
//база данных для уведомлении о резком скачке биткойна
func addUserChangeCostDB(chatId int,changeCost float64,db *sql.DB) error {
	_, err := db.Exec("insert into notif_change_cost (chatId, Change_cost) values (?,?)",
		chatId,changeCost)
	return err
}
func allChatIdChangeCostDB(changeCost float64,db *sql.DB) ([]userCost,error) {
	us:=make([]userCost,0)
	rows, err := db.Query("select ChatId ,Change_cost from notif_change_cost where Change_cost < ?",changeCost)
	if err != nil {
		return us,err
	}

	i:=0
	for rows.Next(){
		us=append(us,userCost{0,0})
		err = rows.Scan(&us[i].ChatId, &us[i].cost)
		i++
		if err != nil{
			return us,err
		}
	}
	//удаление всех кого касался запрос
    err = rows.Close()
	if err != nil {
		return us,err
	}
	_, err = db.Exec("delete from notif_change_cost where Change_cost < ?",changeCost)
	return us,err
}
func deleteUserChatIdChangeCostDB(chatId int,db *sql.DB) error {
	_, err := db.Exec("delete from notif_change_cost where chatId = ?",chatId)
	return err
}