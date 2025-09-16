package stunning

import (
	"encoding/binary"
	"fmt"
	"net"
)

type AttributeType uint16

const (
	MAPPED_ADDRESS          = 0x0001
	LEGACY_RESPONSE_ADDRESS = 0x0002
	LEGACY_CHANGE_REQUEST   = 0x0003
	LEGACY_SOURCE_ADDRESS   = 0x0004
	LEGACY_CHANGED_ADDRESS  = 0x0005
	USERNAME                = 0x0006
	LEGACY_PASSWORD         = 0x0007
	MESSAGE_INTEGRITY       = 0x0008
	ERROR_CODE              = 0x0009
	UNKNOWN_ATTRIBUTES      = 0x000A
	LEGACY_REFLECTED_FROM   = 0x000B
	REALM                   = 0x0014
	NONCE                   = 0x0015
	XOR_MAPPED_ADDRESS      = 0x0020
)

type PayloadT interface {
	Length() uint16
	Encode() []byte
}

type Attribute struct {
	Type     AttributeType
	Length   uint16
	RawValue []byte
	Value    PayloadT
}

func DecodeAttributesFromBytes(rawAttributes []byte) ([]Attribute, error) {
	var attributes []Attribute
	offset := 0

	for offset+4 <= len(rawAttributes) { // Need at least 4 bytes for attr header
		// Type (2 bytes)
		aType := binary.BigEndian.Uint16(rawAttributes[offset : offset+2])

		// Length (2 bytes)
		aLength := binary.BigEndian.Uint16(rawAttributes[offset+2 : offset+4])

		// Check if we have enough bytes for the value
		if offset+4+int(aLength) > len(rawAttributes) {
			return nil, fmt.Errorf("insufficient data for attribute value")
		}

		// Get the value (without padding)
		payloadWithoutPadding := rawAttributes[offset+4 : offset+4+int(aLength)]

		value, err := DecodeAttributeValue(AttributeType(aType), payloadWithoutPadding)
		if err != nil {
			return nil, err
		}

		attributes = append(attributes, Attribute{
			Type:     AttributeType(aType),
			Length:   aLength,
			Value:    value,
			RawValue: rawAttributes[offset : offset+4+int(aLength)], // Include header
		})

		// Move offset: 4 (header) + length + padding
		offset += 4 + int(aLength)

		// STUN attributes are padded to 4-byte boundaries
		padding := (4 - (int(aLength) % 4)) % 4
		offset += padding
	}

	return attributes, nil
}
func DecodeAttributeValue(aType AttributeType, rawValue []byte) (PayloadT, error) {
	switch aType {
	case MAPPED_ADDRESS:
		return MappedAddressFromBytes(rawValue), nil
	}

	return nil, nil
}

func (a Attribute) Encode() []byte {
	data := []byte{}
	data = binary.BigEndian.AppendUint16(data, uint16(a.Type))
	data = binary.BigEndian.AppendUint16(data, a.Length)

	// add padding to the payload
	rawValue := a.Value.Encode()
	var paddedValue []byte
	if len(rawValue)%4 != 0 {
		// Add padding
	}

	data = append(data, paddedValue...)
	return data
}

type WithAttribute func() Attribute

type MappedAddress struct {
	Family  MappedAddressFamily
	Port    uint16
	Address net.IP
}

type MappedAddressFamily uint16

const (
	MappedAddressFamilyIPv4 = MappedAddressFamily(0x0001)
	MappedAddressFamilyIPv6 = MappedAddressFamily(0x0002)
)

func (ma MappedAddress) Length() uint16 {
	return 0
}

func (ma MappedAddress) Encode() []byte {
	var data []byte
	data = binary.BigEndian.AppendUint16(data, uint16(ma.Family))
	data = binary.BigEndian.AppendUint16(data, ma.Port)
	data = append(data, ma.Address...)

	return data
}

func MappedAddressFromBytes(value []byte) MappedAddress {

	family := binary.BigEndian.Uint16(value[:2])
	port := binary.BigEndian.Uint16(value[2:4])
	address := value[4:]

	return MappedAddress{
		Family:  MappedAddressFamily(family),
		Port:    port,
		Address: address,
	}
}

func WithMappedAddress(val *MappedAddress) func() Attribute {
	return func() Attribute {
		return Attribute{
			Type:   MAPPED_ADDRESS,
			Length: val.Length(),
			Value:  *val,
		}
	}
}
