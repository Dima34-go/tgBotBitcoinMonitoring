package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strconv"
	"strings"
)
type infoBitcoin struct {
	Cost float64
	changeCostPr float64
	changeCostUSD float64
	isIncrease bool
}
//Парсится страница https://ru.investing.com/crypto/bitcoin
//получение ссылок на новости
func parsNews() ([]string,error){
	prefix:="https://ru.investing.com"
	//Количество считываемых новостей
	quatityNews:=3
	var allNews []string
	var err error
	i:=0
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://ru.investing.com/crypto/bitcoin/news"},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("div.js-articles-wrapper.largeTitle.js-news-items").Find("a.img").Each( func( size int , s *goquery.Selection) {
				if i<quatityNews{
					url,ok:=s.Attr("href")
					if !ok{
						err=errors.New("ошибка выделения ссылок, проверьте доступность сайта")
					}
					if url[0:4]!="http"{
						url=prefix+url
					}
					allNews=append(allNews,url)
					i++
				}
			})
		},
	}).Start()
	return allNews,err
}
//получение ссылок на аналитику
func parsAnalysis() ([]string,error){
	prefix:="https://ru.investing.com"
	//Количество считываемых аналитик
	quatityAnalysis:=3
	var allAnalysis []string
	var err error
	i:=0
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://ru.investing.com/crypto/bitcoin/analysis"},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("div.js-content-wrapper").Find("a.img").Each( func( size int , s *goquery.Selection) {
				if i<quatityAnalysis{
					url,ok:=s.Attr("href")
					if url[0:4]!="http"{
						url=prefix+url
					}
					if !ok{
						err=errors.New("ошибка выделения ссылок, проверьте доступность сайта")
					}
					allAnalysis=append(allAnalysis,url)
					i++
				}
			})
		},
	}).Start()
	return allAnalysis,err
}
//получение всей информации о цене с сайта, а именно цены,роста цены, роста в процентах
func parsAllInfoAboutBitcoin() (infoBitcoin,error)  {
	var costBitcoinStr,changeCostPrStr,changeCostUSDStr string
	var isIncrease  bool
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://ru.investing.com/crypto/bitcoin"},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// два варианта развития событий, либо стоимость поднимается, либо опускается
			//случай падения цены
			costBitcoinStr=r.HTMLDoc.Find("span#last_last.pid-1057391-last").Text()
			changeCostPrStr=r.HTMLDoc.Find("span.arial_20.pid-1057391-pcp.parentheses.redFont").Text()
			changeCostUSDStr=r.HTMLDoc.Find("span.arial_20.pid-1057391-pc.redFont").Text()
			if changeCostPrStr==""{//случай подорожания
				changeCostPrStr=r.HTMLDoc.Find("span.arial_20.pid-1057391-pcp.parentheses.greenFont").Text()
				changeCostUSDStr=r.HTMLDoc.Find("span.arial_20.pid-1057391-pc.greenFont").Text()
				isIncrease=true
			}
		},
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
	//получены необработанные данные, нужно превратить их в цифры
	//обработка стоимости биткойна
	costBitcoinStr=strings.ReplaceAll(costBitcoinStr,".","")
	costBitcoinStr=strings.ReplaceAll(costBitcoinStr,",",".")
	costBitcoin,err:=strconv.ParseFloat(costBitcoinStr,64)
	if err!=nil{
		return infoBitcoin{},err
	}
	//обработка изменения цены
	changeCostUSDStr=strings.ReplaceAll(changeCostUSDStr,"+","")
	//changeCostUSDStr=strings.ReplaceAll(changeCostUSDStr,"-","")
	changeCostUSDStr=strings.ReplaceAll(changeCostUSDStr,".","")
	changeCostUSDStr=strings.ReplaceAll(changeCostUSDStr,",",".")
	changeCostUSD,err:=strconv.ParseFloat(changeCostUSDStr,64)
	if err!=nil{
		return infoBitcoin{},err
	}
	//обработка изменения цены в процентах
	changeCostPrStr=strings.ReplaceAll(changeCostPrStr,"+","")
	changeCostPrStr=strings.ReplaceAll(changeCostPrStr,"%","")
	//changeCostUSDStr=strings.ReplaceAll(changeCostUSDStr,"-","")
	changeCostPrStr=strings.ReplaceAll(changeCostPrStr,".","")
	changeCostPrStr=strings.ReplaceAll(changeCostPrStr,",",".")
	changeCostPr,err:=strconv.ParseFloat(changeCostPrStr,64)
	if err!=nil{
		return infoBitcoin{},err
	}
	return infoBitcoin{costBitcoin,changeCostPr,changeCostUSD,isIncrease},nil
}