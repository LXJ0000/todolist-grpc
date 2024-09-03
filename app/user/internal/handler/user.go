package handler

import (
	"context"
	"errors"

	"github.com/LXJ0000/todolist-grpc-user/domain"
	service "github.com/LXJ0000/todolist-grpc-user/internal/service/pb"
	"github.com/LXJ0000/todolist-grpc-user/pkg/orm"
	"gorm.io/gorm"
)

type UserService struct {
	orm orm.Database
	service.UnimplementedUserServiceServer
}

func NewUserService(orm orm.Database) *UserService {
	return &UserService{orm: orm}
}

func (u *UserService) Login(ctx context.Context, req *service.UserRequest) (rsp *service.UserDetailResponse, err error) {
	rsp = &service.UserDetailResponse{}
	user, err := u.getUserInfo(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rsp.Code = 0
			return rsp, errors.New("用户不存在")
		}
		rsp.Code = 1
		return rsp, err
	}
	userModel := &service.UserModel{
		UserID:   uint32(user.ID),
		UserName: user.UserName,
		NickName: user.NickName,
	}
	rsp = &service.UserDetailResponse{
		UserModel: userModel,
		Code:      0,
	}
	return rsp, nil
}

func (u *UserService) Register(ctx context.Context, req *service.UserRequest) (rsp *service.UserDetailResponse, err error) {
	rsp = &service.UserDetailResponse{}
	user, err := u.createUser(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			rsp.Code = 0
			return rsp, errors.New("用户已存在")
		}
		rsp.Code = 1
		return rsp, err
	}
	userModel := &service.UserModel{
		UserID:   uint32(user.ID),
		UserName: user.UserName,
		NickName: user.NickName,
	}
	rsp = &service.UserDetailResponse{
		UserModel: userModel,
		Code:      0,
	}
	return rsp, nil
}

func (u *UserService) getUserInfo(ctx context.Context, req *service.UserRequest) (domain.User, error) {
	filter := map[string]interface{}{
		"user_name": req.UserName,
		"password":  req.Password, // TODO: 加密
	}
	var user domain.User
	err := u.orm.FindOne(ctx, &domain.User{}, filter, &user)
	return user, err
}

func (u *UserService) createUser(ctx context.Context, req *service.UserRequest) (domain.User, error) {
	var user domain.User
	user.UserName = req.UserName
	user.Password = req.Password // TODO 加密
	user.NickName = req.NickName
	u.orm.Insert(ctx, &domain.User{}, &user)
	return user, nil
}
