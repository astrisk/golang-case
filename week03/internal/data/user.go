package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"user/internal/biz"
	"user/internal/model"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func (u userRepo) CreateUser(ctx context.Context, user *biz.User) error {
	// TODO getUser

	return u.data.db.Model(&model.User{}).Create(&model.User{
		Id:   0,
		Name: user.Name,
	}).Error
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
