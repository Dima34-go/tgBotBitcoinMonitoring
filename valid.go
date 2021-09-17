package main

import "strconv"

func validAndPrepare(costStr string) (float64,bool){
	cost,err:=strconv.ParseFloat(costStr,64)
	if err!=nil{
		return 0,false
	}else {
		return cost,true
	}
}