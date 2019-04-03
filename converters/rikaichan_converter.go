package converters

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"
)

type RikaichanConverter struct {
	reader *Reader
}

func NewRikaichanConverter(parser Parser) *RikaichanConverter {
	return &RikaichanConverter{
		reader: NewReader(parser),
	}
}

func (r *RikaichanConverter) ConvertToEdict(dirName string) error {
	wordList, err := r.reader.ReadAllDirectory(dirName)
	if err != nil {
		return errors.Wrap(err, "Fail to read file")
	}

	splitWord := 999
	timeStart := time.Now()
	edictChan := make(chan *edictData, len(wordList) / splitWord + 1)
	fmt.Println("WordList length is " + strconv.Itoa(len(wordList)))

	for i := 0; i < len(wordList); i += splitWord {
		list := wordList[i:min(i+splitWord, len(wordList))]
		fmt.Println("Transferring data from " + strconv.Itoa(i) + " to " + strconv.Itoa(min(i+splitWord, len(wordList))))
		go func(list []*Word) {
			edictString, dictIndex := convert(list)
			edictChan <- &edictData{edictString, dictIndex}
		}(list)
	}

	edict := ""
	dictIndex := make(map[string][] int)
	receiveCount := 0

	for {
		select {
		case edictData := <- edictChan:
			fmt.Println("Receive edictData from dict " + strconv.Itoa(receiveCount))
			lenEdict := utf8.RuneCountInString(edict)
			for k, v := range edictData.DictIndex {
				for i := 0; i < len(v); i++ {
					v[i] += lenEdict
				}
				dictIndex[k] = append(dictIndex[k], v...)
			}
			edict += edictData.Edict
			fmt.Println("Done execute edictData from dict " + strconv.Itoa(receiveCount))
			receiveCount++
		}
		if receiveCount == len(wordList) / splitWord + 1 {
			fmt.Print("Close editchan ")
			close(edictChan)
			break
		}
	}

	fmt.Println("Convert done after " + fmt.Sprintf("%f", time.Now().Sub(timeStart).Seconds()))
	fmt.Println()

	dictFile, err := os.Create("dict.dat")
	if err != nil {
		return errors.Wrap(err, "Cannot create the file")
	}
	defer dictFile.Close()

	w := bufio.NewWriter(dictFile)
	fmt.Fprint(w, edict)
	err = w.Flush() // Don't forget to flush!
	if err != nil {
		return errors.Wrap(err, "Cannot flush")
	}

	indexFile, err := os.Create("dict.idx")
	if err != nil {
		return errors.Wrap(err, "Cannot create the file")
	}
	defer indexFile.Close()
	w2 := bufio.NewWriter(indexFile)
	// sort key

	keys := make([]string, 0, len(dictIndex))
	for k := range dictIndex {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		printString := k
		for _, i := range dictIndex[k] {
			printString += "," + strconv.Itoa(i)
		}
		printString += "\n"
		fmt.Fprint(w2, printString)
	}
	err = w2.Flush() // Don't forget to flush!
	if err != nil {
		return errors.Wrap(err, "Cannot flush")
	}

	return nil
}

func convert(wordList []*Word) (string, map[string][] int) {
	wordChan := make(chan *wordEdict, len(wordList))

	for _, word := range wordList {
		go func(word *Word) {
			wordString := word.ToEdictFormat()


			if word.Kanji != "" {
				wordChan <- &wordEdict{wordString, []string{word.Kanji, word.Kana}}
			} else {
				wordChan <- &wordEdict{wordString, []string{word.Kana}}
			}
		}(word)

	}

	edictString := ""
	dictIndex := make(map[string][] int)
	receiveCount := 0
	for {
		select {
		case dict := <-wordChan:
			lenEdict := utf8.RuneCountInString(edictString)
			for _, t := range dict.Title {
				dictIndex[t] = append(dictIndex[t], lenEdict)
			}
			wordEdict := dict.WordString
			edictString += wordEdict + "\n"
			receiveCount++
		}
		if receiveCount == len(wordList) {
			close(wordChan)
			fmt.Println("Close word channel")
			return edictString, dictIndex
		}
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

type wordEdict struct {
	WordString string
	Title     []string
}

type edictData struct {
	Edict string
	DictIndex map[string][] int
}