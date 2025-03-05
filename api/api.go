/*
Copyright (c) 2019 Dell Inc, or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/akutz/gournal"
	"github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
)

const (
	headerKeyContentType                  = "Content-Type"
	headerValContentTypeJSON              = "application/json"
	headerValContentTypeBinaryOctetStream = "binary/octet-stream"
	headerKeyContentLength                = "Content-Length"
	defaultVolumesPath                    = "/ifs/volumes"
	defaultVolumesPathPermissions         = "0777"
	defaultIgnoreUnresolvableHosts        = false
	headerISISessToken                    = "Cookie"
	headerISICSRFToken                    = "X-CSRF-Token" // #nosec, G101  False positive since no hard-coded CSRF token
	headerISIReferer                      = "Referer"
	isiSessCsrfToken                      = "Set-Cookie"
	authTypeBasic                         = 0
	authTypeSessionBased                  = 1
)

var (
	debug, _     = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
	errNewClient = errors.New("missing endpoint, username, or password")
)

// Client is an API client.
type Client interface {
	// Do sends an HTTP request to the OneFS API.
	Do(
		ctx context.Context,
		method, path, id string,
		params OrderedValues,
		body, resp interface{}) error

	// DoWithHeaders sends an HTTP request to the OneFS API.
	DoWithHeaders(
		ctx context.Context,
		method, path, id string,
		params OrderedValues, headers map[string]string,
		body, resp interface{}) error

	// Get sends an HTTP request using the GET method to the OneFS API.
	Get(
		ctx context.Context,
		path, id string,
		params OrderedValues, headers map[string]string,
		resp interface{}) error

	// Post sends an HTTP request using the POST method to the OneFS API.
	Post(
		ctx context.Context,
		path, id string,
		params OrderedValues, headers map[string]string,
		body, resp interface{}) error

	// Put sends an HTTP request using the PUT method to the OneFS API.
	Put(
		ctx context.Context,
		path, id string,
		params OrderedValues, headers map[string]string,
		body, resp interface{}) error

	// Delete sends an HTTP request using the DELETE method to the OneFS API.
	Delete(
		ctx context.Context,
		path, id string,
		params OrderedValues, headers map[string]string,
		resp interface{}) error

	// APIVersion returns the API version.
	APIVersion() uint8

	// User returns the user name used to access the OneFS API.
	User() string

	// Group returns the group name used to access the OneFS API.
	Group() string

	// VolumesPath returns the client's configured volumes path.
	VolumesPath() string

	// VolumePath returns the path to a volume with the provided name.
	VolumePath(name string) string

	// SetAuthToken sets the Auth token/Cookie for the HTTP client
	SetAuthToken(token string)

	// SetCSRFToken sets the Auth token for the HTTP client
	SetCSRFToken(csrf string)

	// SetReferer sets the Referer header
	SetReferer(referer string)

	// GetAuthToken gets the Auth token/Cookie for the HTTP client
	GetAuthToken() string

	// GetCSRFToken gets the CSRF token for the HTTP client
	GetCSRFToken() string

	// GetReferer gets the Referer header
	GetReferer() string
}

type client struct {
	http                    *http.Client
	hostname                string
	username                string
	groupname               string
	password                string
	volumePath              string
	volumePathPermissions   string
	ignoreUnresolvableHosts bool
	apiVersion              uint8
	apiMinorVersion         uint8
	verboseLogging          VerboseType
	sessionCredentials      session
	authType                uint8
}

type session struct {
	sessionCookies string
	sessionCSRF    string
	referer        string
}

type setupConnection struct {
	Services []string `json:"services"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

type VerboseType uint

const (
	VerboseHigh   VerboseType = 0
	VerboseMedium VerboseType = 1
	VerboseLow    VerboseType = 2
)

type apiVerResponse struct {
	Latest *string `json:"latest"`
}

// Error is an API error.
type Error struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// JSONError is a JSON response with one or more errors.
type JSONError struct {
	StatusCode int
	Err        []Error `json:"errors"`
}

// HTMLError is an HTML response with one or more errors.
type HTMLError struct {
	StatusCode int
	Message    string
}

// ClientOptions are options for the API client.
type ClientOptions struct {
	// Insecure is a flag that indicates whether or not to supress SSL errors.
	Insecure bool

	// VolumesPath is the location on the Isilon server where volumes are
	// stored.
	VolumesPath string

	// VolumesPathPermissions is the directory permissions for VolumesPath
	VolumesPathPermissions string

	// IgnoreUnresolvableHosts is the unresolvable hosts param from platform
	IgnoreUnresolvableHosts bool

	// Timeout specifies a time limit for requests made by this client.
	Timeout time.Duration
}

// New returns a new API client.
func New(
	ctx context.Context,
	hostname, username, password, groupname string,
	verboseLogging uint, authType uint8,
	opts *ClientOptions,
) (Client, error) {
	if hostname == "" || username == "" || password == "" {
		return nil, errNewClient
	}

	if authType != authTypeBasic && authType != authTypeSessionBased {
		log.Warn(ctx, "AuthType can be 0 or 1. Setting it to default value 0")
		authType = authTypeBasic
	}

	c := &client{
		hostname:                hostname,
		username:                username,
		groupname:               groupname,
		password:                password,
		volumePath:              defaultVolumesPath,
		volumePathPermissions:   defaultVolumesPathPermissions,
		ignoreUnresolvableHosts: defaultIgnoreUnresolvableHosts,
		verboseLogging:          VerboseType(verboseLogging),
		authType:                authType,
	}

	c.http = &http.Client{}

	if opts != nil {
		if opts.VolumesPath != "" {
			c.volumePath = opts.VolumesPath
		}

		if opts.VolumesPathPermissions != "" {
			c.volumePathPermissions = opts.VolumesPathPermissions
		}

		if opts.IgnoreUnresolvableHosts {
			c.ignoreUnresolvableHosts = opts.IgnoreUnresolvableHosts
		}
		if opts.Timeout != 0 {
			c.http.Timeout = opts.Timeout
		}

		log.Debug(ctx, "opts.Insecure : '%v'", opts.Insecure)

		if opts.Insecure {
			c.http.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // #nosec,G402
					MinVersion:         tls.VersionTLS12,
					MaxVersion:         tls.VersionTLS13,
					CipherSuites:       GetSecuredCipherSuites(),
				},
			}
		} else {
			pool, err := x509.SystemCertPool()
			if err != nil {
				return nil, err
			}
			c.http.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{ //nolint:gosec,G402
					RootCAs:            pool,
					InsecureSkipVerify: false,
					MinVersion:         tls.VersionTLS12,
					MaxVersion:         tls.VersionTLS13,
					CipherSuites:       GetSecuredCipherSuites(),
				},
			}
		}
	}

	if c.authType == authTypeSessionBased {
		_ = c.authenticate(ctx, username, password, hostname)
	}
	resp := &apiVerResponse{}
	if err := c.Get(ctx, "/platform/latest", "", nil, nil, resp); err != nil &&
		!strings.HasPrefix(err.Error(), "json: ") {
		return nil, err
	}

	if resp.Latest != nil {
		s := *resp.Latest
		c.apiMinorVersion = 0
		if i := strings.Index(s, "."); i != -1 {
			ms := s[i+1:]
			m, err := strconv.ParseUint(ms, 10, 8)
			if err != nil {
				return nil, err
			}
			c.apiMinorVersion = uint8(m)
			s = s[:i]
		}
		i, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}
		c.apiVersion = uint8(i)
	} else {
		c.apiVersion = 2
	}

	if c.apiVersion < 3 {
		return nil, errors.New("OneFS releases older than 8.0 are no longer supported")
	}

	return c, nil
}

func (c *client) Get(
	ctx context.Context,
	path, id string,
	params OrderedValues, headers map[string]string,
	resp interface{},
) error {
	return c.executeWithRetryAuthenticate(
		ctx, http.MethodGet, path, id, params, headers, nil, resp)
}

func (c *client) Post(
	ctx context.Context,
	path, id string,
	params OrderedValues, headers map[string]string,
	body, resp interface{},
) error {
	return c.executeWithRetryAuthenticate(
		ctx, http.MethodPost, path, id, params, headers, body, resp)
}

func (c *client) Put(
	ctx context.Context,
	path, id string,
	params OrderedValues, headers map[string]string,
	body, resp interface{},
) error {
	return c.executeWithRetryAuthenticate(
		ctx, http.MethodPut, path, id, params, headers, body, resp)
}

func (c *client) Delete(
	ctx context.Context,
	path, id string,
	params OrderedValues, headers map[string]string,
	resp interface{},
) error {
	return c.executeWithRetryAuthenticate(
		ctx, http.MethodDelete, path, id, params, headers, nil, resp)
}

func (c *client) Do(
	ctx context.Context,
	method, path, id string,
	params OrderedValues,
	body, resp interface{},
) error {
	return c.executeWithRetryAuthenticate(ctx, method, path, id, params, nil, body, resp)
}

func beginsWithSlash(s string) bool {
	return s[0] == '/'
}

func endsWithSlash(s string) bool {
	return s[len(s)-1] == '/'
}

func (c *client) DoWithHeaders(
	ctx context.Context,
	method, uri, id string,
	params OrderedValues, headers map[string]string,
	body, resp interface{},
) error {
	return doWithHeadersFunc(c, ctx, method, uri, id, params, headers, body, resp)
}

var doWithHeadersFunc = func(c *client, ctx context.Context, method string, uri string, id string, params OrderedValues, headers map[string]string, body, resp interface{}) error {
	res, _, err := c.DoAndGetResponseBody(
		ctx, method, uri, id, params, headers, body)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logrus.Printf("Error closing HTTP response: %s", err.Error())
		}
	}()
	logResponse(ctx, res, c.verboseLogging)

	// parse the response
	switch {
	case res == nil:
		return nil
	case res.StatusCode >= 200 && res.StatusCode <= 299:
		if resp == nil {
			return nil
		}
		dec := json.NewDecoder(res.Body)
		if err = dec.Decode(resp); err != nil && err != io.EOF {
			return err
		}
	default:
		return parseJSONHTMLError(res)
	}

	return nil
}

func (c *client) DoAndGetResponseBody(
	ctx context.Context,
	method, uri, id string,
	params OrderedValues, headers map[string]string,
	body interface{},
) (*http.Response, bool, error) {
	return doAndGetResponseBodyFunc(c, ctx, method, uri, id, params, headers, body)
}

var doAndGetResponseBodyFunc = func(
	c *client,
	ctx context.Context,
	method, uri, id string,
	params OrderedValues, headers map[string]string,
	body interface{},
) (*http.Response, bool, error) {
	var (
		err                   error
		req                   *http.Request
		res                   *http.Response
		ubf                   = &bytes.Buffer{}
		lid                   = len(id)
		luri                  = len(uri)
		hostnameEndsWithSlash = endsWithSlash(c.hostname)
		uriBeginsWithSlash    = beginsWithSlash(uri)
		uriEndsWithSlash      = endsWithSlash(uri)
	)

	ubf.WriteString(c.hostname)

	if !hostnameEndsWithSlash && (luri > 0 || lid > 0) {
		ubf.WriteString("/")
	}

	if luri > 0 {
		if uriBeginsWithSlash {
			ubf.WriteString(uri[1:])
		} else {
			ubf.WriteString(uri)
		}
		if !uriEndsWithSlash {
			ubf.WriteString("/")
		}
	}

	if lid > 0 {
		ubf.WriteString(id)
	}

	// add parameters to the URI
	if len(params) > 0 {
		ubf.WriteByte('?')
		if err := params.EncodeTo(ubf); err != nil {
			return nil, false, err
		}
	}

	u, err := url.Parse(ubf.String())
	if err != nil {
		return nil, false, err
	}

	var isContentTypeSet bool

	// marshal the message body (assumes json format)
	if body != nil {
		if r, ok := body.(io.ReadCloser); ok {
			req, err = http.NewRequest(method, u.String(), r)
			defer func() {
				if err := r.Close(); err != nil {
					logrus.Printf("Error closing HTTP response: %s", err.Error())
				}
			}()
			if v, ok := headers[headerKeyContentType]; ok {
				req.Header.Set(headerKeyContentType, v)
			} else {
				req.Header.Set(
					headerKeyContentType, headerValContentTypeBinaryOctetStream)
			}
			isContentTypeSet = true
			// Avoid chunked encoding
			if _, ok := headers[headerKeyContentLength]; ok {
				req.TransferEncoding = []string{"native"}
			}
		} else {
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			if err = enc.Encode(body); err != nil {
				return nil, false, err
			}
			req, err = http.NewRequest(method, u.String(), buf)
			if v, ok := headers[headerKeyContentType]; ok {
				req.Header.Set(headerKeyContentType, v)
			} else {
				req.Header.Set(headerKeyContentType, headerValContentTypeJSON)
			}
			isContentTypeSet = true
		}
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
	}

	if err != nil {
		return nil, false, err
	}

	if !isContentTypeSet {
		isContentTypeSet = req.Header.Get(headerKeyContentType) != ""
	}

	// add headers to the request
	if len(headers) > 0 {
		for header, value := range headers {
			if header == headerKeyContentType && isContentTypeSet {
				continue
			}
			req.Header.Add(header, value)
		}
	}

	if c.authType == authTypeBasic {
		req.SetBasicAuth(c.username, c.password)
	} else {
		if c.GetAuthToken() != "" {
			req.Header.Set(headerISISessToken, c.GetAuthToken())
			req.Header.Set(headerISIReferer, c.GetReferer())
			req.Header.Set(headerISICSRFToken, c.GetCSRFToken())
		}
	}

	var (
		logReqBuf = &bytes.Buffer{}
	)

	if debug {
		log.Info(ctx, "Setting log level to debug in gopowerscale")
		ctx = context.WithValue(
			ctx,
			log.LevelKey(),
			log.DebugLevel)
	}

	logRequest(ctx, logReqBuf, req, c.verboseLogging)
	log.Debug(ctx, logReqBuf.String())

	// send the request
	req = req.WithContext(ctx)
	if res, err = c.http.Do(req); err != nil {
		return nil, debug, err
	}

	return res, debug, err
}

func (c *client) APIVersion() uint8 {
	return c.apiVersion
}

func (c *client) User() string {
	return c.username
}

func (c *client) Group() string {
	return c.groupname
}

func (c *client) VolumesPath() string {
	return c.volumePath
}

func (c *client) VolumePath(volumeName string) string {
	return path.Join(c.volumePath, volumeName)
}

func (err *JSONError) Error() string {
	return err.Err[0].Message
}

func (err *HTMLError) Error() string {
	return err.Message
}

func (c *client) SetAuthToken(cookie string) {
	c.sessionCredentials.sessionCookies = cookie
}

func (c *client) SetCSRFToken(csrf string) {
	c.sessionCredentials.sessionCSRF = csrf
}

func (c *client) SetReferer(referer string) {
	c.sessionCredentials.referer = referer
}

func (c *client) GetAuthToken() string {
	return c.sessionCredentials.sessionCookies
}

func (c *client) GetCSRFToken() string {
	return c.sessionCredentials.sessionCSRF
}

func (c *client) GetReferer() string {
	return c.sessionCredentials.referer
}

func parseJSONHTMLError(r *http.Response) error {
	// check the content type of the response
	if r.Header.Get("Content-Type") == "text/html" {
		htmlError := &HTMLError{}
		// decode HTML error response
		doc, err := goquery.NewDocumentFromReader(r.Body)
		if err != nil {
			return err
		}
		htmlError.StatusCode = r.StatusCode

		if h1 := doc.Find("h1"); h1 != nil {
			htmlError.Message = h1.Text()
		}

		if strings.TrimSpace(htmlError.Message) == "" && doc.Find("title") != nil {
			htmlError.Message = doc.Find("title").Text()
		}
		return htmlError
	}
	jsonErr := &JSONError{}
	// decode JSON error response
	err := json.NewDecoder(r.Body).Decode(jsonErr)
	if err != nil {
		return err
	}
	jsonErr.StatusCode = r.StatusCode
	if len(jsonErr.Err) > 0 && jsonErr.Err[0].Message == "" {
		jsonErr.Err[0].Message = r.Status
	}
	return jsonErr
}

// Authenticate make a REST API call [/session/1/session] to PowerScale to authenticate the given credentials.
// The response contains the session Cookie, X-CSRF-Token and the client uses it for further communication.
func (c *client) authenticate(ctx context.Context, username string, password string, endpoint string) error {
	return authenticateFunc(c, ctx, username, password, endpoint)
}

var authenticateFunc = func(c *client, ctx context.Context, username string, password string, endpoint string) error {
	headers := make(map[string]string, 1)
	headers[headerKeyContentType] = headerValContentTypeJSON
	data := &setupConnection{Services: []string{"platform", "namespace"}, Username: username, Password: password}
	resp, _, err := c.DoAndGetResponseBody(ctx, http.MethodPost, "/session/1/session", "", nil, headers, data)
	if err != nil {
		return fmt.Errorf("Authentication error: %v", err)
	}

	if resp != nil {
		log.Debug(ctx, "Authentication response code: %d", resp.StatusCode)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				logrus.Printf("Error closing HTTP response: %s", err.Error())
			}
		}()

		switch {
		case resp.StatusCode == 201:
			{
				log.Debug(ctx, "Authentication successful")
			}
		case resp.StatusCode == 401:
			{
				log.Debug(ctx, "Response Code %v", resp)
				return fmt.Errorf("authentication failed. unable to login to powerscale. verify username and password")
			}
		default:
			return fmt.Errorf("authenticate error. response-")
		}

		headerRes := strings.Join(resp.Header.Values(isiSessCsrfToken), " ")

		startIndex, endIndex, matchStrLen := FetchValueIndexForKey(headerRes, "isisessid=", ";")
		if startIndex < 0 || endIndex < 0 {
			return fmt.Errorf("Session ID not retrieved")
		}
		c.SetAuthToken(headerRes[startIndex : startIndex+matchStrLen+endIndex])

		startIndex, endIndex, matchStrLen = FetchValueIndexForKey(headerRes, "isicsrf=", ";")
		if startIndex < 0 || endIndex < 0 {
			log.Warn(ctx, "Anti-CSRF Token not retrieved")
		} else {
			c.SetCSRFToken(headerRes[startIndex+matchStrLen : startIndex+matchStrLen+endIndex])
		}

		c.SetReferer(endpoint)
	} else {
		log.Error(ctx, "Authenticate error: Nil response received")
	}
	return nil
}

// executeWithRetryAuthenticate re-authenticates when session credentials become invalid due to time-out or requests exceed.
// it retries the same operation after performing authentication.
func (c *client) executeWithRetryAuthenticate(ctx context.Context, method, uri string, id string, params OrderedValues, headers map[string]string, body, resp interface{}) error {
	err := c.DoWithHeaders(ctx, method, uri, id, params, headers, body, resp)
	if c.authType == authTypeBasic {
		return err
	}
	if err == nil {
		log.Debug(ctx, "Execution successful on Method: %v, URI: %v", method, uri)
		return nil
	}

	switch e := err.(type) {
	case *JSONError:
		if e.StatusCode == 401 {
			log.Debug(ctx, "Authentication failed. Trying to re-authenticate")
			if err := c.authenticate(ctx, c.username, c.password, c.hostname); err != nil {
				return fmt.Errorf("authentication failure due to: %v", err)
			}
			return c.DoWithHeaders(ctx, method, uri, id, params, headers, body, resp)
		}
		log.Error(ctx, "Error in response. Method:%s URI:%s Error: %v JSON Error: %+v", method, uri, err, e)

	case *HTMLError:
		if e.StatusCode == 401 {
			log.Debug(ctx, "Authentication failed. Trying to re-authenticate")
			if err := c.authenticate(ctx, c.username, c.password, c.hostname); err != nil {
				return fmt.Errorf("authentication failure due to: %v", err)
			}
			return c.DoWithHeaders(ctx, method, uri, id, params, headers, body, resp)
		}
		log.Error(ctx, "Error in response. Method:%s URI:%s Error: %v HTML Error: %+v", method, uri, err, e)
	default:
		log.Error(ctx, "Error is not a type of \"*JSONError or *HTMLError\". Error:", err)
	}
	return err
}

func FetchValueIndexForKey(l string, match string, sep string) (int, int, int) {
	startIndex, endIndex := -1, -1
	if strings.Contains(l, match) {
		startIndex = strings.Index(l, match)
		if startIndex != -1 && sep != "" {
			endIndex = strings.Index(l[startIndex+len(match):], sep)
		}
	}
	return startIndex, endIndex, len(match)
}

// GetSecuredCipherSuites returns a set of secure cipher suites.
func GetSecuredCipherSuites() (suites []uint16) {
	securedSuite := tls.CipherSuites()
	for _, v := range securedSuite {
		suites = append(suites, v.ID)
	}
	return suites
}
