package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
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
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
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
		Id:          u.Id,
		Email:       u.Email,
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
		Email:       u.Email,
		Password:    u.Password,
		Description: u.Description,
		Nickname:    u.Nickname,
		Birthday:    time.UnixMilli(u.Birthday),
		Gender:      u.Gender,
	}
}
