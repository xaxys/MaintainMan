package util

import (
	"strconv"
	"strings"
)

type PermSet struct {
	data map[string]interface{}
}

func NewPermSet() *PermSet {
	return &PermSet{
		data: make(map[string]interface{}),
	}
}

func (p *PermSet) seperate(key string) ([]string, string, bool) {
	positive := !strings.HasPrefix(key, "-")
	key = strings.TrimLeft(key, "-")
	parts := strings.Split(key, ".")
	last := "@"
	if parts[len(parts)-1] == "*" {
		last = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	} else if _, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
		last = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}
	return parts, last, positive
}

func (p *PermSet) Add(keys ...string) *PermSet {
	for _, v := range keys {
		p.add(v)
	}
	return p
}

func (p *PermSet) add(key string) *PermSet {
	parts, last, positive := p.seperate(key)
	data := p.data
	for _, v := range parts {
		if _, ok := data[v]; !ok {
			data[v] = make(map[string]interface{})
		}
		data = data[v].(map[string]interface{})
	}
	data[last] = positive
	return p
}

func (p *PermSet) Delete(keys ...string) *PermSet {
	for _, v := range keys {
		p.delete(v)
	}
	return p
}

func (p *PermSet) delete(key string) *PermSet {
	parts, last, _ := p.seperate(key)
	data := p.data
	for _, v := range parts {
		if _, ok := data[v]; !ok {
			break
		}
		data = data[v].(map[string]interface{})
	}
	delete(data, last)
	return p
}

func (p *PermSet) Has(key string) bool {
	positive, _ := p.Find(key)
	return positive
}

func (p *PermSet) Find(key string) (positive, found bool) {
	if strings.HasPrefix(key, "-") {
		positive, found := p.Find(strings.TrimLeft(key, "-"))
		return !positive, found
	}
	parts, last, _ := p.seperate(key)
	data := p.data
	for _, v := range parts {
		if data["*"] != nil {
			positive = true
			found = true
		}
		if data[v] == nil {
			return
		}
		data = data[v].(map[string]interface{})
	}
	if num, err := strconv.Atoi(last); err == nil {
		for k, v := range data {
			if knum, kerr := strconv.Atoi(k); kerr == nil {
				if knum >= num && v.(bool) {
					positive = true
				}
				if knum == num {
					found = true
				}
			}
		}
	} else if data[last] != nil {
		positive = data[last].(bool)
		found = true
	}
	return
}
