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
	"reflect"
	"sort"
	"time"

	"github.com/roberChen/ghip/controller"
	"github.com/spf13/viper"
)

//go:embed statics
var statics embed.FS

// GenHostListsHTML generates host list as an HTML node with data from ipctrl, to be embed in index.md
func GenHostListsHTML(ipctrl *controller.IPController) []byte {
	t, err := template.New("host.go.tmpl").ParseFS(statics, "statics/host.go.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	d := new(Display)
	d.Time = time.Now().String()
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
	for _, urlname := range ipctrl.Sequence {
		ip, ok := localsmap[urlname]
		if !ok {
			continue
		}
		d.Locals = append(d.Locals, struct {
			URL string
			IP  string
		}{
			URL: urlname,
			IP:  ip,
		})
	}
	for _, urlname := range ipctrl.Sequence {
		ips, ok := ipaddrmap[urlname]
		if !ok {
			continue
		}
		d.IPADDRCOM = append(d.IPADDRCOM, struct {
			URL string
			IPS []string
		}{
			URL: urlname,
			IPS: ips,
		})
	}
	out := new(bytes.Buffer)
	if err := t.Execute(out, d); err != nil {
		log.Fatalln(err)
	}
	return out.Bytes()
}

// IndexPageHandler writes index page when handling the http request, requires ipctrl when update is needed
func IndexPageHandler(w http.ResponseWriter, ipctrl *controller.IPController) error {
	c, err := os.Open("cache/index.html")
	if err == nil {
		if _, err := io.Copy(w, c); err != nil {
			return fmt.Errorf("IndexPage write: %s", err)
		}
		return nil
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

// GenIndexPage will generate a index.html file according to index.md and cache/host.html, and write to cache/index.html, also return it
func GenIndexPage(ipctrl *controller.IPController) ([]byte, error) {
	log.Printf("updating index.html")
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
	log.Printf("updating host.html")
	if _, err := ipctrl.Updates(); err != nil {
		log.Fatalf("Update failed: %s", err)
	}
	d := GenHostListsHTML(ipctrl)
	if err := os.WriteFile("cache/host.html", d, os.ModePerm); err != nil {
		log.Fatalf("updating host.html: %s",err)
	}
	return d
}

// CheckConfigUpdate will check whether config url lists has been updated, if is, then refetch
func CheckConfigUpdate(ipctrl *controller.IPController, config *viper.Viper) {
	seq := []string{}
	confseq := []string{}
	copy(seq, ipctrl.Sequence)
	copy(confseq, config.GetStringSlice("urls.list"))
	sort.Strings(seq)
	sort.Strings(confseq)
	if reflect.DeepEqual(seq, confseq) {
		return
	}
	UpdateHostFile(ipctrl)
	if _, err := GenIndexPage(ipctrl); err != nil {
		log.Printf("check config update: %s", err)
	}
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
	Time string
}
