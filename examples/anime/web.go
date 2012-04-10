package main

import (
	"github.com/zozor/hgui"
)

func main() {
	
	//Setup top Left Cell
	frame := hgui.NewFrame()
	listen := hgui.NewSelect(20, false, []hgui.Style{{"width":"150px"}})
	
	frame.Add(
		hgui.Html("<h3><center>Anime</center></h3>"),
		listen,
	)
	
	//Setup top Right cell
	txStyle := hgui.Style{
		"width":"200px",
	}
	
	navn := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	udlaant := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	antal := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	set := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	lager := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	dato := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	gudemappe := hgui.NewTextinput("", hgui.TextType_Text, txStyle)
	kommentar := hgui.NewTextarea("", hgui.Style{"width":"300px", "height":"130px"})
	
	innertable := hgui.NewTable(nil,
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Navn")),
			hgui.NewCell(false, 1, 1, navn),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Antal Afsnit")),
			hgui.NewCell(false, 1, 1, antal),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Sete Afsnit")),
			hgui.NewCell(false, 1, 1, set),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Lager")),
			hgui.NewCell(false, 1, 1, lager),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Dato")),
			hgui.NewCell(false, 1, 1, dato),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Gudemappe side")),
			hgui.NewCell(false, 1, 1, gudemappe),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, hgui.Html("Udlånt til")),
			hgui.NewCell(false, 1, 1, udlaant),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 2, 1, hgui.Html("Kommentar"), hgui.Style{"text-align": "left"}),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 2, 1, kommentar),
		),
		
	)
	
	//setup down left cell
	sortlist := hgui.NewSelect(0, false, nil,
		hgui.NewOptions("Sort by anime", "sort by id", "sort by dato")...
	)
	
	//setup down right cell
	buttonframe := hgui.NewFrame()
	buttonframe.Add(
		hgui.NewButton("New", nil, func() {
			//do stuff
		}),
		hgui.NewButton("Save", nil, func() {
			//do stuff
		}),
		hgui.NewButton("Reload", nil, func() {
			//do stuff
		}),
		hgui.NewButton("Delete", nil, func() {
			//do stuff
		}),
	)
	
	//Setup Table
	table := hgui.NewTable(
		[]hgui.Style{
			{"border-width": "1px", "border-style":"solid", "margin": "auto", "background-color":"white"},
		},
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, frame, hgui.Style{"width":"200px", "text-align":"center"}),
			hgui.NewCell(false, 1, 1, innertable, hgui.Style{"vertical-align":"top"}),
		),
		hgui.NewRow(nil,
			hgui.NewCell(false, 1, 1, sortlist, hgui.Style{"text-align":"center"}),
			hgui.NewCell(false, 1, 1, buttonframe, hgui.Style{"text-align":"right"}),
		),
	
	)
	hgui.Topframe.Add(table)
	hgui.Topframe.AddStyle(hgui.Css_bgcolor("green"))
	hgui.StartServer(600, 450, 20000, "Anime Manager")
} 
