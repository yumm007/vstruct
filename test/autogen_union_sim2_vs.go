// Code generated by vstruct; DO NOT EDIT.

package teststruct

import (
	"bytes"
	"encoding/binary"
	"github.com/yumm007/gohash"
)


func (u *UnionSim2)Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := u.encodeToBuffer(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (u *UnionSim2)encodeToBuffer(buf *bytes.Buffer) error {
	if err := binary.Write(buf, binary.LittleEndian, &u.Len); err != nil {
		return err
	}
	if err := u.Arr.encodeToBuffer(buf); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, gohash.Crc16ccitt(buf.Bytes())); err != nil {
		return err
	}

	return nil
}

func (u *UnionSim2)Decode(payload []byte) error {
	buf := bytes.NewBuffer(payload)
	return u.decodeFromBuffer(buf)
}

func (u *UnionSim2)decodeFromBuffer(buf *bytes.Buffer) error {
	if err := binary.Read(buf, binary.LittleEndian, &u.Len); err != nil {
		return err
	}
	if err := u.Arr.decodeFromBuffer(buf); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &u.Crc); err != nil {
		return err
	}

	return nil
}
