package payment

import (
	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	moneyaccount "github.com/iota-uz/iota-sdk/modules/finance/domain/aggregates/money_account"
	"time"
)

func NewWithID(
	id uint,
	amount float64,
	transactionID uint,
	comment string,
	account *moneyaccount.Account,
	createdBy *user.User,
	date, accountingPeriod, createdAt, updatedAt time.Time,
) Payment {
	return &payment{
		id:               id,
		amount:           amount,
		account:          account,
		transactionID:    transactionID,
		transactionDate:  date,
		accountingPeriod: accountingPeriod,
		user:             createdBy,
		comment:          comment,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

func New(
	amount float64,
	transactionID uint,
	comment string,
	account *moneyaccount.Account,
	createdBy *user.User,
	date, accountingPeriod time.Time,
) Payment {
	return NewWithID(
		0,
		amount,
		transactionID,
		comment,
		account,
		createdBy,
		date,
		accountingPeriod,
		time.Now(),
		time.Now(),
	)
}

type payment struct {
	id               uint
	amount           float64
	transactionID    uint
	transactionDate  time.Time
	accountingPeriod time.Time
	comment          string
	account          *moneyaccount.Account
	user             *user.User
	createdAt        time.Time
	updatedAt        time.Time
}

func (p *payment) ID() uint {
	return p.id
}

func (p *payment) SetID(id uint) {
	p.id = id
}

func (p *payment) Amount() float64 {
	return p.amount
}

func (p *payment) SetAmount(a float64) {
	p.amount = a
	p.updatedAt = time.Now()
}

func (p *payment) TransactionID() uint {
	return p.transactionID
}

func (p *payment) TransactionDate() time.Time {
	return p.transactionDate
}

func (p *payment) SetTransactionDate(t time.Time) {
	p.transactionDate = t
	p.updatedAt = time.Now()
}

func (p *payment) AccountingPeriod() time.Time {
	return p.accountingPeriod
}

func (p *payment) SetAccountingPeriod(t time.Time) {
	p.accountingPeriod = t
	p.updatedAt = time.Now()
}

func (p *payment) Comment() string {
	return p.comment
}

func (p *payment) SetComment(s string) {
	p.comment = s
	p.updatedAt = time.Now()
}

func (p *payment) Account() *moneyaccount.Account {
	return p.account
}

func (p *payment) User() *user.User {
	return p.user
}

func (p *payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *payment) UpdatedAt() time.Time {
	return p.updatedAt
}