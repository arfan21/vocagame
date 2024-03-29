package server

import (
	"github.com/arfan21/vocagame/internal/middleware"
	productctrl "github.com/arfan21/vocagame/internal/product/controller"
	productrepo "github.com/arfan21/vocagame/internal/product/repository"
	productsvc "github.com/arfan21/vocagame/internal/product/service"
	transactionctrl "github.com/arfan21/vocagame/internal/transaction/controller"
	transactionrepo "github.com/arfan21/vocagame/internal/transaction/repository"
	transactionsvc "github.com/arfan21/vocagame/internal/transaction/service"
	userctrl "github.com/arfan21/vocagame/internal/user/controller"
	userrepo "github.com/arfan21/vocagame/internal/user/repository"
	usersvc "github.com/arfan21/vocagame/internal/user/service"
	walletctrl "github.com/arfan21/vocagame/internal/wallet/controller"
	walletrepo "github.com/arfan21/vocagame/internal/wallet/repository"
	walletsvc "github.com/arfan21/vocagame/internal/wallet/service"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userRepoRedis := userrepo.NewRedis(s.dbRedis)
	userSvc := usersvc.New(userRepo, userRepoRedis)
	userCtrl := userctrl.New(userSvc)

	productRepo := productrepo.New(s.db, s.db)
	productSvc := productsvc.New(productRepo)
	productCtrl := productctrl.New(productSvc)

	walletRepo := walletrepo.New(s.db, s.db)
	walletSvc := walletsvc.New(walletRepo)
	walletCtrl := walletctrl.New(walletSvc)

	transactionRepo := transactionrepo.New(s.db, s.db)
	transactionSvc := transactionsvc.New(transactionRepo, walletSvc, productSvc)
	transactionCtrl := transactionctrl.New(transactionSvc)

	s.RoutesCustomer(api, userCtrl)
	s.RoutesProduct(api, productCtrl)
	s.RoutesWallet(api, walletCtrl)
	s.RoutesTransaction(api, transactionCtrl)
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
	productV1.Delete("/:productId", middleware.JWTAuth, ctrl.Delete)
}

func (s Server) RoutesWallet(route fiber.Router, ctrl *walletctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	walletV1 := v1.Group("/wallets")
	walletV1.Post("", middleware.JWTAuth, ctrl.Create)
	walletV1.Get("", middleware.JWTAuth, ctrl.GetByUserID)
}

func (s Server) RoutesTransaction(route fiber.Router, ctrl *transactionctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	transactionV1 := v1.Group("/transactions")
	transactionV1.Post("/deposit", middleware.JWTAuth, ctrl.CreateDepositTransaction)
	transactionV1.Post("/withdraw", middleware.JWTAuth, ctrl.CreateWithdrawTransaction)
	transactionV1.Get("/wallet", middleware.JWTAuth, ctrl.GetHistoryWalletByUserID)
	transactionV1.Post("/checkout", middleware.JWTAuth, ctrl.Checkout)
	transactionV1.Get("/:transactionId", middleware.JWTAuth, ctrl.GetByID)
}
