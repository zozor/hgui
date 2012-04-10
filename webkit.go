package hgui

/*
#cgo pkg-config: webkit-1.0
#cgo pkg-config: gtk+-2.0

#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <glib-object.h>
#include <stdlib.h>

static inline gchar* to_gcharptr(const char* s) { return (gchar*)s; }
static inline void free_string(char* s) { free(s); }

static GtkWindow* to_GtkWindow(GtkWidget* w) { return GTK_WINDOW(w); }
static GtkContainer* to_GtkContainer(GtkWidget* w) { return GTK_CONTAINER(w); }

static void loadHtmlString(GtkWidget* widget, gchar* pcontent, gchar* pbase_uri) {
	webkit_web_view_load_html_string(WEBKIT_WEB_VIEW(widget), pcontent, pbase_uri);
}

static void connect_destroy(GtkWidget* window) {
	gtk_signal_connect(GTK_OBJECT(window), "destroy", GTK_SIGNAL_FUNC(gtk_exit), NULL);
}
*/
import "C"
import "fmt"


func startGui(width, height int, title string, port int) {
	C.gtk_init(nil, nil); //gtk.Init(nil)
	
	window := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	C.gtk_window_set_title(C.to_GtkWindow(window), C.to_gcharptr(C.CString(title)))
	C.connect_destroy(window)

	vbox := C.gtk_hbox_new(0, 1)

	webview := C.webkit_web_view_new()
	
	C.gtk_container_add(C.to_GtkContainer(vbox), webview);

	embed := `
<iframe width="100%" height="100%" frameborder="0" scrolling="no" marginheight="0" marginwidth="0" src="http://127.0.0.1:`+fmt.Sprintf("%d", port)+`"></iframe>
`
	loadHtmlString(webview, embed, ".")	

	C.gtk_container_add(C.to_GtkContainer(window), vbox)
	C.gtk_widget_set_size_request(window, C.gint(width), C.gint(height))
	
	C.gtk_widget_show(vbox);
	C.gtk_widget_show(window); //Window.ShowAll()
    C.gtk_widget_show(webview);

	/*
	This only matters if proxy is stupid!
	proxy := os.Getenv("HTTP_PROXY")
	if len(proxy) > 0 {
		ptr := C.CString(uri)
		C.proxyshit(ptr)
		C.free(ptr)
	}
	*/
	
	C.gtk_main(); //gtk.GtkMain()
}

func loadHtmlString(webview *C.GtkWidget, content, base_uri string) {
	pcontent := C.CString(content)
	defer C.free_string(pcontent)
	pbase_uri := C.CString(base_uri)
	defer C.free_string(pbase_uri)
	C.loadHtmlString(webview, C.to_gcharptr(pcontent), C.to_gcharptr(pbase_uri))
}
