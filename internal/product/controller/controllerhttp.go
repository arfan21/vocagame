package productctrl

import (
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/internal/product"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/exception"
	"github.com/arfan21/vocagame/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ControllerHTTP struct {
	svc product.Service
}

func New(svc product.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create Product
// @Description Create Product
// @Tags Product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.ProductCreateRequest true "Payload Create Product Request"
// @Success 201 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/products [post]
func (ctrl ControllerHTTP) Create(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "invalid or expired token",
		})
	}

	var req model.ProductCreateRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	uuidUserID, err := uuid.Parse(claims.Subject)
	exception.PanicIfNeeded(err)
	req.UserID = uuidUserID

	err = ctrl.svc.Create(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusCreated,
	})
}

// @Summary Get Products
// @Description Get Products
// @Tags Product
// @Accept json
// @Produce json
// @Param page query string true "Page"
// @Param limit query string true "Limit"
// @Param name query string false "Name of product"
// @Param owner_id query string false "Owner ID"
// @Success 200 {object} pkgutil.HTTPResponse{data=pkgutil.PaginationResponse{data=model.GetProductResponse}}
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/products [get]
func (ctrl ControllerHTTP) GetProducts(c *fiber.Ctx) error {
	reqQuery := model.GetListProductRequest{}
	err := c.QueryParser(&reqQuery)
	exception.PanicIfNeeded(err)

	res, err := ctrl.svc.GetProducts(c.UserContext(), reqQuery)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusOK,
		Data: res,
	})
}
