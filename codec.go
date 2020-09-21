package codec

import (
	"bytes"
	"encoding/binary"
	"image/color"

	"github.com/faiface/pixel"
	"gopkg.in/yaml.v2"
)

// AnimationData contains data
// about animation frames, their
// durations and the name of the
// spritesheet.
type AnimationData struct {
	spritesheet string
	frames      []pixel.Rect
	durations   []int32
}

// AnimationMeta is animation metadata
// read from the YAML file.
type AnimationMeta struct {
	Name        string   `yaml:"name"`
	Spritesheet string   `yaml:"spritesheet"`
	Frames      [][2]int `yaml:"frames"`
}

// SpritesheetMeta is spritesheet metadata
// read from the YAML file.
type SpritesheetMeta struct {
	Width  int
	Height int
}

// ReadAnimationsData deserializes the animations
// metadata stored as YAML.
func ReadAnimationsData(data []byte) ([]AnimationMeta, error) {
	animations := make([]AnimationMeta, 0)
	err := yaml.Unmarshal(data, &animations)

	if err != nil {
		return nil, err
	}

	return animations, nil
}

// ReadSpritesheetsData deserializes the spritesheets
// metadata stored as YAML.
func ReadSpritesheetsData(data []byte) (map[string]SpritesheetMeta, error) {
	spritesheets := make(map[string]SpritesheetMeta, 0)
	err := yaml.Unmarshal(data, &spritesheets)

	if err != nil {
		return nil, err
	}

	return spritesheets, nil
}

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

// AnimationDataToBytes converts the animation data
// to a byte array.
func AnimationDataToBytes(anim *AnimationData) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// Write the name of the spritesheet.
	err := binary.Write(buffer, binary.BigEndian, int32(len(anim.spritesheet)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(anim.spritesheet)

	if err != nil {
		return nil, err
	}

	// Write the frame data.
	err = binary.Write(buffer, binary.BigEndian, int32(len(anim.frames)))

	if err != nil {
		return nil, err
	}

	for _, frame := range anim.frames {
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
	for _, duration := range anim.durations {
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
		spritesheet: spritesheetName,
		frames:      frames,
		durations:   durations,
	}, nil
}
