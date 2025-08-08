package watchdog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	messagebus "github.com/ferretcode/switchyard/scheduler/internal/message_bus"
	"github.com/ferretcode/switchyard/scheduler/pkg/types"
	"github.com/redis/go-redis/v9"
)

type WatchdogService struct {
	Logger            *slog.Logger
	Config            *types.Config
	RedisConn         *redis.Client
	MessageBusService *messagebus.MessageBusService
	Context           context.Context
}

func NewWatchdogService(logger *slog.Logger, config *types.Config, redisConn *redis.Client, messageBusService *messagebus.MessageBusService, context context.Context) WatchdogService {
	return WatchdogService{
		Logger:            logger,
		Config:            config,
		RedisConn:         redisConn,
		MessageBusService: messageBusService,
		Context:           context,
	}
}

func (w *WatchdogService) WatchStuckJobs() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		w.checkStuckJobs()
	}
}

func (w *WatchdogService) checkStuckJobs() {
	now := time.Now().Unix()

	w.Logger.Info("executing stuck jobs watchdog")

	cutoff := now - int64(w.Config.WorkerStuckJobThreshold.Seconds())

	jobIds, err := w.RedisConn.ZRangeByScore(w.Context, "jobs:pending", &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%d", cutoff),
	}).Result()
	if err != nil {
		w.Logger.Error("error fetching stuck jobs", "err", err)
		return
	}

	for _, jobId := range jobIds {
		jobKey := fmt.Sprintf("jobs:%s", jobId)

		status, err := w.RedisConn.HGet(w.Context, jobKey, "status").Result()
		if err != nil {
			w.Logger.Error("error reading job status", "err", err)
			continue
		}

		if status == "ok" || status == "error" {
			w.RedisConn.ZRem(w.Context, "jobs:pending", jobId)
			continue
		}

		w.Logger.Info("job has not been processed yet, retrying", "job-id", jobId)

		retryCountStr, err := w.RedisConn.HGet(w.Context, jobKey, "retry_count").Result()
		retryCount := 0
		if err == nil {
			fmt.Sscanf(retryCountStr, "%d", &retryCount)
		}

		if retryCount >= w.Config.WorkerMaxJobRetries {
			w.RedisConn.HSet(w.Context, jobKey, "status", "error", "updated_at", now)
			w.RedisConn.ZRem(w.Context, "jobs:pending", jobId)
			w.Logger.Error("job marked as failed after max retries", "job-id", jobId)
			continue
		}

		jobName, err := w.RedisConn.HGet(w.Context, jobKey, "job_name").Result()
		if err != nil {
			w.Logger.Error("error reading job name", "err", err)
			continue
		}

		jobContextString, err := w.RedisConn.HGet(w.Context, jobKey, "job_context").Result()
		if err != nil {
			w.Logger.Error("error reading job context", "err", err)
			continue
		}

		jobContext := make(map[string]any)

		if err := json.Unmarshal([]byte(jobContextString), &jobContext); err != nil {
			w.Logger.Error("error decoding job context", "err", err)
			continue
		}

		err = w.MessageBusService.SendRetryJobMessage(jobName, jobContext, jobId)
		if err != nil {
			w.Logger.Error("error re-scheduling job", "err", err)
			continue
		}

		err = w.RedisConn.HIncrBy(w.Context, jobKey, "retry_count", 1).Err()
		if err != nil {
			w.Logger.Error("error updating retry count", "err", err)
			continue
		}

		err = w.RedisConn.HSet(w.Context, jobKey, "updated_at", now).Err()
		if err != nil {
			w.Logger.Error("error updating retry updated field", "err", err)
			continue
		}
	}
}
