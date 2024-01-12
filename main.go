package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getRatesMap() map[string]string {
	requestURL := "https://api.coinbase.com/v2/exchange-rates?currency=USD"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	var valueMap map[string]map[string]map[string]string //lose currency value
	json.Unmarshal(resBody, &valueMap)

	rates := valueMap["data"]["rates"]
	return rates
}

func ConvertFromMap(amount float64, rates map[string]string) (map[string]interface{}, error) {

	if amount <= 0 {
		return nil, errors.New("Negative Amount now allowed")
	}
	if _, ok := rates["BTC"]; !ok {
		return nil, errors.New("Missing BTC value")
	}
	if _, ok := rates["ETH"]; !ok {
		return nil, errors.New("Missing ETH value")
	}

	btcString := rates["BTC"]
	ethString := rates["ETH"]
	btc, _ := strconv.ParseFloat(btcString, 64)
	eth, _ := strconv.ParseFloat(ethString, 64)

	resultMap := make(map[string]interface{})
	resultMap["BTC"] = calc(amount, .7, btc)
	resultMap["ETH"] = calc(amount, .3, eth)
	resultMap["timestamp"] = time.Now()
	return resultMap, nil
}

func ToJson(resultMap map[string]interface{}) (string, error) {
	jsonBytes, jsonErr := json.MarshalIndent(resultMap, "", "   ")
	if jsonErr != nil {
		return "", jsonErr //should use a pointer?
	}
	json := string(jsonBytes)
	return fmt.Sprintf("%v", json), nil
}

func calc(amount float64, percentage float64, coin float64) float64 {
	amt := amount * percentage * coin
	return amt
}

func main() {
	log.Println("----stephen leonard----")
	argsWithProg := os.Args
	if len(argsWithProg) == 0 {
		panic("Starting Dollar amount needed as program argument")
	}
	amountString := os.Args[1]
	fmt.Println("Converting for: [$" + amountString + "]")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	var btc = "0.0"
	var eth = "0.0"

	for {
		ratesMap := getRatesMap()
		if ratesMap["BTC"] != btc || ratesMap["ETH"] != eth {
			resultsMap, convertError := ConvertFromMap(amount, ratesMap)
			if convertError != nil {
				fmt.Print("E")
			} else {
				result, _ := ToJson(resultsMap)
				fmt.Println("\n" + result)
				btc = ratesMap["BTC"]
				eth = ratesMap["ETH"]
			}
		} else {
			fmt.Print(".")
		}
		time.Sleep(2 * time.Second)
	}
}
