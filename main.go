package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/roberChen/ghip/conf"
	"github.com/roberChen/ghip/controller"
	"github.com/roberChen/ghip/cronjobs"
	"github.com/roberChen/ghip/database"
	"github.com/roberChen/ghip/frontend"
)

var (
	conffile *string
	help     *bool
	gen      *bool
)

func init() {
	conffile = flag.String("conf", "ghip.toml", "config file")
	help = flag.Bool("help", false, "help page")
	gen = flag.Bool("gen", false, "generate config file with specific config file path, specified by conffile")

	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *gen {
		if err := conf.GenDefaultConfig(*conffile); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}
	c := conf.ReadConfig(*conffile)
	conf.ParseConfigs(c)
	db := database.Open(c)
	ipctrl, err := controller.NewIPController(db, c)
	if err != nil {
		panic(err)
	}
	frontend.CheckConfigUpdate(ipctrl, c)
	cronjobs.Trigger(func() error {
		log.Println("cron job")
		d := frontend.GenHostListsHTML(ipctrl)
		if err := os.WriteFile("cache/host.html", d, os.ModePerm); err != nil {
			log.Fatalf("updating host.html: %s", err)
		}
		_, err := frontend.GenIndexPage(ipctrl)
		return err
	})
	if err := cronjobs.Initial(c, ipctrl); err != nil {
		log.Fatalf("cron init failed: [%s]: %s", c.GetString("cron.spec"), err)
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] from %s\n", r.Method, r.RemoteAddr)
		if r.Method != http.MethodGet {
			fmt.Fprintf(rw, "invalid method: %s", r.Method)
			log.Printf("invalid method: %s", r.Method)
			return
		}
		err := frontend.IndexPageHandler(rw, ipctrl)
		if err != nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(rw, "inner error")
			log.Printf("E: %s", err)
		}
	})

	log.Fatalln(http.ListenAndServe(":"+c.GetString("server.port"), nil))
}
