// Package dyndns updates dynamic DNS hostnames. It can be used with services
// that support the DNS Update API:
//
// http://dyn.com/support/developers/api/
package dyndns

import (
	"bufio"
	"net"
	"net/http"
	"strings"
)

// UserAgent identifies the client in update requests.
var UserAgent = "go-dyndns/0.0 (github.com/jayschwa/go-dyndns)"

// A Service represents a dynamic DNS service and its account credentials.
type Service struct {
	URL, Username, Password string
}

// Update sends a request to the service to change the hostname to ip.
// If ip is nil, the update server will use the client's IP address.
// It returns the updated IP address on success and an error, if any.
func (s Service) Update(hostname string, ip net.IP) (net.IP, error) {

	// Prepare HTTP request.
	url := s.URL + "?hostname=" + hostname
	if ip != nil {
		url += "&myip=" + ip.String()
		ip = nil // ip is reused for output.
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.Username, s.Password)
	req.Header.Add("User-Agent", UserAgent)

	// Execute the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response.
	buf := bufio.NewReader(resp.Body)
	code, _ := buf.ReadString(' ')
	code = strings.TrimSpace(code)
	info, _ := buf.ReadString(0)
	if code == "good" || code == NoChange.Code {
		ip = net.ParseIP(info)
	}
	err = errors[code]
	if err == nil && code != "good" {
		err = &Error{"invalid response code", code}
	}
	return ip, err
}

// errors maps return code text to an error.
var errors = make(map[string]error)

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

// Update protocol response codes.
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
