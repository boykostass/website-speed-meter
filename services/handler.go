package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	apperror "pinger/app_error"
	"pinger/logger"
	"time"
)

const (
	sitesURL      = "/sites"
	sitesStatsURL = "/sites/stats"
)

type handler struct {
	logger   *logger.Logger
	repQuery RepQuery
	database Database
}

func NewHandler(logger *logger.Logger, database Database, repQuery RepQuery) Handler {
	return &handler{
		logger:   logger,
		repQuery: repQuery,
		database: database,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, sitesURL, apperror.Middleware(h.GetInfo))
	router.HandlerFunc(http.MethodGet, sitesStatsURL, apperror.Middleware(h.GetStats))
}

func (h *handler) GetInfo(w http.ResponseWriter, r *http.Request) error {
	var lastCheck time.Time
	var answer SiteInfo
	var found bool

	site := r.URL.Query().Get("search")
	u, err := url.Parse(site)
	if err != nil {
		return err
	}
	u.Scheme = "https"
	res, err := h.repQuery.FindOne(context.TODO(), site)
	if err == nil {
		for _, member := range res {
			lastCheck, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprint(member.Date+" "+member.Time), time.Local)
			if time.Since(lastCheck) <= (24 * time.Hour) {
				answer = member
				found = true
				break
			}
		}
	}
	if found == false {
		delay, perf, err := GetPageSpeed(h.logger, context.TODO(), u.String())
		if err != nil {
			return err
		}
		dateNow := time.Now().Format("2006-01-02")
		timeNow := time.Now().Format("15:04:05")
		_, err = h.database.Query(context.TODO(), "INSERT INTO list (site, date, time, delay, performance) VALUES ($1, $2, $3, $4, $5)", site, dateNow, timeNow, delay.String(), fmt.Sprint(perf))
		if err != nil {
			return err
		}
		answer.Site = site
		answer.Date = time.Now().Format("2006-01-02")
		answer.Time = time.Now().Format("15:04:05")
		answer.Delay = delay.String()
		answer.Performance = fmt.Sprint(perf)
	}

	allBytes, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}

func (h *handler) GetStats(w http.ResponseWriter, r *http.Request) error {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	all, err := h.repQuery.FindAll(context.TODO())
	if err != nil {
		return err
	}

	temp := make(map[string][]time.Duration)
	stats := make(map[string]float64)
	fromTime, err := time.ParseInLocation("2006-01-02", from, time.Local)
	if err != nil {
		return err
	}
	toTime, err := time.ParseInLocation("2006-01-02", to, time.Local)
	if err != nil {
		return err
	}
	if fromTime.Sub(toTime) > 0 {
		return errors.New("incorrect data entered")
	}
	for _, one := range all {
		oneDateTime, err := time.ParseInLocation("2006-01-02", one.Date, time.Local)
		if err != nil {
			return err
		}
		if oneDateTime.Sub(fromTime) >= 0 && toTime.Sub(oneDateTime) < 0 {
			fmt.Println(oneDateTime.Sub(fromTime))
			fmt.Println(toTime.Sub(oneDateTime))
			continue
		}
		delayDur, err := time.ParseDuration(one.Delay)
		if err != nil {
			return err
		}

		temp[one.Site] = append(temp[one.Site], delayDur)
	}

	for key, one := range temp {
		var sum float64

		for _, val := range one {
			sum = sum + val.Seconds()
		}
		stats[key] = sum / float64(len(one))
	}

	allBytes, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}
