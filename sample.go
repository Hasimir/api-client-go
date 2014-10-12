package main

import (
	"fmt"
	"github.com/martinbajalan/ice3x"
)

func main() {
	var apikey string = "Your API Key"
	var privatekey string = "Your Private Key"
	ice3x.SetAuth(apikey, privatekey)
	
	// trade history sample 
	trades, err := ice3x.TradeHistory("ZAR", "BTC", 10, 1)
	if err != nil {
		fmt.Printf("Can't get trade history : %s\n", err)
		return
	}
	fmt.Printf("results : %v\n", trades)
}
