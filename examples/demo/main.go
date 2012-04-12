package main

import (
	"github.com/zozor/hgui"
	"fmt"
)

func main() {
	/*Test of Resources*/
	hgui.SetResource(map[string][]byte{
		"/pics/test.jpg": BILLEDE,
	})
	
	/*Test of Tables*/
	toptable := hgui.NewTable([]hgui.Style{{
		"border-width": "2px",
		"border-style": "solid",
		"margin": "auto",
	}})
	
	list := hgui.NewSelect(20, false, nil)
	list.SetStyle(hgui.Style{
		"width": "150px",
	})
	
	
	name := hgui.NewTextinput("", hgui.TextType_Text)
	afsnit := hgui.NewTextinput("", hgui.TextType_Text)
	set := hgui.NewTextinput("", hgui.TextType_Text)
	
	formtable := hgui.NewTable(nil)
	formtable.Addrows(
		hgui.NewRow(
			[]hgui.Style{{"background-color":"red"}},
			hgui.NewCell(false, 0,0, hgui.NewLabel("Name")),
			hgui.NewCell(false, 0,0, name),
		),
		hgui.NewRow(
			[]hgui.Style{{"background-color":"green"}},
			hgui.NewCell(false, 0,0, hgui.NewLabel("Afsnit")),
			hgui.NewCell(false, 0,0, afsnit),
		),
		hgui.NewRow(
			[]hgui.Style{{"background-color":"red"}},
			hgui.NewCell(false, 0,0, hgui.NewLabel("Set")),
			hgui.NewCell(false, 0,0, set),
		),
	)
	
	cellstyle := hgui.Style{
		"vertical-align": "top",
	}
	
	toptable.Addrows(
		hgui.NewRow(
			nil,
			hgui.NewCell(false, 0,0, list, cellstyle, hgui.Style{"width":"150px"}),
			hgui.NewCell(false, 0,0, formtable, cellstyle),
		),
	)
	
	/*Test of TextArea, textinput, radio and check box*/
	ramme := hgui.NewFrame()
	textinput := hgui.NewTextinput("", hgui.TextType_Text)
	textarea := hgui.NewTextarea("", hgui.Style{"width":"200px", "height":"100px"})
	textlabel := hgui.NewLabel("")
	radio1 := hgui.NewRadioCheckbox(true, "hej")
	radio2 := hgui.NewRadioCheckbox(true, "hej")
	radio1.SetEvent(hgui.Evt_onclick, func() {
		textlabel.SetValue(textinput.Value())
	})
	radio2.SetEvent(hgui.Evt_onclick, func() {
		textlabel.SetValue(textarea.Value())
	})
	radio3 := hgui.NewRadioCheckbox(false, "")
	
	knap1 := hgui.NewButton("Check", nil, func() {
		radio3.Check()
	})
	knap2 := hgui.NewButton("Uncheck", nil, func() {
		radio3.Uncheck()
	})
	knap3 := hgui.NewButton("GetCheck", nil, func() {
		if radio3.Checked() {
			textarea.SetValue("Checkbox is checked")
		} else {
			textarea.SetValue("Checkbox is not checked")
		}
	})
	knap4 := hgui.NewButton("Check Radio", nil, func() {
		radio2.Check()
	})
	
	ramme.Add(
		textinput, hgui.Html("<br/>"),
		textarea, hgui.Html("<br/>"),
		radio1,hgui.Html("Pressing here gets input from textinput <br/>"), 
		radio2, hgui.Html("Pressing here gets input from textarea <br/>"),
		radio3, hgui.Html("<br/>"),	
		knap1, knap2, knap3, knap4, hgui.Html("<br/>"),
		
		textlabel,
	)
	/*Test of Fieldtype, ordered lists and unordered lists*/
	ramme2 := hgui.NewFrame()
	ol := hgui.NewList(true, nil,
		hgui.NewListItem("monster!"),
		hgui.NewListItem("monster!"),
	)
	ul := hgui.NewList(false, nil,
		hgui.NewListItem("monster!"),
		hgui.NewListItem("monster!"),
	)
	ramme2.Add(ol, ul)
	fieldset := hgui.NewFieldset("test!", ramme2)
	ramme3 := hgui.NewFrame()
	
	ramme3.Add(fieldset, hgui.Html("<br/>"),
		hgui.NewButton("Clear list", nil, func() {
			ol.SetList()
		}),
		hgui.NewButton("Big list", nil, func() {
			ol.SetList(
				hgui.NewListItem("monster 1"),
				hgui.NewListItem("monster 2"),
				hgui.NewListItem("monster 1"),
				hgui.NewListItem("monster 2"),
				hgui.NewListItem("monster 1"),
				hgui.NewListItem("monster 2"),
				hgui.NewListItem("monster 1"),
				hgui.NewListItem("monster 2"),
			)
		}),hgui.Html("<br/>"),
		hgui.NewLink("#", hgui.NewImage("/pics/test.jpg")),
	)
	
	/*Test of Select input and Gauge*/
	ramme4 := hgui.NewFrame()
	select1 := hgui.NewSelect(10, true, nil,
		hgui.NewOptions("aaaaaa","bbbbbb","cccccc","dddddd")...
	)
	select2 := hgui.NewSelect(20, false, nil,
		hgui.NewOptions("aaaaaa","bbbbbb","cccccc","dddddd")...
	)
	select3 := hgui.NewSelect(0, false, nil,
		hgui.NewOptions("aaaaaa","bbbbbb","cccccc","dddddd")...
	)
	text8 := hgui.NewLabel("")
	
	gauge := hgui.NewGauge(10, 100, "red")
	
	ramme4.Add(select1, select2,select3,hgui.Html("<br/>"), text8,hgui.Html("<br/>"),
		hgui.NewButton("Get Multiple", nil, func(){
			_, v := select1.Selected()
			text8.SetValue(fmt.Sprint(v))
		}),
		hgui.NewButton("Get Single", nil, func(){
			v, _ := select2.Selected()
			text8.SetValue(v)
		}),
		hgui.NewButton("Set Big", nil, func(){
			select2.SetOptions(
				hgui.NewOptions("aaaaaaaaa","aaaaaaaaa","aaaaaaaaa","aaaaaaaaa","aaaaaaaaa","aaaaaaaaa","aaaaaaaaa")...,
			)
		}),
		hgui.NewButton("Clear!", nil, func(){
			select2.SetOptions()
		}),hgui.Html("<br/>"),
		
		gauge,hgui.Html("<br/>"),
		hgui.NewButton("Change Gauge!", nil, func(){
			gauge.SetValue(gauge.Value()+5)
		}),
	)
	
	
	
	
	
	/*Test of Buttons, widget type*/
	ramme.Hide()
	ramme3.Hide()
	ramme4.Hide()
		
	hgui.Topframe.Add(
		hgui.NewButton("Frame 1", nil, func() {
			toptable.Show()
			ramme.Hide()
			ramme3.Hide()
			ramme4.Hide()
		}),
		
		hgui.NewButton("Frame 2", nil, func() {
			toptable.Hide()
			ramme.Show()
			ramme3.Hide()
			ramme4.Hide()
		}),
		
		hgui.NewButton("Frame 4", nil, func() {
			toptable.Hide()
			ramme.Hide()
			ramme3.Show()
			ramme4.Hide()
		}),
		
		hgui.NewButton("Frame 5", nil, func() {
			toptable.Hide()
			ramme.Hide()
			ramme3.Hide()
			ramme4.Show()
		}),
		
		hgui.Html("<br/>"),
		
		hgui.NewButton("Remove Style", nil, func() {
			hgui.Topframe.RemoveStyle(hgui.Style{"background-color":"red"})
		}),
		
		hgui.NewButton("Add Style", nil, func() {
			hgui.Topframe.AddStyle(hgui.Style{"background-color":"red"})
		}),
		
		hgui.NewButton("Set Attribute", nil, func() {
			hgui.Topframe.SetAttribute("bgcolor", "green")
		}),
		
		hgui.NewButton("Remove Attribute", nil, func() {
			hgui.Topframe.RemoveAttribute("bgcolor")
		}),
		
		toptable,
		ramme,
		ramme3,hgui.Html("<br/>"),hgui.Html("<br/>"),
		ramme4,
	)
	hgui.StartServer(800,600, "Demo!")
}
