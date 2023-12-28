package servants

import (
	"context"
	"time"

	"github.com/alimy/freecar/app/user/infras/mysql"
	"github.com/alimy/freecar/idle/auto/rpc/base"
	"github.com/alimy/freecar/idle/auto/rpc/blob"
	"github.com/alimy/freecar/idle/auto/rpc/user"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/paseto"
)

var (
	_ user.UserService = (*userSrv)(nil)
)

// userSrv implements the last service interface defined in the IDL.
type userSrv struct {
	OpenIDResolver
	EncryptManager
	AdminMysqlManager
	UserMysqlManager
	BlobManager
	TokenGenerator
}

// TokenGenerator creates token.
type TokenGenerator interface {
	CreateToken(claims *paseto.StandardClaims) (token string, err error)
}

// OpenIDResolver resolves an authorization code
// to an open id.
type OpenIDResolver interface {
	Resolve(code string) string
}

type EncryptManager interface {
	EncryptPassword(code string) string
}

type UserMysqlManager interface {
	CreateUser(user *mysql.User) (*mysql.User, error)
	GetUserByOpenId(openId string) (*mysql.User, error)
	GetUserByAccountId(aid string) (*mysql.User, error)
	GetSomeUsers() ([]*mysql.User, error)
	GetAllUsers() ([]*mysql.User, error)
	UpdateUser(user *mysql.User) error
	DeleteUser(aid string) error
}

type AdminMysqlManager interface {
	GetAdminByAccountId(aid string) (*mysql.Admin, error)
	GetAdminByName(name string) (*mysql.Admin, error)
	UpdateAdminPassword(aid string, password string) error
}

// BlobManager defines the Anti Corruption Layer
// for get blob logic.
type BlobManager interface {
	GetBlobURL(ctx context.Context, req *blob.GetBlobURLRequest, callOptions ...callopt.Option) (*blob.GetBlobURLResponse, error)
	CreateBlob(ctx context.Context, req *blob.CreateBlobRequest, callOptions ...callopt.Option) (*blob.CreateBlobResponse, error)
}

// Login implements the userSrv interface.
func (s *userSrv) Login(_ context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp = new(user.LoginResponse)
	// Resolve code to openID.
	openID := s.OpenIDResolver.Resolve(req.Code)
	if openID == "" {
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("bad open id"))
		return resp, nil
	}

	usr, err := s.UserMysqlManager.GetUserByOpenId(openID)
	if err != nil {
		if err != errno.RecordNotFound {
			klog.Error("get user by open id err", err)
			resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr)
			return resp, nil
		}
		usr, err = s.UserMysqlManager.CreateUser(&mysql.User{OpenID: openID})
		if err != nil {
			klog.Error("create user err", err)
			resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr)
			return resp, nil
		}
	}

	now := time.Now()
	resp.Token, err = s.TokenGenerator.CreateToken(&paseto.StandardClaims{
		ID:        usr.ID,
		Issuer:    consts.Issuer,
		Audience:  consts.User,
		IssuedAt:  now,
		NotBefore: now,
		ExpiredAt: now.Add(consts.ThirtyDays),
	})
	if err != nil {
		klog.Error("create token error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}

	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// AdminLogin implements the userSrv interface.
func (s *userSrv) AdminLogin(_ context.Context, req *user.AdminLoginRequest) (resp *user.AdminLoginResponse, err error) {
	resp = new(user.AdminLoginResponse)
	admin, err := s.AdminMysqlManager.GetAdminByName(req.Username)
	if err != nil {
		klog.Error("get password by name err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("login error"))
		return resp, nil
	}
	cryPassword := s.EncryptPassword(req.Password)
	if admin.Password != cryPassword {
		klog.Infof("%s login err", req.Username)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("wrong username or password"))
		return resp, nil
	}

	now := time.Now()
	resp.Token, err = s.TokenGenerator.CreateToken(&paseto.StandardClaims{
		ID:        admin.ID,
		Issuer:    consts.Issuer,
		Audience:  consts.Admin,
		IssuedAt:  now,
		NotBefore: now,
		ExpiredAt: now.Add(consts.ThirtyDays),
	})
	if err != nil {
		klog.Error("create token error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}

	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// ChangeAdminPassword implements the userSrv interface.
func (s *userSrv) ChangeAdminPassword(_ context.Context, req *user.ChangeAdminPasswordRequest) (resp *user.ChangeAdminPasswordResponse, err error) {
	resp = new(user.ChangeAdminPasswordResponse)
	admin, err := s.AdminMysqlManager.GetAdminByAccountId(req.AccountId)
	if err != nil {
		klog.Error("get password by aid err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("change password error"))
		return resp, nil
	}
	cryPassword := s.EncryptManager.EncryptPassword(req.OldPassword)
	if admin.Password != cryPassword {
		klog.Infof("%s change password err", admin.Username)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("wrong password"))
		return resp, nil
	}
	err = s.UpdateAdminPassword(req.AccountId, req.NewPassword_)
	if err != nil {
		klog.Error("update password err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("change password error"))
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// UploadAvatar implements the userSrv interface.
func (s *userSrv) UploadAvatar(ctx context.Context, req *user.UploadAvatarRequset) (resp *user.UploadAvatarResponse, err error) {
	resp = new(user.UploadAvatarResponse)
	aid := req.AccountId
	br, err := s.BlobManager.CreateBlob(ctx, &blob.CreateBlobRequest{
		AccountId:           aid,
		UploadUrlTimeoutSec: int32(10 * time.Second.Seconds()),
	})
	if err != nil {
		klog.Error("cannot create blob", err)
		resp.BaseResp = utils.BuildBaseResp(errno.BlobSrvErr)
		return resp, nil
	}

	if err = s.UserMysqlManager.UpdateUser(&mysql.User{
		ID:           req.AccountId,
		AvatarBlobId: br.Id,
	}); err != nil {
		if err == errno.RecordNotFound {
			return nil, errno.RecordNotFound
		}
		klog.Error("update user blob id error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("upload avatar error\""))
		return resp, nil
	}

	resp.BaseResp = utils.BuildBaseResp(nil)
	resp.UploadUrl = br.UploadUrl
	return resp, nil
}

// UpdateUser implements the userSrv interface.
func (s *userSrv) UpdateUser(_ context.Context, req *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	resp = new(user.UpdateUserResponse)
	err = s.UserMysqlManager.UpdateUser(&mysql.User{
		ID:          req.AccountId,
		PhoneNumber: req.PhoneNumber,
		Username:    req.Username,
	})
	if err != nil {
		if err == errno.RecordNotFound {
			resp.BaseResp = utils.BuildBaseResp(errno.RecordNotFound)
			return resp, nil
		}
		klog.Error("update user error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr)
		return resp, nil
	}

	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// GetUser implements the userSrv interface.
func (s *userSrv) GetUser(ctx context.Context, req *user.GetUserRequest) (resp *user.GetUserInfoResponse, err error) {
	resp = new(user.GetUserInfoResponse)
	u, err := s.UserMysqlManager.GetUserByAccountId(req.AccontId)
	if err != nil {
		if err == errno.RecordNotFound {
			resp.BaseResp = utils.BuildBaseResp(errno.RecordNotFound)
			return resp, nil
		}
		klog.Error("get user by accountId err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr)
		return resp, nil
	}
	resp.UserInfo = &base.UserInfo{
		AccountId:   u.ID,
		Username:    u.Username,
		PhoneNumber: u.PhoneNumber,
		AvatarUrl:   "",
		Balance:     u.Balance,
	}
	if u.AvatarBlobId != "" {
		res, err := s.BlobManager.GetBlobURL(ctx, &blob.GetBlobURLRequest{
			Id:         u.AvatarBlobId,
			TimeoutSec: int32(5 * time.Second.Seconds()),
		})
		if err != nil {
			klog.Error("get blob url err", err)
			resp.BaseResp = utils.BuildBaseResp(errno.BlobSrvErr)
			return resp, nil
		}
		resp.UserInfo.AvatarUrl = res.Url
	}

	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// AddUser implements the userSrv interface.
func (s *userSrv) AddUser(ctx context.Context, req *user.AddUserRequest) (resp *user.AddUserResponse, err error) {
	resp = new(user.AddUserResponse)
	_, err = s.UserMysqlManager.CreateUser(&mysql.User{
		ID:           req.AccountId,
		PhoneNumber:  req.PhoneNumber,
		AvatarBlobId: req.AvatarBlobId,
		Username:     req.Username,
		OpenID:       req.OpenId,
	})
	if err != nil {
		if err == errno.RecordAlreadyExist {
			klog.Error("add user error", err)
			resp.BaseResp = utils.BuildBaseResp(errno.RecordAlreadyExist)
			return resp, nil
		}
		klog.Error("add user error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr)
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// DeleteUser implements the userSrv interface.
func (s *userSrv) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (resp *user.DeleteUserResponse, err error) {
	resp = new(user.DeleteUserResponse)
	err = s.UserMysqlManager.DeleteUser(req.AccountId)
	if err != nil {
		klog.Error("delete user error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("delete user err"))
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}

// GetSomeUsers implements the userSrv interface.
func (s *userSrv) GetSomeUsers(ctx context.Context, req *user.GetSomeUsersRequest) (resp *user.GetSomeUsersResponse, err error) {
	resp = new(user.GetSomeUsersResponse)
	users, err := s.UserMysqlManager.GetSomeUsers()
	if err != nil {
		klog.Error("get users error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr)
		return resp, nil
	}
	var uInfos []*base.User
	for _, u := range users {
		var uInfo base.User
		uInfo.Username = u.Username
		uInfo.AccountId = u.ID
		uInfo.PhoneNumber = u.PhoneNumber
		uInfo.AvatarBlobId = u.AvatarBlobId
		uInfo.OpenId = u.OpenID
		uInfos = append(uInfos, &uInfo)
	}

	resp.BaseResp = utils.BuildBaseResp(nil)
	resp.Users = uInfos
	return resp, nil
}

// GetAllUsers implements the userSrv interface.
func (s *userSrv) GetAllUsers(ctx context.Context, req *user.GetAllUsersRequest) (resp *user.GetAllUsersResponse, err error) {
	resp = new(user.GetAllUsersResponse)
	users, err := s.UserMysqlManager.GetAllUsers()
	if err != nil {
		klog.Error("get users error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("get users error"))
		return resp, nil
	}
	var uInfos []*base.User
	for _, u := range users {
		var uInfo base.User
		uInfo.Username = u.Username
		uInfo.AccountId = u.ID
		uInfo.PhoneNumber = u.PhoneNumber
		uInfo.AvatarBlobId = u.AvatarBlobId
		uInfo.OpenId = u.OpenID
		uInfos = append(uInfos, &uInfo)
	}
	resp.BaseResp = utils.BuildBaseResp(nil)
	resp.Users = uInfos
	return resp, nil
}

// Pay implements the userSrv interface.
func (s *userSrv) Pay(ctx context.Context, req *user.PayRequest) (resp *user.PayResponse, err error) {
	resp = new(user.PayResponse)
	var u *mysql.User
	if u, err = s.UserMysqlManager.GetUserByAccountId(req.AccountId); err != nil {
		if err == errno.RecordNotFound {
			resp.BaseResp = utils.BuildBaseResp(errno.RecordNotFound)
		} else {
			klog.Error("get user error", err)
			resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("get user error"))
		}
		return resp, nil
	}
	u.Balance -= req.FeeCent
	if err = s.UserMysqlManager.UpdateUser(u); err != nil {
		klog.Error("update user error", err)
		resp.BaseResp = utils.BuildBaseResp(errno.UserSrvErr.WithMessage("update user error"))
		return resp, nil
	}
	resp.BaseResp = utils.BuildBaseResp(nil)
	return resp, nil
}
