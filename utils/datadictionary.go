package utils

type DataDict struct {
	FirstNames  []string            `json:"first_names"`
	MiddleNames []string            `json:"middle_names"`
	LastNames   []string            `json:"last_names"`
	Fields      map[string][]string `json:"fields"`
}

var TheDataDict = DataDict{
	FirstNames:  []string{},
	MiddleNames: []string{},
	LastNames:   []string{},
	Fields:      map[string][]string{},
}

func (d *DataDict) Update(key string, value string) {
	switch key {
	case "FirstName":
		{
			d.FirstNames = append(d.FirstNames, value)
		}
	case "MiddleName":
		{
			d.MiddleNames = append(d.MiddleNames, value)
		}
	case "LastName":
		{
			d.LastNames = append(d.LastNames, value)
		}
	default:
		{
			d.Fields[key] = append(d.Fields[key], value)
		}
	}

}
