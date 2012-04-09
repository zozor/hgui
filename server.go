 /*
A simple gui toolkit to communicate between html/javascript and go server.
On execution it opens a browser connecting it to the server.
*/
package hgui

import (
	"net/http"
	"os"
	"time"
	"fmt"
)

var firstTimeRequest = true

//Allow the user to open the gui in more tabs. This is not recommended and for DEBUG ONLY
var AllowMoreRequests = false

var resources = map[string][]byte{}
func SetResource(files map[string][]byte) {
	resources = files
}

var handlers = map[string]func() {}
var Topframe = &frame{newWidget(), make([]HTMLer, 0, 20), true}

//Killing the server when not active
var pingChannel chan bool

func dieCounter() {
	for {
		select {
		case <-pingChannel:
			continue
		case <-time.After(10e9):
			os.Exit(0)
		}
	}
}

//This starts the server with the address addr. should be localhost:23192 (or some other port)
func StartServer(port int) { //"127.0.0.1:3939"
	http.Handle("/", http.HandlerFunc(requests))
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	
	openBrowser(addr)

	pingChannel = make(chan bool)
	go dieCounter()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("You need to wait atleast 10 seconds before you start this program again")
	}
}

func requests(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/events" {
		pingChannel <- true
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
		if !firstTimeRequest && !AllowMoreRequests {
			w.Write([]byte("Refrehed? Restart server.<br/>New tab? don't do that :)"))
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
