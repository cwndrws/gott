package gott

import (
	"testing"
)

func TestForReality(t *testing.T) {
	if true == false {
		t.Fatal("RUNNNNNNNN!!!!!")
	}
}

func TestFixedHeaderEncodeDecode(t *testing.T) {
	fh := FixedHeader{
		MessageType: PINGREQ,
		Dup: true,
		Qos: 1,
		Retain: true,
		Remaining: 321,
	}
	b := fh.Bytes()
	testFH, bl := FixedHeaderFromBytes(b)
	if testFH != fh {
		t.Error("fixed headers did not match")
	}
	if bl != 3 {
		t.Error("Incorrect byte length")
	}
}

func TestLsb(t *testing.T) {
	l := lsb(65535)
	if l != uint8(255) {
		t.Errorf("lsb fail %d\n", l)
	}
}

func TestMsb(t *testing.T) {
	m := msb(65535)
	if m != uint8(255) {
		t.Error("msb fail")
	}

	m = msb(255)
	if m != uint8(0) {
		t.Error("msb fail")
	}
}
