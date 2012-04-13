package hgui

/*
#cgo pkg-config: webkit-1.0
#cgo pkg-config: gtk+-2.0

#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <glib-object.h>
#include <stdlib.h>

GtkWidget* window;
GtkWidget* webview;

static inline gchar* to_gcharptr(const char* s) { return (gchar*)s; }
static inline void free_string(char* s) { free(s); }

static GtkWindow* to_GtkWindow(GtkWidget* w) { return GTK_WINDOW(w); }
static GtkContainer* to_GtkContainer(GtkWidget* w) { return GTK_CONTAINER(w); }

static void loadHtmlString(GtkWidget* widget, gchar* pcontent, gchar* pbase_uri) {
	webkit_web_view_load_html_string(WEBKIT_WEB_VIEW(widget), pcontent, pbase_uri);
}

static void connect_destroy(GtkWidget* window) {
	gtk_signal_connect(GTK_OBJECT(window), "destroy", G_CALLBACK(gtk_exit), NULL);
}

static void loadUri(GtkWidget *widget, gchar* uri) {
	webkit_web_view_load_uri(WEBKIT_WEB_VIEW(widget), uri);
}

static void scriptEvent() {
	webkit_web_view_execute_script(WEBKIT_WEB_VIEW(webview), "GetEvents();");
}

static void _emit_script() {
	g_signal_emit_by_name(G_OBJECT(webview), "send-script");
}


static GtkWidget* _new_webkit() {
	
	GtkWidget* ww = webkit_web_view_new();
	
	g_signal_new("send-script",
             G_TYPE_OBJECT, G_SIGNAL_RUN_FIRST,
             0, NULL, NULL,
             g_cclosure_marshal_VOID__POINTER,
             G_TYPE_NONE, 1, G_TYPE_POINTER);
    
	gtk_signal_connect(GTK_OBJECT(ww), "send-script", G_CALLBACK(scriptEvent), NULL);
	
	return ww;
}

*/
import "C"
import "fmt"


func startGui(width, height int, title string, port int) {
	C.gtk_init(nil, nil); //gtk.Init(nil)
	
	window := C.window
	window = C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	C.gtk_window_set_title(C.to_GtkWindow(window), C.to_gcharptr(C.CString(title)))
	C.connect_destroy(window)

	vbox := C.gtk_hbox_new(0, 1)
	
	C.webview = C._new_webkit()
	
	C.gtk_container_add(C.to_GtkContainer(vbox), C.webview);

	C.loadUri(C.webview, C.to_gcharptr(C.CString(fmt.Sprintf("http://127.0.0.1:%d", port))))

	C.gtk_container_add(C.to_GtkContainer(window), vbox)
	C.gtk_widget_set_size_request(window, C.gint(width), C.gint(height))
	
	C.gtk_widget_show(vbox);
	C.gtk_widget_show(window); //Window.ShowAll()
    C.gtk_widget_show(C.webview);

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

func jsGetEvents() {
	C._emit_script()
}

//=============================================
//  Statusbar  //
//=============================================

//=============================================
//  Window Options  //
//=============================================
