package prisma

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
)

type CronJobRepo struct {
	q *PrismaClient
}

func NewCronJobRepo(prismaClient *PrismaClient) *CronJobRepo {
	return &CronJobRepo{
		q: prismaClient,
	}
}

func (c *CronJobModel) toCronJob() internal.CronJob {
	return internal.CronJob{
		ID:      c.ID,
		Name:    c.Name,
		Crontab: c.Crontab,
		Tags:    c.Tags,
	}
}

func (c *CronJobStatusModel) toCronJobStatus() internal.CronJobStatus {
	return internal.CronJobStatus{
		ID:       c.ID,
		JobID:    c.JobID,
		Status:   c.Status,
		Message:  c.Message,
		Duration: c.Duration,
	}
}

func (c *CronJobRepo) Create(ctx context.Context, params internal.CreateCronJobParams) (internal.CronJob, error) {
	defer newSentrySpan(ctx, "CronJobRepo.Create").Finish()

	cronJob, err := c.q.CronJob.CreateOne(
		CronJob.ID.Set(params.ID),
		CronJob.Name.Set(params.Name),
		CronJob.Crontab.Set(params.Crontab),
		CronJob.Tags.Set(params.Tags),
	).Exec(ctx)
	if err != nil {
		return internal.CronJob{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to create cron job")
	}

	return cronJob.toCronJob(), nil
}

func (c *CronJobRepo) Upsert(ctx context.Context, params internal.CreateCronJobParams) (internal.CronJob, error) {
	defer newSentrySpan(ctx, "CronJobRepo.Upsert").Finish()

	cronJob, err := c.q.CronJob.UpsertOne(
		CronJob.ID.Equals(params.ID),
	).Create(
		CronJob.ID.Set(params.ID),
		CronJob.Name.Set(params.Name),
		CronJob.Crontab.Set(params.Crontab),
		CronJob.Tags.Set(params.Tags),
	).Update(
		CronJob.Name.Set(params.Name),
		CronJob.Crontab.Set(params.Crontab),
		CronJob.Tags.Set(params.Tags),
	).Exec(ctx)
	if err != nil {
		return internal.CronJob{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to upsert cron job")
	}

	return cronJob.toCronJob(), nil
}

func (c *CronJobRepo) Find(ctx context.Context, id string) (internal.CronJob, error) {
	defer newSentrySpan(ctx, "CronJobRepo.Find").Finish()

	cronJob, err := c.q.CronJob.FindUnique(
		CronJob.ID.Equals(id),
	).Exec(ctx)
	if err != nil {
		return internal.CronJob{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find cron job")
	}

	return cronJob.toCronJob(), nil
}

func (c *CronJobRepo) FindAll(ctx context.Context) ([]internal.CronJob, error) {
	defer newSentrySpan(ctx, "CronJobRepo.FindAll").Finish()

	cronJobs, err := c.q.CronJob.FindMany().Exec(ctx)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find all cron jobs")
	}

	var res []internal.CronJob
	for i := range cronJobs {
		res = append(res, cronJobs[i].toCronJob())
	}

	return res, nil
}

func (c *CronJobRepo) CreateStatus(ctx context.Context, params internal.CreateCronJobStatusParams) (internal.CronJobStatus, error) {
	defer newSentrySpan(ctx, "CronJobRepo.CreateStatus").Finish()

	cronJobStatus, err := c.q.CronJobStatus.CreateOne(
		CronJobStatus.Status.Set(params.Status),
		CronJobStatus.Message.Set(params.Message),
		CronJobStatus.Duration.Set(params.Duration),
		CronJobStatus.CronJob.Link(
			CronJob.ID.Equals(params.JobID),
		),
	).Exec(ctx)
	if err != nil {
		return internal.CronJobStatus{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to create cron job status")
	}

	return cronJobStatus.toCronJobStatus(), nil
}

func (c *CronJobRepo) UpdateStatus(ctx context.Context, params internal.UpdateCronJobStatusParams) (internal.CronJobStatus, error) {
	defer newSentrySpan(ctx, "CronJobRepo.UpdateStatus").Finish()

	cronJobStatus, err := c.q.CronJobStatus.FindUnique(
		CronJobStatus.ID.Equals(params.ID),
	).Update(
		CronJobStatus.Status.Set(params.Status),
		CronJobStatus.Message.Set(params.Message),
		CronJobStatus.Duration.Set(params.Duration),
	).Exec(ctx)
	if err != nil {
		return internal.CronJobStatus{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to update cron job status")
	}

	return cronJobStatus.toCronJobStatus(), nil
}

func (c *CronJobRepo) Delete(ctx context.Context, id string) error {
	defer newSentrySpan(ctx, "CronJobRepo.Delete").Finish()

	_, err := c.q.CronJob.FindUnique(
		CronJob.ID.Equals(id),
	).Delete().Exec(ctx)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "failed to delete cron job")
	}

	return nil
}
