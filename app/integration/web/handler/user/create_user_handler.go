package user

import (
	"net/http"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"podGopher/core/port/inbound/user"
	"podGopher/integration/web/auth"
	"podGopher/integration/web/handler"

	"github.com/gin-gonic/gin"
)

type CreateUserHandler struct {
	route     *handler.Route
	port      user.CreateUserPort
	adminAuth auth.AdminAuth
}

func (h *CreateUserHandler) Authorize(context *gin.Context) {
	username, password, ok := context.Request.BasicAuth()
	if !ok || !h.adminAuth.IsValid(username, password) {
		_ = context.AbortWithError(http.StatusForbidden, domainError.NewAuthorizationError())
	}
}

type CreateUserRequestDto struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type userResponseDto struct {
	Id        string           `json:"id" binding:"required"`
	Username  string           `json:"username" binding:"required"`
	ShowRoles []model.ShowRole `json:"showRoles" binding:"required"`
}

func (h *CreateUserHandler) GetRoute() *handler.Route {
	return h.route
}

func NewCreateUserHandler(portMap inbound.PortMap, adminAuth auth.AdminAuth) *CreateUserHandler {
	return &CreateUserHandler{
		route: &handler.Route{
			Method: http.MethodPost,
			Path:   "/admin/show/:showId/user",
		},
		port:      portMap[inbound.CreateUser].(user.CreateUserPort),
		adminAuth: adminAuth,
	}
}

func (h *CreateUserHandler) Handle(context *gin.Context) {
	var request *CreateUserRequestDto
	if err := context.BindJSON(&request); err != nil {
		context.Abort()
		return
	}

	h.handleCreateUser(context, request)
}

func (h *CreateUserHandler) handleCreateUser(context *gin.Context, request *CreateUserRequestDto) {
	createUserCommand := &user.CreateUserCommand{ShowId: context.Param("showId"), Username: request.Username, Role: request.Role}
	if createdUser, err := h.port.CreateUser(createUserCommand); err != nil {
		_ = context.Error(err)
	} else {
		responseDto := userResponseDto{Id: createdUser.Id, Username: createdUser.Username, ShowRoles: createdUser.ShowRoles}
		context.JSON(http.StatusCreated, responseDto)
	}
}
