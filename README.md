Simple GUI toolkit with go
------------------------------

Install
===========================

`go get github.com/zozor/hgui`


The idea
===========================

This is a gui toolkit that relies on HTML, CSS, and javascript. But as a user of this toolkit you do not have to worry about that
although it helps understanding it.

The compiled program will when run, start a server at 127.0.0.1:port, and runs the standard browser with that address.
(it is, at the moment, up to the user of the package to set the port. The server fails if the port is taken and exits).

This means the gui is rendered in the browser. You cannot open another tab with the page, and you cannot refresh
(doing so makes the server unreachable, and it will die in 10 seconds)

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
			hgui.StartServer(20000)
		}

Features
===========================
It should support windows (confirm this for me)

#### Widgets

- Frames
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

#### Other
- Styling with css
- Raw javascript
- Resources (it actually does not allow anything else)