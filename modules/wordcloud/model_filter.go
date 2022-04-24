package wordcloud

type Filter interface {
	IsLegal(word WordJson) bool
} 

type FilterWithDict struct {
	Dictionary map[string]any
}

func (filter *FilterWithDict) Has(str string) bool {
	_, ok := filter.Dictionary[str]
	return ok
}

func (filter *FilterWithDict) IsLegal(word WordJson) bool {
	return !filter.Has(word.Content) && len(word.Content) > 3
}

func NewWordFilter(dict map[string]any) *FilterWithDict {
	return &FilterWithDict{
		Dictionary: dict,
	}
}
