package hgui

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/mattn/go-webkit/webkit"
import "fmt"

func startGui(width, height int, title string, port int) {
	gtk.Init(nil)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle(title)
	window.Connect("destroy", gtk.MainQuit)

	vbox := gtk.VBox(false, 1)

	/*swin := gtk.ScrolledWindow(nil, nil)
	swin.SetShadowType(gtk.GTK_SHADOW_IN)
	swin.SetPolicy(gtk.GTK_POLICY_ALWAYS, gtk.GTK_POLICY_ALWAYS)
	*/
	webview := webkit.WebView()
	
	//swin.Add(webview)

	//vbox.Add(swin)
	vbox.Add(webview)

	embed := `
<iframe width="100%" height="100%" frameborder="0" scrolling="no" marginheight="0" marginwidth="0" src="http://127.0.0.1:`+fmt.Sprintf("%d", port)+`"></iframe>
`
	webview.LoadHtmlString(fmt.Sprintf(embed, port), ".")
	

	window.Add(vbox)
	window.SetSizeRequest(width, height)
	window.ShowAll()

	proxy := os.Getenv("HTTP_PROXY")
	if len(proxy) > 0 {
		soup_uri := webkit.SoupUri(proxy)
		webkit.GetDefaultSession().Set("proxy-uri", soup_uri)
		soup_uri.Free()
	}
	gtk.Main()
}
