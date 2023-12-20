package profile

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/alimy/freecar/idle/auto/rpc/base"
	"github.com/alimy/freecar/idle/auto/rpc/profile"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/alimy/freecar/library/core/id"
	"github.com/bytedance/sonic"
)

// Manager defines a profile manager.
type Manager struct {
	ProfileService profileservice.Client
}

// Verify verifies account identity.
func (m *Manager) Verify(c context.Context, aid id.AccountID) (id.IdentityID, error) {
	nilID := id.IdentityID("")
	resp, err := m.ProfileService.GetProfile(c, &profile.GetProfileRequest{AccountId: aid.String()})
	if err != nil {
		return nilID, fmt.Errorf("cannot get profile: %v", err)
	}
	if resp.Profile.IdentityStatus != base.IdentityStatus_VERIFIED {
		return nilID, fmt.Errorf("invalid identity status")
	}

	b, err := sonic.Marshal(resp.Profile.Identity)
	if err != nil {
		return nilID, fmt.Errorf("cannot marshal identity:%v", err)
	}
	return id.IdentityID(base64.StdEncoding.EncodeToString(b)), nil
}
