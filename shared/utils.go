package shared

import "log"

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
