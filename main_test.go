package main

import (
	"testing"
)

func TestGetUSDPrice(t *testing.T) {
	usdPrice, err := getUSDPrice()
	if err != nil {
		t.Errorf("getUSDPrice() error: %v", err)
	}

	if usdPrice == "" {
		t.Error("getUSDPrice() returned an empty string")
	}
}
