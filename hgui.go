/*
hgui is a gui toolkit for webgui

Except for Style type, never, EVER use the widgets without initializing them, it probably won't work if you do that.
All types that must be initialized has a New function (f.eks. NewFrame, NewButton...).
*/
package hgui

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"runtime"
)


//=============================================
//  Variable Declartions  //
//=============================================
var (
	events     = make(chan event, 1000) //event channel
	replyid    unique                   //id for replies
	replyqueue = make(chan event, 1000) //replies

	ids unique //Id for widgets

)

//=============================================
//  evtType Declartions  //
//=============================================
type evtType string

var (
	Evt_onclick     = evtType("onclick")
	Evt_onabort     = evtType("onabort")
	Evt_onblur      = evtType("onblur")
	Evt_onchange    = evtType("onchange")
	Evt_ondblclick  = evtType("ondblclick")
	Evt_onerror     = evtType("onerror")
	Evt_onfocus     = evtType("onfocus")
	Evt_onkeydown   = evtType("onkeydown")
	Evt_onkeypress  = evtType("onkeypress")
	Evt_onkeyup     = evtType("onkeyup")
	Evt_onload      = evtType("onload")
	Evt_onmousedown = evtType("onmousedown")
	Evt_onmousemove = evtType("onmousemove")
	Evt_onmouseout  = evtType("onmouseout")
	Evt_onmouseover = evtType("onmouseover")
	Evt_onmouseup   = evtType("onmouseup")
	Evt_onreset     = evtType("onreset")
	Evt_onresize    = evtType("onresize")
	Evt_onselect    = evtType("onselect")
	Evt_onsubmit    = evtType("onsubmit")
	Evt_onunload    = evtType("onunload")
)

//=============================================
//  TextType Declartions  //
//=============================================
type texttype string

var (
	TextType_Password = texttype("password")
	TextType_Text     = texttype("text")
	TextType_Hidden   = texttype("hidden")
)

//=============================================
//  css Declartions  //
//=============================================

func Css_bgcolor(color string) Style {
	return Style{"background-color":color}
}

var (
	Css_bgcolor_black = Style{"background-color": "black"}
)

//=============================================
//  Core  //
//=============================================
type event struct {
	id         string
	javascript string
	reply      chan string
}

//Makes an event that can be send to the browser
func Event(js string, reply chan string) event {
	return event{replyid.New(""), js, reply}
}

type jsEvent struct {
	Id         string
	Javascript string
	Reply      bool
}

type reply struct {
	Id    string
	Reply string
}

func eventPoll(w http.ResponseWriter) {
	buf := make([]jsEvent, 0, 10)
	for {
		select {
		case evt := <-events:
			if evt.reply == nil {
				buf = append(buf, jsEvent{evt.id, evt.javascript, false})
			} else {
				buf = append(buf, jsEvent{evt.id, evt.javascript, true})
				replyqueue <- evt
			}
		default:
			goto done
		}
	}
done:

	out, err := json.Marshal(buf)
	if err != nil {
		w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(out)
}

func eventReply(r reply) {
	for v := range replyqueue {
		if r.Id == v.id {
			v.reply <- r.Reply
			replyid.Remove(r.Id)
			return
		}
		replyqueue <- v
	}	
}

//Sends and Event to the browser
func SendEvent(js string, reply chan string) {
	events <- Event(js, reply)
}

//=============================================
//  Practical  //
//=============================================
func jq(id, property string) string {
	return fmt.Sprintf(`$("#%s").%s`, id, property)
}

func escape(s string) string {
	s = strings.Replace(s, `"`, `\"`, -1)
	s = strings.Replace(s, `'`, `\'`, -1)
	return s
}

//Threadsafe unique list of strings
type unique struct {
	list  []string
	mutex sync.Mutex
}

func (i *unique) New(prefix string) string {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if i.list == nil {
		i.list = make([]string, 0, 100)
	}

	var id string
	for {
		id = prefix + strconv.Itoa(rand.Int())
		for _, v := range i.list {
			if v == id {
				continue
			}
		}
		break
	}
	i.list = append(i.list, id)
	return id
}

func (i *unique) Remove(id string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if i.list == nil {
		return
	}

	for j, v := range i.list {
		if v == id {
			i.list = append(i.list[:j], i.list[j+1:]...)
		}
	}
}

//=============================================
//  Style  //
//=============================================

//CSS map
type Style map[string]string

func (s Style) Marshal() string {
	buf := make([]string, 0, len(s))
	for key, value := range s {
		buf = append(buf, key+":"+escape(value))
	}
	return strings.Join(buf, ";")
}

//Creates a Style type from a string
func UnmarshalStyle(style string) (Style, error) {
	out := Style{}
	for _, v := range strings.Split(style, ";") {
		kv := strings.Split(v, ":")
		if len(kv) != 2 {
			return Style{}, errors.New("Something wrong with input: " + style)
		}
		out[kv[0]] = kv[1]
	}
	return out, nil
}

//Add more css to your style
func (s Style) AddStyle(n Style) {
	for k, v := range n {
		s[k] = v
	}
}

//Remove some css from your style
func (s Style) RemoveStyle(n Style) {
	for k, _ := range n {
		delete(s, k)
	}
}

//=============================================
//  Widget  //
//=============================================

type HTMLer interface {
	HTML() string
}

type widget struct {
	id    string
	style Style
}

func newWidget(styles ...Style) *widget {
	style := Style{}
	for _, s := range styles {
		style.AddStyle(s)
	}
	w := &widget{id:ids.New("id")}
	w.SetStyle(style)
	runtime.SetFinalizer(w, func(last *widget) {
		ids.Remove(last.id)
		events <- Event(jq(last.id, "remove()"), nil)
	})
	return w
}

func (w *widget) ID() string {
	return w.id
}

//Sets a style for a widget
func (w *widget) SetStyle(style Style) {
	w.style = style
	js := jq(w.id, `attr("style", "`+style.Marshal()+`")`)
	events <- Event(js, nil)
}

//returns the style of a widget
func (w *widget) Style() Style {
	return w.style
}

//Adds a style to existing style on the widget
func (w *widget) AddStyle(n Style) {
	w.style.AddStyle(n)
	w.SetStyle(w.style)
}

//Removes a style from the widget, specified by the n style
func (w *widget) RemoveStyle(n Style) {
	w.style.RemoveStyle(n)
	w.SetStyle(w.style)
}

//Hides the widget
func (w *widget) Hide() {
	events <- Event(jq(w.id, "hide()"), nil)
}

//Shows the widget
func (w *widget) Show() {
	events <- Event(jq(w.id, "show()"), nil)
}

//This is set ONLY on the client side, the server does not record this. Don't set style with this
func (w *widget) SetAttribute(attribute, value string) {
	events <- Event(jq(w.id, "attr('"+attribute+"','"+value+"')"), nil)
}

//Not everything can be removed, see $.removeAttr() in jquery for details
//The use of frame.Flip, removes everything made through this call.
//If used with flip, it should be called afterwards
func (w *widget) RemoveAttribute(attr string) {
	events <- Event(jq(w.id, "removeAttr('"+attr+"')"), nil)
}

//This is set ONLY on the client side, the server does not record this.
//The use of frame.Flip, removes everything made through this call
//If used with flip, it should be called afterwards.
func (w *widget) SetEvent(event evtType, action func()) {
	if action == nil {
		return
	}
	handlers[w.id+"."+string(event)] = action
	events <- Event(fmt.Sprintf(`$("#%s").attr("%s", "callHandler('%s.%s')")`, w.id, string(event), w.id, string(event)), nil)
}

//=============================================
//  Frame  //
//=============================================

type Frame struct {
	*widget
	content []HTMLer
	topframe bool
}

//Creates a new container for your widgets
func NewFrame(styles ...Style) *Frame {
	return &Frame{newWidget(styles...), make([]HTMLer, 0, 20), false}
}

//Add widget to your frame
func (f *Frame) Add(widget ...HTMLer) {
	f.content = append(f.content, widget...)
}

//Resets the content of the frame. Events on widgets should be reassigned after using this.
//The function has it's usefulness, but should seldom be used.
func (f *Frame) Flip() {
	buf := make([]string, len(f.content))
	for i, v := range f.content {
		buf[i] = v.HTML()
	}
	events <- Event(jq(f.id, `html("`+escape(strings.Join(buf, ""))+`")`), nil)
}

func (f *Frame) HTML() string {
	buf := make([]string, len(f.content))
	for i, v := range f.content {
		buf[i] = v.HTML()+"\n"
	}
	if f.topframe {
		return fmt.Sprintf(`<body id="%s" style="%s">%s</body>`, f.id, f.style.Marshal(), strings.Join(buf, ""))
	}
	return fmt.Sprintf(`<div id="%s" style="%s">%s</div>`, f.id, f.style.Marshal(), strings.Join(buf, ""))
}

//=============================================
//  Label  //
//=============================================

type Label struct {
	*widget
	value string
}

//Creates a new label with a value
func NewLabel(value string, styles ...Style) *Label {
	return &Label{newWidget(styles...), value}
}

//Gets the value of the label
func (l *Label) Value() string {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").html()`, l.id), reply)
	events <- evt
	return <-evt.reply
}

//Set the value of the Label
func (l *Label) SetValue(s string) {
	events <- Event(fmt.Sprintf(`$("#%s").html("%s")`, l.id, escape(s)), nil)
}

func (l *Label) HTML() string {
	return fmt.Sprintf(`<span id="%s" style="%s">%s</span>`, l.id, l.style.Marshal(), l.value)
}

//=============================================
//  Button  //
//=============================================

type Button struct {
	*widget
	value  string
	action func()
}

//Creates a new button, with a caption and a callback
func NewButton(value string, styles []Style, action func()) *Button {
	b := &Button{newWidget(styles...), value, action}
	if action != nil {
		handlers[b.id+".onclick"] = action
	}
	return b
}

func (b *Button) HTML() string {
	return fmt.Sprintf(`<input type="button" value="%s" id="%s" onclick="callHandler('%s.onclick');" style="%s" />`,
		b.value, b.id, b.id, b.style.Marshal())
}

//=============================================
//  Table  //
//=============================================

type Table struct {
	*widget
	rows []*Row
}

//Creates a grid to put stuff into
func NewTable(styles []Style, r ...*Row) *Table {
	return &Table{newWidget(styles...), r}
}

func (t *Table) Addrows(r ...*Row) {
	t.rows = append(t.rows, r...)
}

func (t *Table) HTML() (html string) {
	html = `<table id="` + t.id + `" style="`+t.style.Marshal()+`">`
	for _, v := range t.rows {
		html += v.HTML()+"\n"
	}
	html += "</table>"
	return
}


//you put this into your table/grid. It's a row.
type Row struct {
	*widget
	cells []*Cell
}

func NewRow(styles []Style, c ...*Cell) *Row {
	return &Row{newWidget(styles...), c}
}

func (r *Row) AddCells(c ...*Cell) {
	r.cells = append(r.cells, c...)
}

func (r *Row) HTML() (html string) {
	html = `<tr id="` + r.id + `" style="`+r.style.Marshal()+`">`
	for _, v := range r.cells {
		html += v.HTML()+"\n"
	}
	html += "</tr>"
	return
}

type Cell struct {
	*widget
	header  bool
	colspan int
	rowspan int
	content HTMLer
}

//Does this need introduction? Colspan tells the cell how many columns to span across, and rowspan how many rows..
func NewCell(header bool, colspan, rowspan int, content HTMLer, styles ...Style) *Cell {
	return &Cell{newWidget(styles...), header, colspan, rowspan, content}
}

func (t *Cell) HTML() (html string) {
	if t.header {
		html = fmt.Sprintf(`<th id="%s" colspan="%d" rowspan="%d" style="%s">`, t.id, t.colspan, t.rowspan, t.style.Marshal())
	} else {
		html = fmt.Sprintf(`<td id="%s" colspan="%d" rowspan="%d" style="%s">`, t.id, t.colspan, t.rowspan, t.style.Marshal())
	}
	html += t.content.HTML()
	if t.header {
		html += "</th>"
	} else {
		html += "</td>"
	}

	return
}

//=============================================
//  Textinput  //
//=============================================

type Textinput struct {
	*widget
	value string
	_type string
}

//Creates a inpup field for text
func NewTextinput(value string, ttype texttype,styles ...Style) *Textinput {
	return &Textinput{newWidget(styles...), value, string(ttype)}
}

//Grabs the value of the textinput
func (t *Textinput) Value() string {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").val()`, t.id), reply)
	events <- evt
	return <-evt.reply
}

//Sets the value of the input
func (t *Textinput) SetValue(s string) {
	events <- Event(fmt.Sprintf(`$("#%s").val("%s")`, t.id, escape(s)), nil)
}

func (t *Textinput) HTML() string {
	return fmt.Sprintf(`<input type="%s" id="%s" value="%s" style="%s"/>`, t._type, t.id, t.value, t.style.Marshal())
}

//=============================================
//  Textarea  //
//=============================================

type Textarea struct {
	*widget
	value string
}

//Multiline text input
func NewTextarea(value string, styles ...Style) *Textarea {
	return &Textarea{newWidget(styles...), value}
}

func (t *Textarea) Value() string {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").val()`, t.id), reply)
	events <- evt
	return <-evt.reply
}

func (t *Textarea) SetValue(s string) {
	events <- Event(fmt.Sprintf(`$("#%s").text("%s")`, t.id, escape(s)), nil)
}

func (t *Textarea) HTML() string {
	return `<textarea id="` + t.id + `" style="`+t.style.Marshal()+`">`+t.value+`</textarea>`
}

//=============================================
//  Radiobuttons checkboxes  //
//=============================================

type Radiocheckbox struct {
	*widget
	group   string
	radiobox bool
}

//Creates either new radiobox or checkbox.
//Checkboxes are not affected by the grouping
func NewRadioCheckbox(radiobox bool, group string, styles ...Style) *Radiocheckbox {
	return &Radiocheckbox{newWidget(styles...), group, radiobox}
}

//Get the state of a radiobox/checkbox
func (t *Radiocheckbox) Checked() bool {
	reply := make(chan string)
	evt := Event(`reply = $("#`+t.id+`").prop("checked")`, reply)
	events <- evt
	if <-reply == "true" {
		return true
	}
	return false
}

//Checks the checkbox/radiobox
func (t *Radiocheckbox) Check() {
	events <- Event(jq(t.id, `prop("checked", "checked")`), nil)
}

//Unchecks the checkbox/radiobox
func (t *Radiocheckbox) Uncheck() {
	events <- Event(jq(t.id, `prop("checked", false)`), nil)
}

func (t *Radiocheckbox) HTML() string {
	if t.radiobox {
		return fmt.Sprintf(`<input type="radio" id="%s" name="%s" style="%s"/>`, t.id, t.group, t.style.Marshal())
	}
	return fmt.Sprintf(`<input type="checkbox" id="%s" style="%s"/>`, t.id, t.style.Marshal())
}

//=============================================
//  Image  //
//=============================================

type Image struct {
	*widget
	src string
}

//New image...
func NewImage(src string, styles ...Style) *Image {
	return &Image{newWidget(styles...), src}
}

func (i *Image) HTML() string {
	return fmt.Sprintf(`<img id="%s" src="%s" style="%s"/>`, i.id, i.src, i.style.Marshal())
}

//=============================================
//  Lists  //
//=============================================

//List are bullet points or numbered lists.
type List struct {
	*widget
	items []*Listitem
	ordered bool
}

func NewList(ordered bool, styles []Style, items ...*Listitem) *List {
	return &List{newWidget(styles...), items, ordered}
}

func (l *List) SetList(items ...*Listitem) {
	l.items = items
	html := ""
	for _, v := range l.items {
		html += v.HTML()
	}
	events <- Event(jq(l.id, `html("`+escape(html)+`")`), nil)
}

func (l *List) HTML() (html string) {
	if l.ordered {
		html = fmt.Sprintf(`<ol id="%s" style="%s">`, l.id, l.style.Marshal())
	} else {
		html = fmt.Sprintf(`<ul id="%s" style="%s">`, l.id, l.style.Marshal())
	}
	
	for _, v := range l.items {
		html += v.HTML()+"\n"
	}
	
	if l.ordered {
		html += "</ol>"
	} else {
		html += "</ul>"
	}
	return
}

type Listitem struct {
	*widget
	value string
}

func NewListItem(value string, styles ...Style) *Listitem {
	return &Listitem{newWidget(styles...), value}
}

func (l *Listitem) HTML() string {
	return fmt.Sprintf(`<li id="%s">%s</li>`, l.id, l.value)
}

//=============================================
//  Links  //
//=============================================

type Link struct {
	*widget
	href  string
	value HTMLer
}

func NewLink(href string, value HTMLer, styles ...Style) *Link {
	return &Link{newWidget(styles...), href, value}
}

func (l *Link) HTML() string {
	return fmt.Sprintf(`<a href="%s" id="%s" style="%s">%s</a>`, l.href, l.id, l.style.Marshal(), l.value.HTML())
}

//=============================================
//  Fieldset  //
//=============================================

type Fieldset struct {
	legend  string
	content HTMLer
}

func NewFieldset(caption string, content HTMLer) *Fieldset {
	return &Fieldset{caption, content}
}

func (f *Fieldset) HTML() string {
	return `<fieldset><legend>` + f.legend + `</legend>` + f.content.HTML() + `</fieldset>`
}

//=============================================
//  Select  //
//=============================================

type Selectform struct {
	*widget
	options  []*Option
	size     int
	multiple bool
}

//Create new combobox, multiselection og list item.
func NewSelect(size int, multiple bool, styles []Style, options ...*Option) *Selectform {
	return &Selectform{newWidget(styles...), options, size, multiple}
}

func (s *Selectform) HTML() (html string) {
	if s.multiple {
		html = fmt.Sprintf(`<select id="%s" size="%d" style="%s" multiple>`, s.id, s.size, s.style.Marshal())
	} else {
		html = fmt.Sprintf(`<select id="%s" size="%d" style="%s">`, s.id, s.size, s.style.Marshal())
	}
	for _, v := range s.options {
		html += v.HTML()+"\n"
	}
	html += "</select>"
	return
}

func (s *Selectform) SetOptions(o ...*Option) {
	s.options = o
	html := ""
	for _, v := range s.options {
		html += v.HTML()
	}
	events <- Event(jq(s.id, "html('"+escape(html)+"')"), nil)
}

func (s *Selectform) Selected() (string, []string) {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").val(); if (reply == null) {reply = ""}`, s.id), reply)
	events <- evt
	if s.multiple {
		return "", strings.Split(<-evt.reply, ",")
	}
	return <-evt.reply, nil
}

type Option struct {
	*widget
	value string
	text  string
}

func NewOption(value, text string) *Option {
	out := &Option{newWidget(), value, text}
	return out
}

func NewOptions(values ...string) []*Option {
	buf := make([]*Option, len(values))
	for i, v := range values {
		buf[i] = NewOption(v, v)
	}
	return buf
}

func (o *Option) HTML() string {
	return fmt.Sprintf(`<option id="%s" value="%s" style="%s">%s</option>`, o.id, o.value, o.style.Marshal(), o.text)
}

//=============================================
//  Select  //
//=============================================

type Modal struct {
	*widget
	content HTMLer
	width int
	height int
}

func NewModal(width, height int) *Modal {
	s := Style{
		"display": "none",
		"position": "absolute",
		"left": "0px",
		"top": "0px",
		"width":"100%",
		"height":"100%",
		"text-align":"center",
		"z-index": "1000",
		"background": "rgba(0, 0, 0, 0.6)",
	}
	return &Modal{newWidget(s), Html(""), width, height}
}

func (m *Modal) SetContent(content HTMLer) {
	m.content = content
}

func (m *Modal) HTML() string {
	underlaystyle := Style{
		"width": fmt.Sprintf("%dpx", m.width),
		"height": fmt.Sprintf("%dpx", m.height),
		"margin": "20% auto",
		"background-color": "white",
		"border":"1px solid #000",
		"padding":"15px",
		"text-align":"center",
	}
	return fmt.Sprintf(`<div style="%s" id="%s"><div style="%s">%s</div></div>`,
	m.style.Marshal(), m.id, underlaystyle.Marshal(), m.content.HTML())
}

//=============================================
//  Slider  //
//=============================================

//=============================================
//  Gauge  //
//=============================================

type Gauge struct {
	*widget
	value int
	width int
	color string
}

func NewGauge(value, width int, color string) *Gauge {
	s := Style{"border":"black solid 1px", "vertical-align": "middle"}
	g := &Gauge{newWidget(s), value, width, color}
	g.SetValue(value)
	return g
}

func (g *Gauge) SetValue(pct int) {
	switch {
	case pct > 100: pct = 100
	case pct < 0: pct = 0
	}
	js := fmt.Sprintf(`
	canvas = document.getElementById("%s");
	ctx = canvas.getContext("2d");  
  
    ctx.fillStyle = "%s";  
    ctx.fillRect (0, 0, %d, 20);
	`, g.id, g.color, g.width*pct/100)
	g.value = pct
	events <- Event(js, nil)
}

func (g *Gauge) Value() int {
	return g.value
}

func (g *Gauge) HTML() string {
	return fmt.Sprintf(`<canvas id="%s" width="%d" height="20" style="%s"></canvas>`, g.id, g.width, g.style)
}


//=============================================
//  Treecontrol  //
//=============================================

//=============================================
//  Spincontrol  //
//=============================================
/*
Da fuck am I gonna do here D:
type Spinctrl struct {
	*widget
	value int
	min int
	max int
	step int
}

func NewSpinctrl(value, min, max, step int) *Spinctrl { //&#x25B2; UP --- &#x25BC; DOWN
	return &Spinctrl{newWidget(), value, min, max, step}
}

func (s *Spinctrl) HTML() string {
	html := fmt.Sprintf(`
	<input type="text" value="%s" id="%s-tekst"/>
	<input type="button" value="&#x25B2;" id="%s-op" onclick="" style="vertical-align:super;height:5px;"/>
	<input type="button" value="&#x25BC;" id="%s-ned" onclick="" style="vertical-align:sub;height:5px;"/>
	`)
}
*/
//=============================================
//  Misc  //
//=============================================

type texthmler struct {
	value string
}

func (t texthmler) HTML() string {
	return t.value
}

func Html(value string) HTMLer {
	return texthmler{value}
}

//Javascript alert
func Alert(s string) {
	events <- Event("alert("+escape(s)+");", nil)
}

