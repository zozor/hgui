 /*
A simple gui toolkit to communicate between html/javascript and go server.
On execution it opens a browser connecting it to the server.
*/
package hgui

import (
	"net/http"
	"net"
	"os"
	"fmt"
	"math/rand"
)

var firstTimeRequest = true
var DEBUG = false


var resources = map[string][]byte{}
//When you compile a file, be it image, or page or whatever, to a []byte, it can be used with this map.
//when the page is requested on the server, fx. /img/cat.jpg, it will write the bytes in 
//		hgui.SetResource(map[string][]byte{"/img/cat.jpg", catpicvar})
//back to the client.
func SetResource(files map[string][]byte) {
	resources = files
}

var handlers = map[string]func() {}
var Topframe = &Frame{newWidget(), make([]HTMLer, 0, 20), true}

//This starts the server on the specified port. It also runs the mainloop for webkit.
//It also takes width and heigh + a title for the window to appear in.
func StartServer(width, height int, title string) { //"127.0.0.1"
	http.Handle("/", http.HandlerFunc(requests))
	
	var port int
	//find port
	for {
		port = rand.Intn(40000)+10000
		ln, err := net.Listen("tcp",fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		err = ln.Close()
		if err != nil {
			panic(err)
		}
		break
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	
	
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	
	if DEBUG {
		fmt.Println(addr)
		c := make(chan int)
		<-c
	}
	startGui(width, height, title, port)
}

func requests(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/events" {
		eventPoll(w)
		return
	}
	
	println(req.URL.Path)
	if req.URL.Path == "/reply" {
		q := req.URL.Query()
		eventReply(reply{q.Get("Id"), q.Get("Reply")})
		return
	}
	
	if req.URL.Path == "/handler" {
		q := req.URL.Query()
		if f, ok := handlers[q.Get("id")]; ok {
			f()
		}
		return
	}
	
	if req.URL.Path == "/" {
		if !firstTimeRequest {
			os.Exit(0)
			return
		}
		firstTimeRequest = false
		w.Write(head())
		w.Write([]byte(Topframe.HTML()))
		w.Write(bottom())
		return
	}
	
	if req.URL.Path == "/js" { //<script type="text/javascript" src="/webgui"></script>
		w.Header().Set("Content-Type", "text/javascript")
		w.WriteHeader(http.StatusOK)
		w.Write(fileJQuery)
		w.Write([]byte("\n\n"))
		w.Write(filecorejs)
		return
	}
		
	b, ok := resources[req.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page Not Found - 404"))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func head() []byte {
	return []byte(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
<script type="text/javascript" src="js"/></script>
</head>
`)
}
func bottom() []byte {
	return []byte(`
</html>
`)
}
