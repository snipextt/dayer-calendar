package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func catchError() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}

func Init() {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println(err)
	}
	s := gocron.NewScheduler(location)

	s.Every(1).Day().At("00:00").Do(timedoctorCron)

	s.StartAsync()
}
