package wordcloud

import (
	"github.com/go-ego/gse"
)

var (
	seg gse.Segmenter
)

type WordCollector struct {
	wordSet []WordJson
}

func NewWordCollectorWithStr(str string) *WordCollector {
	words := seg.Pos(str, true)
	words = seg.TrimPos(words)
	words = seg.TrimWithPos(words, "x", "m", "eng")
	wordSet := newWordSet(words)
	return NewWordCollectorWithSet(wordSet)
}

func NewWordCollectorWithSet(wordSet []WordJson) *WordCollector {
	return &WordCollector{
		wordSet: wordSet,
	}
}

func accumulateWords(words []*WordJson) []*WordJson {
	length := len(words)
	ans := make([]*WordJson, 0)
	classMap := make(map[string]string)
	wordCounter := make(map[string]uint)
	for i := 0; i < length; i++ {
		if _, ok := classMap[words[i].Content]; !ok {
			wordCounter[words[i].Content] = words[i].Count
			classMap[words[i].Content] = words[i].WordClass
		} else {
			wordCounter[words[i].Content] += words[i].Count
		}
	}

	for word, count := range wordCounter {
		ans = append(ans, &WordJson{
			Content:   word,
			WordClass: classMap[word],
			Count:     count,
		})
	}

	return ans
}

func newWordSet(words []gse.SegPos) (wordSet []WordJson) {
	counter := map[gse.SegPos]uint{}
	for _, word := range words {
		counter[word]++
	}
	for word, count := range counter {
		wordSet = append(wordSet, WordJson{
			Content:   word.Text,
			WordClass: word.Pos,
			Count:     count,
		})
	}
	return
}

func (wc *WordCollector) Filter(filter Filter) *WordCollector {
	newWordSet := make([]WordJson, 0)
	for _, word := range wc.wordSet {
		if filter.IsLegal(word) {
			newWordSet = append(newWordSet, word)
		}
	}
	wc.wordSet = newWordSet
	return wc
}

func (wc *WordCollector) ToSlice() []WordJson {
	return wc.wordSet
}
