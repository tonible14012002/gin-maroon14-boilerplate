package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/userutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewUserRepositoryParams struct {
	Store *store.DBStore
	Cfg   config.Config
}

func NewUserRepository(params NewUserRepositoryParams) ports.UserRepository {
	return &UserRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundById(id)
		}

		return nil, domain.ErrDatabaseQuery
	}

	return userutils.TransformUserModelToDomain(user), nil
}

func (r *UserRepository) GetUserByPkID(ctx context.Context, pkId int64) (*domain.User, *domain.Error) {
	cachedUser := r.store.Cache().GetUser(pkId)
	if cachedUser != nil {
		return cachedUser, nil
	}

	var userModel model.User
	err := r.store.DB().Where("pkid = ?", pkId).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, domain.ErrDatabaseQuery
	}

	user := userutils.TransformUserModelToDomain(userModel)

	// go func() {
	// 	r.store.Cache().SetUser(user, time.Hour)
	// }()

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFoundByEmail(email)
		}

		return nil, domain.ErrDatabaseQuery
	}

	return userutils.TransformUserModelToDomain(user), nil
}

func (r *UserRepository) GetOrCreateUserByEmail(ctx context.Context, email string, salt string) (*domain.User, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDatabaseQuery
		}

		user = model.User{
			Email: email,
			Salt:  salt,
		}

		err = r.store.DB().Create(&user).Error
		if err != nil {
			return nil, domain.ErrDatabaseQuery
		}
	}

	return userutils.TransformUserModelToDomain(user), nil
}

func (r *UserRepository) CreateUserWithGoogleInfo(ctx context.Context, email, salt, firstName, lastName, avatar string) (*domain.User, *domain.Error) {
	user := model.User{
		Email:      email,
		Salt:       salt,
		FirstName:  firstName,
		LastName:   lastName,
		Avatar:     avatar,
		OauthGmail: email,
	}

	err := r.store.DB().Create(&user).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return userutils.TransformUserModelToDomain(user), nil
}

func (r *UserRepository) SetUserPassword(ctx context.Context, pkID int64, hashedPassword string) *domain.Error {
	// FIXME: Add password hashing
	err := r.store.DB().Model(&model.User{}).Where("pkid = ?", pkID).Update("password", hashedPassword).Error
	if err != nil {
		return domain.ErrDatabaseMutation
	}

	return nil
}

func (r *UserRepository) CheckPassword(ctx context.Context, email, rawPassword string, hasher ports.Hasher) (bool, *domain.Error) {
	var user model.User
	err := r.store.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, domain.ErrUserNotFoundByEmail(email)
		}

		return false, domain.ErrDatabaseQuery
	}

	valid := hasher.Compare(rawPassword, *user.Password, user.Salt)

	return valid, nil
}

func (r *UserRepository) UpdateUserInfo(ctx context.Context, PkID int64, firstName, lastName, avatar string) (*domain.User, *domain.Error) {
	var user = model.User{
		FirstName: firstName,
		LastName:  lastName,
		Avatar:    avatar,
	}
	err := r.store.DB().Model(&model.User{}).Where("pkid = ?", PkID).Updates(&user).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return userutils.TransformUserModelToDomain(user), nil
}

func (r *UserRepository) SetUserActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.User, *domain.Error) {
	var user model.User

	err := r.store.DB().Model(&user).Clauses(clause.Returning{}).Where("pkid = ?", pkID).Update("activated_at", activatedAt).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return userutils.TransformUserModelToDomain(user), nil
}
