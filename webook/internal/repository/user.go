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

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, userCache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: userCache,
	}
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

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
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//打印日志
		}
	}()
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

//func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
//	return r.dao.Update(ctx, dao.User{
//		Id:          u.Id,
//		Birthday:    u.Birthday,
//		Gender:      u.Gender,
//		Description: u.Description,
//		Nickname:    u.Nickname,
//	})
//}

func (r *UserRepository) UpdateNonZeroFields(ctx context.Context,
	user domain.User) error {
	return r.dao.UpdateById(ctx, r.toEntity(user))
}

func (r *UserRepository) toEntity(u domain.User) dao.User {
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
		Password:    u.Password,
		Birthday:    u.Birthday.UnixMilli(),
		Description: u.Description,
		Nickname:    u.Nickname,
		Gender:      u.Gender,
	}
}

func (r *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		Email:       u.Email.String,
		Phone:       u.Phone.String,
		Password:    u.Password,
		Description: u.Description,
		Nickname:    u.Nickname,
		Birthday:    time.UnixMilli(u.Birthday),
		Gender:      u.Gender,
	}
}
