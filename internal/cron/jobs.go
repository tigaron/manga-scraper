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
	cronjobs, err := s.repo.FindAll(context.Background())
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to find all cron jobs")
	}

	for i := range cronjobs {
		if cronjobs[i].Tags == "skip" {
			return nil
		}

		switch cronjobs[i].Name {
		case "scrape-series-list":
			err = s.createNewJob(scheduler, cronjobs[i].Crontab, cronjobs[i].Name, s.scrapeSeriesList, cronjobs[i].ID)
		case "scrape-series-detail":
			err = s.createNewJob(scheduler, cronjobs[i].Crontab, cronjobs[i].Name, s.scrapeSeriesDetail, cronjobs[i].ID)
		case "scrape-chapters-list":
			err = s.createNewJob(scheduler, cronjobs[i].Crontab, cronjobs[i].Name, s.scrapeChaptersList, cronjobs[i].ID)
		case "scrape-chapters-detail":
			err = s.createNewJob(scheduler, cronjobs[i].Crontab, cronjobs[i].Name, s.scrapeChaptersDetail, cronjobs[i].ID)
		}
	}

	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create job")
	}

	return nil
}

func (s *Cron) createNewJob(scheduler gocron.Scheduler, crontab, name string, jobFunc func(), prevJobID string) error {
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
					Message:  jobName,
					Duration: 0,
				})
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				s.logger.Info("Job completed", zap.String("jobID", jobID.String()), zap.String("jobName", jobName), zap.Duration("duration", time.Duration(s.cronMonitor.time[jobName][s.cronMonitor.counter[jobName]])))
				_, _ = s.repo.CreateStatus(context.Background(), internal.CreateCronJobStatusParams{
					JobID:    jobID.String(),
					Status:   "completed",
					Message:  jobName,
					Duration: s.cronMonitor.time[jobName][s.cronMonitor.counter[jobName]].Seconds(),
				})
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				s.logger.Error("Job failed", zap.String("jobID", jobID.String()), zap.String("jobName", jobName), zap.Error(err))
				_, _ = s.repo.CreateStatus(context.Background(), internal.CreateCronJobStatusParams{
					JobID:    jobID.String(),
					Status:   "failed",
					Message:  jobName,
					Duration: s.cronMonitor.time[jobName][s.cronMonitor.counter[jobName]].Seconds(),
				})
			}),
		),
	)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create job")
	}

	_, err = s.repo.Create(context.Background(), internal.CreateCronJobParams{
		ID:      job.ID().String(),
		Name:    job.Name(),
		Crontab: crontab,
		Tags:    "",
	})
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "Failed to create new job")
	}

	return s.repo.Delete(context.Background(), prevJobID)
}

// TODO:
// - Implement cron job to update chapter count 2x a day
// - Implement cron job to update latest chapter 2x a day
// - Implement cron job to update series index 2x a day
