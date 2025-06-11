package main

import (
	"net"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	if hostFound && err != nil {
		return err
	}

	if hostFound && err == nil {
		return nil
	}

	return goph.AddKnownHost(host, remote, key, "")
}
