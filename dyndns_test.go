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
	ip, err := Service{DynDNS, username, password}.Update(hostname, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
}

func TestBadAuth(t *testing.T) {
	ip, err := Service{DynDNS, "bogus", password}.Update(hostname, nil)
	if err != ErrAuth {
		t.Error(err)
	}
	t.Log(ip)
	ip, err = Service{DynDNS, username, "bogus"}.Update(hostname, nil)
	if err != ErrAuth {
		t.Error(err)
	}
	t.Log(ip)
}

func TestBadDomain(t *testing.T) {
	ip, err := Service{DynDNS, username, password}.Update("bogus", nil)
	if err != ErrDomain {
		t.Error(err)
	}
	t.Log(ip)
}

func TestNoHost(t *testing.T) {
	ip, err := Service{DynDNS, username, password}.Update("bogus.com", nil)
	if err != ErrNoHost {
		t.Error(err)
	}
	t.Log(ip)
}
