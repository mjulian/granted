package browsers

import (
	"errors"
	"os/exec"
	"runtime"
)

const (
	ChromeKey        string = "CHROME"
	FirefoxKey       string = "FIREFOX"
	FirefoxStdoutKey string = "FIREFOX_STDOUT"
	EdgeKey          string = "EDGE"
	BraveKey         string = "BRAVE"
	StdoutKey        string = "STDOUT"
	ChromiumKey      string = "CHROMIUM"
)

// A few default paths to check for the browser
var ChromePathMac = []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"}
var ChromePathLinux = []string{`/usr/bin/google-chrome`, `/../../mnt/c/Program Files/Google/Chrome/Application/chrome.exe`, `/../../mnt/c/Program Files (x86)/Google/Chrome/Application/chrome.exe`}
var ChromePathWindows = []string{`\Program Files\Google\Chrome\Application\chrome.exe`, `\Program Files (x86)\Google\Chrome\Application\chrome.exe`}

var BravePathMac = []string{"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"}
var BravePathLinux = []string{`/usr/bin/brave-browser`, `/../../mnt/c/Program Files/BraveSoftware/Brave-Browser/Application/brave.exe`}
var BravePathWindows = []string{`\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`}

var EdgePathMac = []string{"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"}
var EdgePathLinux = []string{`/usr/bin/edge`, `/../../mnt/c/Program Files (x86)/Microsoft/Edge/Application/msedge.exe`}
var EdgePathWindows = []string{`\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`}

var FirefoxPathMac = []string{"/Applications/Firefox.app/Contents/MacOS/firefox"}
var FirefoxPathLinux = []string{`/usr/bin/firefox`, `/../../mnt/c/Program Files/Mozilla Firefox/firefox.exe`}
var FirefoxPathWindows = []string{`\Program Files\Mozilla Firefox\firefox.exe`}

var ChromiumPathMac = []string{"/Applications/Chromium.app/Contents/MacOS/Chromium"}
var ChromiumPathLinux = []string{`/usr/bin/chromium`, `/../../mnt/c/Program Files/Chromium/chromium.exe`}
var ChromiumPathWindows = []string{`\Program Files\Chromium\chromium.exe`}

func ChromePathDefaults() ([]string, error) {
	//check linuxpath for binary install
	path, err := exec.LookPath("google-chrome-stable")
	if err != nil {
		path, err = exec.LookPath("google-chrome")
		if err == nil {
			return []string{path}, nil
		}
	}
	if err == nil {
		return []string{path}, nil
	}
	switch runtime.GOOS {
	case "windows":
		return ChromePathWindows, nil
	case "darwin":
		return ChromePathMac, nil
	case "linux":
		return ChromePathLinux, nil
	default:
		return nil, errors.New("os not supported")
	}
}
func BravePathDefaults() ([]string, error) {
	//check linuxpath for binary install
	path, err := exec.LookPath("brave")
	if err == nil {
		return []string{path}, nil
	}
	switch runtime.GOOS {
	case "windows":
		return BravePathWindows, nil
	case "darwin":
		return BravePathMac, nil
	case "linux":
		return BravePathLinux, nil
	default:
		return nil, errors.New("os not supported")
	}
}
func EdgePathDefaults() ([]string, error) {
	//check linuxpath for binary install
	path, err := exec.LookPath("edge")
	if err == nil {
		return []string{path}, nil
	}
	switch runtime.GOOS {
	case "windows":
		return EdgePathWindows, nil
	case "darwin":
		return EdgePathMac, nil
	case "linux":
		return EdgePathLinux, nil
	default:
		return nil, errors.New("os not supported")
	}
}
func FirefoxPathDefaults() ([]string, error) {
	//check linuxpath for binary install
	path, err := exec.LookPath("firefox")
	if err == nil {
		return []string{path}, nil
	}
	switch runtime.GOOS {
	case "windows":
		return FirefoxPathWindows, nil
	case "darwin":
		return FirefoxPathMac, nil
	case "linux":
		return FirefoxPathLinux, nil
	default:
		return nil, errors.New("os not supported")
	}
}

func ChromiumPathDefaults() ([]string, error) {
	//check linuxpath for binary install
	path, err := exec.LookPath("chromium")
	if err == nil {
		return []string{path}, nil
	}
	switch runtime.GOOS {
	case "windows":
		return ChromiumPathWindows, nil
	case "darwin":
		return ChromiumPathMac, nil
	case "linux":
		return ChromiumPathLinux, nil
	default:
		return nil, errors.New("os not supported")
	}
}
