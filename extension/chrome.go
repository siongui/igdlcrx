package main

import (
	"github.com/gopherjs/gopherjs/js"
	. "github.com/siongui/godom"
)

var Chrome = &ChromeExtension{Object: js.Global.Get("chrome")}

// search gopherjs struct tag
type ChromeExtension struct {
	*js.Object
	Runtime   *js.Object `js:"runtime"`
	Tabs      *js.Object `js:"tabs"`
	Downloads *js.Object `js:"downloads"`
}

func GetElementInElement(element *Object, selector string) (elm *Object, ok bool) {
	elms := element.QuerySelectorAll(selector)
	if len(elms) == 1 {
		elm = elms[0]
		ok = true
		return
	}
	return
}
