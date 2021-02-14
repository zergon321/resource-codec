package codec

import (
	"bytes"
	"encoding/gob"
)

// EncodeTag encodes the tag as a byte sequence.
func EncodeTag(tag []string) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(tag)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DecodeTag decodes the byte sequence as a tag.
func DecodeTag(data []byte) ([]string, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	var tag []string
	err := decoder.Decode(&tag)

	if err != nil {
		return nil, err
	}

	return tag, nil
}
