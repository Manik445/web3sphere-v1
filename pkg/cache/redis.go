package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
)

// RedisClient wraps the go-redis client with convenience methods.
type RedisClient struct {
	Client *redis.Client
	log    *logger.Logger
}

// New creates a new Redis client.
func New(cfg *configs.RedisConfig, log *logger.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Redis connected successfully")
	return &RedisClient{Client: client, log: log}, nil
}

// Close shuts down the Redis connection.
func (r *RedisClient) Close() {
	if err := r.Client.Close(); err != nil {
		r.log.Errorf("Failed to close Redis: %v", err)
	} else {
		r.log.Info("Redis connection closed")
	}
}

// --- Cache Abstraction ---

// Set stores a value in the cache with a TTL.
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}
	return r.Client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from the cache.
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key from the cache.
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists.
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.Client.Exists(ctx, key).Result()
	return n > 0, err
}

// --- Distributed Lock ---

// AcquireLock acquires a distributed lock.
func (r *RedisClient) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
}

// ReleaseLock releases a distributed lock.
func (r *RedisClient) ReleaseLock(ctx context.Context, key string) error {
	return r.Client.Del(ctx, "lock:"+key).Err()
}

// --- OTP Storage ---

// StoreOTP saves a one-time password.
func (r *RedisClient) StoreOTP(ctx context.Context, identifier string, otp string, ttl time.Duration) error {
	return r.Client.Set(ctx, "otp:"+identifier, otp, ttl).Err()
}

// GetOTP retrieves an OTP.
func (r *RedisClient) GetOTP(ctx context.Context, identifier string) (string, error) {
	return r.Client.Get(ctx, "otp:"+identifier).Result()
}

// DeleteOTP removes an OTP.
func (r *RedisClient) DeleteOTP(ctx context.Context, identifier string) error {
	return r.Client.Del(ctx, "otp:"+identifier).Err()
}

// IncrementOTPAttempts increments the OTP attempt counter.
func (r *RedisClient) IncrementOTPAttempts(ctx context.Context, identifier string, maxAttempts int, window time.Duration) (int64, error) {
	key := "otp_attempts:" + identifier
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.Client.Expire(ctx, key, window)
	}
	return count, nil
}

// --- JWT Blacklist ---

// BlacklistToken adds a token to the blacklist.
func (r *RedisClient) BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	return r.Client.Set(ctx, "blacklist:"+tokenID, "1", ttl).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted.
func (r *RedisClient) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	return r.Exists(ctx, "blacklist:"+tokenID)
}

// --- Session Storage ---

// StoreSession saves a session.
func (r *RedisClient) StoreSession(ctx context.Context, userID string, sessionID string, data interface{}, ttl time.Duration) error {
	return r.Set(ctx, "session:"+userID+":"+sessionID, data, ttl)
}

// GetSession retrieves a session.
func (r *RedisClient) GetSession(ctx context.Context, userID string, sessionID string, dest interface{}) error {
	return r.Get(ctx, "session:"+userID+":"+sessionID, dest)
}

// DeleteSession removes a session.
func (r *RedisClient) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	return r.Delete(ctx, "session:"+userID+":"+sessionID)
}

// DeleteAllSessions removes all sessions for a user.
func (r *RedisClient) DeleteAllSessions(ctx context.Context, userID string) error {
	pattern := "session:" + userID + ":*"
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.Client.Del(ctx, keys...).Err()
	}
	return nil
}

// --- Rate Limiter ---

// CheckRateLimit checks if a request is within the rate limit.
func (r *RedisClient) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, int64, error) {
	rateKey := "rate:" + key
	count, err := r.Client.Incr(ctx, rateKey).Result()
	if err != nil {
		return false, 0, err
	}
	if count == 1 {
		r.Client.Expire(ctx, rateKey, window)
	}
	remaining := int64(limit) - count
	if remaining < 0 {
		remaining = 0
	}
	return count <= int64(limit), remaining, nil
}

// --- Feature Flags ---

// SetFeatureFlag sets a feature flag.
func (r *RedisClient) SetFeatureFlag(ctx context.Context, flag string, enabled bool) error {
	val := "0"
	if enabled {
		val = "1"
	}
	return r.Client.HSet(ctx, "feature_flags", flag, val).Err()
}

// GetFeatureFlag checks if a feature flag is enabled.
func (r *RedisClient) GetFeatureFlag(ctx context.Context, flag string) (bool, error) {
	val, err := r.Client.HGet(ctx, "feature_flags", flag).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
