// Package dyndns updates dynamic DNS hostnames.
package dyndns

import (
	"bufio"
	"net"
	"net/http"
	"strings"
)

// URL specifies where to send update requests.
var URL = "https://members.dyndns.org/nic/update"

// UserAgent identifies the client in update requests.
var UserAgent = "go-dyndns/0.0 (github.com/jayschwa/go-dyndns)"

// errors maps return code text to an error.
var errors = make(map[string]error)

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

// Update protocol errors.
type Error struct {
	Code, Description string
}

// NewError returns a new Error from a return code and description.
func NewError(code, description string) *Error {
	err := &Error{code, description}
	errors[code] = err
	return err
}

// Error satisfies the built-in error interface.
func (e *Error) Error() string {
	str := "dyndns: " + e.Code
	if len(e.Description) > 0 {
		str += ": " + e.Description
	}
	return str
}

// Update protocol return codes.
//
// http://dyn.com/support/developers/api/return-codes/
var (
	NoChange = NewError("nochg", "no settings changed")

	// Account errors.
	ErrAuth    = NewError("badauth", "bad username or password")
	ErrDonator = NewError("!donator", "option available only to credited users")

	// Hostname errors.
	ErrDomain  = NewError("notfqdn", "hostname is not a fully-qualified domain name")
	ErrNoHost  = NewError("nohost", "hostname does not exist in this account")
	ErrNumHost = NewError("numhost", "too many hosts")
	ErrAbuse   = NewError("abuse", "hostname blocked for update abuse")

	// User agent errors.
	ErrAgent = NewError("badagent", "bad user agent or http method")

	// Server errors.
	ErrDns = NewError("dnserror", "dns error")
	Err911 = NewError("911", "server problem or scheduled maintenance")
)
