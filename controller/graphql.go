package controller

import (
	"fmt"
	"maintainman/model"
	"maintainman/service"

	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":           &graphql.Field{Type: graphql.Int},
			"name":         &graphql.Field{Type: graphql.String},
			"display_name": &graphql.Field{Type: graphql.String},
			"role_id":      &graphql.Field{Type: graphql.Int},
			"role": &graphql.Field{
				Type: roleType,
			},
		},
	},
)

var roleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Role",
		Fields: graphql.Fields{
			"id":           &graphql.Field{Type: graphql.Int},
			"name":         &graphql.Field{Type: graphql.String},
			"display_name": &graphql.Field{Type: graphql.String},
			"discription":  &graphql.Field{Type: graphql.String},
			"perms": &graphql.Field{
				Type: graphql.NewList(permissionType),
			},
		},
	},
)

var permissionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Permission",
		Fields: graphql.Fields{
			"id":           &graphql.Field{Type: graphql.Int},
			"name":         &graphql.Field{Type: graphql.String},
			"display_name": &graphql.Field{Type: graphql.String},
			"discription":  &graphql.Field{Type: graphql.String},
			"default":      &graphql.Field{Type: graphql.Boolean},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, idOK := p.Args["id"].(int)
					name, nameOK := p.Args["name"].(string)
					if idOK && nameOK {
						return nil, fmt.Errorf("Can not query with both \"id\" and \"name\"")
					}
					var result *model.ApiJson
					if idOK {
						result = service.GetUserInfoByID(uint(id))
					} else if nameOK {
						result = service.GetUserInfoByName(name)
					} else {
						return nil, fmt.Errorf("No query key specified")
					}
					if !result.Status {
						return nil, fmt.Errorf(result.Msg)
					} else {
						return result.Data, nil
					}
				},
			},
			"role": &graphql.Field{
				Type: roleType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name, nameOK := p.Args["name"].(string)
					if !nameOK {
						return nil, fmt.Errorf("No query key specified")
					}
					result := service.GetRoleByName(name)
					if !result.Status {
						return nil, fmt.Errorf(result.Msg)
					} else {
						return result.Data, nil
					}
				},
			},
			"permission": &graphql.Field{
				Type: permissionType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name, nameOK := p.Args["name"].(string)
					if !nameOK {
						return nil, fmt.Errorf("No query key specified")
					}
					result := service.GetPermission(name)
					if !result.Status {
						return nil, fmt.Errorf(result.Msg)
					} else {
						return result.Data, nil
					}
				},
			},
		},
	})

// var mutationType = graphql.NewObject(graphql.ObjectConfig{
//     Name: "Mutation",
//     Fields: graphql.Fields{
//         /* Create new product item
//         http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
//         */
//         "create": &graphql.Field{
//             Type:        productType,
//             Description: "Create new product",
//             Args: graphql.FieldConfigArgument{
//                 "name": &graphql.ArgumentConfig{
//                     Type: graphql.NewNonNull(graphql.String),
//                 },
//                 "info": &graphql.ArgumentConfig{
//                     Type: graphql.String,
//                 },
//                 "price": &graphql.ArgumentConfig{
//                     Type: graphql.NewNonNull(graphql.Float),
//                 },
//             },
//             Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//                 rand.Seed(time.Now().UnixNano())
//                 product := Product{
//                     ID:    int64(rand.Intn(100000)), // generate random ID
//                     Name:  params.Args["name"].(string),
//                     Info:  params.Args["info"].(string),
//                     Price: params.Args["price"].(float64),
//                 }
//                 products = append(products, product)
//                 return product, nil
//             },
//         },

//         /* Update product by id
//            http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
//         */
//         "update": &graphql.Field{
//             Type:        productType,
//             Description: "Update product by id",
//             Args: graphql.FieldConfigArgument{
//                 "id": &graphql.ArgumentConfig{
//                     Type: graphql.NewNonNull(graphql.Int),
//                 },
//                 "name": &graphql.ArgumentConfig{
//                     Type: graphql.String,
//                 },
//                 "info": &graphql.ArgumentConfig{
//                     Type: graphql.String,
//                 },
//                 "price": &graphql.ArgumentConfig{
//                     Type: graphql.Float,
//                 },
//             },
//             Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//                 id, _ := params.Args["id"].(int)
//                 name, nameOk := params.Args["name"].(string)
//                 info, infoOk := params.Args["info"].(string)
//                 price, priceOk := params.Args["price"].(float64)
//                 product := Product{}
//                 for i, p := range products {
//                     if int64(id) == p.ID {
//                         if nameOk {
//                             products[i].Name = name
//                         }
//                         if infoOk {
//                             products[i].Info = info
//                         }
//                         if priceOk {
//                             products[i].Price = price
//                         }
//                         product = products[i]
//                         break
//                     }
//                 }
//                 return product, nil
//             },
//         },

//         /* Delete product by id
//            http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
//         */
//         "delete": &graphql.Field{
//             Type:        productType,
//             Description: "Delete product by id",
//             Args: graphql.FieldConfigArgument{
//                 "id": &graphql.ArgumentConfig{
//                     Type: graphql.NewNonNull(graphql.Int),
//                 },
//             },
//             Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//                 id, _ := params.Args["id"].(int)
//                 product := Product{}
//                 for i, p := range products {
//                     if int64(id) == p.ID {
//                         product = products[i]
//                         // Remove from product list
//                         products = append(products[:i], products[i+1:]...)
//                     }
//                 }

//                 return product, nil
//             },
//         },
//     },
// })

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
		// Mutation: mutationType,
	},
)

func GetGraphQL(ctx iris.Context) {
	query := ctx.URLParam("query")
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	var response *model.ApiJson
	if result.HasErrors() {
		var errs []error
		for _, err := range result.Errors {
			errs = append(errs, err.OriginalError())
		}
		response = model.ErrorQueryDatabase(errs...)
	} else {
		response = model.Success(result.Data, "查询成功")
	}
	ctx.StatusCode(response.Code)
	ctx.JSON(response)
}
