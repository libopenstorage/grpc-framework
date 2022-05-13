/*
Generic role manager
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

// For a full Role Manager example see github.com/libopenstorage/openstorage/pkg/role/sdkserviceapi.go

package role

import (
	"context"
	"fmt"

	grpcutil "github.com/libopenstorage/grpc-framework/pkg/grpc/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SystemAdminRoleName = "system.admin"
	SystemGuestRoleName = "system.guest"
)

type GenericRule struct {
	Services []string
	Apis     []string
}

type GenericRoleManager struct{}

type DefaultRole struct {
	Rules []*GenericRule
}

var (
	// DefaultRoles are the default roles to load on system startup
	// Should be prefixed by `system.` to avoid collisions
	DefaultRoles = map[string]*DefaultRole{
		// system:admin role can run any command
		SystemAdminRoleName: &DefaultRole{
			Rules: []*GenericRule{
				&GenericRule{
					Services: []string{"*"},
					Apis:     []string{"*"},
				},
			},
		},

		// system:guest role is used for any unauthenticated user.
		// They can only use standard volume lifecycle commands.
		SystemGuestRoleName: &DefaultRole{
			Rules: []*GenericRule{
				&GenericRule{
					Services: []string{"!*"},
					Apis:     []string{"!*"},
				},
			},
		},
	}
)

// Verify determines if the role has access to `fullmethod`
func (r *GenericRoleManager) Verify(ctx context.Context, roles []string, fullmethod string) error {

	// Check all roles
	for _, role := range roles {
		if defaultRole, ok := DefaultRoles[role]; ok {
			if err := r.VerifyRules(defaultRole.Rules,
				"", /* this would be the default root path of the APIs */
				fullmethod); err == nil {
				return nil
			}
		}
	}

	return status.Errorf(codes.PermissionDenied, "Access denied to roles: %+s", roles)
}

// VerifyRules checks if the rules authorize use of the API called `fullmethod`
func (r *GenericRoleManager) VerifyRules(rules []*GenericRule, rootPath, fullmethod string) error {

	reqService, reqApi := grpcutil.GetMethodInformation(rootPath, fullmethod)

	// Look for denials first
	for _, rule := range rules {
		for _, service := range rule.Services {
			// if the service is denied, then return here
			if DenyRule(service, reqService) {
				return fmt.Errorf("access denied to service by role")
			}

			// If there is a match to the service now check the apis
			if MatchRule(service, reqService) {
				for _, api := range rule.Apis {
					if DenyRule(api, reqApi) {
						return fmt.Errorf("access denied to api by role")
					}
				}
			}
		}
	}

	// Look for permissions
	for _, rule := range rules {
		for _, service := range rule.Services {
			if MatchRule(service, reqService) {
				for _, api := range rule.Apis {
					if MatchRule(api, reqApi) {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("no accessible rule to authorize access found")
}
