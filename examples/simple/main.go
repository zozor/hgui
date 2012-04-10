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
	hgui.StartServer(800, 600, 20000, "Simple program")
}
