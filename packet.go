package gott

import (
	"bytes"
)

const (
	CONNECT     = uint8(1)
	CONNACK     = uint8(2)
	PUBLISH     = uint8(3)
	PUBACK      = uint8(4)
	PUBREC      = uint8(5)
	PUBREL      = uint8(6)
	PUBCOMP     = uint8(7)
	SUBSCRIBE   = uint8(8)
	SUBACK      = uint8(9)
	UNSUBSCRIBE = uint8(10)
	UNSUBACK    = uint8(11)
	PINGREQ     = uint8(12)
	PINGRESP    = uint8(13)
	DISCONNECT  = uint8(14)
)

// StaticHeader holds all of the data that is found
// in the mqtt static header
type FixedHeader struct {
	MessageType uint8
	Dup         bool
	Qos         uint8
	Retain      bool
	Remaining   int
}


// *************** VARIABLE HEADERS ******************** //

// VariableHeader is an interface that all of
// the different types of variable headers
type VariableHeader interface {
	Type() string
	Bytes() []byte
}


// ConnectHeader holds all of the data for the
// variable header for CONNECT messages
type ConnectHeader struct {
	ProtoName    string
	ProtoVersion uint8
	CleanSession bool
	Will         bool
	WillQos      uint8
	WillRetain   bool
	Pass         bool
	User         bool
	KeepAlive    int
}

// ConnackHeader holds all of the data for the
// variable header for CONNACK messages
type ConnackHeader struct {
	ReturnCode uint8
}

// PublishHeader holds all of the data for the
// variable header for PUBLISH messages
type PublishHeader struct {
	Topic string
	MessageID int
}

type MessageIDHeader struct {
	MessageID int
}

type EmptyVariableHeader struct{}

// ************ PAYLOADS ******************//

// Payload is an interface that can be any kind of data
// you want to send as long as you can write
// it to bytes before we send it on the wire
type Payload interface {
	Bytes() []byte
	MessageType() uint8
}

// PayloadBuffer is an alias of []byte so we can attach
// functions to received payloads
type PayloadBuffer []byte

type EmptyPayload struct{}

type ConnectPayload struct {
	ClientID string
	WillTopic string
	WillMessage string
	Username string
	Password string
}

type PublishPayload struct {
	Message string
}

type Subscription struct {
	Topic string
	Qos uint8
}

type SubscribePayload struct {
	Subscriptions []Subscription
}

type SubackPayload struct {
	Qoss []uint8
}

type UnsubscribePayload struct {
	Subscriptions []Subscription
}


// ************** MESSAGE ******************//

// Message holds everything that a message can be
// message is the construct that the api will mainly
// deal with messages as the construct with which to
// attach functionality
type Message struct {
	FixedHeader    FixedHeader
	VariableHeader VariableHeader
	Payload        Payload
}


/********************ENCODING**************************/


// Bytes writes all the data in a message to bytes
// this is to prepare the message for sending
func (m Message) Bytes() []byte {
	bytesToReturn := make([]byte, 0)
	vh := m.VariableHeader.Bytes()
	pl := m.Payload.Bytes()
	m.FixedHeader.Remaining = len(vh) + len(pl)
	bytesToReturn = append(bytesToReturn, m.FixedHeader.Bytes()...)
	bytesToReturn = append(bytesToReturn, vh...)
	bytesToReturn = append(bytesToReturn, pl...)
	return bytesToReturn
}

// Bytes writes all the data in the fixed header
// as defined here: http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#fixed-header
func (f FixedHeader) Bytes() []byte {
	bytesToReturn := make([]byte, 0)
	var byte1 byte

	// Set Message type
	byte1 |= f.MessageType

	// Set Dup Flag
	if f.Dup {
		byte1 |= (1 << 4)
	}

	// Set Qos
	byte1 |= (f.Qos << 5)

	// Set Retain
	if f.Retain {
		byte1 |= (1 << 7)
	}

	bytesToReturn = append(bytesToReturn, byte1)

	// Set Remaining length
	lengthBytes := EncodeRemainingLength(f.Remaining)

	bytesToReturn = append(bytesToReturn, lengthBytes...)
	return bytesToReturn
}

func (f FixedHeader) BytesWithBuffer() []byte {
	var byte1 byte

	// Set Message type
	byte1 |= f.MessageType

	// Set Dup Flag
	if f.Dup {
		byte1 |= (1 << 4)
	}

	// Set Qos
	byte1 |= (f.Qos << 5)

	// Set Retain
	if f.Retain {
		byte1 |= (1 << 7)
	}

	lengthBytes := EncodeRemainingLengthBuffer(f.Remaining)
	numOfBytes := 1 + len(lengthBytes)
	buf := new(bytes.Buffer)
	buf.Grow(numOfBytes)
	err := buf.WriteByte(byte1)
	if err != nil {
		panic(err)
	}
	_, err = buf.Write(lengthBytes)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
// ************ VARIABLE HEADER ENCODING ************//

// Bytes writes all the data in the variable header for
// CONNECT messages
func (c ConnectHeader) Bytes() []byte {
	bytesToReturn := make([]byte, 0)
	ProtoNameBytes := MqttUTF8Bytes(c.ProtoName)
	ProtoVersionByte := byte(c.ProtoVersion)

	FlagByte := byte(0)

	if c.CleanSession {
		FlagByte |= (1 << 1)
	}

	if c.Will {
		FlagByte |= (1 << 2)
	}

	FlagByte |= (c.WillQos << 3)

	if c.WillRetain {
		FlagByte |= (1 << 5)
	}

	if c.Pass {
		FlagByte |= (1 << 6)
	}

	if c.User {
		FlagByte |= (1 << 7)
	}
	msb := msb(c.KeepAlive)
	lsb := lsb(c.KeepAlive)

	bytesToReturn = append(bytesToReturn, ProtoNameBytes...)
	bytesToReturn = append(bytesToReturn, ProtoVersionByte, FlagByte, msb, lsb)

	return bytesToReturn
}

// Type returns a string of the type of variable header
func (c ConnectHeader) Type() string {
	return "CONNECT"
}

// Bytes returns the bytes of the variable header
// for the CONNACK message
func (c ConnackHeader) Bytes() []byte {
	empty := byte(0)
	retCodeByte := byte(c.ReturnCode)
	return []byte{empty, retCodeByte}
}

// Tyoe returns a string of the type of variable header
func (c ConnackHeader) Type() string {
	return "CONNACK"
}

// Bytes returns the bytes of the variable header
// for the PUBLISH message
func (p PublishHeader) Bytes() []byte {
	topicBytes := MqttUTF8Bytes(p.Topic)
	if p.MessageID != 0 {
		msb := msb(p.MessageID)
		lsb := lsb(p.MessageID)
		topicBytes = append(topicBytes, msb, lsb)
	}
	return topicBytes
}

// Type returns the string of the type of variable header
func (p PublishHeader) Type() string {
	return "PUBLISH"
}

// Bytes returns an empty byte slice
func (e EmptyVariableHeader) Bytes() []byte {
	return []byte{}
}

// Type returns the string "EMPTY"
func (e EmptyVariableHeader) Type() string {
	return "EMPTY"
}

func (m MessageIDHeader) Bytes() []byte {
	msb := msb(m.MessageID)
	lsb := lsb(m.MessageID)
	return []byte{msb, lsb}
}

func (m MessageIDHeader) Type() string {
	return "MESSAGE"
}

// EncodeRemainingLength encodes an int into the
// Remaining length encoding format as defined in the spec
func EncodeRemainingLength(length int) []byte {
	if length == 0 {
		panic("INVALID REMAINING LENGTH")
	}
	encodedLength := make([]byte, 0)
	for length > 0 {
		mod := length % 128
		digit := uint8(mod)
		length = length / 128
		if length > 0 {
			digit |= 0x80
		}
		encodedLength = append(encodedLength, digit)
	}
	return encodedLength
}

func EncodeRemainingLengthBuffer(length int) []byte {
	if length == 0 {
		panic("INVALID REMAINING LENGTH")
	}
	var numOfBytes int;
	switch {
	case length < 128 :
		numOfBytes = 1
	case length < 16384:
		numOfBytes = 2
	case length < 2097152:
		numOfBytes = 3
	default:
		numOfBytes = 4
	}
	buf := new(bytes.Buffer)
	buf.Grow(numOfBytes)
	for length > 0 {
		mod := length % 128
		digit := uint8(mod)
		length = length / 128
		if length > 0 {
			digit |= 0x80
		}
		err := buf.WriteByte(digit)
		if err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}

// **************** PAYLOAD ENCODING ******************** //

// Bytes Returns the []byte value of the
// PayloadBuffer
func (p PayloadBuffer) Bytes() []byte {
	return []byte(p)
}

func (p PayloadBuffer) MessageType() uint8 {
	return uint8(0)
}

func (c ConnectPayload) Bytes() []byte {
	bytesToReturn := make([]byte, 0)

	clientIDBytes := MqttUTF8Bytes(c.ClientID)
	bytesToReturn = append(bytesToReturn, clientIDBytes...)

	if c.WillTopic != "" {
		willTopicBytes := MqttUTF8Bytes(c.WillTopic)
		bytesToReturn = append(bytesToReturn, willTopicBytes... )
	}

	if c.WillMessage != "" {
		willMessageBytes := MqttUTF8Bytes(c.WillMessage)
		bytesToReturn = append(bytesToReturn, willMessageBytes... )
	}

	if c.Username != "" {
		usernameBytes := MqttUTF8Bytes(c.Username)
		bytesToReturn = append(bytesToReturn, usernameBytes... )
	}

	if c.Password != "" {
		passwordBytes := MqttUTF8Bytes(c.Password)
		bytesToReturn = append(bytesToReturn, passwordBytes... )
	}

	return bytesToReturn
}

func (c ConnectPayload) MessageType() uint8 {
	return CONNECT
}

/********************DECODING****************************/

// FixedHeaderFromBytes takes the first few bytes from
// an incoming packet and parses them into the FixedHeader
// Returns the fixed header and the number of bytes
// That were parsed
func FixedHeaderFromBytes(b []byte) (FixedHeader, int) {
	var messageType uint8
	messageType |= b[0]
	messageType &^= (15 << 4)

	dup := b[0]&(1<<4) > 0

	var qos uint8
	qos |= (b[0] >> 5)
	qos &^= (63 << 2)

	retain := b[0]&(1<<7) > 0

	remaining, byteLength := DecodeRemainingLength(b[1:])

	fh := FixedHeader{
		MessageType: messageType,
		Dup:         dup,
		Qos:         qos,
		Retain:      retain,
		Remaining:   remaining,
	}
	return fh, byteLength + 1
}


// ********** VARIABLE HEADER DECODING *********//

// ConnectHeaderFromBytes takes the first few bytes from
// an incoming packet who's message type is CONNECT and
// parses them into the ConnectHeader
func ConnectHeaderFromBytes(b []byte) (ConnectHeader, int) {

	protoNameLength := strLen(b[0], b[1])
	protoName := string(b[2 : protoNameLength+2])

	protoVersion := uint8(b[protoNameLength+2])

	flagByte := b[protoNameLength+3]

	clean := flagByte&(1<<1) > 0
	will := flagByte&(1<<2) > 0

	var willQos uint8
	willQos |= (flagByte >> 3)
	willQos &^= (63 << 2)

	willRetain := flagByte&(1<<5) > 0
	pass := flagByte&(1<<6) > 0
	user := flagByte&(1<<7) > 0

	keepAlive := strLen(b[protoNameLength+4], b[protoNameLength+5])

	ch := ConnectHeader{
		ProtoName:    protoName,
		ProtoVersion: protoVersion,
		CleanSession: clean,
		Will:         will,
		WillQos:      willQos,
		WillRetain:   willRetain,
		Pass:         pass,
		User:         user,
		KeepAlive:    keepAlive,
	}

	return ch, protoNameLength + 6
}

// TODO make this a real function
func ConnackHeaderFromBytes(b []byte) (ConnackHeader, int) {
	return ConnackHeader{}, 0
}

// TODO make this a real function
func PublishHeaderFromBytes(b []byte) (PublishHeader, int) {
	return PublishHeader{}, 0
}

// DecodeRemainingLength decodes the encoded remaining
// length from an incoming packet .Returns the remaining
// length and how many bytes were parsed.
func DecodeRemainingLength(b []byte) (int, int) {
	multiplier := 1
	value := 0
	cur := 0
	last := 0
	for cur < len(b) && (b[last]&128) != 0 {
		value += int(b[cur]&127) * multiplier
		multiplier *= 128
		last = cur
		cur++
	}
	return value, cur + 1
}

// VariableHeaderFromBytes parses bytes from an incoming
// packet and returns a VariableHeader and how many
// bytes were parsed
func VariableHeaderFromBytes(b []byte, messageType uint8) (VariableHeader, int) {

	switch messageType {
	case CONNECT:
		return ConnectHeaderFromBytes(b)
	case CONNACK:
		return ConnackHeaderFromBytes(b)
	case PUBLISH:
		return PublishHeaderFromBytes(b)
	default:
		return EmptyVariableHeader{}, 0
	}
}

// MessageFromBytes parses an incoming packet and returns
// a message
func MessageFromBytes(b []byte) Message {
	fixedHeader, fixedLength := FixedHeaderFromBytes(b)
	variableHeader, variableLength := VariableHeaderFromBytes(b[fixedLength:], fixedHeader.MessageType)
	payload := b[fixedLength+variableLength:]
	return Message{
		FixedHeader:    fixedHeader,
		VariableHeader: variableHeader,
		Payload:        PayloadBuffer(payload),
	}
}

/***************** HELPERS ********************/

// lsb takes a 16 bit int and returns the least
// significant byte
func lsb(i int) byte {
	if i > 65535 {
		panic("NUMBER IS TOO LARGE MUST BE ABLE TO FIT IN 16 bit unsigned")
	}
	lsb := i % 256
	return uint8(lsb)
}

// msb takes a 16 bit int and returns the most
// significant byte
func msb(i int) byte {
	if i > 65535 {
		panic("NUMBER IS TOO LARGE MUST BE ABLE TO FIT IN 16 bit unsigned")
	}
	msb := i / 256
	return uint8(msb)
}

// strLen takes the msb and lsb bytes and returns
// an int of those two multiplied
func strLen(msb, lsb byte) int {
	if msb == uint8(0) {
		return int(lsb)
	} else {
		return (int(msb) * 256) + int(lsb)
	}
	return 0
}

// MqttUTF8Bytes takes a string and encodes it like:
// http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#utf-8
func MqttUTF8Bytes(s string) []byte {
	msb := msb(len(s))
	lsb := lsb(len(s))

	b := []byte{msb, lsb}
	b = append(b, []byte(s)...)
	return b
}
