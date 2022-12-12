/*
Package ownership manages access to resources
Copyright 2019 Portworx

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
package ownership

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/libopenstorage/grpc-framework/pkg/auth"
)

func TestOwnershipAccessType(t *testing.T) {
	tests := []struct {
		acl       Ownership_AccessType
		request   Ownership_AccessType
		permitted bool
	}{
		{
			acl:       Ownership_READ,
			request:   Ownership_READ,
			permitted: true,
		},
		{
			acl:       Ownership_READ,
			request:   Ownership_WRITE,
			permitted: false,
		},
		{
			acl:       Ownership_READ,
			request:   Ownership_ADMIN,
			permitted: false,
		},
		{
			acl:       Ownership_WRITE,
			request:   Ownership_READ,
			permitted: true,
		},
		{
			acl:       Ownership_WRITE,
			request:   Ownership_WRITE,
			permitted: true,
		},
		{
			acl:       Ownership_WRITE,
			request:   Ownership_ADMIN,
			permitted: false,
		},
		{
			acl:       Ownership_ADMIN,
			request:   Ownership_READ,
			permitted: true,
		},
		{
			acl:       Ownership_ADMIN,
			request:   Ownership_WRITE,
			permitted: true,
		},
		{
			acl:       Ownership_ADMIN,
			request:   Ownership_ADMIN,
			permitted: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t,
			test.acl.isAccessPermitted(test.request),
			test.permitted,
			fmt.Sprintf("acl:%v req:%v p:%v\n", test.acl, test.request, test.permitted))
	}

}

func TestOwnershipIsPermitted(t *testing.T) {
	tests := []struct {
		owner      *Ownership
		user       *auth.UserInfo
		accessType Ownership_AccessType
		permitted  bool
	}{
		{
			// no owner set, so it is a public volume
			owner:     &Ownership{},
			permitted: true,
		},
		{
			// no owner set, so it is a public volume
			owner: &Ownership{
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{},
				},
			},
			permitted: true,
		},
		{
			// no owner set, so it is a public volume
			owner: &Ownership{
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"somegroup": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
			},
			user: &auth.UserInfo{
				Username: "notme",
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
			},
			user: &auth.UserInfo{
				Username: "me",
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
			},
			user: &auth.UserInfo{
				Username: "me",
			},
			accessType: Ownership_WRITE,
			permitted:  true,
		},
		{
			owner: &Ownership{
				Owner: "me",
			},
			user: &auth.UserInfo{
				Username: "me",
			},
			accessType: Ownership_ADMIN,
			permitted:  true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{},
				},
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{},
					Collaborators: map[string]Ownership_AccessType{
						"*": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{},
					Collaborators: map[string]Ownership_AccessType{
						"*": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{},
				},
			},
			accessType: Ownership_WRITE,
			permitted:  false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"*": Ownership_ADMIN,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"*": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{},
				},
			},
			accessType: Ownership_WRITE,
			permitted:  false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"*": Ownership_ADMIN,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group2": Ownership_ADMIN,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"*": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"*": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			accessType: Ownership_ADMIN,
			permitted:  false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
						"notme": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
						"notme": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2"},
				},
			},
			accessType: Ownership_ADMIN,
			permitted:  false,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"group1", "group2", AdminGroup},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group3": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"*"},
				},
			},
			permitted: true,
		},
		{
			owner: &Ownership{
				Owner: "me",
			},
			user: &auth.UserInfo{
				Username: "notme",
				Claims: auth.Claims{
					Groups: []string{"*"},
				},
			},
			permitted: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t,
			test.owner.IsPermitted(test.user, test.accessType),
			test.permitted,
			fmt.Sprintf("Owner:%v\nUser:%v\nPermitted:%v\n", test.owner, test.user, test.permitted))
	}
}

func TestOwnershipUpdate(t *testing.T) {

	tests := []struct {
		owner     *Ownership
		update    *Ownership
		result    *Ownership
		user      *auth.UserInfo
		expectErr bool
	}{
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Owner: "user2",
				Acls: &Ownership_AccessControl{
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "user1",
				Claims: auth.Claims{
					Groups: []string{"group"},
				},
			},
			expectErr: true,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Acls: &Ownership_AccessControl{
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "user1",
				Claims: auth.Claims{
					Groups: []string{"group"},
				},
			},
			expectErr: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"user1": Ownership_READ,
						"user2": Ownership_READ,
						"user3": Ownership_READ,
					},
				},
			},
			user: &auth.UserInfo{
				Username: "user1",
				Claims: auth.Claims{
					Groups: []string{"group"},
				},
			},
			expectErr: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Owner: "user2",
				Acls:  &Ownership_AccessControl{},
			},
			result: &Ownership{
				Owner: "user2",
				Acls:  &Ownership_AccessControl{},
			},
			user: &auth.UserInfo{
				Claims: auth.Claims{
					Groups: []string{"*"},
				},
			},
			expectErr: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Acls: &Ownership_AccessControl{},
			},
			result: &Ownership{
				Owner: "user1",
				Acls:  &Ownership_AccessControl{},
			},
			user: &auth.UserInfo{
				Claims: auth.Claims{
					Groups: []string{"*"},
				},
			},
			expectErr: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Acls: &Ownership_AccessControl{},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			user: &auth.UserInfo{
				Username: "anotheruser",
			},
			expectErr: true,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"anotheruser": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_READ,
					},
					Collaborators: map[string]Ownership_AccessType{
						"anotheruser": Ownership_ADMIN,
					},
				},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_READ,
					},
					Collaborators: map[string]Ownership_AccessType{
						"anotheruser": Ownership_ADMIN,
					},
				}},
			user: &auth.UserInfo{
				Username: "anotheruser",
			},
			expectErr: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_READ,
					},
				},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_READ,
					},
				}},
			user: &auth.UserInfo{
				Username: "anotheruser",
				Claims: auth.Claims{
					Groups: []string{"group1"},
				},
			},
			expectErr: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			update: &Ownership{
				Owner: "anotheruser",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
						"group2": Ownership_READ,
					},
				},
			},
			result: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				}},
			user: &auth.UserInfo{
				Username: "anotheruser",
				Claims: auth.Claims{
					Groups: []string{"group1"},
				},
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		err := test.owner.Update(test.update, test.user)
		if test.expectErr {
			assert.Error(t, err, fmt.Sprintf("%v | %v", test.owner, test.result))
		} else {
			assert.NoError(t, err, fmt.Sprintf("%v | %v", test.owner, test.result))
		}
		assert.True(t, reflect.DeepEqual(test.owner, test.result), fmt.Sprintf("%v | %v", test.owner, test.result))
	}
}

func TestOwnershipIsMatch(t *testing.T) {
	tests := []struct {
		owner   *Ownership
		match   *Ownership
		isMatch bool
	}{
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				},
			},
			match: &Ownership{
				Owner: "user2",
				Acls:  &Ownership_AccessControl{},
			},
			isMatch: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				},
			},
			isMatch: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				},
			},
			match: &Ownership{
				Owner: "user2",
			},
			isMatch: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				},
			},
			match: &Ownership{
				Owner: "user2",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
				},
			},
			isMatch: true,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"one":   Ownership_READ,
						"two":   Ownership_READ,
						"three": Ownership_READ,
					},
				},
			},
			match: &Ownership{
				Owner: "user2",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"one": Ownership_ADMIN,
					},
				},
			},
			isMatch: false,
		},
		{
			owner: &Ownership{
				Owner: "user1",
				Acls: &Ownership_AccessControl{
					Groups: map[string]Ownership_AccessType{
						"group1": Ownership_ADMIN,
					},
					Collaborators: map[string]Ownership_AccessType{
						"one":   Ownership_READ,
						"two":   Ownership_READ,
						"three": Ownership_READ,
					},
				},
			},
			match: &Ownership{
				Owner: "user2",
				Acls: &Ownership_AccessControl{
					Collaborators: map[string]Ownership_AccessType{
						"three": Ownership_READ,
					},
				},
			},
			isMatch: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.owner.IsMatch(test.match), test.isMatch)
	}
}

func TestOwnershipIsPermittedByContext(t *testing.T) {
	o := &Ownership{
		Owner: "user1",
	}
	ctxUser1 := auth.ContextSaveUserInfo(context.Background(), &auth.UserInfo{
		Username: "user1",
	})
	ctxUser2 := auth.ContextSaveUserInfo(context.Background(), &auth.UserInfo{
		Username: "user2",
	})

	assert.True(t, o.IsPermittedByContext(ctxUser1, Ownership_READ))
	assert.False(t, o.IsPermittedByContext(ctxUser2, Ownership_READ))

	var oNil *Ownership
	assert.True(t, oNil.IsPermittedByContext(ctxUser1, Ownership_READ))
	assert.True(t, oNil.IsPermittedByContext(ctxUser2, Ownership_READ))
}

func TestOwnershipIsAdminByContext(t *testing.T) {
	adminctx := auth.ContextSaveUserInfo(context.Background(), &auth.UserInfo{
		Username: "admin",
		Claims: auth.Claims{
			Groups: []string{AdminGroup, "another group"},
		},
	})
	userctx := auth.ContextSaveUserInfo(context.Background(), &auth.UserInfo{
		Username: "user",
	})
	guestctx := auth.ContextSaveUserInfo(context.Background(), auth.NewGuestUser())
	assert.True(t, IsAdminByContext(context.Background()))
	assert.True(t, IsAdminByContext(adminctx))
	assert.False(t, IsAdminByContext(userctx))
	assert.False(t, IsAdminByContext(guestctx))
}
