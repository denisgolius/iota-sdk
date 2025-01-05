package persistence

import (
	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/currency"
	corepersistence "github.com/iota-uz/iota-sdk/modules/core/infrastructure/persistence"
	"github.com/iota-uz/iota-sdk/modules/finance/domain/aggregates/expense"
	category "github.com/iota-uz/iota-sdk/modules/finance/domain/aggregates/expense_category"
	moneyaccount "github.com/iota-uz/iota-sdk/modules/finance/domain/aggregates/money_account"
	"github.com/iota-uz/iota-sdk/modules/finance/domain/aggregates/payment"
	"github.com/iota-uz/iota-sdk/modules/finance/domain/entities/counterparty"
	"github.com/iota-uz/iota-sdk/modules/finance/domain/entities/transaction"
	"github.com/iota-uz/iota-sdk/modules/finance/infrastructure/persistence/models"
	"github.com/iota-uz/iota-sdk/pkg/mapping"
)

func toDBTransaction(entity *transaction.Transaction) *models.Transaction {
	return &models.Transaction{
		ID:                   entity.ID,
		Amount:               entity.Amount,
		Comment:              entity.Comment,
		AccountingPeriod:     entity.AccountingPeriod,
		TransactionDate:      entity.TransactionDate,
		DestinationAccountID: entity.DestinationAccountID,
		OriginAccountID:      entity.OriginAccountID,
		TransactionType:      string(entity.TransactionType),
		CreatedAt:            entity.CreatedAt,
	}
}

func toDomainTransaction(dbTransaction *models.Transaction) (*transaction.Transaction, error) {
	_type, err := transaction.NewType(dbTransaction.TransactionType)
	if err != nil {
		return nil, err
	}

	return &transaction.Transaction{
		ID:                   dbTransaction.ID,
		Amount:               dbTransaction.Amount,
		TransactionType:      _type,
		Comment:              dbTransaction.Comment,
		AccountingPeriod:     dbTransaction.AccountingPeriod,
		TransactionDate:      dbTransaction.TransactionDate,
		DestinationAccountID: dbTransaction.DestinationAccountID,
		OriginAccountID:      dbTransaction.OriginAccountID,
		CreatedAt:            dbTransaction.CreatedAt,
	}, nil
}

func toDBPayment(entity payment.Payment) (*models.Payment, *models.Transaction) {
	dbTransaction := &models.Transaction{
		ID:                   entity.TransactionID(),
		Amount:               entity.Amount(),
		Comment:              entity.Comment(),
		AccountingPeriod:     entity.AccountingPeriod(),
		TransactionDate:      entity.TransactionDate(),
		OriginAccountID:      nil,
		DestinationAccountID: &entity.Account().ID,
		TransactionType:      string(transaction.Deposit),
		CreatedAt:            entity.CreatedAt(),
	}
	dbPayment := &models.Payment{
		ID:             entity.ID(),
		TransactionID:  entity.TransactionID(),
		CounterpartyID: entity.CounterpartyID(),
		CreatedAt:      entity.CreatedAt(),
		UpdatedAt:      entity.UpdatedAt(),
	}
	return dbPayment, dbTransaction
}

// TODO: populate user && account
func toDomainPayment(dbPayment *models.Payment, dbTransaction *models.Transaction) (payment.Payment, error) {
	t, err := toDomainTransaction(dbTransaction)
	if err != nil {
		return nil, err
	}
	return payment.NewWithID(
		dbPayment.ID,
		t.Amount,
		t.ID,
		dbPayment.CounterpartyID,
		t.Comment,
		&moneyaccount.Account{ID: *t.DestinationAccountID}, //nolint:exhaustruct
		&user.User{}, //nolint:exhaustruct
		t.TransactionDate,
		t.AccountingPeriod,
		dbPayment.CreatedAt,
		dbPayment.UpdatedAt,
	), nil
}

func toDBExpenseCategory(entity *category.ExpenseCategory) *models.ExpenseCategory {
	return &models.ExpenseCategory{
		ID:               entity.ID,
		Name:             entity.Name,
		Description:      &entity.Description,
		Amount:           entity.Amount,
		AmountCurrencyID: string(entity.Currency.Code),
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

func toDomainExpenseCategory(dbCategory *models.ExpenseCategory) (*category.ExpenseCategory, error) {
	return &category.ExpenseCategory{
		ID:          dbCategory.ID,
		Name:        dbCategory.Name,
		Description: mapping.Value(dbCategory.Description),
		Amount:      dbCategory.Amount,
		Currency:    currency.Currency{Code: currency.Code(dbCategory.AmountCurrencyID)}, //nolint:exhaustruct
		CreatedAt:   dbCategory.CreatedAt,
		UpdatedAt:   dbCategory.UpdatedAt,
	}, nil
}

func toDomainMoneyAccount(dbAccount *models.MoneyAccount) (*moneyaccount.Account, error) {
	currencyEntity, err := corepersistence.ToDomainCurrency(dbAccount.Currency)
	if err != nil {
		return nil, err
	}
	return &moneyaccount.Account{
		ID:            dbAccount.ID,
		Name:          dbAccount.Name,
		AccountNumber: dbAccount.AccountNumber,
		Balance:       dbAccount.Balance,
		Currency:      *currencyEntity,
		Description:   dbAccount.Description,
		CreatedAt:     dbAccount.CreatedAt,
		UpdatedAt:     dbAccount.UpdatedAt,
	}, nil
}

func toDBMoneyAccount(entity *moneyaccount.Account) *models.MoneyAccount {
	return &models.MoneyAccount{
		ID:                entity.ID,
		Name:              entity.Name,
		AccountNumber:     entity.AccountNumber,
		Balance:           entity.Balance,
		BalanceCurrencyID: string(entity.Currency.Code),
		Currency:          corepersistence.ToDBCurrency(&entity.Currency),
		Description:       entity.Description,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}

func toDomainExpense(dbExpense *models.Expense, dbTransaction *models.Transaction) (*expense.Expense, error) {
	return &expense.Expense{
		ID:               dbExpense.ID,
		Amount:           -1 * dbTransaction.Amount,
		Account:          moneyaccount.Account{ID: *dbTransaction.OriginAccountID}, //nolint:exhaustruct
		Category:         category.ExpenseCategory{ID: dbExpense.CategoryID},       //nolint:exhaustruct
		Comment:          dbTransaction.Comment,
		TransactionID:    dbExpense.TransactionID,
		AccountingPeriod: dbTransaction.AccountingPeriod,
		Date:             dbTransaction.TransactionDate,
		CreatedAt:        dbExpense.CreatedAt,
		UpdatedAt:        dbExpense.UpdatedAt,
	}, nil
}

func toDBExpense(entity *expense.Expense) (*models.Expense, *transaction.Transaction) {
	domainTransaction := &transaction.Transaction{
		ID:                   entity.TransactionID,
		Amount:               -1 * entity.Amount,
		Comment:              entity.Comment,
		AccountingPeriod:     entity.AccountingPeriod,
		TransactionDate:      entity.Date,
		OriginAccountID:      &entity.Account.ID,
		DestinationAccountID: nil,
		TransactionType:      transaction.Withdrawal,
		CreatedAt:            entity.CreatedAt,
	}
	dbExpense := &models.Expense{
		ID:            entity.ID,
		CategoryID:    entity.Category.ID,
		TransactionID: entity.TransactionID,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
	return dbExpense, domainTransaction
}

func toDomainCounterparty(dbRow *models.Counterparty) (counterparty.Counterparty, error) {
	partyType, err := counterparty.NewType(dbRow.Type)
	if err != nil {
		return nil, err
	}
	legalType, err := counterparty.NewLegalType(dbRow.LegalType)
	if err != nil {
		return nil, err
	}
	return counterparty.NewWithID(
		dbRow.ID,
		dbRow.TIN,
		dbRow.Name,
		partyType,
		legalType,
		dbRow.LegalAddress,
		dbRow.CreatedAt,
		dbRow.UpdatedAt,
	), nil
}

func toDBCounterparty(entity counterparty.Counterparty) (*models.Counterparty, error) {
	return &models.Counterparty{
		ID:           entity.ID(),
		TIN:          entity.TIN(),
		Name:         entity.Name(),
		Type:         string(entity.Type()),
		LegalType:    string(entity.LegalType()),
		LegalAddress: entity.LegalAddress(),
		CreatedAt:    entity.CreatedAt(),
		UpdatedAt:    entity.UpdatedAt(),
	}, nil
}