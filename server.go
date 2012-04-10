 /*
A simple gui toolkit to communicate between html/javascript and go server.
On execution it opens a browser connecting it to the server.
*/
package hgui

import (
	"net/http"
	"os"
	"fmt"
)

var firstTimeRequest = true

var resources = map[string][]byte{}
func SetResource(files map[string][]byte) {
	resources = files
}

var handlers = map[string]func() {}
var Topframe = &frame{newWidget(), make([]HTMLer, 0, 20), true}

//This starts the server with the address addr. should be localhost:23192 (or some other port)
func StartServer(width, height, port int, title string) { //"127.0.0.1:3939"
	http.Handle("/", http.HandlerFunc(requests))
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
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
	return []byte(`
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
