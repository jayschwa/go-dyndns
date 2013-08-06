package dyndns

import (
	"bufio"
	"net"
	"net/http"
	"strings"
)

const URL = "https://members.dyndns.org/nic/update"

var UserAgent = "Go DynDNS package v0.0.0 jay@jayschwa.net"

var errors = make(map[string]error)

type Error struct {
	Code, Description string
}

func (e *Error) Error() string {
	str := "dyndns: " + e.Code
	if len(e.Description) > 0 {
		str += ": " + e.Description
	}
	return str
}

func NewError(code, description string) error {
	err := &Error{code, description}
	errors[code] = err
	return err
}

var (
	ErrNoChange = NewError("nochg", "no settings changed")

	ErrAuth    = NewError("badauth", "bad username or password")
	ErrDonator = NewError("!donator", "option available only to credited users")

	ErrDomain  = NewError("notfqdn", "hostname is not a fully-qualified domain name")
	ErrNoHost  = NewError("nohost", "hostname does not exist in this account")
	ErrNumHost = NewError("numhost", "too many hosts")
	ErrAbuse   = NewError("abuse", "hostname blocked for update abuse")

	ErrAgent = NewError("badagent", "bad user agent or http method")

	ErrDns = NewError("dnserror", "dns error")
	Err911 = NewError("911", "problem or scheduled maintenance")
)

func Update(username, password, hostname string, ip net.IP) (net.IP, error) {
	url := URL + "?hostname=" + hostname
	if ip != nil {
		url += "&myip=" + ip.String()
		ip = nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := bufio.NewReader(resp.Body)
	code, _ := buf.ReadString(' ')
	code = strings.TrimSpace(code)
	info, _ := buf.ReadString(0)
	if code == "good" || code == "nochg" {
		ip = net.ParseIP(info)
	}
	err = errors[code]
	if err == nil && code != "good" {
		err = &Error{"invalid response", code}
	}
	return ip, err
}
