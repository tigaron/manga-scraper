package cron

import (
	"context"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Cron) scheduleJobs(scheduler gocron.Scheduler) error {
	err := s.createNewJob(scheduler, "0 0 * * *", "scrape-series-list", s.scrapeSeriesList)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create scrape-series-list job")
	}

	err = s.createNewJob(scheduler, "0 0 * * *", "scrape-series-detail", s.scrapeSeriesDetail)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create scrape-series-detail job")
	}

	err = s.createNewJob(scheduler, "0 0,12 * * *", "scrape-chapters-list", s.scrapeChaptersList)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create scrape-chapters-list job")
	}

	err = s.createNewJob(scheduler, "0 0,12 * * *", "scrape-chapters-detail", s.scrapeChaptersDetail)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create scrape-chapters-detail job")
	}

	return nil
}

func (s *Cron) createNewJob(scheduler gocron.Scheduler, crontab, name string, jobFunc func()) error {
	job, err := scheduler.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(jobFunc),
		gocron.WithName(name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				s.logger.Info("Job started", zap.String("jobID", jobID.String()), zap.String("jobName", jobName))
				_, _ = s.repo.CreateStatus(context.Background(), internal.CreateCronJobStatusParams{
					JobID:    jobID.String(),
					Status:   "started",
					Message:  "",
					Duration: 0,
				})
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				s.logger.Info("Job completed", zap.String("jobID", jobID.String()), zap.String("jobName", jobName), zap.Duration("duration", time.Duration(s.cronMonitor.time[jobName][s.cronMonitor.counter[jobName]])))
				_, _ = s.repo.CreateStatus(context.Background(), internal.CreateCronJobStatusParams{
					JobID:    jobID.String(),
					Status:   "completed",
					Message:  "",
					Duration: s.cronMonitor.time[jobName][s.cronMonitor.counter[jobName]].Seconds(),
				})
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				s.logger.Error("Job failed", zap.String("jobID", jobID.String()), zap.String("jobName", jobName), zap.Error(err))
				_, _ = s.repo.CreateStatus(context.Background(), internal.CreateCronJobStatusParams{
					JobID:    jobID.String(),
					Status:   "failed",
					Message:  "",
					Duration: s.cronMonitor.time[jobName][s.cronMonitor.counter[jobName]].Seconds(),
				})
			}),
		),
	)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create job")
	}

	_, err = s.repo.Upsert(context.Background(), internal.CreateCronJobParams{
		ID:      job.ID().String(),
		Name:    job.Name(),
		Crontab: crontab,
		Tags:    "",
	})
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to upsert job")
	}

	return nil
}

// TODO:
// - Implement cron job to update chapter count 2x a day
// - Implement cron job to update latest chapter 2x a day
// - Implement cron job to update series index 2x a day
