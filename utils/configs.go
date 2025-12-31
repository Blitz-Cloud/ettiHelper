package utils

import "log"

var Log = InitLogger(log.Default(), DEBUG, log.LstdFlags|log.Llongfile, "debug.txt")
