package main

import (
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	doc := js.Global.Get("document")
	println("Hello, browser console!", doc)
}

func mySearch() {
	doc := js.Global.Get("document")
	input := doc.Get("myInput")
	println("myInput=", input)
}
