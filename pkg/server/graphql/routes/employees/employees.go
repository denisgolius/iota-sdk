package employees

import (
	"github.com/graphql-go/graphql"
	"github.com/iota-agency/iota-erp/models"
	"github.com/iota-agency/iota-erp/sdk/graphql/adapters"
	"gorm.io/gorm"
)

func Queries(db *gorm.DB) []*graphql.Field {
	return adapters.DefaultQueries(db, &models.Employee{}, "employee", "employees")
}

func Mutations(db *gorm.DB) []*graphql.Field {
	return adapters.DefaultMutations(db, &models.Employee{}, "employee")
}
