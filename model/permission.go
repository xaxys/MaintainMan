package model

import (
	"strings"
)

type PermissionJson struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type PermissionSet struct {
	data map[string]interface{}
}

func NewPermissionSet() *PermissionSet {
	return &PermissionSet{
		data: make(map[string]interface{}),
	}
}

func (p *PermissionSet) seperate(key string) ([]string, string, bool) {
	positive := !strings.HasPrefix(key, "-")
	key = strings.TrimLeft(key, "-")
	parts := strings.Split(key, ".")
	last := "@"
	if parts[len(parts)-1] == "*" {
		last = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}
	return parts, last, positive
}

func (p *PermissionSet) Add(keys ...string) *PermissionSet {
	for _, v := range keys {
		p.add(v)
	}
	return p
}

func (p *PermissionSet) add(key string) *PermissionSet {
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

func (p *PermissionSet) Delete(keys ...string) *PermissionSet {
	for _, v := range keys {
		p.delete(v)
	}
	return p
}

func (p *PermissionSet) delete(key string) *PermissionSet {
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

func (p *PermissionSet) Has(key string) bool {
	positive, _ := p.Find(key)
	return positive
}

func (p *PermissionSet) Find(key string) (positive, found bool) {
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
	if data[last] != nil {
		positive = data[last].(bool)
		found = true
	}
	return
}
