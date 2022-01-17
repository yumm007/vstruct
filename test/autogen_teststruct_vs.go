// Code generated by vstruct; DO NOT EDIT.

package teststruct

import (
	"bytes"
	"encoding/binary"
	"github.com/yumm007/gohash"
)


func (s *Simples)encodeToBuffer(buf *bytes.Buffer) error {
	if err := binary.Write(buf, binary.LittleEndian, &s.Id); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, &s.NameLen); err != nil {
		return err
	}
	for i := 0; i < int(s.NameLen); i++ {
		if err := binary.Write(buf, binary.LittleEndian, &s.Name); err != nil {
			return err
		}
	}

	return nil
}

func (s *Simples)Encode(buf *bytes.Buffer) ([]byte, error) {
	if buf == nil {
		buf = new(bytes.Buffer)
	} else {
		buf.Reset()
	}
	if err := s.encodeToBuffer(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Simples)decodeFromBuffer(buf *bytes.Buffer) error {
	if err := binary.Read(buf, binary.LittleEndian, &s.Id); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &s.NameLen); err != nil {
		return err
	}
	for i := 0; i < int(s.NameLen); i++ {
		var ele uint8
		if err := binary.Read(buf, binary.LittleEndian, &s.Name); err != nil {
			return err
		}
		s.Name = append(s.Name, ele)
	}

	return nil
}

func (s *Simples)Decode(payload []byte) error {
	buf := bytes.NewBuffer(payload)
	return s.decodeFromBuffer(buf)
}

func (u *UnionSim)encodeToBuffer(buf *bytes.Buffer) error {
	if err := binary.Write(buf, binary.LittleEndian, &u.Len); err != nil {
		return err
	}
	for i := 0; i < int(u.Len); i++ {
		if err := u.Arr[i].encodeToBuffer(buf); err != nil {
			return err
		}
	}
	if err := binary.Write(buf, binary.LittleEndian, &u.Crc); err != nil {
		return err
	}

	return nil
}

func (u *UnionSim)Encode(buf *bytes.Buffer) ([]byte, error) {
	if buf == nil {
		buf = new(bytes.Buffer)
	} else {
		buf.Reset()
	}
	if err := u.encodeToBuffer(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (u *UnionSim)decodeFromBuffer(buf *bytes.Buffer) error {
	if err := binary.Read(buf, binary.LittleEndian, &u.Len); err != nil {
		return err
	}
	for i := 0; i < int(u.Len); i++ {
		var ele Simples
		if err := ele.decodeFromBuffer(buf); err != nil {
			return err
		}
		u.Arr = append(u.Arr, ele)
	}
	if err := binary.Read(buf, binary.LittleEndian, &u.Crc); err != nil {
		return err
	}

	return nil
}

func (u *UnionSim)Decode(payload []byte) error {
	buf := bytes.NewBuffer(payload)
	return u.decodeFromBuffer(buf)
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

func (u *UnionSim2)Encode(buf *bytes.Buffer) ([]byte, error) {
	if buf == nil {
		buf = new(bytes.Buffer)
	} else {
		buf.Reset()
	}
	if err := u.encodeToBuffer(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func (u *UnionSim2)Decode(payload []byte) error {
	buf := bytes.NewBuffer(payload)
	return u.decodeFromBuffer(buf)
}