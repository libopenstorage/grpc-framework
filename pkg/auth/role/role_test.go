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

func TestGenericRoleVerifyRules(t *testing.T) {

	tests := []struct {
		denied     bool
		fullmethod string
		rules      []*Rule
		roles      []string
	}{
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"!enumerate"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			rules: []*Rule{
				&Rule{
					Services: []string{"!openstorage.api.OpenStorageVolumes"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"openstorage.api.OpenStorageVolumes"},
					Apis:     []string{"*", "!create"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"openstorage.api.OpenStorageVolumes"},
					Apis:     []string{"!create", "*"},
				},
			},
		},
		{
			// Denials have more priority
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"openstorage.api.OpenStorageVolumes"},
					Apis:     []string{"!*", "create"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"*", "!*"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"*", "!*"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"openstorage.api.OpenStorageVolumes"},
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
			rules:      []*Rule{},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{"openstorage.api.OpenStorageFutureService"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{"openstorage.api.OpenStorageFutureService"},
					Apis:     []string{"anothercall"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"anothercall"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{
						"openstorage.api.OpenStorageCluster",
						"openstorage.api.OpenStorageVolume",
						"openstorage.api.openStoragefutureservice",
					},
					Apis: []string{"somecallinthefuture"},
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

	r := NewDefaultGenericRoleManager()
	for i, test := range tests {
		var err error
		if i == 14 {
			if len(test.roles) != 0 {
				err = r.Verify(context.Background(), test.roles, test.fullmethod)
			} else {
				err = r.VerifyRules(test.rules, "", test.fullmethod)
			}

			if test.denied {
				assert.NotNil(t, err, fmt.Sprintf("%v:%d", test, i))
			} else {
				assert.Nil(t, err, fmt.Sprintf("%v:%d", test, i))
			}

		}
		if len(test.roles) != 0 {
			err = r.Verify(context.Background(), test.roles, test.fullmethod)
		} else {
			err = r.VerifyRules(test.rules, "", test.fullmethod)
		}

		if test.denied {
			assert.NotNil(t, err, fmt.Sprintf("%v:%d", test, i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("%v:%d", test, i))
		}
	}
}

func TestGenericRoleVerifyRulesWithTag(t *testing.T) {

	tests := []struct {
		denied     bool
		fullmethod string
		rules      []*Rule
		roles      []string
	}{
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"!enumerate"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			rules: []*Rule{
				&Rule{
					Services: []string{"!volumes"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"volumes"},
					Apis:     []string{"*", "!create"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"volumes"},
					Apis:     []string{"!create", "*"},
				},
			},
		},
		{
			// Denials have more priority
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"volumes"},
					Apis:     []string{"!*", "create"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"*", "!*"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"*", "!*"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			rules: []*Rule{
				&Rule{
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
			rules:      []*Rule{},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{"futureservice"},
					Apis:     []string{"*"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{"futureservice"},
					Apis:     []string{"anothercall"},
				},
			},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"anothercall"},
				},
			},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			rules: []*Rule{
				&Rule{
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

	r := NewGenericRoleManager("openstorage.api.OpenStorage", DefaultRoles)
	for _, test := range tests {
		var err error
		if len(test.roles) != 0 {
			err = r.Verify(context.Background(), test.roles, test.fullmethod)
		} else {
			err = r.VerifyRules(test.rules, r.tag, test.fullmethod)
		}

		if test.denied {
			assert.NotNil(t, err, test.fullmethod, fmt.Sprintf("%v", test))
		} else {
			assert.Nil(t, err, test.fullmethod)
		}
	}
}

func TestGenericRoleVerifyRulesWithTagAndCustomRoles(t *testing.T) {

	customRoles := map[string]*Role{
		"admin": &Role{
			Name: "admin",
			Rules: []*Rule{
				&Rule{
					Services: []string{"*"},
					Apis:     []string{"*"},
				},
			},
		},

		"user": &Role{
			Name: "user",
			Rules: []*Rule{
				&Rule{
					Services: []string{"futureservice"},
					Apis:     []string{"somecallinthefuture", "anothercall"},
				},
				&Rule{
					Services: []string{"volumes"},
					Apis:     []string{"*"},
				},
			},
		},
	}

	tests := []struct {
		denied     bool
		fullmethod string
		roles      []string
	}{
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Enumerate",
			roles:      []string{"qa"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			roles:      []string{"user"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			roles:      []string{"admin", "user"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageVolumes/Create",
			roles:      []string{"qa", "user"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			roles:      []string{"admin"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			roles:      []string{"user"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			roles:      []string{"qa", "user"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/SomeCallInTheFuture",
			roles:      []string{"qa"},
		},
		{
			denied:     false,
			fullmethod: "/openstorage.api.OpenStorageFutureService/AnotherCall",
			roles:      []string{"qa", "user"},
		},
		{
			denied:     true,
			fullmethod: "/openstorage.api.OpenStorageFutureService/AdminOnly",
			roles:      []string{"qa", "user"},
		},
	}

	r := NewGenericRoleManager("openstorage.api.OpenStorage", customRoles)
	for _, test := range tests {
		err := r.Verify(context.Background(), test.roles, test.fullmethod)

		if test.denied {
			assert.NotNil(t, err, test.fullmethod, fmt.Sprintf("%v", test))
		} else {
			assert.Nil(t, err, test.fullmethod)
		}
	}
}
