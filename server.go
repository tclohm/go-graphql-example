package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Venue struct {
	ID int
	Name string
	Owners []Patron
}

type Patron struct {
	ID int
	Firstname string
	Lastname string
}


func populate() []Venue {
	patron := &Patron{ ID:1, Firstname: "Taylor", Lastname: "Lohman" }
	venue := Venue{
		ID: 1,
		Name: "Jones Hollywood",
		Owners: []Patron{*patron},
	}
	var venues []Venue
	venues = append(venues, venue)

	return venues
}

func main() {
	fmt.Println("---starting---")
	venues := populate()
	// ID int
	// Firstname string
	// Lastname string
	// Owned []Venue
	// Visited []Venue

	var patronType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Patron",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"Firstname": &graphql.Field{
					Type: graphql.String,
				},
				"Lastname": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	// ID int
	// Name string
	// Owner Patron
	// Visitors []Patron

	var venueType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Venue",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Owners": &graphql.Field{
					Type: graphql.NewList(patronType),
				},
			},
		},
	)


	// fields := graphql.Fields{
	// 	"hello": &graphql.Field{
	// 		Type: graphql.String,
	// 		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
	// 			return "World", nil
	// 		},
	// 	},
	// }

	fields := graphql.Fields{
		"venue": &graphql.Field{
			Type: venueType,
			Description: "Get Venue by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int) // type cast to int
				if ok {
					for _, venue := range venues {
						if int(venue.ID) == id {
							return venue, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type: graphql.NewList(venueType),
			Description: "Get Full Venue List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return venues, nil
			},
		},
	}

	// define object config
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// define schema config
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	// create schema
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("Failed to create Graphql Schema, err %v", err)
	}

	// query := `
	// 	{
	// 		list {
	// 			id
	// 			Name
	// 		}
	// 	}
	// `
	query := `
	{
		venue(id: 1) {
			Name
			Owners {
				Firstname
				Lastname
			}
		}
	}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute graphql operation, errors: %+v", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}