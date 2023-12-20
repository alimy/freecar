package pay

import (
	"context"

	"github.com/alimy/freecar/idle/auto/rpc/user"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/tools"

	"github.com/alimy/freecar/idle/auto/rpc/user/userservice"
	"github.com/alimy/freecar/library/core/id"
)

// Manager defines a car manager.
type Manager struct {
	UserClient userservice.Client
}

// NewManager creates a new pay manager.
func NewManager(c userservice.Client) *Manager {
	return &Manager{
		UserClient: c,
	}
}

// Pay pays for the trip.
func (m *Manager) Pay(ctx context.Context, aid id.AccountID, feeCent int32) error {
	resp, err := m.UserClient.Pay(ctx, &user.PayRequest{
		AccountId: aid.String(),
		FeeCent:   feeCent,
	})
	if err != nil {
		return errno.RPCUserSrvErr
	}
	return tools.ParseBaseResp(resp.BaseResp)
}
