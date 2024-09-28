package weather

import (
    "github.com/graphql-go/graphql"
)

// Define LocationType GraphQL object
var LocationType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Location",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "latitude": &graphql.Field{
                Type: graphql.String,
            },
            "longitude": &graphql.Field{
                Type: graphql.String,
            },
        },
    },
)

// RootQuery definition
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
    Name: "RootQuery",
    Fields: graphql.Fields{
        "locations": &graphql.Field{
            Type:        graphql.NewList(LocationType),
            Description: "Get all locations",
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                db := params.Context.Value("db").(Database)
                lc := params.Context.Value("lc").(LocationController)

                locations, err := lc.GetLocations(db)
                if err != nil {
                    return nil, err
                }

                return locations, nil
            },
        },
    },
})

	// RootMutation definition
	var RootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"addLocation": &graphql.Field{
				Type:        LocationType,
				Description: "Add a new location",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"latitude": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"longitude": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db := params.Context.Value("db").(Database)
					lc := params.Context.Value("lc").(LocationController)

					var name string
					if nameArg, ok := params.Args["name"]; ok && nameArg != nil {
						name = nameArg.(string)
					} else {
						name = ""
					}
					latitude := params.Args["latitude"].(string)
					longitude := params.Args["longitude"].(string)

					// Create the location object
					location := Location{
						Name:      name,
						Latitude:  latitude,
						Longitude: longitude,
					}

					err := lc.AddLocation(db, location)
					if err != nil {
						return nil, err
					}
	
					return location, nil
				},
			},
			"deleteLocation": &graphql.Field{
				Type:        LocationType,
				Description: "Delete a location by ID",
				Args: graphql.FieldConfigArgument{
					"latitude": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"longitude": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db := params.Context.Value("db").(Database)
					lc := params.Context.Value("lc").(LocationController)

					latitude := params.Args["latitude"].(string)
					longitude := params.Args["longitude"].(string)

					err := lc.DeleteLocation(db, latitude, longitude)

					if err != nil {
						return nil, err
					}
					return nil, nil
				},
			},
		},
	})


// Define the GraphQL schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: RootQuery,
	Mutation: RootMutation,
})
