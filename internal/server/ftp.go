package server

import (
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	ftp "goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
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
		return nil, fmt.Errorf("FTP port not specified")
	}
	if conf.Username == "" {
		return nil, fmt.Errorf("FTP user name not specified")
	}
	if conf.Password == "" {
		return nil, fmt.Errorf("FTP user password not specified")
	}
	if len(conf.Password) < 20 {
		log.Warn("FTP user password length is less than 20 symbols, this can be security issue")
	}
	ip := net.ParseIP(conf.PublicIP)
	if conf.PublicIP == "" {
		return nil, fmt.Errorf("FTP's public IP not specified")
	} else if ip.IsUnspecified() || ip.IsPrivate() {
		return nil, fmt.Errorf("bad FTP's public IP specification")
	}

	driver, err := file.NewDriver(viper.GetString("data-dir"))
	if err != nil {
		return nil, fmt.Errorf("failed create driver for data directory: %s", err)
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
		RateLimit:    1_000_000,
		PublicIP:     conf.PublicIP,
		PassivePorts: "30000-30020",
	})
	if err != nil {
		return nil, err
	}
	return &FTPServer{server: srv}, nil
}

func (s *FTPServer) Run() error {
	return s.server.ListenAndServe()
}
