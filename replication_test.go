package goisilon_test

/*
Copyright (c) 2021-2023 Dell Inc, or its subsidiaries.

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

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dell/goisilon"
	"github.com/dell/goisilon/api"
	"github.com/stretchr/testify/suite"
)

type ReplicationTestSuite struct {
	suite.Suite
	localClient    *goisilon.Client
	remoteClient   *goisilon.Client
	localEndpoint  string
	remoteEndpoint string
}

func (suite *ReplicationTestSuite) SetupSuite() {
	lc, err := goisilon.NewClientWithArgs(
		context.Background(),
		"https://10.225.111.21:8080",
		true,
		1,
		"admin",
		"",
		"dangerous",
		"/ifs/data/test-goisilon",
		"0777", false, 0)
	if err != nil {
		panic(err)
	}
	suite.localClient = lc
	suite.localEndpoint = "10.225.111.21"

	rc, err := goisilon.NewClientWithArgs(
		context.Background(),
		"https://10.225.111.70:8080",
		true,
		1,
		"admin",
		"",
		"dangerous",
		"/ifs/data/test-goisilon",
		"0777", false, 0)
	if err != nil {
		panic(err)
	}
	suite.remoteClient = rc
	suite.remoteEndpoint = "10.225.111.70"
}

func (suite *ReplicationTestSuite) TearDownSuite() {
}

func (suite *ReplicationTestSuite) TestUnplannedFailoverScenario() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	volumeName := "replicated"

	// *** SIMULATE CREATE_VOLUME CALL *** //

	// Create volume that would serve as VG
	volume, err := suite.localClient.CreateVolume(ctx, volumeName)
	suite.NoError(err)
	suite.NotNil(volume)

	// defer func() {
	// 	err := suite.localClient.DeleteVolume(ctx, volumeName)
	// 	suite.NoError(err)
	// }()

	res, err := suite.localClient.GetVolume(context.Background(), "", volumeName)
	suite.NoError(err)
	fmt.Println("local", res)

	err = suite.localClient.CreatePolicy(ctx,
		volumeName,
		300,
		"/ifs/data/test-goisilon/replicated",
		"/ifs/data/test-goisilon/replicated",
		suite.remoteEndpoint,
		"",
		true)
	// suite.NoError(err)

	p, err := suite.localClient.GetPolicyByName(ctx, volumeName)
	suite.NoError(err)
	suite.NotNil(p)
	fmt.Println("local policy", p)

	err = suite.localClient.WaitForPolicyLastJobState(ctx, volumeName, goisilon.FINISHED)
	suite.NoError(err)

	// defer func() {
	// 	err := suite.localClient.DeletePolicy(ctx, volumeName)
	// 	suite.NoError(err)
	// }()

	// *** SIMULATE EXECUTE_ACTION UNPLANNED_FAILOVER CALL ***

	err = suite.remoteClient.BreakAssociation(ctx, volumeName)
	suite.NoError(err)

	// *** SIMULATE EXECUTE_ACTION REPROTECT CALL ***
	// In driver EXECUTE_ACTION reprotect will be called on another side, but here we just talk to remote client

	local := suite.remoteClient
	remote := suite.localClient

	pp, err := remote.GetPolicyByName(ctx, volumeName)
	suite.NoError(err)

	if pp.Enabled {
		// Disable policy on remote
		err = remote.DisablePolicy(ctx, volumeName)
		suite.NoError(err)

		err = remote.WaitForPolicyEnabledFieldCondition(ctx, volumeName, false)
		suite.NoError(err)

		// Run reset on the policy
		err = remote.ResetPolicy(ctx, volumeName)
		suite.NoError(err)

		// Create policy on local (actually get it before creating it)
		err = local.CreatePolicy(ctx,
			volumeName,
			300,
			"/ifs/data/test-goisilon/replicated",
			"/ifs/data/test-goisilon/replicated",
			suite.localEndpoint,
			"",
			true)
		suite.NoError(err)

		err = local.WaitForPolicyEnabledFieldCondition(ctx, volumeName, true)
		suite.NoError(err)

	} else {
		err = local.EnablePolicy(ctx, volumeName)
		suite.NoError(err)

		err = local.WaitForPolicyEnabledFieldCondition(ctx, volumeName, true)
		suite.NoError(err)
	}

	tp, err := remote.GetTargetPolicyByName(ctx, volumeName)
	if err != nil {
		if e, ok := err.(*api.JSONError); ok {
			if e.StatusCode != 404 {
				suite.NoError(err)
			}
		}
	}

	if tp != nil {
		err = remote.DisallowWrites(ctx, volumeName)
		suite.NoError(err)
	}
}

func (suite *ReplicationTestSuite) TestReplication() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	volumeName := "replicated"

	// *** SIMULATE CREATE_VOLUME CALL *** //

	// Create volume that would serve as VG
	volume, err := suite.localClient.CreateVolume(ctx, volumeName)
	suite.NoError(err)
	suite.NotNil(volume)

	// defer func() {
	// 	err := suite.localClient.DeleteVolume(ctx, volumeName)
	// 	suite.NoError(err)
	// }()

	res, err := suite.localClient.GetVolume(context.Background(), "", volumeName)
	suite.NoError(err)
	fmt.Println("local", res)

	err = suite.localClient.CreatePolicy(ctx,
		volumeName,
		300,
		"/ifs/data/test-goisilon/replicated",
		"/ifs/data/test-goisilon/replicated",
		suite.remoteEndpoint,
		"",
		true)
	suite.NoError(err)

	p, err := suite.localClient.GetPolicyByName(ctx, volumeName)
	suite.NoError(err)
	suite.NotNil(p)
	fmt.Println("local policy", p)

	err = suite.localClient.WaitForPolicyLastJobState(ctx, volumeName, goisilon.FINISHED)
	suite.NoError(err)

	// defer func() {
	// 	err := suite.localClient.DeletePolicy(ctx, volumeName)
	// 	suite.NoError(err)
	// }()

	// *** SIMULATE EXECUTE_ACTION FAILOVER CALL ***

	err = suite.localClient.SyncPolicy(ctx, volumeName)
	suite.NoError(err)

	// Create Remote Policy
	err = suite.remoteClient.CreatePolicy(ctx,
		volumeName,
		300,
		"/ifs/data/test-goisilon/replicated",
		"/ifs/data/test-goisilon/replicated",
		suite.remoteEndpoint,
		"",
		false)
	suite.NoError(err)

	rp, err := suite.remoteClient.GetPolicyByName(ctx, volumeName)
	suite.NoError(err)
	suite.NotNil(rp)
	suite.Equal(rp.Enabled, false)
	fmt.Println("remote policy", rp)

	err = suite.remoteClient.WaitForPolicyLastJobState(ctx, volumeName, goisilon.UNKNOWN)
	suite.NoError(err)

	// defer func() {
	// 	err := suite.remoteClient.DeletePolicy(ctx, volumeName)
	// 	suite.NoError(err)
	// }()

	// Allow writes on remote
	err = suite.remoteClient.AllowWrites(ctx, volumeName)
	suite.NoError(err)

	// Disable policy on local
	err = suite.localClient.DisablePolicy(ctx, volumeName)
	suite.NoError(err)

	err = suite.localClient.WaitForPolicyEnabledFieldCondition(ctx, volumeName, false)
	suite.NoError(err)

	// Disable writes on local (if we can)
	tp, err := suite.localClient.GetTargetPolicyByName(ctx, volumeName)
	if err != nil {
		if e, ok := err.(*api.JSONError); ok {
			if e.StatusCode != 404 {
				suite.NoError(err)
			}
		}
	}

	fmt.Println("local target policy", tp)

	if tp != nil {
		err = suite.localClient.DisallowWrites(ctx, volumeName)
		suite.NoError(err)
	}

	// *** SIMULATE EXECUTE_ACTION REPROTECT CALL ***
	// In driver EXECUTE_ACTION reprotect will be called on another side, but here we just talk to remote client
	err = suite.remoteClient.EnablePolicy(ctx, volumeName)
	suite.NoError(err)

	err = suite.remoteClient.WaitForPolicyEnabledFieldCondition(ctx, volumeName, true)
	suite.NoError(err)
}

func TestReplicationSuite(_ *testing.T) {
	//	suite.Run(t, new(ReplicationTestSuite))
}
