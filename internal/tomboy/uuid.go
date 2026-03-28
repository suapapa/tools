package tomboy

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
)

// UUID is RFC4122 v4 UUID
type UUID [16]byte

// MakeUUID makes a UUID
func MakeUUID() (u UUID) {
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}

	u[6] = (u[6] & 0x0f) | 0x40 // version 4
	u[8] = (u[8] & 0xbf) | 0x80 // RFC4122

	return
}

func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// UnmarshalXML implements xml.Unmarshaler
func (u *UUID) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v []byte
	d.DecodeElement(&v, &start)
	b := u[:]

	return decodeUUID(b, v)
}

// MarshalXMLAttr implements xml.MarshalerAttr
func (u UUID) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr.Name = name
	attr.Value = u.String()

	return
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr
func (u *UUID) UnmarshalXMLAttr(attr xml.Attr) error {
	v := []byte(attr.Value)
	b := u[:]

	return decodeUUID(b, v)
}

func decodeUUID(b, v []byte) error {
	for _, byteGroup := range []int{8, 4, 4, 4, 12} {
		if v[0] == '-' {
			v = v[1:]
		}

		_, err := hex.Decode(b[:byteGroup/2], v[:byteGroup])
		if err != nil {
			return err
		}

		v = v[byteGroup:]
		b = b[byteGroup/2:]
	}

	return nil
}
