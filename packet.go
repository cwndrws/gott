package gott

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

// VariableHeader holds all of the data that can be found
// in the mqtt variable header
// TODO maybe make this an interface since
// It can look different a lot of places
type VariableHeader struct{}

// Payload is an interface that can be any kind of data
// you want to send as long as you can write
// it to bytes before we send it on the wire
type Payload interface {
	Bytes() []byte
}

// Message holds everything that a message can be
// message is the construct that the api will mainly
// deal with messages as the construct with which to
// attach functionality
type Message struct {
	FixedHeader    FixedHeader
	VariableHeader VariableHeader
	Payload        Payload
}

// Bytes writes all the data in a message to bytes
// this is to prepare the message for sending
func (m Message) Bytes() []byte {
	bytesToReturn := make([]byte, 0)
	bytesToReturn = append(bytesToReturn, m.FixedHeader.Bytes()...)
	bytesToReturn = append(bytesToReturn, m.VariableHeader.Bytes()...)
	bytesToReturn = append(bytesToReturn, m.Payload.Bytes()...)
	return bytesToReturn
}

/********************ENCODING**************************/

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

// Bytes writes the variable header to a byte slice
func (v VariableHeader) Bytes() []byte {
	return []byte{}
}

/********************DECODING****************************/

func FixedHeaderFromBytes(b []byte) FixedHeader {
	var messageType uint8
	messageType |= b[0]
	messageType &^= (15 << 4)

	dup := b[0] & (1 << 4) > 0


	var qos uint8
	qos |= (b[0] >> 5)
	qos &^= (63 << 2)

	retain := b[0] & (1 << 7) > 0

	remaining := DecodeRemainingLength(b[1:])

	return FixedHeader{
		MessageType: messageType,
		Dup: dup,
		Qos: qos,
		Retain: retain,
		Remaining: remaining,
	}
}

func DecodeRemainingLength(b []byte) int {
	multiplier := 1
	value := 0
	cur := 0
	for cur < len(b) {
		value += int(b[cur] & 127) * multiplier
		multiplier *= 128
		cur++
	}
	return value
}
