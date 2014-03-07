package gott

import (
	"reflect"
	"testing"
	"log"
)

func TestForReality(t *testing.T) {
	if true == false {
		t.Fatal("RUNNNNNNNN!!!!!")
	}
}

func TestFixedHeaderEncodeDecode(t *testing.T) {
	fh := FixedHeader{
		MessageType: PINGREQ,
		Dup:         true,
		Qos:         1,
		Retain:      true,
		Remaining:   321,
	}
	b := fh.Bytes()
	testFH, bl := FixedHeaderFromBytes(b)
	if testFH != fh {
		t.Error("fixed headers did not match")
	}
	if bl != 4 {
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

func TestLengthFuncs(t *testing.T) {
	testInt := 256
	msb := msb(testInt)
	lsb := lsb(testInt)
	newInt := strLen(msb, lsb)
	if newInt != testInt {
		t.Error("fail")
	}
}

func TestConnectPacket(t *testing.T) {
	fh := FixedHeader{
		MessageType: CONNECT,
	}
	ch := ConnectHeader{
		ProtoName:    "MQIsdp",
		ProtoVersion: 3,
		CleanSession: true,
		Will:         true,
		WillQos:      1,
		WillRetain:   false,
		Pass:         true,
		User:         true,
		KeepAlive:    10,
	}
	pl := PayloadBuffer("This is the payload")
	message := Message{
		FixedHeader:    fh,
		VariableHeader: ch,
		Payload:        pl,
	}
	b := message.Bytes()
	testMessage := MessageFromBytes(b)
	if !reflect.DeepEqual(testMessage, message) {
		t.Error("encoding and decoding CONNECT packet failed")
	}
}

func TestConnectMessage(t *testing.T) {
	con := NewConnect(false, false, false, 0, "", "", 10, "", "")
	log.Printf("MESSAGE: %+v\n", con)
}
