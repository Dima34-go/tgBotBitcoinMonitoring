package main

import (
	"fmt"
	"time"
)

func updateInfoAboutBitcoin(updateRate time.Duration,bitcoinNow *infoBitcoin,errInfoBitcoinPars *error){
	for ;;{
		time.Sleep(updateRate)
		*(bitcoinNow),*errInfoBitcoinPars = parsAllInfoAboutBitcoin()
		if *errInfoBitcoinPars!=nil{
			fmt.Println("Проблема с парсингом информации о биткойне: ",*errInfoBitcoinPars)
		}
	}
}
func updateNewsAnalysis(updateRate time.Duration,allAnalysis []string,allNews []string,errNews *error,errAnalysis *error){
	for ;;{
		time.Sleep(updateRate)
		allAnalysis,*errAnalysis=parsAnalysis()
		if *errAnalysis!=nil{
			fmt.Println(*errAnalysis)
		}
		allNews,*errNews=parsNews()
		if *errNews!=nil{
			fmt.Println(*errNews)
		}
	}
}