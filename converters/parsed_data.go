package converters

type ParsedData struct {
	ListWord                 []*Word
	BelongToPreviousFileWord *Word
}

func EmptyParsedData() *ParsedData {
	return &ParsedData{
		ListWord:                 []*Word{},
		BelongToPreviousFileWord: nil,
	}
}

func (p *ParsedData) MergeLastWord(word *Word) {
	p.ListWord[len(p.ListWord) -1].Merge(word)
}