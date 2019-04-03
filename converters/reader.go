package converters

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Reader struct {
	Parser Parser
}

func NewReader(parser Parser) *Reader {
	return &Reader{
		Parser: parser,
	}
}

func (r *Reader) ReadAllDirectory(dirName string) ([]*Word, error) {
	timeStart := time.Now()
	files, _ := ioutil.ReadDir(dirName)
	if len(files) == 0 {
		return nil, errors.New("Directory is empty")
	}

	errChan := make(chan error, len(files))
	byteDataChan := make(chan *fileContent, len(files))

	for _, file := range files {
		go func(file os.FileInfo) {
			regex := regexp.MustCompile(`\D`)
			fileIndex, _ := strconv.Atoi(regex.ReplaceAllString(file.Name(), ""))

			path := dirName + "/" + file.Name()
			fmt.Println(path)
			defer fmt.Println("Done " + path + " " + strconv.Itoa(fileIndex) + "after " + fmt.Sprintf("%f", time.Now().Sub(timeStart).Seconds()))

			jsonDict, err := os.Open(path)
			defer jsonDict.Close()
			if err != nil {
				errChan <- errors.Wrap(err, "Cannot open the file: "+path)
			}
			byteData, err := ioutil.ReadAll(jsonDict)
			if err != nil {
				errChan <- errors.Wrap(err, "Error reading file: "+path)
			}
			parsedData, err := r.Parser.Parse(byteData)
			if err != nil {
				errChan <- errors.Wrap(err, "Error parsing file " + path)
			}
			byteDataChan <- &fileContent{fileIndex, parsedData}
		}(file)
	}

	contentMap := make(map[int]*ParsedData)
	for {
		select {
		case e := <-errChan:
			return nil, e
		case content := <-byteDataChan:
			contentMap[content.Index] = content.Data
		}
		if len(contentMap) == len(files) {
			close(errChan)
			close(byteDataChan)
			var listWord []*Word
			for i := 1; i < len(files) ; i++ {
				if previousFileWord := contentMap[i].BelongToPreviousFileWord; previousFileWord != nil {
					contentMap[i-1].MergeLastWord(previousFileWord)
				}
				listWord = append(listWord, contentMap[i-1].ListWord...)
			}
			listWord = append(listWord, contentMap[len(contentMap)-1].ListWord...)
			fmt.Println("Done Read file")
			fmt.Println("Finish reading files after " + fmt.Sprintf("%f", time.Now().Sub(timeStart).Seconds()))
			return listWord, nil
		}
	}
}

type fileContent struct {
	Index int
	Data  *ParsedData
}
