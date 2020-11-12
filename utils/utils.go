package utils

import (
	"log"
	"os"
)

//GetEnvParameter env parameter
func GetEnvParameter(name string) string {
	parameter := os.Getenv(name)
	if parameter == "" {
		log.Fatalf("Missed required env parameter '%s'", name)
	}

	return parameter
}
