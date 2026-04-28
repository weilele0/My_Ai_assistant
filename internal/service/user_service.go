package service

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"My_AI_Assistant/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务层
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建 UserService
func NewUserService(userRepo *repository.UserRepository) *UserService {

	return &UserService{userRepo: userRepo}
} //实例化userservice结构体，并装入数据库工具，就是传入repository结构体

// Register 用户注册
func (s *UserService) Register(req model.RegisterRequest) error {
	// 1. 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return errors.New("用户名已存在")
	}
	//慢哈希算法加密密码
	// 2. 密码加密（使用 bcrypt）
	/////////////          结构体         第一个参数密码，并强转成byte类型切片  加密强度
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. 创建用户对象
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword), //把加密的密码转换成字符串
	}

	// 4. 保存到数据库
	return s.userRepo.Create(user)
}

// Login 用户登录
func (s *UserService) Login(req model.LoginRequest) (token string, userInfo *model.User, err error) {
	// 1. 根据用户名查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 2. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 3. 生成 JWT Token
	token, err = jwt.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}
