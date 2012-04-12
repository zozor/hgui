Simple GUI toolkit with go
------------------------------

Install
===========================

You will need to install gtk-dev files.

`go get github.com/zozor/hgui`


The idea
===========================

This is a gui toolkit that relies on HTML, CSS, and javascript. But as a user of this toolkit you do not have to worry about that
although it helps understanding it.

The compiled program will, when you run it, start a server at 127.0.0.1:randomport, and connect to it with gtk-webkit.

To make a simple interface, no knowledge of html, javascript or css is needed. For styling some CSS knowledge is needed, but not much.

Here is a simple program. An advanced one can be found in examples


	package main

	import "github.com/zozor/hgui"

	func main() {
		label := hgui.NewLabel("Simple label")
		input := hgui.NewTextinput("", hgui.TextType_Text)
		hgui.Topframe.Add(
			input, 
			hgui.NewButton("New Text!", nil, func() {
				label.SetValue(input.Value())
				label.SetStyle(hgui.Style{"color":"blue"})
			}),
			hgui.Html("<br/>"), 
			label,
		)
		hgui.StartServer(800, 600, "Simple program!")
	}

Features
===========================
#### Widgets

- Frames / Container
- Tables
- Radio and Check boxes
- Fieldset
- Labels
- Text input
- Textarea
- Selectform
- Buttons
- Links
- Images
- Modal Dialogs

#### Other
- Styling with css
- Raw javascript
- Resources (it actually does not allow anything else)

Issues
===========================
Gtk sucks. It spams my face with errors.