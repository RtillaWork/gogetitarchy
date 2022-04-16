package musician

import (
	"encoding/json"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"log"
	"strings"
	"time"
)

type DataDict struct {
	LastModified time.Time           `json:"last_modified"`
	Fields       map[string][]string `json:"fields"`
	KeyStats     map[string]int      `json"keys_tats"`
	ValuesStats  map[string]int      `json"values_tats"`
}

var TheDataDict DataDict

func init() {
	TheDataDict = DataDict{
		LastModified: time.Now(),

		Fields: map[string][]string{
			"FIRSTNAMES":   []string{},
			"MIDDLENAMES":  []string{},
			"LASTNAMES":    []string{},
			"MISCELLANEAS": []string{},
			"DATESBEGIN":   []string{},
			"DATESEND":     []string{},
			"DATESOTHER":   []string{},
			"DATESCSV":     []string{},
		},

		KeyStats:    make(map[string]int),
		ValuesStats: make(map[string]int),
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
		fname := utils.NormalizeKey(mv.FName)
		mname := utils.NormalizeKey(mv.MName)
		lname := utils.NormalizeKey(mv.LName)

		// commented out, we want to count the unassigned names too
		//if fname != Defaults.FName {
		firstnames[fname]++
		// if mname != Defaults.MName {
		middlenames[mname]++
		// if lname != Defaults.LName {
		lastnames[lname]++

		////Fields[key]value
		for key, val := range mv.Fields {
			//log.Printf("KEY   %#v            VALUES  %#v\n", key, val)
			k := utils.NormalizeKey(key)
			v := utils.NormalizeValue(val)
			keys[k]++
			values[v]++
			valueskey[v] = k
		}
		log.Printf("\n\nVALUESKEYS %#v\n\n", valueskey)
		//utils.WaitForKeypress()
	}

	// then adds all unique keys and values to TheDataDict
	for k, _ := range firstnames {
		TheDataDict.Fields["FIRSTNAMES"] = append(TheDataDict.Fields["FIRSTNAMES"], k)
	}
	log.Printf("\n\nfirstnames %#v\n\n", TheDataDict.Fields["FIRSTNAMES"])

	for k, _ := range middlenames {
		TheDataDict.Fields["MIDDLENAMES"] = append(TheDataDict.Fields["MIDDLENAMES"], k)
	}
	log.Printf("\n\nmiddlenames %#v\n\n", TheDataDict.Fields["MIDDLENAMES"])

	for k, _ := range lastnames {
		TheDataDict.Fields["LASTNAMES"] = append(TheDataDict.Fields["LASTNAMES"], k)
	}
	log.Printf("\n\nlastnames %#v\n\n", TheDataDict.Fields["LASTNAMES"])

	for key, _ := range keys {
		for val, keyofv := range valueskey {
			if key == keyofv {
				TheDataDict.Fields[key] = append(TheDataDict.Fields[key], val)
			}
		}
	}

	TheDataDict.LastModified = time.Now()
	TheDataDict.KeyStats = keys
	TheDataDict.ValuesStats = values
	json.Marshal(TheDataDict)
	log.Printf("%#v", TheDataDict)
	utils.WaitForKeypress()
}

//
func UpdateTheDict(key string, value string) {
	TheDataDict.Update(strings.ToUpper(key), strings.ToLower(value))
}

func (d *DataDict) ToJson() string {
	jsoned, err := json.Marshal(*d)
	errors.FailOn(err, "DataDict::ToJson json.Marshal")
	return fmt.Sprintf("%s", string(jsoned))
}
