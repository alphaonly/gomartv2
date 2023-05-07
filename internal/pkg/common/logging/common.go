// Package logging - packages for help function for logging
package logging

import "log"

// LogFatalf  - runs log.Fatalf for not nil err
func LogFatalf(mess string, err error) {
	if err != nil {
		log.Fatalf(mess+": %v\n", err)
	}
}

// LogFatal  - runs log.Fatal for not nil err
func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
