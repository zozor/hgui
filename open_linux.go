package hgui

import "os/exec"

func openBrowser(addr string) {
	cmd := exec.Command("x-www-browser", addr)
	cmd.Start()
} 
