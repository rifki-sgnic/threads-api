package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists         = errors.New("username or email already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService interface {
	Register(in RegisterInput) (*models.User, error)
	Login(in LoginInput) (string, error)
}

type authService struct {
	users     repository.UserRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthService(users repository.UserRepository, secret string, ttl time.Duration) AuthService {
	return &authService{
		users:     users,
		jwtSecret: []byte(secret),
		jwtTTL:    ttl,
	}
}

func (s *authService) Register(in RegisterInput) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: string(hash),
	}

	if err := s.users.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(in LoginInput) (string, error) {
	user, err := s.users.FindByEmail(in.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(user.ID)
}

func (s *authService) generateToken(userID uint) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
