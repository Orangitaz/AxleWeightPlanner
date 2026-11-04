package main

import (
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

//go:embed index.html
var page string

const head = "<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"><style>:root{color-scheme:light dark}body{margin:0;font:14px/1.45 system-ui,sans-serif}</style></head><body>"
const tail = "<script>setInterval(function(){fetch('/ping').catch(function(){})},4000);fetch('/ping');</script></body></html>"

func main() {
	var pings atomic.Int64
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		pings.Add(1)
		w.WriteHeader(204)
	})
	// Saved setups live in setups.json next to the program.
	mux.HandleFunc("/api/setups", func(w http.ResponseWriter, r *http.Request) {
		exe, err := os.Executable()
		if err != nil {
			http.Error(w, "no path", 500)
			return
		}
		dir := filepath.Dir(exe)
		if d := os.Getenv("AXLE_DATA_DIR"); d != "" {
			dir = d
		}
		file := filepath.Join(dir, "setups.json")
		switch r.Method {
		case "GET":
			b, err := os.ReadFile(file)
			if err != nil {
				b = []byte("[]")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		case "PUT":
			b, err := io.ReadAll(io.LimitReader(r.Body, 20<<20))
			if err != nil {
				http.Error(w, "bad body", 400)
				return
			}
			tmp := file + ".tmp"
			if os.WriteFile(tmp, b, 0644) != nil || os.Rename(tmp, file) != nil {
				http.Error(w, "cannot write next to the program", 500)
				return
			}
			w.WriteHeader(204)
		default:
			http.Error(w, "method", 405)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, head+page+tail)
	})
	ln, err := net.Listen("tcp", "127.0.0.1:47813") // fixed port keeps saved settings between runs
	if err != nil {
		// already running? just open another window on it
		c := http.Client{Timeout: 2 * time.Second}
		if resp, e := c.Get("http://127.0.0.1:47813/ping"); e == nil {
			resp.Body.Close()
			openApp("http://127.0.0.1:47813/")
			return
		}
		ln, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			os.Exit(1)
		}
	}
	url := "http://" + ln.Addr().String() + "/"
	go http.Serve(ln, mux)
	openApp(url)
	// Quit once the window has been closed: no heartbeat for 4 minutes. This counts 2 second ticks instead of
	// comparing clock times, so a laptop that slept for an hour does not make the program quit on wake-up.
	time.Sleep(60 * time.Second)
	seen, idle := pings.Load(), 0
	for {
		time.Sleep(2 * time.Second)
		if n := pings.Load(); n != seen {
			seen, idle = n, 0
			continue
		}
		if idle++; idle > 120 {
			return
		}
	}
}

func openApp(url string) {
	if runtime.GOOS == "windows" {
		pf := []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"), os.Getenv("LocalAppData")}
		rel := []string{`\Microsoft\Edge\Application\msedge.exe`, `\Google\Chrome\Application\chrome.exe`}
		for _, r := range rel {
			for _, p := range pf {
				if p == "" {
					continue
				}
				exe := strings.TrimRight(p, `\`) + r
				if _, err := os.Stat(exe); err == nil {
					if exec.Command(exe, "--app="+url, "--window-size=1400,950").Start() == nil {
						return
					}
				}
			}
		}
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		return
	}
	for _, name := range []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "microsoft-edge", "microsoft-edge-stable", "brave-browser"} {
		if p, err := exec.LookPath(name); err == nil {
			if exec.Command(p, "--app="+url, "--window-size=1400,950").Start() == nil {
				return
			}
		}
	}
	for _, name := range []string{"xdg-open", "gio", "sensible-browser"} {
		if p, err := exec.LookPath(name); err == nil {
			args := []string{url}
			if name == "gio" {
				args = []string{"open", url}
			}
			if exec.Command(p, args...).Start() == nil {
				return
			}
		}
	}
}
