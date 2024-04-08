package service

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"testing"
)

//[]*ast.Field{
//					{
//						Name: &ast.Name{
//							Value: "id",
//						},
//					},
//					{
//						Name: &ast.Name{
//							Value: "company",
//						},
//						SelectionSet: &ast.SelectionSet{
//							Selections: []ast.Selection{
//								&ast.Field{
//									Name: &ast.Name{
//										Value: "id",
//									},
//								},
//								&ast.Field{
//									Name: &ast.Name{
//										Value: "name",
//									},
//								},
//								&ast.Field{
//									Name: &ast.Name{
//										Value: "logo",
//									},
//									SelectionSet: &ast.SelectionSet{
//										Selections: []ast.Selection{
//											&ast.Field{
//												Name: &ast.Name{
//													Value: "id",
//												},
//											},
//										},
//									},
//								},
//							},
//						},
//					},
//				}

func TestResolveToQuery(t *testing.T) {
	t.Run("Test ResolveToQuery", func(t *testing.T) {
		resolveParams := graphql.ResolveParams{
			Info: graphql.ResolveInfo{
				FieldName: "employees",
				FieldASTs: []*ast.Field{
					{
						SelectionSet: &ast.SelectionSet{
							Selections: []ast.Selection{
								&ast.Field{
									Name: &ast.Name{
										Value: "id",
									},
								},
								&ast.Field{
									Name: &ast.Name{
										Value: "company",
									},
									SelectionSet: &ast.SelectionSet{
										Selections: []ast.Selection{
											&ast.Field{
												Name: &ast.Name{
													Value: "id",
												},
											},
											&ast.Field{
												Name: &ast.Name{
													Value: "name",
												},
											},
											&ast.Field{
												Name: &ast.Name{
													Value: "logo",
												},
												SelectionSet: &ast.SelectionSet{
													Selections: []ast.Selection{
														&ast.Field{
															Name: &ast.Name{
																Value: "id",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		model := &Model{
			Table: "employees",
			Pk: &Field{
				Name: "id",
			},
			Fields: []*Field{
				{
					Name: "company_id",
					Type: Integer,
					Association: &Association{
						To: &Model{
							Table: "companies",
							Pk: &Field{
								Name: "id",
							},
							Fields: []*Field{
								{
									Name: "name",
									Type: CharacterVarying,
								},
								{
									Name: "logo_id",
									Type: Integer,
									Association: &Association{
										To: &Model{
											Table: "images",
											Pk: &Field{
												Name: "id",
											},
										},
									},
								},
							},
						},
						Column: "id",
						As:     "company",
					},
				},
			},
		}
		q := ResolveToQuery(resolveParams, model)
		stmt, _, err := q.ToSQL()
		if err != nil {
			t.Error(err)
		}
		expected := `SELECT "employees"."id", "companies"."id" AS "company.id", "companies"."name" AS "company.name", "images"."id" AS "company.logo.id" FROM employees LEFT JOIN companies ON employees.company_id = companies.id LEFT JOIN images AS company_logo ON companies.logo_id = company_logo.id`
		if stmt != expected {
			t.Errorf("Expected %s\nGot %s", expected, stmt)
		}
	})
}
