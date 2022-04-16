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
	KeyStats     map[string]int      `json:"keys_stats"`
	ValuesStats  map[string]int      `json:"values_stats"`
	ValuesKey    map[string]string   `json:"values_key"`
}

var TheDataDict DataDict

func init() {
	TheDataDict = DataDict{
		LastModified: time.Now(),

		Fields: map[string][]string{
			"RAWNAMES":     []string{},
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
		ValuesKey:   make(map[string]string),
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
	// collect all keys and values TheDataDict
	for _, m := range musiciansmap {
		for k, v := range m.Fields {
			log.Printf("KEY   %#v            VALUES  %#v\n", k, v)
			TheDataDict.KeyStats[k]++
			TheDataDict.ValuesStats[v]++
			_, ok := TheDataDict.Fields[k]
			if ok {
				TheDataDict.Fields[k] = append(TheDataDict.Fields[k], v)
			} else {
				TheDataDict.Fields[k] = []string{v}
			}
		}
	}

	// collect all rawnames TheDataDict
	for _, m := range musiciansmap {
		//log.Printf("KEY   %#v            VALUES  %#v\n", k, v)
		TheDataDict.KeyStats["RAWNAMES"]++
		TheDataDict.ValuesStats[m.RawName]++

		_, ok := TheDataDict.Fields["RAWNAMES"]
		if ok {
			TheDataDict.Fields["RAWNAMES"] = append(TheDataDict.Fields["RAWNAMES"], m.RawName)
		} else {
			TheDataDict.Fields["RAWNAMES"] = []string{m.RawName}
		}
	}

	// collect all FNames TheDataDict
	for _, m := range musiciansmap {
		//log.Printf("KEY   %#v            VALUES  %#v\n", k, v)
		if m.FName != Defaults.FName {
			TheDataDict.KeyStats["FIRSTNAMES"]++
			TheDataDict.ValuesStats[m.FName]++
			_, ok := TheDataDict.Fields["FIRSTNAMES"]
			if ok {
				TheDataDict.Fields["FIRSTNAMES"] = append(TheDataDict.Fields["FIRSTNAMES"], m.FName)
			} else {
				TheDataDict.Fields["FIRSTNAMES"] = []string{m.FName}
			}
		}
	}

	// collect all MNames TheDataDict
	for _, m := range musiciansmap {
		//log.Printf("KEY   %#v            VALUES  %#v\n", k, v)
		if m.MName != Defaults.MName {
			TheDataDict.KeyStats["MIDDLENAMES"]++
			TheDataDict.ValuesStats[m.MName]++
			_, ok := TheDataDict.Fields["MIDDLENAMES"]
			if ok {
				TheDataDict.Fields["MIDDLENAMES"] = append(TheDataDict.Fields["MIDDLENAMES"], m.MName)
			} else {
				TheDataDict.Fields["MIDDLENAMES"] = []string{m.MName}
			}
		}
	}

	// collect all LNames TheDataDict
	for _, m := range musiciansmap {
		//log.Printf("KEY   %#v            VALUES  %#v\n", k, v)
		if m.LName != Defaults.LName {
			TheDataDict.KeyStats["LASTNAMES"]++
			TheDataDict.ValuesStats[m.LName]++
			_, ok := TheDataDict.Fields["LASTNAMES"]
			if ok {
				TheDataDict.Fields["LASTNAMES"] = append(TheDataDict.Fields["LASTNAMES"], m.LName)
			} else {
				TheDataDict.Fields["LASTNAMES"] = []string{m.LName}
			}
		}
	}

	TheDataDict.LastModified = time.Now()
	TheDataDict.ToJson()
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

//// BuildTheDataDict build maps of lists of key-values apprearing in Musician struct fields
//func BuildTheDataDict(musiciansmap MusiciansMap) {
//	firstnames := make(map[string]int)
//	middlenames := make(map[string]int)
//	lastnames := make(map[string]int)
//
//	keys := make(map[string]int)
//	values := make(map[string]int)
//	valueskey := make(map[string]string)
//
//	// collect all first, middle and last names... to TheDataDict
//	// ...then collect all keys and values
//	for _, mv := range musiciansmap {
//		//// fname, mname, lname
//		log.Printf("Musician RAWNAME: %#v", mv.RawName)
//		fname := utils.NormalizeKey(mv.FName)
//		mname := utils.NormalizeKey(mv.MName)
//		lname := utils.NormalizeKey(mv.LName)
//
//		// commented out, we want to count the unassigned names too
//		//if fname != Defaults.FName {
//		firstnames[fname]++
//		// if mname != Defaults.MName {
//		middlenames[mname]++
//		// if lname != Defaults.LName {
//		lastnames[lname]++
//
//		////Fields[key]value
//		for key, val := range mv.Fields {
//			log.Printf("KEY   %#v            VALUES  %#v\n", key, val)
//			//k := utils.NormalizeKey(key)
//			//v := utils.NormalizeValue(val)
//			//keys[k]++
//			//values[v]++
//			//valueskey[v] = k
//			keys[key]++
//			values[val]++
//			valueskey[val] = key
//
//			c, ok := keys[key]
//			if ok {
//				if c == 1 {
//					log.Printf("KEY: %#v", key)
//				}
//
//			} else {
//				log.Printf("NEW KEY: %#v", key)
//
//			}
//
//			c, ok = values[val]
//			if ok {
//				if c == 1 {
//					log.Printf("VALUE: %#v", val)
//				}
//
//			} else {
//				log.Printf("NEW VALUE: %#v", val)
//
//			}
//		}
//		//log.Printf("\n\nVALUESKEYS %#v\n\n", valueskey)
//		utils.WaitForKeypress()
//	}
//
//	// then adds all unique keys and values to TheDataDict
//	for k, _ := range firstnames {
//		TheDataDict.Fields["FIRSTNAMES"] = append(TheDataDict.Fields["FIRSTNAMES"], k)
//	}
//	//log.Printf("\n\nfirstnames %#v\n\n", TheDataDict.Fields["FIRSTNAMES"])
//
//	for k, _ := range middlenames {
//		TheDataDict.Fields["MIDDLENAMES"] = append(TheDataDict.Fields["MIDDLENAMES"], k)
//	}
//	//log.Printf("\n\nmiddlenames %#v\n\n", TheDataDict.Fields["MIDDLENAMES"])
//
//	for k, _ := range lastnames {
//		TheDataDict.Fields["LASTNAMES"] = append(TheDataDict.Fields["LASTNAMES"], k)
//	}
//	//log.Printf("\n\nlastnames %#v\n\n", TheDataDict.Fields["LASTNAMES"])
//
//	for key, _ := range keys {
//		for val, keyofv := range valueskey {
//			if key == keyofv {
//				_, ok := TheDataDict.Fields[key]
//				if ok {
//					TheDataDict.Fields[key] = append(TheDataDict.Fields[key], val)
//					log.Printf("KEY   %#v            VALUES  %#v\n", key, val)
//				} else {
//					TheDataDict.Fields[key] = []string{val}
//					log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!! CREATING   KEY   %#v      FOR       VALUES  %#v\n", key, val)
//				}
//
//			}
//		}
//	}
//
//	TheDataDict.LastModified = time.Now()
//	TheDataDict.KeyStats = keys
//	TheDataDict.ValuesStats = values
//	TheDataDict.ToJson()
//	log.Printf("%#v", TheDataDict)
//	utils.WaitForKeypress()
//}
