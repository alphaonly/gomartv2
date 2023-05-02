package logging

import "log"

// Log
func LogFatalf(mess string, err error) {
	if err != nil {
		log.Fatalf(mess+": %v\n", err)
	}
}

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
