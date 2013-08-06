package dyndns

import (
	"testing"
)

const (
	hostname = "test.dyndns.org"
	username = "test"
	password = "test"
)

func TestGood(t *testing.T) {
	ip, err := Update(username, password, hostname, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
}

func TestBadAuth(t *testing.T) {
	ip, err := Update("bogus", password, hostname, nil)
	if err != ErrAuth {
		t.Error(err)
	}
	t.Log(ip)
	ip, err = Update(username, "bogus", hostname, nil)
	if err != ErrAuth {
		t.Error(err)
	}
	t.Log(ip)
}

func TestBadDomain(t *testing.T) {
	ip, err := Update(username, password, "test", nil)
	if err != ErrDomain {
		t.Error(err)
	}
	t.Log(ip)
}

func TestNoHost(t *testing.T) {
	ip, err := Update(username, password, "test.com", nil)
	if err != ErrNoHost {
		t.Error(err)
	}
	t.Log(ip)
}
