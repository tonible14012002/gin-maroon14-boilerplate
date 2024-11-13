package ports

import (
	"context"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, *domain.Error)
	GetUserByPkID(ctx context.Context, pkID int64) (*domain.User, *domain.Error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error)
	GetOrCreateUserByEmail(ctx context.Context, email, salt string) (*domain.User, *domain.Error)
	CreateUserWithGoogleInfo(ctx context.Context, email, salt, firstName, lastName, avatar string) (*domain.User, *domain.Error)
	SetUserPassword(ctx context.Context, PkID int64, hashedPassword string) *domain.Error
	CheckPassword(ctx context.Context, email, rawPassword string, hasher Hasher) (bool, *domain.Error)
	UpdateUserInfo(ctx context.Context, PkID int64, firstName, lastName, avatar string) (*domain.User, *domain.Error)
	SetUserActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.User, *domain.Error)
}

type OrganizationRepository interface {
	GetOrgMembers(ctx context.Context, pkID int64) ([]domain.OrganizationMember, *domain.Error)
	GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByName(ctx context.Context, ownerPkID int64, name string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByPkID(ctx context.Context, ownerPkID, pkId int64) (*domain.Organization, *domain.Error)
	GetOrgsByUserPkID(ctx context.Context, usePkID int64) ([]*domain.Organization, *domain.Error)
	GetOrgMemberByEmail(ctx context.Context, orgPkID int64, email string) (*domain.OrganizationMember, *domain.Error)
	GetOrgMemberByUserPkID(ctx context.Context, orgPkID int64, userPkID int64) (*domain.OrganizationMember, *domain.Error)
	CreateOrg(ctx context.Context, userPkID int64, name, description, avatar string) (*domain.Organization, *domain.Error)
	AddMemberToOrg(ctx context.Context, orgPkID int64, userPkID *int64, role string) (*domain.OrganizationMember, *domain.Error)
	SetOrgMemberActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.OrganizationMember, *domain.Error)
}

type PageRepository interface {
	List(ctx context.Context, query domain.PageListQuery) ([]domain.Page, *domain.Error)
	Update(ctx context.Context, pagePkID int64, page domain.PageUpdateInput) (*domain.Page, *domain.Error)
	Move(ctx context.Context, pagePkID int64, parentPagePkID *int64) (*domain.Page, *domain.Error)
	CreatePage(ctx context.Context, page domain.PageInput) (*domain.Page, *domain.Error)
	GetByID(ctx context.Context, pageID string) (*domain.Page, *domain.Error)
	UpdateContent(ctx context.Context, pagePkID int64, content domain.DocumentInput) (*domain.Page, *domain.Error)
	Archive(ctx context.Context, pagePkID int64) (*domain.Page, *domain.Error)
}

type OrganizationInviteRepository interface {
	CreateInvite(ctx context.Context, organizationPkId int64, userPkId int64) (*domain.OrganizationInvite, *domain.Error)
	UpdateInvite(ctx context.Context, invite model.OrganizationInvite) (*domain.OrganizationInvite, *domain.Error)
	GetInviteByID(ctx context.Context, inviteID string) (*domain.OrganizationInvite, *domain.Error)
}
