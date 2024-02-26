package user

import (
	"context"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		timeout:    time.Duration(2) * time.Second,
		Repository: repository,
	}
}
func (s *service) CreateUser(c context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// TODO: hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserResponse{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) Login(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserResponse{}, err
	}

	if err := util.CheckPassword(req.Password, u.Password); err != nil {
		return &LoginUserResponse{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserResponse{}, err
	}

	return &LoginUserResponse{
		accessToken: ss,
		ID:          strconv.Itoa(int(u.ID)),
		Username:    u.Username,
	}, nil

}

const (
	secretKey = "secret"
)

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
