package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
	ErrRecordNotFound        = gorm.ErrRecordNotFound
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context,
		user domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo   repository.UserRepository
	logger logger.LoggerV1
}

func NewUserService(repo repository.UserRepository, logger logger.LoggerV1) UserService {
	return &userService{
		repo: repo,
		//logger: logger.Named("user"),
		logger: logger, //注入了，但是又没完全注入
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, gorm.ErrRecordNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Id: u.Id, Email: u.Email, Birthday: u.Birthday, Gender: u.Gender, Description: u.Description}, nil
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	// UpdateNicknameAndXXAnd
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil {
		return u, err
	}
	// 这里会有主从延迟的坑
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenID)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}
	u = domain.User{
		WechatInfo: wechatInfo,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}
	// 这里会有主从延迟的坑
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenID)
}
