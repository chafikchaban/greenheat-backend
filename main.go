
package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/graphql-go/handler"
    "github.com/chafikchaban/greenheat-backend/weather"
  )

func main() {
    r := gin.Default()
    lc := weather.LocationController{}
	db := weather.BootstrapDatabase("./")

    lc.InitializeLocations(db)

    // Create GraphQL handler
    h := handler.New(&handler.Config{
        Schema:   &weather.Schema,
        Pretty:   true,
        GraphiQL: true,
    })

    // GraphQL endpoint
    r.POST("/graphql", func(c *gin.Context) {
        ctx := c.Request.Context()
        ctx = context.WithValue(ctx, "db", db)
        ctx = context.WithValue(ctx, "lc", lc)
        h.ContextHandler(ctx, c.Writer, c.Request)
    })

    // GraphiQL endpoint for testing
    r.GET("/graphiql", func(c *gin.Context) {
        h.ContextHandler(c.Request.Context(), c.Writer, c.Request)
    })

    r.Run()
}