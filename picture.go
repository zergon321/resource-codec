package codec

import (
	"bytes"
	"encoding/binary"
	"image/color"

	"github.com/faiface/pixel"
)

// PictureDataToBytes converts the picture data to
// a byte array.
func PictureDataToBytes(picData *pixel.PictureData) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// Write number of picture pixels into the buffer.
	err := binary.Write(buffer, binary.BigEndian, int32(len(picData.Pix)))

	if err != nil {
		return nil, err
	}

	// Write picture pixels into the buffer.
	for _, pix := range picData.Pix {
		err = binary.Write(buffer, binary.BigEndian, pix.R)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, pix.G)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, pix.B)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, pix.A)

		if err != nil {
			return nil, err
		}
	}

	// Write picture stride into the buffer.
	err = binary.Write(buffer, binary.BigEndian, int32(picData.Stride))

	if err != nil {
		return nil, err
	}

	// Write picture rect into the buffer.
	err = binary.Write(buffer, binary.BigEndian, picData.Rect.Min.X)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, picData.Rect.Min.Y)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, picData.Rect.Max.X)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, picData.Rect.Max.Y)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// PictureDataFromBytes restores the picture data from the byte array.
func PictureDataFromBytes(data []byte) (*pixel.PictureData, error) {
	buffer := bytes.NewBuffer(data)

	// Read the length of the pixel array.
	var picLength int32
	err := binary.Read(buffer, binary.BigEndian, &picLength)

	if err != nil {
		return nil, err
	}

	// Read the RGBA pixel data.
	var r, g, b, a uint8
	pixels := make([]color.RGBA, 0)

	for i := 0; i < int(picLength); i++ {
		err = binary.Read(buffer, binary.BigEndian, &r)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &g)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &b)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &a)

		if err != nil {
			return nil, err
		}

		rgbaPixel := color.RGBA{r, g, b, a}
		pixels = append(pixels, rgbaPixel)
	}

	// Read the picture stride.
	var picStride int32
	err = binary.Read(buffer, binary.BigEndian, &picStride)

	if err != nil {
		return nil, err
	}

	// Read the picture bounds.
	var minX, minY, maxX, maxY float64

	err = binary.Read(buffer, binary.BigEndian, &minX)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &minY)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &maxX)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &maxY)

	if err != nil {
		return nil, err
	}

	// Assemble the picture.
	picture := &pixel.PictureData{
		Pix:    pixels,
		Stride: int(picStride),
		Rect:   pixel.R(minX, minY, maxX, maxY),
	}

	return picture, nil
}
