package dl_cbf

import (
	"fmt"
	"reflect"
)

type Fingerprint interface {
	IsFingerprint()
	String() string
}

// Support for 16,32,64-bit fp
type fingerprint16 uint16
type fingerprint32 uint32
type fingerprint64 uint64

func (_ fingerprint16) IsFingerprint() {}
func (_ fingerprint32) IsFingerprint() {}
func (_ fingerprint64) IsFingerprint() {}

func (fp16 fingerprint16) String() string {
	return fmt.Sprintf("fp{%d}", fp16)
}

func (fp32 fingerprint32) String() string {
	return fmt.Sprintf("fp{%d}", fp32)
}

func (fp64 fingerprint64) String() string {
	return fmt.Sprintf("fp{%d}", fp64)
}

func CompareFingerprints(fp1, fp2 Fingerprint) uint8 {

	if fp1 == nil {
		if fp2 == nil {
			return 1
		}
		return 0
	}

	if fp2 == nil {
		return 2
	}

	if reflect.TypeOf(fp1) != reflect.TypeOf(fp2) {
		panic("Cannot compare fingerprints of different types")
	}

	switch fp1.(type) {
	case fingerprint16:
		if fp1.(fingerprint16) == fp2.(fingerprint16) {
			return 1
		}
		if fp1.(fingerprint16) < fp2.(fingerprint16) {
			return 0
		}
		return 2
	case fingerprint32:
		if fp1.(fingerprint32) == fp2.(fingerprint32) {
			return 1
		}
		if fp1.(fingerprint32) < fp2.(fingerprint32) {
			return 0
		}
		return 2
	case fingerprint64:
		if fp1.(fingerprint64) == fp2.(fingerprint64) {
			return 1
		}
		if fp1.(fingerprint64) < fp2.(fingerprint64) {
			return 0
		}
		return 2
	default:
		panic("Unsupported fingerprint type")
	}
}
