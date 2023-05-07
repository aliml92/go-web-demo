package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	db "conduit/db/sqlc"
)

type userResponse struct {
	User struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
		Token    string  `json:"token"`
	} `json:"user"`
}

func newUserResponse(user *db.User) *userResponse {
	resp := new(userResponse)
	resp.User.Username = user.Username
	resp.User.Email = user.Email
	resp.User.Bio = user.Bio
	resp.User.Image = user.Image
	token, _ := GenerateJWT(user.ID)
	resp.User.Token = token
	return resp
}

type userRegisterReq struct {
	User struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`
}

func (r *userRegisterReq) bind(c *gin.Context, p *db.CreateUserParams) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(r.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.ID = generateID()
	p.Username = r.User.Username
	p.Email = r.User.Email
	p.Password = string(hashed)
	return nil
}

func (s *Server) RegisterUser(c *gin.Context) { // TODO: POST /users - RegisterUser
	var (
		req userRegisterReq
		p   db.CreateUserParams
	)
	if err := req.bind(c, &p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidationError(err))
		return
	}
	user, err := s.store.CreateUser(c, p)
	if err != nil {
		if apiErr := convertToApiErr(err); apiErr != nil {
			c.JSON(http.StatusUnprocessableEntity, NewValidationError(apiErr))
			return
		}
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusCreated, newUserResponse(user))
}

type userLoginReq struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`
}

func (s *Server) LoginUser(c *gin.Context) { // TODO: POST /users/login - LoginUser
	var req userLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, NewValidationError(err))
		return
	}
	u, err := Nullable(s.store.GetUserByEmail(c.Request.Context(), req.User.Email))
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewValidationError(err))
		return
	}
	if u == nil {
		c.JSON(http.StatusForbidden, NewValidationError(ErrAccessForbidden))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.User.Password)); err != nil {
		c.JSON(http.StatusForbidden, NewValidationError(ErrAccessForbidden))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(u))
}
