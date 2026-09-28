/*
Copyright (c) 2026 Dell Inc. or its subsidiaries. All Rights Reserved.

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

package apiv27

import "github.com/stretchr/testify/mock"

// anyArgs matches any set of arguments for the api.Client methods used in this
// package. The api.Client Do/Get/Post/Put methods take seven parameters, so
// seven mock.Anything entries are required.
var anyArgs = []interface{}{
	mock.Anything,
	mock.Anything,
	mock.Anything,
	mock.Anything,
	mock.Anything,
	mock.Anything,
	mock.Anything,
}

// stringPtr returns a pointer to the given string. It exists because the
// OneFS API models use *string for optional fields, and Go does not allow
// taking the address of a string literal directly.
func stringPtr(s string) *string {
	return &s
}
