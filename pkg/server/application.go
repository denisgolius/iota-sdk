package server

import (
	services2 "github.com/iota-agency/iota-sdk/modules/finance/services"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/event"
	"github.com/iota-agency/iota-sdk/pkg/infrastructure/persistence"
	"github.com/iota-agency/iota-sdk/pkg/services"
	"gorm.io/gorm"
)

func ConstructApp(db *gorm.DB) application.Application {
	eventPublisher := event.NewEventPublisher()
	app := application.New(db, eventPublisher)
	moneyAccountService := services2.NewMoneyAccountService(
		persistence.NewMoneyAccountRepository(),
		eventPublisher,
	)

	app.RegisterService(services.NewUserService(persistence.NewUserRepository(), eventPublisher))
	app.RegisterService(services.NewSessionService(persistence.NewSessionRepository(), eventPublisher))
	app.RegisterService(services.NewAuthService(app))
	app.RegisterService(services.NewRoleService(persistence.NewRoleRepository(), eventPublisher))
	app.RegisterService(services2.NewPaymentService(
		persistence.NewPaymentRepository(), eventPublisher, moneyAccountService,
	))
	app.RegisterService(services.NewProjectStageService(persistence.NewProjectStageRepository(), eventPublisher))
	app.RegisterService(services2.NewCurrencyService(persistence.NewCurrencyRepository(), eventPublisher))
	app.RegisterService(services2.NewExpenseCategoryService(
		persistence.NewExpenseCategoryRepository(),
		eventPublisher,
	))
	app.RegisterService(services.NewPositionService(persistence.NewPositionRepository(), eventPublisher))
	app.RegisterService(services.NewEmployeeService(persistence.NewEmployeeRepository(), eventPublisher))
	app.RegisterService(services.NewAuthLogService(persistence.NewAuthLogRepository(), eventPublisher))
	app.RegisterService(services.NewPromptService(persistence.NewPromptRepository(), eventPublisher))
	app.RegisterService(services2.NewExpenseService(
		persistence.NewExpenseRepository(), eventPublisher, moneyAccountService,
	))
	app.RegisterService(services.NewProjectService(persistence.NewProjectRepository(), eventPublisher))

	app.RegisterService(services.NewEmbeddingService(app))
	app.RegisterService(services.NewDialogueService(persistence.NewDialogueRepository(), app))
	app.RegisterService(moneyAccountService)
	return app
}