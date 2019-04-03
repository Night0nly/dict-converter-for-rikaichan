package converters

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type RikaichanDictParser struct {}

func NewRikaichanDictParser() *RikaichanDictParser {
	return &RikaichanDictParser{}
}

func (r *RikaichanDictParser) Parse(data []byte) (*ParsedData, error){
	var dataArray [][]string
	err := json.Unmarshal(data, &dataArray)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing file")
	}
	fmt.Println("Length of the dict: " + strconv.Itoa(len(dataArray)))

	parsedData := EmptyParsedData()
	previousDefinitionIndexInArray := 0
	previousDefinition := ""

	for i, element := range dataArray {
		if len(element) == 1 {
			definition := strings.Replace(element[0], "<br", "/", 1)
			if i == previousDefinitionIndexInArray {
				previousDefinition += definition
				previousDefinitionIndexInArray++
			} else {
				parsedData.ListWord[len(parsedData.ListWord) - 1].Definition += definition
			}
		} else if len(element) == 3 {
			parsedData.ListWord = append(parsedData.ListWord, NewWord(element[0], element[1], element[2]))
		} else {
			fmt.Print("Wrong format. Length of element is " + strconv.Itoa(len(element)) + " ")
			fmt.Println(element)
		}
	}

	if previousDefinition != "" {
		parsedData.BelongToPreviousFileWord = NewWord("", "", previousDefinition)
	}

	return parsedData, nil
}

