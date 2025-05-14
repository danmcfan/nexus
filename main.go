package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/danmcfan/nexus/internal/database"
)

const (
	numClients    = 100
	numUsers      = 200
	numProperties = 20_000
	pageSize      = 1000
)

func main() {
	log.Println("Nexus is running...")

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
	router.Static("/static", "./static")
	router.LoadHTMLGlob("./static/*.html")

	router.GET("/", func(c *gin.Context) {
		properties, err := queries.ListPropertiesWithFilter(ctx, database.ListPropertiesWithFilterParams{
			Offset: 0,
			Limit:  pageSize,
		})
		if err != nil {
			log.Println(err)
		}

		if len(properties) == 0 {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Properties": properties,
				"NextPage":   nil,
				"Last":       nil,
			})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Properties": properties[0 : len(properties)-1],
			"NextPage":   "1",
			"Last":       properties[len(properties)-1],
		})
	})

	router.GET("/properties/", func(c *gin.Context) {
		filter := c.Query("filter")
		page := c.Query("page")
		if page == "" {
			page = "0"
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			log.Fatal(err)
		}
		properties, err := queries.ListPropertiesWithFilter(ctx, database.ListPropertiesWithFilterParams{
			Name:    "%" + filter + "%",
			Address: "%" + filter + "%",
			Offset:  int64(pageInt * pageSize),
			Limit:   pageSize,
		})
		if err != nil {
			log.Println(err)
		}

		if len(properties) == 0 {
			c.HTML(http.StatusOK, "rows.html", gin.H{
				"Properties": properties,
				"NextPage":   nil,
				"Last":       nil,
				"Filter":     filter,
			})
			return
		}
		nextPage := strconv.Itoa(pageInt + 1)
		if len(properties) < pageSize {
			nextPage = ""

		}

		if len(properties) < 10 {
			for _, property := range properties {
				log.Println(property)
			}
		}

		c.HTML(http.StatusOK, "rows.html", gin.H{
			"Properties": properties[0 : len(properties)-1],
			"NextPage":   nextPage,
			"Last":       properties[len(properties)-1],
			"Filter":     filter,
		})
	})

	router.Run(":8080")
}

func randomString(n int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
