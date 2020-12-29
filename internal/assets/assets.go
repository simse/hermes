package assets

import "github.com/gobuffalo/packr/v2"

// Box contains assets embedded in the binary
var Box *packr.Box

// InitBox sets up the asset box
func InitBox() {
	Box = packr.New("assets", "./assets")
}
