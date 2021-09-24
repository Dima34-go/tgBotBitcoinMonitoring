package main

import (
	"fmt"
	"time"
)

func updateInfoAboutBitcoin(updateRate time.Duration,bitcoinNow *infoBitcoin){
	for ;;{
		time.Sleep(updateRate)
		var err error
		*(bitcoinNow),err = parsAllInfoAboutBitcoin()
		if err!=nil{
			fmt.Println("Проблема с парсингом информации о биткойне: ",err)
		}
	}
}
func updateNewsAnalysis(updateRate time.Duration,allAnalysis []string,allNews []string){
	for ;;{
		time.Sleep(updateRate)
		var err error
		allAnalysis,err=parsAnalysis()
		if err!=nil{
			fmt.Println(err)
		}
		allNews,err=parsNews()
		if err!=nil{
			fmt.Println(err)
		}
	}
}