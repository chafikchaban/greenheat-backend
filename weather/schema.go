package weather

import (
	"fmt"
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

var HourlyDataType = graphql.NewObject(graphql.ObjectConfig{
    Name: "HourlyData",
    Fields: graphql.Fields{
        "time": &graphql.Field{
            Type: graphql.NewList(graphql.String),
        },
        "temperature2m": &graphql.Field{
            Type: graphql.NewList(graphql.Float),
        },
        "cloudCover": &graphql.Field{
            Type: graphql.NewList(graphql.Int),
        },
        "windSpeed80m": &graphql.Field{
            Type: graphql.NewList(graphql.Float),
        },
        "uvIndex": &graphql.Field{
            Type: graphql.NewList(graphql.Float),
        },
    },
})

var DailyDataType = graphql.NewObject(graphql.ObjectConfig{
    Name: "DailyData",
    Fields: graphql.Fields{
        "time": &graphql.Field{
            Type: graphql.NewList(graphql.String),
        },
        "temperature2mMax": &graphql.Field{
            Type: graphql.NewList(graphql.Float),
        },
        "temperature2mMin": &graphql.Field{
            Type: graphql.NewList(graphql.Float),
        },
    },
})

var UnitsType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Units",
    Fields: graphql.Fields{
        "time": &graphql.Field{
            Type: graphql.String,
        },
        "temperature2m": &graphql.Field{
            Type: graphql.String,
        },
        "temperature2mMax": &graphql.Field{
            Type: graphql.String,
        },
        "temperature2mMin": &graphql.Field{
            Type: graphql.String,
        },
        "cloud_cover": &graphql.Field{
            Type: graphql.String,
        },
        "wind_speed_80m": &graphql.Field{
            Type: graphql.String,
        },
        "uv_index": &graphql.Field{
            Type: graphql.String,
        },
    },
})

var WeatherInfoType = graphql.NewObject(graphql.ObjectConfig{
    Name: "WeatherInfo",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.String,
        },
        "locationName": &graphql.Field{
            Type: graphql.String,
        },
        "latitude": &graphql.Field{
            Type: graphql.String,
        },
        "longitude": &graphql.Field{
            Type: graphql.String,
        },
        "temperature": &graphql.Field{
            Type: graphql.Float,
        },
        "maxTemperature": &graphql.Field{
            Type: graphql.Float,
        },
        "minTemperature": &graphql.Field{
            Type: graphql.Float,
        },
        "cloudCoverage": &graphql.Field{
            Type: graphql.Float,
        },
        "windSpeed": &graphql.Field{
            Type: graphql.Float,
        },
        "uvIndex": &graphql.Field{
            Type: graphql.Float,
        },
        "daily": &graphql.Field{
            Type: DailyDataType,
        },
        "hourly": &graphql.Field{
            Type: HourlyDataType,
        },
        "weatherCode": &graphql.Field{
            Type: graphql.Int,
        },
        "units": &graphql.Field{
            Type: UnitsType,
        },
        "dailyUnits": &graphql.Field{
            Type: UnitsType,
        },
        "hourlyUnits": &graphql.Field{
            Type: UnitsType,
        },
    },
})

// RootQuery definition
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
    Name: "RootQuery",
    Fields: graphql.Fields{
        "locations": &graphql.Field{
            Type: graphql.NewList(LocationType),
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
        "WeatherForecast": &graphql.Field{
            Type: WeatherInfoType,
            Description: "Get weather forecast for a specific location",
            Args: graphql.FieldConfigArgument{
                "locationID": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
            },
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                db := params.Context.Value("db").(Database)
                wc := params.Context.Value("wc").(WeatherController)

                locationID := params.Args["locationID"].(interface{})

                var location Location
                if err := db.d.Read("locations", locationID.(string), &location); err != nil {
                    return nil, fmt.Errorf("location not found: %v", err)
                }

                weatherData, err := wc.FetchWeatherForecast(db, location)
                if err != nil {
                    return nil, fmt.Errorf("could not fetch weather data: %v", err)
                }
                return weatherData, nil
            },
        },
		"weatherForLocations": &graphql.Field{
            Type: graphql.NewList(WeatherInfoType),
            Description: "Get current weather for a list of locations",
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                db := params.Context.Value("db").(Database)
                wc := params.Context.Value("wc").(WeatherController)
                lc := params.Context.Value("lc").(LocationController)

                weatherData, err := wc.FetchWeatherForLocations(db, lc)
                if err != nil {
                    return nil, fmt.Errorf("could not fetch weather data: %v", err)
                }
                return weatherData, nil
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
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					db := params.Context.Value("db").(Database)
					lc := params.Context.Value("lc").(LocationController)

					id := params.Args["id"].(string)

					err := lc.DeleteLocation(db, id)

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
