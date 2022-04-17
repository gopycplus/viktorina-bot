package utils

import "log"

func Check(err error) {
	if err != nil {
		//panic(err)
		log.Println(err)
	}
}
