package health

import (
	"context"
	"fmt"
	"time"

	"charm.land/log/v2"
	"gorm.io/gorm"
)

const CheckTimeout = time.Second * 3

type Status = string

const (
	StatusOK          = Status("ok")
	StatusUnavailable = Status("unavailable")
)

type Health struct {
	Dependencies map[string]Status
}

func (h *Health) Status() Status {
	for _, status := range h.Dependencies {
		if status == StatusUnavailable {
			return StatusUnavailable
		}
	}

	return StatusOK
}

type checker interface {
	Check(ctx context.Context) error
}

type Service struct {
	logger *log.Logger

	checkers map[string]checker
}

func (s *Service) GetHealth(ctx context.Context) Health {
	ctx, cancel := context.WithTimeout(ctx, CheckTimeout)
	defer cancel()

	health := Health{
		map[string]Status{},
	}

	for name, check := range s.checkers {
		err := check.Check(ctx)
		if err != nil {
			health.Dependencies[name] = StatusUnavailable
			s.logger.Warn("Dependency unavailable", "dependency", name, "error", err)
		} else {
			health.Dependencies[name] = StatusOK
		}
	}

	return health
}

type Option func(*Service)

func New(options ...Option) *Service {
	service := &Service{
		logger:   log.Default(),
		checkers: map[string]checker{},
	}
	for _, opt := range options {
		opt(service)
	}

	return service
}

func WithLogger(logger *log.Logger) Option {
	return func(s *Service) {
		if logger == nil {
			return
		}

		s.logger = logger
	}
}

func WithDB(db *gorm.DB) Option {
	return func(s *Service) {
		s.checkers["db"] = &dbChecker{db: db}
	}
}

type dbChecker struct {
	db *gorm.DB
}

func (c *dbChecker) Check(ctx context.Context) error {
	sql, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("get sql connection: %w", err)
	}

	err = sql.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}
