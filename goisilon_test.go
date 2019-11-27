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
package goisilon

import (
	"context"
	"flag"
	"os"
	"testing"

	log "github.com/akutz/gournal"
	glogrus "github.com/akutz/gournal/logrus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	err        error
	client     *Client
	defaultCtx context.Context
)

func init() {
	defaultCtx = context.Background()
	defaultCtx = context.WithValue(
		defaultCtx,
		log.AppenderKey(),
		glogrus.NewWithOptions(
			logrus.StandardLogger().Out,
			logrus.DebugLevel,
			logrus.StandardLogger().Formatter))
}

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		defaultCtx = context.WithValue(
			defaultCtx,
			log.LevelKey(),
			log.DebugLevel)
	}

	client, err = NewClientWithArgs(
		defaultCtx,
		"https://10.230.24.244:8080",
		true,
		false,
		"admin",
		"",
		"Password123!",
		"/ifs/data/csi_share_1")

	if err != nil {
		log.WithError(err).Panic(defaultCtx, "error creating test client")
	}

	if err != nil {
		log.WithError(err).Panic(defaultCtx, "error creating test client")
	}
	os.Exit(m.Run())
}

func assertLen(t *testing.T, obj interface{}, expLen int) {
	if !assert.Len(t, obj, expLen) {
		t.FailNow()
	}
}

func assertError(t *testing.T, err error) {
	if !assert.Error(t, err) {
		t.FailNow()
	}
}

func assertNoError(t *testing.T, err error) {
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}

func assertNil(t *testing.T, i interface{}) {
	if !assert.Nil(t, i) {
		t.FailNow()
	}
}

func assertNotNil(t *testing.T, i interface{}) {
	if !assert.NotNil(t, i) {
		t.FailNow()
	}
}
