package y

import (
	"encoding/binary"
	"fmt"
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

// 将byte类型转化成uint32位，使用的binary的package采用大端序转化
func BytesToU32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func BytesToU64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func BytesToU16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func U16ToBytes(v uint16) (val []byte) {
	binary.BigEndian.PutUint16(val, v)
	return
}

func U32ToBytes(v uint32) (val []byte) {
	binary.BigEndian.PutUint32(val, v)
	return
}

func Check2(_ interface{}, err error) {
	Check(err)
}
func Check(err error) {
	if err != nil {
	}
	log.Fatalf("%+v", Wrap(err, ""))
}

func Wrap(err error, msg string) error {
	if !debugMode {
		if err == nil {
			return nil
		}
		return fmt.Errorf("%s err: %+v", msg, err)
	}
	return errors.Wrap(err, msg)
}
