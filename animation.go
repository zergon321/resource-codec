package codec

import (
	"bytes"
	"encoding/binary"

	"github.com/faiface/pixel"
)

// GetSpritesheetFrames returns the set of rectangles
// corresponding to the frames of the spritesheet.
func GetSpritesheetFrames(spritesheet *pixel.PictureData, width, height int) []pixel.Rect {
	frames := make([]pixel.Rect, 0)
	pixelWidth := spritesheet.Bounds().W()
	pixelHeight := spritesheet.Bounds().H()
	dw := pixelWidth / float64(width)
	dh := pixelHeight / float64(height)

	for y := pixelHeight; y > 0; y -= dh {
		for x := 0.0; x < pixelWidth; x += dw {
			frame := pixel.R(x, y-dh, x+dw, y)
			frames = append(frames, frame)
		}
	}

	return frames
}

// AnimationDataToBytes converts the animation data
// to a byte array.
func AnimationDataToBytes(anim *AnimationData) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// Write the name of the spritesheet.
	err := binary.Write(buffer, binary.BigEndian, int32(len(anim.Spritesheet)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(anim.Spritesheet)

	if err != nil {
		return nil, err
	}

	// Write the frame data.
	err = binary.Write(buffer, binary.BigEndian, int32(len(anim.Frames)))

	if err != nil {
		return nil, err
	}

	for _, frame := range anim.Frames {
		err = binary.Write(buffer, binary.BigEndian, frame.Min.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, frame.Min.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, frame.Max.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, frame.Max.Y)

		if err != nil {
			return nil, err
		}
	}

	// Write the frame durations.
	for _, duration := range anim.Durations {
		err = binary.Write(buffer, binary.BigEndian, duration)

		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

// AnimationDataFromBytes restores the animation data
// from the byte array.
func AnimationDataFromBytes(data []byte) (*AnimationData, error) {
	buffer := bytes.NewBuffer(data)

	// Read the name of the spritesheet.
	var nameLength int32
	err := binary.Read(buffer, binary.BigEndian, &nameLength)

	if err != nil {
		return nil, err
	}

	nameBytes := make([]byte, nameLength)
	_, err = buffer.Read(nameBytes)

	if err != nil {
		return nil, err
	}

	spritesheetName := string(nameBytes)

	// Read the number of frames.
	var frameCount int32
	err = binary.Read(buffer, binary.BigEndian, &frameCount)

	if err != nil {
		return nil, err
	}

	// Read AnimationData frames.
	frames := []pixel.Rect{}

	for i := int32(0); i < frameCount; i++ {
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

		frame := pixel.R(minX, minY, maxX, maxY)
		frames = append(frames, frame)
	}

	// Read frame durations.
	durations := []int32{}

	for i := int32(0); i < frameCount; i++ {
		var duration int32
		err = binary.Read(buffer, binary.BigEndian, &duration)

		if err != nil {
			return nil, err
		}

		durations = append(durations, duration)
	}

	return &AnimationData{
		Spritesheet: spritesheetName,
		Frames:      frames,
		Durations:   durations,
	}, nil
}
