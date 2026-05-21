package render

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// defaultFace returns the built-in 7×13 bitmap font.
// All text is drawn via the text/v2 adapter wrapping this face.
func defaultFace() font.Face {
	return basicfont.Face7x13
}
