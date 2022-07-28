package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"pinger/internal/config"
	"pinger/logger"
	"pinger/utils"
	"strings"
	"time"
)

type repository struct {
	database Database
	logger   *logger.Logger
}

func (r *repository) FindAll(ctx context.Context) ([]SiteInfo, error) {
	q := `
		SELECT site, date, time, delay, performance FROM list
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.database.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	sites := make([]SiteInfo, 0)
	for rows.Next() {
		var site SiteInfo

		err = rows.Scan(&site.Site, &site.Date, &site.Time, &site.Delay, &site.Performance)
		if err != nil {
			return nil, err
		}

		sites = append(sites, site)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

func (r *repository) FindOne(ctx context.Context, name string) ([]SiteInfo, error) {
	q := `
		SELECT site, date, time, delay, performance FROM list WHERE site = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.database.Query(ctx, q, name)
	if err != nil {
		return nil, err
	}

	sites := make([]SiteInfo, 0)
	for rows.Next() {
		var site SiteInfo

		err = rows.Scan(&site.Site, &site.Date, &site.Time, &site.Delay, &site.Performance)
		if err != nil {
			return nil, err
		}

		sites = append(sites, site)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}

func NewRepository(logger *logger.Logger, database *pgxpool.Pool) RepQuery {
	return &repository{
		database: database,
		logger:   logger,
	}
}
