package driverhttp

import (
	"time"

	"github.com/pivotal-golang/clock"
)

const (
	DefaultBackoffTimeout  = 30 * time.Second
	backoffInitialInterval = 500 * time.Millisecond
	backoffIncrement       = 1.5
)

type Operation func(logger lager.Logger) error

type BackOff interface {
	Retry(func() error) error
}

type exponentialBackOff struct {
	maxElapsedTime time.Duration
	clock          clock.Clock
}

// newExponentialBackOff takes a maximum elapsed time, after which the
// exponentialBackOff stops retrying the operation.
func NewExponentialBackOff(maxElapsedTime time.Duration, clock clock.Clock) BackOff {
	return &exponentialBackOff{
		maxElapsedTime: maxElapsedTime,
		clock:          clock,
	}
}

// Retry takes a retriable operation, and calls it until either the operation
// succeeds, or the retry timeout occurs.
func (b *exponentialBackOff) Retry(logger lager.Logger, operation Operation) error {
	logger = logger.Session("retry")
	logger.Info("start")
	logger.Info("end")

	startTime := b.clock.Now()
	backoffInterval := time.Duration(0)
	backoffExpired := false

	for {
		err := operation(logger)
		if err == nil {
			return nil
		}

		backoffInterval, backoffExpired = b.incrementInterval(startTime, backoffInterval)
		if backoffExpired {
			return err
		}

		b.clock.Sleep(backoffInterval)
	}
}

func (b *exponentialBackOff) incrementInterval(startTime time.Time, currentInterval time.Duration) (nextInterval time.Duration, expired bool) {
	elapsedTime := b.clock.Now().Sub(startTime)

	if elapsedTime > b.maxElapsedTime {
		return 0, true
	}

	switch {
	case currentInterval == 0:
		nextInterval = backoffInitialInterval
	case elapsedTime+backoff(currentInterval) > b.maxElapsedTime:
		nextInterval = time.Millisecond + b.maxElapsedTime - elapsedTime
	default:
		nextInterval = backoff(currentInterval)
	}

	return nextInterval, false
}

func backoff(interval time.Duration) time.Duration {
	return time.Duration(float64(interval) * backoffIncrement)
}
