package api

import (
	"net/http"
	"path"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/organization"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/Stuhub-io/utils/organization_inviteutils"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgService *organization.Service
}

type NewOrganizationHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	OrgService     *organization.Service
}

func UseOrganizationHandler(params NewOrganizationHandlerParams) {
	handler := &OrganizationHandler{
		orgService: params.OrgService,
	}

	router := params.Router.Group("/organization-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())

	router.POST("/create", decorators.CurrentUser(handler.CreateOrganization))
	router.GET("/joined", decorators.CurrentUser(handler.GetJoinedOrgs))
	router.GET("/get-by-slug", decorators.CurrentUser(handler.GetOrgBySlug))
	router.POST("/invite-by-emails", decorators.CurrentUser(handler.InviteMembersByEmail))
	router.GET(path.Join("/invite-details", ":"+organization_inviteutils.InviteIDParam), handler.GetInviteDetails)
	router.POST("/invite-validate", decorators.CurrentUser(handler.ValidateOrgInvitation))
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context, user *domain.User) {
	var body request.CreateOrgBody
	if ok, vr := request.Validate(c, &body); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.orgService.CreateOrganization(organization.CreateOrganizationDto{
		OwnerPkID:   user.PkID,
		Name:        body.Name,
		Description: body.Description,
		Avatar:      body.Avatar,
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *OrganizationHandler) GetJoinedOrgs(c *gin.Context, user *domain.User) {
	data, err := h.orgService.GetJoinedOrgs(user.PkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *OrganizationHandler) GetOrgBySlug(c *gin.Context, user *domain.User) {
	var params request.GetOrgBySlugParams
	if ok, vr := request.Validate(c, &params); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.orgService.GetOrganizationDetailBySlug(params.Slug)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *OrganizationHandler) GetInviteDetails(c *gin.Context) {
	inviteID, valid := organization_inviteutils.GetInviteIDParam(c)
	if !valid {
		response.BindError(c, "inviteID is missing or invalid")
		return
	}

	invite, err := h.orgService.GetInviteDetails(inviteID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, invite)
}

func (h *OrganizationHandler) InviteMembersByEmail(c *gin.Context, user *domain.User) {
	var params request.InviteMembersByEmailParams
	if ok, vr := request.Validate(c, &params); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.orgService.InviteMemberByEmails(organization.InviteMemberByEmailsDto{
		Owner:       user,
		OrgInfo:     params.OrgInfo,
		InviteInfos: params.Infos,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Emails sent")
}

func (h *OrganizationHandler) ValidateOrgInvitation(c *gin.Context, user *domain.User) {
	var params request.ValidateOrgInvitationParams
	if ok, vr := request.Validate(c, &params); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.orgService.ValidateOrgInviteToken(organization.ValidateOrgInviteTokenDto{
		UserPkID: user.PkID,
		Token:    params.Token,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Invitation validated successfully!")
}
