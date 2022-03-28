package utils

import (
	"strings"
	"time"
)

type DataDict struct {
	LastModified time.Time           `json:"last_modified"`
	Fields       map[string][]string `json:"fields"`
}

var TheDataDict = DataDict{
	LastModified: time.Now(),

	Fields: map[string][]string{
		"FIRSTNAMES":  []string{},
		"MIDDLENAMES": []string{},
		"LASTNAMES":   []string{},
	},
}

func (d *DataDict) Update(key string, value string) {
	k := strings.ToUpper(key)
	v := strings.ToUpper(value)
	if _, ok := d.Fields[k]; !ok {
		d.Fields[k] = []string{v}
	} else {
		d.Fields[k] = append(d.Fields[k], v)
	}

}

func UpdateTheDict(key string, value string) {
	TheDataDict.Update(key, value)
}
