package model

import (
	orm "application/pkg/mysql"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// LoginAndRegister TODO: yjy 如果闭包，可以使用 build 建造模式，否则字段要大写，在 json序列化会为空值
type LoginAndRegister struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(256);uniqueIndex;comment:用户名"`
	Password string `json:"password" gorm:"type:varchar(256);comment:密码"`
	NickName string `json:"nickName" gorm:"type:varchar(256);comment:昵称"`
	Phone    string `json:"phone" gorm:"type:varchar(256);comment:手机号"`
	Email    string `json:"email" gorm:"type:varchar(256);comment:邮箱"`
	Sex      string `json:"sex" gorm:"type:varchar(256);comment:性别"`
	Address  string `json:"address" gorm:"type:varchar(256);comment:地址"`
	Age      int    `json:"age" gorm:"type:int;comment:年龄"`
}

func (LoginAndRegister) TableName() string {
	return "login_register"
}

type LoginAndRegisterManager struct {
	orm *orm.Client
}

// CreateCollectionTasksTable 创建数据表
func CreateCollectionTasksTable(ctx context.Context, lar *LoginAndRegisterManager) error {
	tname := new(LoginAndRegister).TableName()
	if lar.orm.DB().Migrator().HasTable(tname) {
		fmt.Println("table already exists")
		return nil
	}
	// 迁移模型
	return lar.orm.DB().AutoMigrate(new(LoginAndRegister))
}

// CreateTable 创建数据表
func (lar *LoginAndRegisterManager) CreateTables(ctx context.Context) (err error) {
	return CreateCollectionTasksTable(ctx, lar)
}

// TODO: yjy 代码位置不对，注意解耦，解耦
type Options func(*LoginAndRegisterManager)

// NewLoginAndRegisterManager 创建一个LoginAndRegisterManager, options选择模式
func NewLoginAndRegisterManager(orm *orm.Client, opts ...Options) (*LoginAndRegisterManager, error) {
	loginAndRegisterManager := &LoginAndRegisterManager{
		orm: orm,
	}
	for _, opt := range opts {
		opt(loginAndRegisterManager)
	}
	return loginAndRegisterManager, nil
}
