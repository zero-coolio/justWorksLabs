package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertNoCurrency(t *testing.T) {
	testMap := make(map[string]string)
	_, err := ConvertFromMap(1000, testMap)
	assert.NotNilf(t, err, "Currencies should be missing")
}

func TestConvert(t *testing.T) {
	testMap := make(map[string]string)
	testMap["BTC"] = "1"
	testMap["ETH"] = "1"

	resultMap, err := ConvertFromMap(1000, testMap)
	assert.Nilf(t, err, "Should be nil error")
	assert.Equal(t, 700.0, resultMap["BTC"], "BTC should be 700")
	assert.Equal(t, 300.0, resultMap["ETH"], "ETH should be 300")

}

func TestConvertNeg(t *testing.T) {
	testMap := make(map[string]string)
	testMap["BTC"] = "1"
	testMap["ETH"] = "1"

	_, err := ConvertFromMap(-1000, testMap)
	assert.NotNil(t, err, "Shouldn't be nil")
}
