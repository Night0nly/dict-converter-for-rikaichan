package converters

import "strings"

type Word struct {
	Kanji      string
	Kana       string
	Definition string
}

func NewWord(kanji, kana, definition string) *Word {
	return &Word{
		Kanji:      kanji,
		Kana:       kana,
		Definition: definition,
	}
}

func (w *Word) ToEdictFormat() string {
	if len(w.Kanji) != 0 {
		return w.Kanji + " " + "[" + w.Kana + "]" + " /" + strings.Trim(w.Definition, "\t") + "/"
	} else {
		if c := []rune(w.Definition)[0]; c >= 0x30a0 && c <= 0x30ff {
			return w.Definition
		}
		return w.Kana + " /" + strings.Trim(w.Definition, "\t") + "/"
	}
}

func (w *Word) Merge(word *Word) {
	w.Kanji += word.Kanji
	w.Kana += word.Kana
	w.Definition += "/" + word.Definition
}