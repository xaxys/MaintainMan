// The MIT License (MIT)

// Copyright (c) 2014 Milan Misak

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package util

// TODO - refactor out a resizing function

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"crypto/sha1"

	"maintainman/bindata"

	"github.com/g4s8/hexcolor"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/math/fixed"
)

const (
	parameterWidth    = "w"
	parameterHeight   = "h"
	parameterCropping = "c"
	parameterGravity  = "g"
	parameterFilter   = "f"
	parameterScale    = "s"

	// CroppingModeExact crops an image exactly to given dimensions
	CroppingModeExact = "e"
	// CroppingModeAll crops an image so that all of it is displayed in a frame of at most given dimensions
	CroppingModeAll = "a"
	// CroppingModePart crops an image so that it fills a frame of given dimensions
	CroppingModePart = "p"
	// CroppingModeKeepScale crops an image so that it fills a frame of given dimensions, keeps scale
	CroppingModeKeepScale = "k"

	GravityNorth     = "n"
	GravityNorthEast = "ne"
	GravityEast      = "e"
	GravitySouthEast = "se"
	GravitySouth     = "s"
	GravitySouthWest = "sw"
	GravityWest      = "w"
	GravityNorthWest = "nw"
	GravityCenter    = "c"

	FilterGrayScale = "grayscale"

	DefaultScale        = 1
	DefaultCroppingMode = CroppingModeExact
	DefaultGravity      = GravityNorthWest
	DefaultFilter       = "none"

	MaxDimension = 10000
)

var (
	transformationNameConfigRe = regexp.MustCompile("^([0-9A-Za-z-]+)$")
)

func isValidTransformationName(name string) bool {
	return transformationNameConfigRe.MatchString(name)
}

// TransformationInfo for config parsing.
type TransformationInfo struct {
	Name    string     `mapstructure:"name"    yaml:"name"`
	Params  string     `mapstructure:"params"  yaml:"params"`
	Texts   []TextInfo `mapstructure:"texts"   yaml:"texts"`
	Default bool       `mapstructure:"default" yaml:"default"`
	Eager   bool       `mapstructure:"eager"   yaml:"eager"`
}

type TextInfo struct {
	Content  string  `mapstructure:"content" yaml:"content"`
	Gravity  string  `mapstructure:"gravity" yaml:"gravity"`
	FontPath string  `mapstructure:"font"    yaml:"font"`
	X        int     `mapstructure:"x-pos"   yaml:"x-pos"`
	Y        int     `mapstructure:"y-pos"   yaml:"y-pos"`
	Size     int     `mapstructure:"size"    yaml:"size"`
	Color    string  `mapstructure:"color"   yaml:"color"`
	Alpha    float64 `mapstructure:"alpha"   yaml:"alpha"`
}

func (t *TransformationInfo) ToTransformation() *Transformation {
	params, err := ParseParameters(t.Params)
	if err != nil && t.Params != "" {
		panic(fmt.Errorf("invalid transformation parameters: %s (%+v)", t.Params, err))
	}

	if !isValidTransformationName(t.Name) {
		panic(fmt.Errorf("invalid transformation name: %s", t.Name))
	}

	texts := []*Text{}
	for _, text := range t.Texts {
		if !isValidGravity(text.Gravity) {
			panic(fmt.Errorf("missing or invalid gravity (transformation %s): %s", t.Name, text.Gravity))
		}
		if text.X < 0 || text.Y < 0 {
			panic(fmt.Errorf("invalid text position (transformation %s): %d,%d", t.Name, text.X, text.Y))
		}
		color, err := hexcolor.Parse(text.Color)
		if err != nil {
			panic(fmt.Errorf("invalid text color (transformation %s): %s (%+v)", t.Name, text.Color, err))
		}

		fontBytes := []byte{}
		var osErr, bindataErr error

		// Try to load font from os filesystem
		if _, err := os.Stat(text.FontPath); os.IsNotExist(err) {
			osErr = fmt.Errorf("font does not exist (transformation %s): %s", t.Name, text.FontPath)
		}
		if osErr == nil {
			fontBytes, err = ioutil.ReadFile(text.FontPath)
			if err != nil {
				panic(fmt.Errorf("loading font failed (transformation %s): %+v", t.Name, err))
			}
		}

		// Try to load font from bindata
		if osErr != nil {
			fontBytes, bindataErr = bindata.Asset(text.FontPath)
		}

		// if not found in both os filesystem and bindata
		if osErr != nil && bindataErr != nil {
			panic(fmt.Errorf("font does not exist in both os file and bindata (transformation %s): osErr: %+v; binDataErr: %+v", t.Name, osErr, bindataErr))
		}

		font, err := freetype.ParseFont(fontBytes)
		if err != nil {
			panic(fmt.Errorf("parsing font failed (transformation %s): %+v", t.Name, err))
		}
		if text.Size <= 1 {
			panic(fmt.Errorf("invalid text size (transformation %s): %d", t.Name, text.Size))
		}
		text := &Text{
			content:  text.Content,
			gravity:  text.Gravity,
			font:     font,
			x:        text.X,
			y:        text.Y,
			size:     text.Size,
			fontPath: text.FontPath,
			color:    color,
		}
		texts = append(texts, text)
	}

	s := &Transformation{
		params: &params,
		texts:  texts,
	}
	s.Hash = s.hash()
	return s
}

// Params is a struct of parameters specifying an image transformation
type Params struct {
	Width    int
	Height   int
	Scale    int
	Cropping string
	Gravity  string
	Filter   string
}

// ToString turns parameters into a unique string for each possible assignment of parameters
func (p Params) ToString() string {
	// 0 as a value for width or height means that it will be calculated
	return fmt.Sprintf("%s_%s,%s_%s,%s_%d,%s_%d,%s_%s,%s_%d", parameterCropping, p.Cropping, parameterGravity, p.Gravity, parameterHeight, p.Height, parameterWidth, p.Width, parameterFilter, p.Filter, parameterScale, p.Scale)
}

// WithScale returns a copy of a Params struct with the scale set to the given value
func (p Params) WithScale(scale int) Params {
	return Params{p.Width, p.Height, scale, p.Cropping, p.Gravity, p.Filter}
}

// Turns a string like "w_400,h_300" and an image path into a Params struct
// The second return value is an error message
// Also validates the parameters to make sure they have valid values
// w = width, h = height
func ParseParameters(parametersStr string) (Params, error) {
	params := Params{0, 0, DefaultScale, DefaultCroppingMode, DefaultGravity, DefaultFilter}
	parts := strings.Split(parametersStr, ",")
	for _, part := range parts {
		keyAndValue := strings.SplitN(part, "_", 2)
		if len(keyAndValue) != 2 {
			return params, fmt.Errorf("invalid parameter: %s", part)
		}
		key := keyAndValue[0]
		value := keyAndValue[1]

		switch key {
		case parameterWidth, parameterHeight:
			value, err := strconv.Atoi(value)
			if err != nil {
				return params, fmt.Errorf("could not parse value for parameter: %q", key)
			}
			if value <= 0 {
				return params, fmt.Errorf("value %d must be > 0: %q", value, key)
			}
			if value > MaxDimension {
				return params, fmt.Errorf("value %d must be <= %d: %q", value, MaxDimension, key)
			}
			if key == parameterWidth {
				params.Width = value
			} else {
				params.Height = value
			}
		case parameterCropping:
			value = strings.ToLower(value)
			if len(value) > 1 {
				return params, fmt.Errorf("value %q must have only 1 character", key)
			}
			if !isValidCroppingMode(value) {
				return params, fmt.Errorf("invalid value for %q", key)
			}
			params.Cropping = value
		case parameterGravity:
			value = strings.ToLower(value)
			if len(value) > 2 {
				return params, fmt.Errorf("value %q must have at most 2 characters", key)
			}
			if !isValidGravity(value) {
				return params, fmt.Errorf("invalid value for %q", key)
			}
			params.Gravity = value
		case parameterFilter:
			value = strings.ToLower(value)
			if !isValidFilter(value) {
				return params, fmt.Errorf("invalid value for %q", key)
			}
			params.Filter = value
		}
	}

	if params.Width == 0 && params.Height == 0 {
		return params, fmt.Errorf("both width and height can't be 0")
	}

	return params, nil
}

func isValidCroppingMode(str string) bool {
	return str == CroppingModeExact || str == CroppingModeAll || str == CroppingModePart || str == CroppingModeKeepScale
}

func isValidGravity(str string) bool {
	return str == GravityNorth || str == GravityNorthEast || str == GravityEast || str == GravitySouthEast || str == GravitySouth || str == GravitySouthWest || str == GravityWest || str == GravityNorthWest || str == GravityCenter
}

func isValidFilter(str string) bool {
	return str == FilterGrayScale
}

func isEasternGravity(str string) bool {
	return str == GravityNorthEast || str == GravityEast || str == GravitySouthEast
}

func isSouthernGravity(str string) bool {
	return str == GravitySouthWest || str == GravitySouth || str == GravitySouthEast
}

// Transformation specifies parameters and a watermark to be used when transforming an image
type Transformation struct {
	params *Params
	// watermark *Watermark
	texts []*Text
	Hash  string
}

func NewTransformation(param *Params) *Transformation {
	s := &Transformation{
		params: param,
	}
	s.Hash = s.hash()
	return s
}

// Watermark specifies a watermark to be applied to an image
// type Watermark struct {
// 	imagePath, gravity string
// 	x, y               int
// }

// Text specifies a text overlay to be applied to an image
type Text struct {
	content  string
	gravity  string
	x        int
	y        int
	size     int
	fontPath string
	font     *truetype.Font
	color    color.Color
}

// FontMetrics defines font metrics for a Text struct as rounded up integers
type FontMetrics struct {
	width, height, ascent, descent float64
}

func (t *Transformation) hash() string {
	hash := t.params.ToString()
	sum := make([]byte, sha1.Size)
	for _, text := range t.texts {
		hash := text.hash()
		for i := range sum {
			sum[i] += hash[i]
		}
	}
	hash += hex.EncodeToString(sum)
	return hash
}

// Turns an image file path and a transformation parameters into a file path combining both.
// It can then be used for file lookups.
// The function assumes that imagePath contains an extension at the end.
func (t *Transformation) createFilePath(imagePath string) (string, error) {
	i := strings.LastIndex(imagePath, ".")
	if i == -1 {
		return "", fmt.Errorf("invalid image path")
	}

	sum := make([]byte, sha1.Size)
	// Watermark
	// if t.watermark != nil {
	// 	hash := t.watermark.hash()
	// 	for i := range sum {
	// 		sum[i] += hash[i]
	// 	}
	// }

	// Texts
	for _, elem := range t.texts {
		hash := elem.hash()
		for i := range sum {
			sum[i] += hash[i]
		}
	}

	extraHash := ""
	if /*t.watermark != nil ||*/ len(t.texts) != 0 {
		extraHash = "--" + hex.EncodeToString(sum)
	}

	return imagePath[:i] + "--" + t.params.ToString() + extraHash + "--" + imagePath[i:], nil
}

// func (w *Watermark) hash() []byte {
// 	h := sha1.New()

// 	io.WriteString(h, w.imagePath)
// 	io.WriteString(h, w.gravity)
// 	io.WriteString(h, strconv.Itoa(w.x))
// 	io.WriteString(h, strconv.Itoa(w.y))

// 	return h.Sum(nil)
// }

func (t *Text) hash() []byte {
	h := sha1.New()
	writeUint := func(i uint32) {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, i)
		h.Write(bs)
	}

	r, g, b, a := t.color.RGBA()
	io.WriteString(h, t.content)
	io.WriteString(h, t.gravity)
	io.WriteString(h, strconv.Itoa(t.x))
	io.WriteString(h, strconv.Itoa(t.y))
	io.WriteString(h, strconv.Itoa(t.size))
	io.WriteString(h, t.fontPath)
	writeUint(r)
	writeUint(g)
	writeUint(b)
	writeUint(a)

	return h.Sum(nil)
}

func (t *Text) getFontMetrics(scale int, content string) FontMetrics {
	// Adapted from: https://code.google.com/p/plotinum/

	// Converts truetype.FUnit to float64
	fUnit2Float64 := float64(t.size) / float64(t.font.FUnitsPerEm())

	width := 0
	prev, hasPrev := truetype.Index(0), false
	for _, ch := range content {
		index := t.font.Index(ch)
		if hasPrev {
			width += int(t.font.Kern(fixed.Int26_6(t.font.FUnitsPerEm()), prev, index))
		}
		width += int(t.font.HMetric(fixed.Int26_6(t.font.FUnitsPerEm()), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	widthFloat := float64(width) * fUnit2Float64 * float64(scale)

	bounds := t.font.Bounds(fixed.Int26_6(t.font.FUnitsPerEm()))
	height := float64(bounds.Max.Y-bounds.Min.Y) * fUnit2Float64 * float64(scale)
	ascent := float64(bounds.Max.Y) * fUnit2Float64 * float64(scale)
	descent := float64(bounds.Min.Y) * fUnit2Float64 * float64(scale)

	return FontMetrics{widthFloat, height, ascent, descent}
}

func TransformCropAndResize(img image.Image, transformation *Transformation, v any) (imgNew image.Image) {
	parameters := transformation.params
	width := parameters.Width
	height := parameters.Height
	gravity := parameters.Gravity
	scale := parameters.Scale

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Scaling factor
	if parameters.Cropping != CroppingModeKeepScale {
		if width*scale <= MaxDimension && height*scale <= MaxDimension {
			width *= scale
			height *= scale
		}
	}

	// Resize and crop
	switch parameters.Cropping {
	case CroppingModeExact:
		imgNew = resize.Resize(uint(width), uint(height), img, resize.Bilinear)
	case CroppingModeAll:
		if float32(width)*(float32(imgHeight)/float32(imgWidth)) > float32(height) {
			// Keep height
			imgNew = resize.Resize(0, uint(height), img, resize.Bilinear)
		} else {
			// Keep width
			imgNew = resize.Resize(uint(width), 0, img, resize.Bilinear)
		}
	case CroppingModePart:
		var croppedRect image.Rectangle
		if float32(width)*(float32(imgHeight)/float32(imgWidth)) > float32(height) {
			// Whole width displayed
			newHeight := int((float32(imgWidth) / float32(width)) * float32(height))
			croppedRect = image.Rect(0, 0, imgWidth, newHeight)
		} else {
			// Whole height displayed
			newWidth := int((float32(imgHeight) / float32(height)) * float32(width))
			croppedRect = image.Rect(0, 0, newWidth, imgHeight)
		}

		topLeftPoint := calculateTopLeftPointFromGravity(gravity, croppedRect.Dx(), croppedRect.Dy(), imgWidth, imgHeight)
		imgDraw := image.NewRGBA(croppedRect)

		draw.Draw(imgDraw, croppedRect, img, topLeftPoint, draw.Src)
		imgNew = resize.Resize(uint(width), uint(height), imgDraw, resize.Bilinear)
	case CroppingModeKeepScale:
		// If passed in dimensions are bigger use those of the image
		if width > imgWidth {
			width = imgWidth
		}
		if height > imgHeight {
			height = imgHeight
		}

		croppedRect := image.Rect(0, 0, width, height)
		topLeftPoint := calculateTopLeftPointFromGravity(gravity, width, height, imgWidth, imgHeight)
		imgDraw := image.NewRGBA(croppedRect)

		draw.Draw(imgDraw, croppedRect, img, topLeftPoint, draw.Src)
		imgNew = imgDraw.SubImage(croppedRect)
	}

	// Filters
	if parameters.Filter == FilterGrayScale {
		bounds := imgNew.Bounds()
		w, h := bounds.Max.X, bounds.Max.Y
		gray := image.NewGray(bounds)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				oldColor := imgNew.At(x, y)
				grayColor := color.GrayModel.Convert(oldColor)
				gray.Set(x, y, grayColor)
			}
		}
		imgNew = gray
	}

	// if transformation.watermark != nil {
	// 	w := transformation.watermark

	// 	var watermarkSrcScaled image.Image
	// 	var watermarkBounds image.Rectangle

	// 	// Try to load a scaled watermark first
	// 	if scale > 1 {
	// 		scaledPath, err := constructScaledPath(w.imagePath, scale)
	// 		if err != nil {
	// 			log.Println("Error:", err)
	// 			return
	// 		}

	// 		watermarkSrc, _, err := loadImage(scaledPath)
	// 		if err != nil {
	// 			log.Println("Error: could not load a watermark", err)
	// 		} else {
	// 			watermarkBounds = watermarkSrc.Bounds()
	// 			watermarkSrcScaled = watermarkSrc
	// 		}
	// 	}

	// 	if watermarkSrcScaled == nil {
	// 		watermarkSrc, _, err := loadImage(w.imagePath)
	// 		if err != nil {
	// 			log.Println("Error: could not load a watermark", err)
	// 			return
	// 		}
	// 		watermarkBounds = image.Rect(0, 0, watermarkSrc.Bounds().Max.X*scale, watermarkSrc.Bounds().Max.Y*scale)
	// 		watermarkSrcScaled = resize.Resize(uint(watermarkBounds.Max.X), uint(watermarkBounds.Max.Y), watermarkSrc, resize.Bilinear)
	// 	}

	// 	bounds := imgNew.Bounds()

	// 	// Make sure we have a transparent watermark if possible
	// 	watermark := image.NewRGBA(watermarkBounds)
	// 	draw.Draw(watermark, watermarkBounds, watermarkSrcScaled, watermarkBounds.Min, draw.Src)

	// 	pt := calculateTopLeftPointFromGravity(w.gravity, watermarkBounds.Dx(), watermarkBounds.Dy(), bounds.Dx(), bounds.Dy())
	// 	pt = pt.Add(getTranslation(w.gravity, w.x*scale, w.y*scale))
	// 	wX := pt.X
	// 	wY := pt.Y

	// 	watermarkRect := image.Rect(wX, wY, watermarkBounds.Dx()+wX, watermarkBounds.Dy()+wY)
	// 	finalImage := image.NewRGBA(bounds)
	// 	draw.Draw(finalImage, bounds, imgNew, bounds.Min, draw.Src)
	// 	draw.Draw(finalImage, watermarkRect, watermark, watermarkBounds.Min, draw.Over)
	// 	imgNew = finalImage.SubImage(bounds)
	// }

	if transformation.texts != nil {
		bounds := imgNew.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, imgNew, image.ZP, draw.Src)

		dpi := float64(72) // Multiply this by scale for a baaad time

		c := freetype.NewContext()
		c.SetDPI(dpi)
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)

		for _, text := range transformation.texts {
			size := float64(text.size * scale)

			c.SetSrc(image.NewUniform(text.color))
			c.SetFont(text.font)
			c.SetFontSize(size)

			content := ProcessString(text.content, v)
			fontMetrics := text.getFontMetrics(scale, content)
			width := int(c.PointToFixed(fontMetrics.width) >> 6)
			height := int(c.PointToFixed(fontMetrics.height) >> 6)

			pt := calculateTopLeftPointFromGravity(text.gravity, width, height, bounds.Dx(), bounds.Dy())
			pt = pt.Add(getTranslation(text.gravity, text.x*scale, text.y*scale))
			x := pt.X
			y := pt.Y + int(c.PointToFixed(fontMetrics.ascent)>>6)

			_, err := c.DrawString(content, freetype.Pt(x, y))
			if err != nil {
				log.Println("Error adding text:", err)
				return
			}
		}

		imgNew = rgba
	}

	return
}

func calculateTopLeftPointFromGravity(gravity string, width, height, imgWidth, imgHeight int) image.Point {
	// Assuming width <= imgWidth && height <= imgHeight
	switch gravity {
	case GravityNorth:
		return image.Point{(imgWidth - width) / 2, 0}
	case GravityNorthEast:
		return image.Point{imgWidth - width, 0}
	case GravityEast:
		return image.Point{imgWidth - width, (imgHeight - height) / 2}
	case GravitySouthEast:
		return image.Point{imgWidth - width, imgHeight - height}
	case GravitySouth:
		return image.Point{(imgWidth - width) / 2, imgHeight - height}
	case GravitySouthWest:
		return image.Point{0, imgHeight - height}
	case GravityWest:
		return image.Point{0, (imgHeight - height) / 2}
	case GravityNorthWest:
		return image.Point{0, 0}
	case GravityCenter:
		return image.Point{(imgWidth - width) / 2, (imgHeight - height) / 2}
	}
	panic("This point should not be reached")
}

// getTranslation returns a point specifying a translation by a given
// horizontal and vertical offset according to gravity
func getTranslation(gravity string, h, v int) image.Point {
	switch gravity {
	case GravityNorth:
		return image.Point{0, v}
	case GravityNorthEast:
		return image.Point{-h, v}
	case GravityEast:
		return image.Point{-h, 0}
	case GravitySouthEast:
		return image.Point{-h, -v}
	case GravitySouth:
		return image.Point{0, -v}
	case GravitySouthWest:
		return image.Point{h, -v}
	case GravityWest:
		return image.Point{h, 0}
	case GravityNorthWest:
		return image.Point{h, v}
	case GravityCenter:
		return image.Point{0, 0}
	}
	panic("This point should not be reached")
}
