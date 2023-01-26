package gclog

import (
	"github.com/arafath-mk/gcstyle"

	"github.com/arafath-mk/gcstyle/wcolor"
)

// var colorVal = wcolor.RGB(95, 95, 95)

var styleKey = gcstyle.Style{
	Color: wcolor.Grey.Clone(),
}

var styleVal = gcstyle.Style{
	Color:  wcolor.Grey.Clone(),
	Darken: true,
}

var styleValErr = gcstyle.Style{
	Color: &wcolor.Red,
}
