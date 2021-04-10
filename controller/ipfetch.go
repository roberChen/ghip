// Package controller contains controllers with control logics
package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/roberChen/ghip/ipget"
	"github.com/roberChen/ghip/module"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IPController struct {
	db      *gorm.DB
	urlmaps map[string]*module.IP
	config  *viper.Viper
}

// NewIPController initialise a ip controller, with cobra settings. It will load urls setted in config files to url maps.
func NewIPController(db *gorm.DB, conf *viper.Viper) (*IPController, error) {
	if db == nil || conf == nil {
		return nil, fmt.Errorf("invalid nil parameters")
	}
	ipctrl := new(IPController)
	ipctrl.db = db
	ipctrl.config = conf
	// insert urlmaps
	ipctrl.urlmaps = make(map[string]*module.IP)
	for num, urlitem := range ipctrl.config.GetStringSlice("urls.list") {
		if _, ok := ipctrl.urlmaps[urlitem]; ok {
			log.Printf("duplicate url `%s` No. %d\n", urlitem, num)
			continue
		}
		ip := new(module.IP)
		if err := ip.LoadFrom(db, urlitem); err != nil {
			return nil, fmt.Errorf("loading data from url `%s` failed: %s", urlitem, err)
		}
		ipctrl.urlmaps[urlitem] = ip
	}
	return ipctrl, nil
}

// GetURLIPsLocal returns ip of url fetched from local
func (ipctrl *IPController) GetURLIPsLocal() map[string]string {
	m := make(map[string]string)
	for urlname, ip := range ipctrl.urlmaps {
		if ip.LocalIP == "" {
			continue
		}
		m[urlname] = ip.LocalIP
	}
	return m
}

// GetURLIPsADDRCOM returns ip of url fetched from ipaddress.com
func (ipctrl *IPController) GetURLIPsADDRCOM() map[string][]string {
	m := make(map[string][]string)
	for urlname, ip := range ipctrl.urlmaps {
		l := []string{}
		if err := json.Unmarshal(ip.ADDRCOMIP, &l); err != nil {
			panic(fmt.Errorf("unexpected ipaddress.com ip list %s: %s", ip.ADDRCOMIP, err))
		}
		if len(l) == 0 {
			continue
		}
		m[urlname] = l
	}

	return m
}

func (ipctrl *IPController) Updates() (bool, error) {
	var updated bool
	log.Printf("updating database [size (%d)]", len(ipctrl.urlmaps))
	for urlname, ip := range ipctrl.urlmaps {
		log.Printf("fetching for url `%s`", urlname)
		localip, err := ipget.IPAddrLocal(urlname, "80")
		if err != nil {
			log.Printf("getting ip for `%s` from local: %s\n", urlname, err)
		}
		ip.IPLocal(localip)
		addrcomip, err := ipget.IPAddrServer(urlname)
		if err != nil {
			log.Printf("getting ip for `%s` from ipaddress.com: %s\n", urlname, err)
		}
		ip.IPADDRCOM(addrcomip)
		up, err := ip.Write(ipctrl.db)
		if err != nil {
			log.Printf("write to database failed: %s\n", err)
		}
		updated = up && updated
	}

	return updated, nil
}
