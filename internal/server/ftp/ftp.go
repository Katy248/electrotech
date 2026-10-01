package ftp

import (
	"errors"
	"fmt"
	"net"
	"os"

	"charm.land/log/v2"
	"github.com/spf13/viper"
	ftp "goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

const StrongPasswordMinLength = 20
const RateLimit = 1_000_000

var (
	ErrInvalidConfig        = errors.New("invalid FTP config")
	ErrPortNotSpecified     = fmt.Errorf("%w: FTP port not specified", ErrInvalidConfig)
	ErrUsernameNotSpecified = fmt.Errorf("%w: FTP user name not specified", ErrInvalidConfig)
	ErrPasswordNotSpecified = fmt.Errorf("%w: FTP user password not specified", ErrInvalidConfig)
	ErrNoPublicIP           = fmt.Errorf("%w: FTP public IP not specified", ErrInvalidConfig)
	ErrBadPublicIP          = fmt.Errorf("%w: bad FTP public IP", ErrInvalidConfig)
)

type Server struct {
	server *ftp.Server
	logger *log.Logger
}
type Config struct {
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	PublicIP string `mapstructure:"public-ip"`
}

func (conf *Config) validate(logger *log.Logger) error {
	if conf.Port == 0 {
		return ErrPortNotSpecified
	}

	if conf.Username == "" {
		return ErrUsernameNotSpecified
	}

	if conf.Password == "" {
		return ErrPasswordNotSpecified
	}

	if len(conf.Password) < StrongPasswordMinLength {
		logger.Warn("FTP user password length is less than 20 symbols, this can be security issue")
	}

	ip := net.ParseIP(conf.PublicIP)
	if conf.PublicIP == "" {
		return ErrNoPublicIP
	} else if ip.IsUnspecified() || ip.IsPrivate() {
		return ErrBadPublicIP
	}

	return nil
}

func NewFTPServer(conf *Config, logger *log.Logger) (*Server, error) {
	if err := conf.validate(logger); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	driver, err := file.NewDriver(viper.GetString("data-dir"))
	if err != nil {
		return nil, fmt.Errorf("failed create driver for data directory: %w", err)
	}

	usr := os.Getenv("USER")

	//nolint:exhaustruct_v5
	srv, err := ftp.NewServer(&ftp.Options{
		Driver: driver,
		Port:   conf.Port,
		Auth: &ftp.SimpleAuth{
			Name:     conf.Username,
			Password: conf.Password,
		},
		Perm:         ftp.NewSimplePerm(usr, usr),
		RateLimit:    RateLimit,
		PublicIP:     conf.PublicIP,
		PassivePorts: "30000-30020",
	})
	if err != nil {
		return nil, fmt.Errorf("create FTP server: %w", err)
	}

	return &Server{server: srv, logger: logger}, nil
}

func (s *Server) Run() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("run FTP server: %w", err)
	}

	return nil
}
