package y

import (
	"github.com/pkg/errors"
	"log"
)

func AssertTrue(b bool) {
	if !b {
		log.Fatalf("%+v", errors.Errorf("Assert failed"))
	}
}

func AssertTruef(b bool, format string, args ...interface{}) {
	if !b {
		log.Fatalf("%+v", errors.Errorf(format, args...))
	}
}
