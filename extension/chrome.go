package main

import (
	"github.com/gopherjs/gopherjs/js"
)

var Chrome = &ChromeExtension{Object: js.Global.Get("chrome")}

// search gopherjs struct tag
type ChromeExtension struct {
	*js.Object
	Runtime   *js.Object `js:"runtime"`
	Tabs      *js.Object `js:"tabs"`
	Downloads *js.Object `js:"downloads"`
}
