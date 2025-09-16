package stunning

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

const MAGIC_COOKIE = 0x2112A442

type MessageMethod uint16

const (
	MethodBinding = 0x0001 // 0b00000001

)

type MessageClass uint8

const (
	ClassRequest         = 0b00
	ClassIndication      = 0b01
	ClassSuccessResponse = 0b10
	ClassErrorResponse   = 0b11
)

type MessageType struct {
	Class  MessageClass
	Method MessageMethod
}

func DecodeMessageTypeFromUint16(encodedMessageType uint16) MessageType {
	classBit0 := GetBit(int(encodedMessageType), 4)
	classBit1 := GetBit(int(encodedMessageType), 8)

	class := 0b00
	class = SetBit(class, 0, classBit0)
	class = SetBit(class, 1, classBit1)

	methodSection1 := encodedMessageType & 0b00000000001111
	methodSection2 := encodedMessageType & 0b00000011100000
	methodSection3 := encodedMessageType & 0b11111000000000

	method := methodSection1 | (methodSection2 >> 1) | (methodSection3 >> 2)

	return MessageType{
		Class:  MessageClass(class),
		Method: MessageMethod(method),
	}
}

func (mt MessageType) ToUint16() uint16 {
	parsedType := 0b000000000000

	// 		0                 1
	// 	2  3  4 5 6 7 8 9 0 1 2 3 4 5
	// +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
	// |M |M |M|M|M|C|M|M|M|C|M|M|M|M|
	// |11|10|9|8|7|1|6|5|4|0|3|2|1|0|
	// +--+--+-+-+-+-+-+-+-+-+-+-+-+-+

	// Because we need to place the bits
	// in that exact order where M is the method
	// and c is the class. We're going to follow the
	// next steps:
	// 1- Insert the class bits into the parsedType
	// 2- separate the sections of the method and
	//    shift them accordingly
	// 3- OR the three parts

	classBit0 := GetBit(int(mt.Class), 0)
	classBit1 := GetBit(int(mt.Class), 1)

	parsedType = SetBit(parsedType, 4, classBit0)
	parsedType = SetBit(parsedType, 8, classBit1)

	section1 := GetBits(int(mt.Method), 0, 4)
	section2 := GetBits(int(mt.Method), 4, 3)
	section3 := GetBits(int(mt.Method), 7, 4)

	// Shift section 2 one bit to the left to accomodate class bit
	section2 = (section2 << 1)
	// Shift section 3 two bits to the left to accomodate section2 shift and class bit
	section3 = (section3 << 2)

	parsedType |= section1 | section2 | section3

	return uint16(parsedType)
}

type Message struct {
	mtype         MessageType
	length        uint16
	cookie        uint32
	transactionId [12]byte
	attributes    []Attribute
	payloadBytes  []byte
}

func NewMessage(
	method MessageMethod,
	class MessageClass,
	attributes ...WithAttribute,
) (Message, error) {

	payload := make([]byte, 0)
	attrs := make([]Attribute, 0)
	for _, attrFn := range attributes {
		attrs = append(attrs, attrFn())
		payload = append(payload, attrFn().Encode()...)
	}

	length := uint16(len(payload))

	transactionId, err := genTransactionId()
	if err != nil {
		return Message{}, err
	}

	return Message{
		mtype: MessageType{
			Method: method,
			Class:  class,
		},
		cookie:        MAGIC_COOKIE,
		length:        length,
		transactionId: transactionId,
		attributes:    attrs,
		payloadBytes:  payload,
	}, nil
}

func DecodeMessageFromBytes(rawMessage []byte) (Message, error) {
	headers := rawMessage[:20]
	rawAttributes := rawMessage[20:]

	mtype := binary.BigEndian.Uint16(headers[:2])
	mlength := binary.BigEndian.Uint16(headers[2:4])
	magicCookie := binary.BigEndian.Uint32(headers[4:8])
	if magicCookie != MAGIC_COOKIE {
		return Message{}, fmt.Errorf("throw error here")
	}

	transactionId := headers[8:]

	attrs, err := DecodeAttributesFromBytes(rawAttributes)
	if err != nil {
		return Message{}, err
	}

	return Message{
		mtype:         DecodeMessageTypeFromUint16(mtype),
		length:        mlength,
		cookie:        MAGIC_COOKIE,
		transactionId: [12]byte(transactionId),
		attributes:    attrs,
		payloadBytes:  rawAttributes,
	}, nil
}

func (m Message) TransactionID() [12]byte {
	return m.transactionId
}

func (m Message) Encode() []byte {
	byteArr := make([]byte, 0)

	byteArr = binary.BigEndian.AppendUint16(byteArr, m.mtype.ToUint16())
	byteArr = binary.BigEndian.AppendUint16(byteArr, m.length)
	byteArr = binary.BigEndian.AppendUint32(byteArr, m.cookie)

	byteArr = append(byteArr, m.transactionId[:]...)
	return byteArr
}

func genTransactionId() ([12]byte, error) {
	var id [12]byte
	_, err := rand.Read(id[:])
	if err != nil {
		return id, err
	}
	return id, nil
}

func (m Message) GetMappedAddress() (MappedAddress, error) {
	var attr Attribute

	for _, at := range m.attributes {
		if at.Type == MAPPED_ADDRESS {
			attr = at
		}
	}

	return attr.Value.(MappedAddress), nil
}
