package server

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
	ErrPortNotSpecified     = errors.New("FTP port not specified")
	ErrUsernameNotSpecified = errors.New("FTP user name not specified")
	ErrPasswordNotSpecified = errors.New("FTP user password not specified")
	ErrNoPublicIP           = errors.New("FTP public IP not specified")
	ErrBadPublicIP          = errors.New("bad FTP public IP specification")
)

type FTPServer struct {
	server *ftp.Server
}

func NewFTPServer() (*FTPServer, error) {
	var conf struct {
		Port     int
		Username string
		Password string
		PublicIP string
	}

	conf.Port = viper.GetInt("ftp.port")
	conf.Username = viper.GetString("ftp.username")
	conf.Password = viper.GetString("ftp.password")
	conf.PublicIP = viper.GetString("ftp.public-ip")

	if conf.Port == 0 {
		return nil, ErrPortNotSpecified
	}

	if conf.Username == "" {
		return nil, ErrUsernameNotSpecified
	}

	if conf.Password == "" {
		return nil, ErrPasswordNotSpecified
	}

	if len(conf.Password) < StrongPasswordMinLength {
		log.Warn("FTP user password length is less than 20 symbols, this can be security issue")
	}

	ip := net.ParseIP(conf.PublicIP)
	if conf.PublicIP == "" {
		return nil, ErrNoPublicIP
	} else if ip.IsUnspecified() || ip.IsPrivate() {
		return nil, ErrBadPublicIP
	}

	driver, err := file.NewDriver(viper.GetString("data-dir"))
	if err != nil {
		return nil, fmt.Errorf("failed create driver for data directory: %w", err)
	}

	usr := os.Getenv("USER")

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

	return &FTPServer{server: srv}, nil
}

func (s *FTPServer) Run() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("run FTP server: %w", err)
	}

	return nil
}
