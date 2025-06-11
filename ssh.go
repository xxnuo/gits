package main

import (
	"strconv"

	"github.com/melbahja/goph"
	"github.com/trzsz/ssh_config"
)

func GetHostConfig(name string) (*goph.Config, error) {
	user := ssh_config.Get(name, "User")
	addr := ssh_config.Get(name, "HostName")

	_port := ssh_config.Get(name, "Port")
	port, err := strconv.ParseUint(_port, 10, 32)
	if err != nil {
		return nil, err
	}

	identityFile, err := NormalizePath(ssh_config.Get(name, "IdentityFile"))
	if err != nil {
		return nil, err
	}
	auth, err := goph.Key(identityFile, "")
	if err != nil {
		return nil, err
	}

	defaultHostkeyCallback, err := goph.DefaultKnownHosts()
	if err != nil {
		return nil, err
	}

	return &goph.Config{
		User:     user,
		Addr:     addr,
		Port:     uint(port),
		Auth:     auth,
		Callback: defaultHostkeyCallback,
	}, nil
}
