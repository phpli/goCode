package repository

import (
	"context"
	"database/sql"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	UpdateNonZeroFields(ctx context.Context,
		user domain.User) error
	FindByWechat(ctx context.Context, openid string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, userCache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: userCache,
	}
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	//if errors.Is(err, cache.ErrKeyNotExist) {
	//
	//}
	// 缓存崩掉的情况，做好内存限流器 雪崩。穿透。击穿
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.toDomain(ue)
	//_ = r.cache.Set(ctx, u)
	//if err != nil {
	//	//打印日志
	//}
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//打印日志
		}
	}()
	return u, nil
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *CachedUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

//func (r *CachedUserRepository) Update(ctx context.Context, u domain.User) error {
//	return r.dao.Update(ctx, dao.User{
//		Id:          u.Id,
//		Birthday:    u.Birthday,
//		Gender:      u.Gender,
//		Description: u.Description,
//		Nickname:    u.Nickname,
//	})
//}

func (r *CachedUserRepository) UpdateNonZeroFields(ctx context.Context,
	user domain.User) error {
	return r.dao.UpdateById(ctx, r.toEntity(user))
}

func (r *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		Password:    u.Password,
		Birthday:    u.Birthday.UnixMilli(),
		Description: u.Description,
		Nickname:    u.Nickname,
		Gender:      u.Gender,
	}
}

func (r *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		Email:       u.Email.String,
		Phone:       u.Phone.String,
		Password:    u.Password,
		Description: u.Description,
		Nickname:    u.Nickname,
		Birthday:    time.UnixMilli(u.Birthday),
		Gender:      u.Gender,
		WechatInfo: domain.WechatInfo{
			OpenID:  u.WechatOpenID.String,
			UnionID: u.WechatUnionID.String,
		},
	}
}
