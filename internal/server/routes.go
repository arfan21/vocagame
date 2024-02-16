package server

import (
	"github.com/arfan21/vocagame/internal/middleware"
	productctrl "github.com/arfan21/vocagame/internal/product/controller"
	productrepo "github.com/arfan21/vocagame/internal/product/repository"
	productsvc "github.com/arfan21/vocagame/internal/product/service"
	userctrl "github.com/arfan21/vocagame/internal/user/controller"
	userrepo "github.com/arfan21/vocagame/internal/user/repository"
	usersvc "github.com/arfan21/vocagame/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userRepoRedis := userrepo.NewRedis(s.dbRedis)
	userSvc := usersvc.New(userRepo, userRepoRedis)
	userCtrl := userctrl.New(userSvc)

	productRepo := productrepo.New(s.db)
	productSvc := productsvc.New(productRepo)
	productCtrl := productctrl.New(productSvc)

	s.RoutesCustomer(api, userCtrl)
	s.RoutesProduct(api, productCtrl)

}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/users")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
	usersV1.Post("/refresh-token", ctrl.RefreshToken)
	usersV1.Post("/logout", middleware.JWTAuth, ctrl.Logout)
}

func (s Server) RoutesProduct(route fiber.Router, ctrl *productctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	productV1 := v1.Group("/products")
	productV1.Post("", middleware.JWTAuth, ctrl.Create)
	productV1.Get("", ctrl.GetProducts)
	productV1.Put("/:productId", middleware.JWTAuth, ctrl.Update)
}
