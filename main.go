package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/danmcfan/nexus/internal"
	"github.com/danmcfan/nexus/internal/components"
	"github.com/danmcfan/nexus/internal/database"
)

//go:embed assets
var embeddedFiles embed.FS

const (
	numClients    = 100
	numUsers      = 200
	numProperties = 20_000
	pageSize      = 1000
)

func main() {
	log.Println("Nexus is running in", internal.Version, "mode...")

	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	for _, file := range []string{"schema/client.sql", "schema/metadata.sql", "schema/property.sql", "schema/user.sql"} {
		schema, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(string(schema))
		if err != nil {
			log.Fatal(err)
		}
	}

	clients := make([]database.Client, numClients)
	clientIDs := make(map[string]bool)
	for i := range numClients {
		clientID := randomString(3)
		_, ok := clientIDs[clientID]
		for ok {
			clientID = randomString(3)
			_, ok = clientIDs[clientID]
		}
		clientIDs[clientID] = true
		client, err := queries.CreateClient(ctx, database.CreateClientParams{
			PkClientID: clientID,
			Name:       randomString(16),
		})
		if err != nil {
			log.Fatal(err)
		}
		clients[i] = client
	}

	users := make([]database.User, numUsers)
	for i := range numUsers {
		user, err := queries.CreateUser(ctx, database.CreateUserParams{
			PkUserID:  uuid.New().String(),
			FirstName: randomString(16),
			LastName:  randomString(16),
		})
		if err != nil {
			log.Fatal(err)
		}
		users[i] = user
	}

	properties := make([]database.Property, numProperties)
	for i := range numProperties {
		fkPointOfContactID := sql.NullString{}
		if rand.Intn(10) == 0 {
			fkPointOfContactID = sql.NullString{
				String: users[rand.Intn(len(users))].PkUserID,
				Valid:  true,
			}
		}

		fkManagerID := sql.NullString{}
		if rand.Intn(10) == 0 {
			fkManagerID = sql.NullString{
				String: users[rand.Intn(len(users))].PkUserID,
				Valid:  true,
			}
		}

		property, err := queries.CreateProperty(ctx, database.CreatePropertyParams{
			PkPropertyID:       uuid.New().String(),
			Name:               randomString(16),
			Address:            randomString(16),
			IsDemo:             rand.Intn(2) == 0,
			FkPointOfContactID: fkPointOfContactID,
			FkManagerID:        fkManagerID,
			FkClientID:         clients[rand.Intn(len(clients))].PkClientID,
		})
		if err != nil {
			log.Fatal(err)
		}
		properties[i] = property
	}

	router := gin.Default()

	if internal.Version == "production" {
		log.Println("Production mode")
		router.StaticFS("/public", http.FS(embeddedFiles))
	} else {
		log.Println("Development mode")
		router.Static("/public/assets", "./assets")
	}

	router.GET("/", func(c *gin.Context) {
		components.HTML().Render(ctx, c.Writer)
	})

	router.GET("/properties/", func(c *gin.Context) {
		filter := c.Query("filter")
		pageParam := c.Query("page")
		if pageParam == "" {
			pageParam = "0"
		}
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			log.Fatal(err)
		}
		properties, err := queries.ListPropertiesWithFilter(ctx, database.ListPropertiesWithFilterParams{
			Name:    "%" + filter + "%",
			Address: "%" + filter + "%",
			Offset:  int64(page * pageSize),
			Limit:   pageSize,
		})
		if err != nil {
			log.Println(err)
		}

		if len(properties) == 0 {
			components.Rows(properties, filter, 0).Render(ctx, c.Writer)
			return
		}
		nextPage := page + 1
		if len(properties) < pageSize {
			nextPage = 0
		}

		components.Rows(properties, filter, nextPage).Render(ctx, c.Writer)
	})

	port := ":8080"
	if internal.Version == "production" {
		port = ":80"
	}
	router.Run(port)
}

func randomString(n int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
