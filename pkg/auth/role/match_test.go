/*
Package role manages roles in Kvdb and provides validation
Copyright 2022 Pure Storage

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
package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchDenyRule(t *testing.T) {

	tests := []struct {
		matchFound bool
		role       string
		s          string
	}{
		{
			matchFound: false,
			role:       "",
			s:          "",
		},
		{
			matchFound: true,
			role:       "!*",
			s:          "test",
		},
		{
			matchFound: true,
			role:       "!test",
			s:          "test",
		},
		{
			matchFound: true,
			role:       "!!!!!!!!!!!!!*******************test****",
			s:          "test",
		},
		{
			matchFound: false,
			role:       "test",
			s:          "test",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.matchFound, DenyRule(test.role, test.s))
	}
}

func TestMatchRule(t *testing.T) {

	tests := []struct {
		matchFound bool
		role       string
		s          string
	}{
		{
			matchFound: false,
			role:       "",
			s:          "",
		},
		{
			matchFound: true,
			role:       "*",
			s:          "test",
		},
		{
			matchFound: true,
			role:       "***********",
			s:          "test",
		},
		{
			matchFound: false,
			role:       "nomatch",
			s:          "test",
		},
		{
			matchFound: false,
			role:       "*nomatch",
			s:          "test",
		},
		{
			matchFound: false,
			role:       "nomatch*",
			s:          "test",
		},
		{
			matchFound: false,
			role:       "*nomatch*",
			s:          "test",
		},
		{
			matchFound: true,
			role:       "*test",
			s:          "thisisatest",
		},
		{
			matchFound: true,
			role:       "this*",
			s:          "thisisatest",
		},
		{
			matchFound: true,
			role:       "*isa*",
			s:          "thisisatest",
		},
		{
			matchFound: false,
			role:       "isa",
			s:          "thisisatest",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.matchFound, MatchRule(test.role, test.s))
	}
}
