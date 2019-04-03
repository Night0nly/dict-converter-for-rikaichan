package main

import (
	"convertedict/converters"
	"fmt"
)

func main() {
	dirPath := "dicts"
	error := converters.NewRikaichanConverter(converters.NewRikaichanDictParser()).ConvertToEdict(dirPath)
	fmt.Println(error)
}