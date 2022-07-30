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
	"context"
	"fmt"
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

func TestGenericRoleVerifyRules(t *testing.T) {

	tests := []struct {
		denied     bool
		fullmethod string
		rules      []*GenericRule
		roles      []string
	}{
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"*"},
					Apis:     []string{"!enumerate"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"!volumes"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"volumes"},
					Apis:     []string{"*", "!create"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"volumes"},
					Apis:     []string{"!create", "*"},
				},
			},
		},
		{
			// Denials have more priority
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"volumes"},
					Apis:     []string{"!*", "create"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"*", "!*"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"*"},
					Apis:     []string{"*", "!*"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"volumes"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			roles:      []string{"system.admin"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			roles:      []string{"system.admin"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules:      []*GenericRule{},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"futureservice"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"futureservice"},
					Apis:     []string{"anothercall"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"*"},
					Apis:     []string{"anothercall"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*GenericRule{
				&GenericRule{
					Services: []string{"cluster", "volume", "futureservice"},
					Apis:     []string{"somecallinthefuture"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFutureEnumerate",
			roles:      []string{"system.admin"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			roles:      []string{"system.view"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageCluster/InspectCurrent",
			roles:      []string{"system.guest"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageNode/Enumerate",
			roles:      []string{"system.guest"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageNode/Inspect",
			roles:      []string{"system.guest"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageCluster/InspectCurrent",
			roles:      []string{"system.user"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageNode/Enumerate",
			roles:      []string{"system.user"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageNode/Inspect",
			roles:      []string{"system.user"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageNode/FutureCall",
			roles:      []string{"system.user"},
		},
		{
			denied:     false,
			fullmethod: "/my.api.GenericCall/FutureCall",
			roles:      []string{"system.admin"},
		},
	}

	r := GenericRoleManager{}
	for _, test := range tests {
		var err error
		if len(test.roles) != 0 {
			err = r.Verify(context.Background(), test.roles, test.fullmethod)
		} else {
			err = r.VerifyRules(test.rules, "openstorage.api.OpenStorage", test.fullmethod)
		}

		if test.denied {
			assert.NotNil(t, err, test.fullmethod, fmt.Sprintf("%v", test))
		} else {
			assert.Nil(t, err, test.fullmethod)
		}
	}
}
