package musician

import (
	"encoding/json"
	"github.com/RtillaWork/gogetitarchy/utils"
	"strings"
	"time"
)

type DataDict struct {
	LastModified time.Time           `json:"last_modified"`
	Fields       map[string][]string `json:"fields"`
}

var TheDataDict DataDict

func init() {
	TheDataDict = DataDict{
		LastModified: time.Now(),

		Fields: map[string][]string{
			"FIRSTNAMES":  []string{},
			"MIDDLENAMES": []string{},
			"LASTNAMES":   []string{},
		},
	}
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

// utility funcs

// BuildTheDataDict build maps of lists of key-values apprearing in Musician struct fields
func BuildTheDataDict(musiciansmap MusiciansMap) {
	firstnames := make(map[string]int)
	middlenames := make(map[string]int)
	lastnames := make(map[string]int)
	keys := make(map[string]int)
	values := make(map[string]int)
	valueskey := make(map[string]string)

	// collect all first, middle and last names... to TheDataDict
	// ...then collect all keys and values
	for _, mv := range musiciansmap {
		//// fname, mname, lname
		fname := strings.ToUpper(strings.TrimSpace(mv.FName))
		mname := strings.ToUpper(strings.TrimSpace(mv.MName))
		lname := strings.ToUpper(strings.TrimSpace(mv.LName))

		if fname != Defaults.FName {

			firstnames[fname]++
		}

		if mname != Defaults.FName {

			middlenames[mname]++
		}

		if lname != Defaults.FName {

			lastnames[lname]++
		}
		////Fields[key]value
		for key, val := range mv.Fields {
			k := strings.ToUpper(strings.TrimSpace(key))
			v := strings.ToUpper(strings.TrimSpace(val))
			keys[k]++
			values[v]++
			valueskey[v] = k
		}
	}

	// then adds all unique keys and values to TheDataDict
	for k, _ := range firstnames {
		TheDataDict.Fields["FIRSTNAMES"] = append(TheDataDict.Fields["FIRSTNAMES"], k)
	}

	for k, _ := range middlenames {
		TheDataDict.Fields["MIDDLENAMES"] = append(TheDataDict.Fields["MIDDLENAMES"], k)
	}

	for k, _ := range lastnames {
		TheDataDict.Fields["LASTNAMES"] = append(TheDataDict.Fields["LASTNAMES"], k)
	}

	for key, _ := range keys {
		for val, keyofv := range valueskey {
			if key == keyofv {
				TheDataDict.Fields[keyofv] = append(TheDataDict.Fields[keyofv], val)
			}
		}
	}
	json.Marshal(TheDataDict)
	utils.WaitForKeypress()
}

//
func UpdateTheDict(key string, value string) {
	TheDataDict.Update(strings.ToUpper(key), strings.ToLower(value))
}