package wordcloud

import (
	"strings"

	"github.com/xaxys/maintainman/core/model"
	"github.com/yanyiwu/gojieba"
	"github.com/xaxys/maintainman/modules/order"
)

type Word struct {
	model.BaseModel
	Content string `gorm:"not null; unique;"`
	WordClass string `gorm:"not null"; size:10;`
	Count int `gorm:"not null"; default:0;`
	OrderId int `gorm:"not null";`
	Order order.Order `gorm:"foreignKey:OrderId"`
}

type WordJson struct {
	Content string `json:"content"`
	WordClass string `json:"wordclass"`
	Count int `json:"count"`
}

type WordCollector struct {
	wordSet []WordJson
}

func NewWordCollectorWithStr(str string) *WordCollector {
	wordAnalyzer := gojieba.NewJieba()
	defer wordAnalyzer.Free()
	words := wordAnalyzer.Tag(str)
	roots, wordClass := getWordClass(words)
    wordCounter := getWordCounter(roots)
	return NewWordCollectorWithSet(generateWordSet(roots, wordClass, wordCounter)) 

}

func NewWordCollectorWithSet(wordSet []WordJson) *WordCollector{
	return &WordCollector{
		wordSet: wordSet,
	}
}

func getWordClass(words []string) ([]string, []string) {
	roots := make([]string, 0)
	wordclass := make([]string, 0)
	for _, word := range words {
		spilt := strings.Split(word, "/")
		roots = append(roots, spilt[0])
		wordclass = append(wordclass, spilt[1])
	}
	return roots, wordclass
}

func getWordCounter(words []string) map[string]int {
	counter := make(map[string]int, 0)
	for _, word := range words {
		_, ok := counter[word]
		if !ok {
			counter[word] = 1
		} else {
			counter[word]++
		}
	}
	return counter
}

func accumulateWords(words []*WordJson) []*WordJson {
	length := len(words)
	ans := make([]*WordJson, 0)
	classMap := make(map[string]string)
	wordCounter := make(map[string]int)
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
			Content: word,
			WordClass: classMap[word],
			Count: count,
		})
	}

	return ans
}

func generateWordSet(roots []string, wordclass []string, wordCounter map[string] int) []WordJson {
	length := len(roots)
	ans := make([]WordJson, 0)
	uniqueMap := make(map[string] bool)
	for i := 0; i < length; i++ {
		if _, ok := wordCounter[roots[i]]; ok {
			if ok := uniqueMap[roots[i]]; !ok {
				uniqueMap[roots[i]] = true
				ans = append(ans, WordJson{
					Content:   roots[i],
					WordClass: wordclass[i],
					Count:     wordCounter[roots[i]],
				})
			}

		}
	}
	return ans
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

type UploadWordsRequest struct {
	Content string
	OrderId uint
}

type GetAllWordsRequest struct {
	model.PageParam
}

type GetWordsByOrderIdRequest struct {
	model.PageParam
	OrderId uint
}

