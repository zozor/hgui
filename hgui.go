/*
hgui is a gui toolkit for webgui
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

func (w *widget) SetStyle(style Style) {
	w.style = style
	js := jq(w.id, `attr("style", "`+style.Marshal()+`")`)
	events <- Event(js, nil)
}

func (w *widget) Style() Style {
	return w.style
}

func (w *widget) AddStyle(n Style) {
	w.style.AddStyle(n)
	w.SetStyle(w.style)
}

func (w *widget) RemoveStyle(n Style) {
	w.style.RemoveStyle(n)
	w.SetStyle(w.style)
}

func (w *widget) Hide() {
	events <- Event(jq(w.id, "hide()"), nil)
}

func (w *widget) Show() {
	events <- Event(jq(w.id, "show()"), nil)
}

//This is set ONLY on the client side, the server does not record this. Don't set style with this
func (w *widget) SetAttribute(attribute, value string) {
	events <- Event(jq(w.id, "attr('"+attribute+"','"+value+"')"), nil)
}

//Not everything can be removed, see $.removeAttr() in jquery for details
func (w *widget) RemoveAttribute(attr string) {
	events <- Event(jq(w.id, "removeAttr('"+attr+"')"), nil)
}

//This is set ONLY on the client side, the server does not record this.
//The use of frame.Flip, removes everything made through this call
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
type frame struct {
	*widget
	content []HTMLer
	topframe bool
}

func NewFrame(styles ...Style) *frame {
	return &frame{newWidget(styles...), make([]HTMLer, 0, 20), false}
}

func (f *frame) Add(widget ...HTMLer) {
	f.content = append(f.content, widget...)
}

func (f *frame) Flip() {
	buf := make([]string, len(f.content))
	for i, v := range f.content {
		buf[i] = v.HTML()
	}
	events <- Event(jq(f.id, `html("`+escape(strings.Join(buf, ""))+`")`), nil)
}

func (f *frame) HTML() string {
	buf := make([]string, len(f.content))
	for i, v := range f.content {
		buf[i] = v.HTML()
	}
	if f.topframe {
		return fmt.Sprintf(`<body id="%s" style="%s">%s</body>`, f.id, f.style.Marshal(), strings.Join(buf, ""))
	}
	return fmt.Sprintf(`<div id="%s" style="%s">%s</div>`, f.id, f.style.Marshal(), strings.Join(buf, ""))
}

//=============================================
//  Label  //
//=============================================
type label struct {
	*widget
	value string
}

func NewLabel(value string, styles ...Style) *label {
	return &label{newWidget(styles...), value}
}

func (l *label) Value() string {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").html()`, l.id), reply)
	events <- evt
	return <-evt.reply
}

func (l *label) SetValue(s string) {
	events <- Event(fmt.Sprintf(`$("#%s").html("%s")`, l.id, escape(s)), nil)
}

func (l *label) HTML() string {
	return fmt.Sprintf(`<span id="%s" style="%s">%s</span>`, l.id, l.style.Marshal(), l.value)
}

//=============================================
//  Button  //
//=============================================
type button struct {
	*widget
	value  string
	action func()
}

func NewButton(value string, styles []Style, action func()) *button {
	b := &button{newWidget(styles...), value, action}
	if action != nil {
		handlers[b.id+".onclick"] = action
	}
	return b
}

func (b *button) HTML() string {
	return fmt.Sprintf(`<input type="button" value="%s" id="%s" onclick="callHandler('%s.onclick');" style="%s"/>`,
		b.value, b.id, b.id, b.style.Marshal())
}

//=============================================
//  Table  //
//=============================================
type table struct {
	*widget
	rows []*row
}

func NewTable(styles []Style, r ...*row) *table {
	return &table{newWidget(styles...), r}
}

func (t *table) Addrows(r ...*row) {
	t.rows = append(t.rows, r...)
}

func (t *table) HTML() (html string) {
	html = `<table id="` + t.id + `" style="`+t.style.Marshal()+`">`
	for _, v := range t.rows {
		html += v.HTML()
	}
	html += "</table>"
	return
}

type row struct {
	*widget
	cells []*cell
}

func NewRow(styles []Style, c ...*cell) *row {
	return &row{newWidget(styles...), c}
}

func (r *row) AddCells(c ...*cell) {
	r.cells = append(r.cells, c...)
}

func (r *row) HTML() (html string) {
	html = `<tr id="` + r.id + `" style="`+r.style.Marshal()+`">`
	for _, v := range r.cells {
		html += v.HTML()
	}
	html += "</tr>"
	return
}

type cell struct {
	*widget
	header  bool
	colspan int
	rowspan int
	content HTMLer
}

func NewCell(header bool, colspan, rowspan int, content HTMLer, styles ...Style) *cell {
	return &cell{newWidget(styles...), header, colspan, rowspan, content}
}

func (t *cell) HTML() (html string) {
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
type textinput struct {
	*widget
	value string
	_type string
}

type texttype string



func NewTextinput(value string, ttype texttype,styles ...Style) *textinput {
	return &textinput{newWidget(styles...), value, string(ttype)}
}

func (t *textinput) Value() string {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").val()`, t.id), reply)
	events <- evt
	return <-evt.reply
}

func (t *textinput) SetValue(s string) {
	events <- Event(fmt.Sprintf(`$("#%s").val("%s")`, t.id, escape(s)), nil)
}

func (t *textinput) HTML() string {
	return fmt.Sprintf(`<input type="%s" id="%s" value="%s" style="%s"/>`, t._type, t.id, t.value, t.style.Marshal())
}

//=============================================
//  Textarea  //
//=============================================
type textarea struct {
	*widget
	value string
}

func NewTextarea(value string, styles ...Style) *textarea {
	return &textarea{newWidget(styles...), value}
}

func (t *textarea) Value() string {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").val()`, t.id), reply)
	events <- evt
	return <-evt.reply
}

func (t *textarea) SetValue(s string) {
	events <- Event(fmt.Sprintf(`$("#%s").text("%s")`, t.id, escape(s)), nil)
}

func (t *textarea) HTML() string {
	return `<textarea id="` + t.id + `" style="`+t.style.Marshal()+`"></textarea>`
}

//=============================================
//  Radiobuttons checkboxes  //
//=============================================
type radiocheckbox struct { //TODO: Lav det færdigt
	*widget
	group   string
	radiobox bool
}

func NewRadioCheckbox(radiobox bool, group string, styles ...Style) *radiocheckbox {
	return &radiocheckbox{newWidget(styles...), group, radiobox}
}

func (t *radiocheckbox) Checked() bool {
	reply := make(chan string)
	evt := Event(`reply = $("#`+t.id+`").prop("checked")`, reply)
	events <- evt
	if <-reply == "true" {
		return true
	}
	return false
}

func (t *radiocheckbox) Check() {
	events <- Event(jq(t.id, `prop("checked", "checked")`), nil)
}

func (t *radiocheckbox) Uncheck() {
	events <- Event(jq(t.id, `prop("checked", false)`), nil)
}

func (t *radiocheckbox) HTML() string {
	if t.radiobox {
		return fmt.Sprintf(`<input type="radio" id="%s" name="%s" style="%s"/>`, t.id, t.group, t.style.Marshal())
	}
	return fmt.Sprintf(`<input type="checkbox" id="%s" style="%s"/>`, t.id, t.style.Marshal())
}

//=============================================
//  Image  //
//=============================================
type image struct {
	*widget
	src string
}

func NewImage(src string, styles ...Style) *image {
	return &image{newWidget(styles...), src}
}

func (i *image) HTML() string {
	return fmt.Sprintf(`<img id="%s" src="%s" style="%s"/>`, i.id, i.src, i.style.Marshal())
}

//=============================================
//  Lists  //
//=============================================
type list struct {
	*widget
	items []*listitem
	ordered bool
}

func NewList(ordered bool, styles []Style, items ...*listitem) *list {
	return &list{newWidget(styles...), items, ordered}
}

func (l *list) SetList(items ...*listitem) {
	l.items = items
	html := ""
	for _, v := range l.items {
		html += v.HTML()
	}
	events <- Event(jq(l.id, `html("`+escape(html)+`")`), nil)
}

func (l *list) HTML() (html string) {
	if l.ordered {
		html = fmt.Sprintf(`<ol id="%s" style="%s">`, l.id, l.style.Marshal())
	} else {
		html = fmt.Sprintf(`<ul id="%s" style="%s">`, l.id, l.style.Marshal())
	}
	
	for _, v := range l.items {
		html += v.HTML()
	}
	
	if l.ordered {
		html += "</ol>"
	} else {
		html += "</ul>"
	}
	return
}

type listitem struct {
	*widget
	value string
}

func NewListItem(value string, styles ...Style) *listitem {
	return &listitem{newWidget(styles...), value}
}

func (l *listitem) HTML() string {
	return fmt.Sprintf(`<li id="%s">%s</li>`, l.id, l.value)
}

//=============================================
//  Links  //
//=============================================
type link struct {
	*widget
	href  string
	value HTMLer
}

func NewLink(href string, value HTMLer, styles ...Style) *link {
	return &link{newWidget(styles...), href, value}
}

func (l link) HTML() string {
	return fmt.Sprintf(`<a href="%s" id="%s" style="%s">%s</a>`, l.href, l.id, l.style.Marshal(), l.value.HTML())
}

//=============================================
//  Fieldset  //
//=============================================
type fieldset struct {
	legend  string
	content HTMLer
}

func NewFieldset(caption string, content HTMLer) *fieldset {
	return &fieldset{caption, content}
}

func (f *fieldset) HTML() string {
	return `<fieldset><legend>` + f.legend + `</legend>` + f.content.HTML() + `</fieldset>`
}

//=============================================
//  Select  //
//=============================================
type selectform struct {
	*widget
	options  []*option
	size     int
	multiple bool
}

func NewSelect(size int, multiple bool, styles []Style, options ...*option) *selectform {
	return &selectform{newWidget(styles...), options, size, multiple}
}

func (s *selectform) HTML() (html string) {
	if s.multiple {
		html = fmt.Sprintf(`<select id="%s" size="%d" style="%s" multiple>`, s.id, s.size, s.style.Marshal())
	} else {
		html = fmt.Sprintf(`<select id="%s" size="%d" style="%s">`, s.id, s.size, s.style.Marshal())
	}
	for _, v := range s.options {
		html += v.HTML()
	}
	html += "</select>"
	return
}

func (s *selectform) SetOptions(o ...*option) {
	s.options = o
	html := ""
	for _, v := range s.options {
		html += v.HTML()
	}
	events <- Event(jq(s.id, "html('"+escape(html)+"')"), nil)
}

func (s *selectform) Selected() (string, []string) {
	reply := make(chan string)
	evt := Event(fmt.Sprintf(`reply = $("#%s").val(); if (reply == null) {reply = ""}`, s.id), reply)
	events <- evt
	if s.multiple {
		return "", strings.Split(<-evt.reply, ",")
	}
	return <-evt.reply, nil
}

type option struct {
	*widget
	value string
	text  string
}

func NewOption(value, text string) *option {
	out := &option{newWidget(), value, text}
	return out
}

func NewOptions(values ...string) []*option {
	buf := make([]*option, len(values))
	for i, v := range values {
		buf[i] = NewOption(v, v)
	}
	return buf
}

func (o *option) HTML() string {
	return fmt.Sprintf(`<option id="%s" value="%s" style="%s">%s</option>`, o.id, o.value, o.style.Marshal(), o.text)
}

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

func Alert(s string) {
	events <- Event("alert("+escape(s)+");", nil)
}

