package fonction

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	// Default setting
	imgColorDefault = "E5E5E5"
	msgColorDefault = "AAAAAA"
	imgWDefault     = 300
	imgHDefault     = 300
	fontSizeDefault = 0

	dpiDefault float64 = 72

	fontfileDefault = "Raleway-Medium.ttf"
	msgDefault      = ""
)

var (
	img   Img
	label Label
)

type Label struct {
	Text     string
	FontSize int
	Color    string
}

type Img struct {
	Width  int
	Height int
	Color  string
	Label  Label
}

func Do(params []string) (*bytes.Buffer, error) {

	var (
		err       error
		imgWidth  = imgWDefault
		imgHeight = imgHDefault
		imgColor  = imgColorDefault
		msgColor  = msgColorDefault
	)

	if len(params) != 0 {
		if len(params) >= 1 {
			imgWidth, err = strconv.Atoi(params[0])
			if err != nil {
				imgWidth = imgWDefault
			}
		}
		if len(params) >= 2 {
			imgHeight, err = strconv.Atoi(params[1])
			if err != nil {
				imgHeight = imgHDefault
			}
		}
		if len(params) >= 3 {
			imgColor = params[2]
		}
		if len(params) >= 4 {
			msgColor = params[3]
		}

		label = Label{Text: msgDefault, FontSize: fontSizeDefault, Color: msgColor}
		img = Img{Width: imgWidth, Height: imgHeight, Color: imgColor, Label: label}
	} else {
		label = Label{Text: msgDefault, FontSize: fontSizeDefault, Color: msgColorDefault}
		img = Img{Width: imgWDefault, Height: imgHDefault, Color: imgColorDefault, Label: label}
	}

	return img.generate()
}

// generate - make the image according to the desired size, color, and text.
func (i Img) generate() (*bytes.Buffer, error) {
	// If there are dimensions and there are no requirements for the Text, we will build the default Text.
	if ((i.Width > 0 || i.Height > 0) && i.Label.Text == "") || i.Label.Text == "" {
		i.Label.Text = fmt.Sprintf("%d x %d", i.Width, i.Height)
	}
	// If there are no parameters for the font size, we will construct it based on the sizes of the image.
	if i.Label.FontSize == 0 {
		i.Label.FontSize = i.Width / 10
		if i.Height < i.Width {
			i.Label.FontSize = i.Height / 5
		}
	}
	// Convert the color from string to color.RGBA.
	clr, err := ToRGBA(i.Color)
	if err != nil {
		return nil, err
	}
	// Create an in-memory image with the desired size.
	m := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
	//Draw a picture:
	// - in the sizes (Bounds)
	// - with color (Uniform - wrapper above color.Color with Image functions)
	// - based on the point (Point) as the base image
	// - fill with color Uniform (draw.Src)
	draw.Draw(m, m.Bounds(), image.NewUniform(clr), image.Point{}, draw.Src)
	// add a text in the picture.
	if err = i.drawLabel(m); err != nil {
		return nil, err
	}
	var im image.Image = m
	// Allocate memory for our data (the bytes of the image)
	buffer := &bytes.Buffer{}
	// Let's encode the image into our allocated memory.
	err = jpeg.Encode(buffer, im, nil)

	return buffer, err
}

// drawLabel - add a text in the picture
func (i *Img) drawLabel(m *image.RGBA) error {
	// Convert string text to RGBA.
	clr, err := ToRGBA(i.Label.Color)
	if err != nil {
		return err
	}
	// Get the font (should work with both latin and cyrillic).
	fontBytes, err := ioutil.ReadFile(fontfileDefault)
	if err != nil {
		return err
	}
	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	// Prepare a Drawer for drawing text on the image.
	d := &font.Drawer{
		Dst: m,
		Src: image.NewUniform(clr),
		Face: truetype.NewFace(fnt, &truetype.Options{
			Size:    float64(i.Label.FontSize),
			DPI:     dpiDefault,
			Hinting: font.HintingNone,
		}),
	}
	//Setting the baseline.
	d.Dot = fixed.Point26_6{
		X: (fixed.I(i.Width) - d.MeasureString(i.Label.Text)) / 2,
		Y: fixed.I((i.Height+i.Label.FontSize)/2 - 12),
	}
	// Directly rendering text to our RGBA image.
	d.DrawString(i.Label.Text)

	return nil
}
