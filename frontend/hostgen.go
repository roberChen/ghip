// Package frontend includes everythings about frontend, including index page
package frontend

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/roberChen/ghip/controller"
)

//go:embed statics
var statics embed.FS

func GenHostListsHTML(ipctrl *controller.IPController) []byte {
	t, err := template.New("host.go.tmpl").ParseFS(statics, "statics/host.go.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	d := new(Display)
	localsmap := ipctrl.GetURLIPsLocal()
	ipaddrmap := ipctrl.GetURLIPsADDRCOM()
	d.Locals = make([]struct {
		URL string
		IP  string
	}, 0)
	d.IPADDRCOM = make([]struct {
		URL string
		IPS []string
	}, 0)
	for urlnames, ip := range localsmap {
		d.Locals = append(d.Locals, struct {
			URL string
			IP  string
		}{
			URL: urlnames,
			IP:  ip,
		})
	}
	for urlnames, ips := range ipaddrmap {
		d.IPADDRCOM = append(d.IPADDRCOM, struct {
			URL string
			IPS []string
		}{
			URL: urlnames,
			IPS: ips,
		})
	}
	out := new(bytes.Buffer)
	if err := t.Execute(out, d); err != nil {
		log.Fatalln(err)
	}
	return out.Bytes()
}

// IndexPage writes index page
func IndexPage(w http.ResponseWriter, ipctrl *controller.IPController) error {
	c, err := os.Open("cache/index.html")
	if err == nil {
		if _, err := io.Copy(w, c); err != nil {
			return fmt.Errorf("IndexPage write: %s", err)
		}
	}
	// write cache
	if err := os.MkdirAll("cache", fs.ModePerm); err != nil {
		return fmt.Errorf("create cache: %s", err)
	}
	b, err := GenIndexPage(ipctrl)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// GenIndexPage will generate a index.html file according to index.md and cache/host.html, and write to cache/index.html, also return
func GenIndexPage(ipctrl *controller.IPController) ([]byte, error) {
	d, err := os.ReadFile("cache/host.html")
	if err != nil {
		d = UpdateHostFile(ipctrl)
	}
	t, err := template.New("index.md").ParseFiles("page/index.md")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, template.HTML(d)); err != nil {
		panic(err)
	}
	b := buf.Bytes()
	if err := os.WriteFile("cache/index.html", b, fs.ModePerm); err != nil {
		return nil, fmt.Errorf("write file failed: %s", err)
	}
	return b, nil
}

// UpdateHostFile will update cache/host.html by refetching
func UpdateHostFile(ipctrl *controller.IPController) []byte {
	if _, err := ipctrl.Updates(); err != nil {
		log.Fatalf("Update failed: %s", err)
	}
	d := GenHostListsHTML(ipctrl)
	if err := os.WriteFile("cache/host.html", d, os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	return d
}

type Display struct {
	Locals []struct {
		URL string
		IP  string
	}
	IPADDRCOM []struct {
		URL string
		IPS []string
	}
}
