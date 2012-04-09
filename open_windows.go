package hgui

import "os/exec"

func openBrowser(addr string) {
	cmd := exec.Command("start", addr)
	cmd.Start()
} 
