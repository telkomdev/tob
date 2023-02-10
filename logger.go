package tob

import (
	"log"
	"os"
)

var Logger = log.New(os.Stderr, "tob => ", log.Ldate|log.Ltime|log.Llongfile)
