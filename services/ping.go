package services

import (
	"context"
	servertiming "github.com/mitchellh/go-server-timing"
	"google.golang.org/api/pagespeedonline/v5"
	"pinger/logger"
	"strings"
	"time"
)

var (
	pageSpeedSvc *pagespeedonline.Service
	ctx          context.Context
)

//
func runPageSpeed(ctx context.Context, url string, insightType string) (int64, error) {
	res, err := pageSpeedSvc.Pagespeedapi.Runpagespeed(url).
		Strategy(strings.ToUpper(insightType)).
		Category("PERFORMANCE").
		Context(ctx).
		Do()

	if err != nil {
		return 0, err
	}

	score := int64(res.LighthouseResult.Categories.Performance.Score.(float64) * 100)

	return score, nil
}

// GetPageSpeed - функция для полученные статистики о сайте
func GetPageSpeed(logger *logger.Logger, ctx context.Context, host string) (time.Duration, int64, error) {
	var err error
	var score int64

	pageSpeedSvc, err = pagespeedonline.NewService(ctx)
	timing := servertiming.FromContext(ctx)
	insightType := "DESKTOP"
	if err != nil {
		logger.Fatalf("%v", err)
		return 0, 0, err
	}

	pagespeedTiming := timing.NewMetric("pageSpeed").
		WithDesc("Google Page Speed Insights API").
		Start()

	score, err = runPageSpeed(
		ctx,
		host,
		insightType,
	)
	if err != nil {
		logger.Fatalf("%v", err)
		return 0, 0, err
	}

	pagespeedTiming.Stop()
	return pagespeedTiming.Duration, score, nil
}
