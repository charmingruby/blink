package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(addr string) (*RedisLock, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisLock{client: client}, nil
}

func (r *RedisLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	result, err := r.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return result, nil
}

func (r *RedisLock) Release(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	return nil
}

func (r *RedisLock) CheckIdempotency(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check idempotency: %w", err)
	}

	return result > 0, nil
}

func (r *RedisLock) MarkIdempotent(ctx context.Context, key string, ttl time.Duration) error {
	if err := r.client.Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf("failed to mark idempotent: %w", err)
	}

	return nil
}

func (r *RedisLock) IncrementRetry(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment retry: %w", err)
	}

	if count == 1 {
		if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
			return count, fmt.Errorf("failed to set TTL: %w", err)
		}
	}

	return count, nil
}

func (r *RedisLock) Close() error {
	return r.client.Close()
}
