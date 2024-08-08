package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				assert.NoError(t, err)
				//这个写法就是insert
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(res)
				//require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				//res := sqlmock.NewResult(3, 1)
				assert.NoError(t, err)
				//这个写法就是insert
				mock.ExpectExec("INSERT INTO .*").WillReturnError(&mysqlDriver.MySQLError{
					Number: 1062,
				})
				//require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				//res := sqlmock.NewResult(3, 1)
				assert.NoError(t, err)
				//这个写法就是insert
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("数据库错误"))
				//require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
			wantErr: errors.New("数据库错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      tt.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
				SkipDefaultTransaction:                   true, //gorm 默认开启事物，跳过默认事物
			})
			assert.NoError(t, err)
			d := NewUserDAO(db)
			err = d.Insert(tt.ctx, tt.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
