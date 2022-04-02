package errors

import "errors"

func FailOn(err error, desc string) {
	const SEP = "\n\n"
	if err != nil {
		//log.Printf(SEP+"ERR: %s"+SEP, err)
		panic(err)
	} else {
		//log.Printf(SEP+"INFO: %s"+SEP, desc)
	}
}

func Assert(ok bool, desc string) {
	const SEP = "\n\n"
	if !ok {
		//log.Printf(SEP+"ERR: FAILED ON NOT OK: %s"+SEP, desc)
		panic(errors.New(desc))
	} else {
		//log.Printf(SEP+"INFO: OK'ed %s"+SEP, desc)
	}
}
