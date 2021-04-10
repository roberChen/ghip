// Package cronjobs contains functions to init cron jobs
package cronjobs

import (
	"log"


	"github.com/roberChen/ghip/controller"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var Crons *cron.Cron
var trigger Triggers

type Triggers []func() error

func Trigger(flist... func()error) {
	trigger = flist
}

func init() {
	Crons = cron.New()
	go Crons.Run()
}

// Initial will init cron tasks. cron spec setted in config file
func Initial(conf *viper.Viper, ipctrl *controller.IPController) (error) {
	return Crons.AddFunc(conf.GetString("cron.spec"), func ()  {
		update, err := ipctrl.Updates()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
		if update {
			log.Println("got update, triggering registered files")
			for id, f:= range trigger {
				if err := f(); err != nil {
					log.Printf("trigger func #%d: %s", id, err)
				}
			}
			
		}
	})
}
