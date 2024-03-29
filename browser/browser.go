package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Attempts to open the given url in the default browser
func Open(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform for browser opening")
	}
	return
}
