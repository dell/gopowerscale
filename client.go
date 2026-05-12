/*
Copyright (c) 2019-2022 Dell Inc, or its subsidiaries.

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
package gopowerscale

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dell/gopowerscale/api"
)

// Client is an Isilon client.
type Client struct {
	// API is the underlying OneFS API client.
	API api.Client
}

// NewClient returns a new Isilon client struct initialized from the environment.
func NewClient(ctx context.Context) (*Client, error) {
	insecure, err := strconv.ParseBool(os.Getenv("GOISILON_INSECURE"))
	if err != nil {
		return nil, err
	}
	ignoreUnresolvableHosts, err := strconv.ParseBool(os.Getenv("GOISILON_UNRESOLVABLE_HOSTS"))
	if err != nil {
		return nil, err
	}
	authType, err := strconv.Atoi(os.Getenv("GOISILON_AUTHTYPE"))
	if err != nil {
		return nil, err
	}
	// #nosec G115
	return NewClientWithArgs(
		ctx,
		os.Getenv("GOISILON_ENDPOINT"),
		insecure,
		1,
		os.Getenv("GOISILON_USERNAME"),
		os.Getenv("GOISILON_GROUP"),
		os.Getenv("GOISILON_PASSWORD"),
		os.Getenv("GOISILON_VOLUMEPATH"),
		os.Getenv("GOISILON_VOLUMEPATH_PERMISSIONS"),
		ignoreUnresolvableHosts,
		uint8(authType),
	)
}

// NewClientWithCAFilePath returns a new Isilon client struct initialized from the environment with CA file.
func NewClientWithCAFilePath(ctx context.Context) (*Client, error) {
	insecure, err := strconv.ParseBool(os.Getenv("GOISILON_INSECURE"))
	if err != nil {
		return nil, err
	}
	ignoreUnresolvableHosts, err := strconv.ParseBool(os.Getenv("GOISILON_UNRESOLVABLE_HOSTS"))
	if err != nil {
		return nil, err
	}
	authType, err := strconv.Atoi(os.Getenv("GOISILON_AUTHTYPE"))
	if err != nil {
		return nil, err
	}
	// #nosec G115
	return NewClientWithArgsWithCAFilePath(
		ctx,
		os.Getenv("GOISILON_ENDPOINT"),
		insecure,
		1,
		os.Getenv("GOISILON_USERNAME"),
		os.Getenv("GOISILON_GROUP"),
		os.Getenv("GOISILON_PASSWORD"),
		os.Getenv("GOISILON_VOLUMEPATH"),
		os.Getenv("GOISILON_VOLUMEPATH_PERMISSIONS"),
		ignoreUnresolvableHosts,
		uint8(authType),
		os.Getenv("GOISILON_CAFILEPATH"),
	)
}

// NewClientWithArgs returns a new Isilon client struct initialized from the supplied arguments.
func NewClientWithArgs(
	ctx context.Context,
	endpoint string,
	insecure bool, verboseLogging uint,
	user, group, pass, volumesPath string, volumesPathPermissions string, ignoreUnresolvableHosts bool, authType uint8,
) (*Client, error) {
	timeout, _ := time.ParseDuration(os.Getenv("GOISILON_TIMEOUT"))

	client, err := api.New(
		ctx, endpoint, user, pass, group, verboseLogging, authType,
		&api.ClientOptions{
			Insecure:                insecure,
			VolumesPath:             volumesPath,
			VolumesPathPermissions:  volumesPathPermissions,
			IgnoreUnresolvableHosts: ignoreUnresolvableHosts,
			Timeout:                 timeout,
		})
	if err != nil {
		return nil, err
	}

	return &Client{client}, err
}

// NewClientWithArgsWithCAFilePath returns a new Isilon client struct initialized from the supplied arguments.
func NewClientWithArgsWithCAFilePath(
	ctx context.Context,
	endpoint string,
	insecure bool, verboseLogging uint,
	user, group, pass, volumesPath string, volumesPathPermissions string, ignoreUnresolvableHosts bool, authType uint8,
	caFilePath string,
) (*Client, error) {
	timeout, _ := time.ParseDuration(os.Getenv("GOISILON_TIMEOUT"))

	client, err := api.New(
		ctx, endpoint, user, pass, group, verboseLogging, authType,
		&api.ClientOptions{
			Insecure:                insecure,
			VolumesPath:             volumesPath,
			VolumesPathPermissions:  volumesPathPermissions,
			IgnoreUnresolvableHosts: ignoreUnresolvableHosts,
			CAFilePath:              caFilePath,
			Timeout:                 timeout,
		})
	if err != nil {
		return nil, err
	}

	return &Client{client}, err
}

// SetCustomHTTPHeaders sets custom HTTP headers that will be sent with every request
func (c *Client) SetCustomHTTPHeaders(headers http.Header) {
	if c.API != nil {
		c.API.SetCustomHTTPHeaders(headers)
	}
}

// SetAuthToken sets the authentication token for the client
func (c *Client) SetAuthToken(token string) {
	if c.API != nil {
		c.API.SetAuthToken(token)
	}
}

// SetCSRFToken sets the CSRF token for the client
func (c *Client) SetCSRFToken(csrf string) {
	if c.API != nil {
		c.API.SetCSRFToken(csrf)
	}
}

// SetReferer sets the referer for the client
func (c *Client) SetReferer(referer string) {
	if c.API != nil {
		c.API.SetReferer(referer)
	}
}

// SetAuthorizationMode sets the global flag to skip authentication for authorization mode
func SetAuthorizationMode(skipAuth bool) {
	api.SetSkipAuthForAuthorization(skipAuth)
}
