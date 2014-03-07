package gott

import (
	"fmt"
)

const (
	PROTONAME = "MQIsdp"
	PROTOVERSION = uint8(3)
)

var (
	MessageIDNum = 0
)

func MessageID() int {
	MessageIDNum++
	return MessageIDNum
}

func ClientID() string {
	// TODO make this an actual random string
	return fmt.Sprintf("TEMPCLIENTID%d", MessageID())
}

func NewConnect(will, willRetain, cleanSession bool, willQos uint8, username, password string, keepAlive int, willTopic, willMessage string) Message {
	fh := FixedHeader{
		MessageType: CONNECT,
	}
	ch := ConnectHeader{
		ProtoName: PROTONAME,
		ProtoVersion: PROTOVERSION,
		CleanSession: cleanSession,
		Will: will,
		WillQos: willQos,
		WillRetain: willRetain,
		Pass: password != "",
		User: username != "",
		KeepAlive: keepAlive,
	}
	cp := ConnectPayload{
		ClientID: ClientID(),
		WillTopic: willTopic,
		WillMessage: willMessage,
		Username: username,
		Password: password,
	}
	connectMessage := Message{
		FixedHeader: fh,
		VariableHeader: ch,
		Payload: cp,
	}
	return connectMessage
}
