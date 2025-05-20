package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"math/rand"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/danmcfan/nexus/internal"
	"github.com/danmcfan/nexus/internal/components"
	"github.com/danmcfan/nexus/internal/database"
)

//go:embed assets schema
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
		schema, err := embeddedFiles.ReadFile(file)
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
			Name:       randomCompany(),
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
			FirstName: randomFirstName(),
			LastName:  randomLastName(),
		})
		if err != nil {
			log.Fatal(err)
		}
		users[i] = user
	}

	properties := make([]database.Property, numProperties)
	propertyIDs := make(map[string]bool)
	for i := range numProperties {
		propertyID := randomString(6)
		_, ok := propertyIDs[propertyID]
		for ok {
			propertyID = randomString(6)
			_, ok = propertyIDs[propertyID]
		}
		propertyIDs[propertyID] = true
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
			PkPropertyID:       propertyID,
			Name:               randomPropertyName(),
			Address:            randomAddress(),
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
		if err != nil {
			log.Fatal(err)
		}
		properties, err := queries.ListPropertiesWithFilter(ctx, database.ListPropertiesWithFilterParams{
			PkPropertyID: "%" + filter + "%",
			Name:         "%" + filter + "%",
			Address:      "%" + filter + "%",
		})
		if err != nil {
			log.Println(err)
		}

		propertiesByClient := make(map[string][]database.ListPropertiesWithFilterRow)
		for _, property := range properties {
			propertiesByClient[property.ClientName.String] = append(propertiesByClient[property.ClientName.String], property)
		}

		clients := make([]string, 0)
		for k, _ := range propertiesByClient {
			clients = append(clients, k)
		}
		sort.Strings(clients)

		components.GroupedRows(clients, propertiesByClient).Render(ctx, c.Writer)
	})

	port := ":8080"
	if internal.Version == "production" {
		port = ":80"
	}
	router.Run(port)
}

func randomCompany() string {
	companies := []string{
		"Acme, Inc.",
		"Globex Corporation",
		"Initech",
		"Soylent Corp",
		"Wayne Enterprises",
		"LexCorp",
		"Daily Planet",
		"Daily Bugle",
	}
	return companies[rand.Intn(len(companies))]
}

func randomPropertyName() string {
	names := []string{
		"First Avenue",
		"Second Street",
		"Third Drive",
		"Fourth Boulevard",
		"Fifth Court",
		"Sixth Place",
		"Seventh Lane",
		"Eighth Avenue",
		"Ninth Street",
	}
	return names[rand.Intn(len(names))]
}

func randomAddress() string {
	addresses := []string{
		"123 Main St",
		"456 Oak Ave",
		"789 Pine Rd",
	}
	return addresses[rand.Intn(len(addresses))]
}

func randomFirstName() string {
	firstNames := []string{
		"John",
		"Jane",
		"Jim",
		"Jill",
		"Jack",
	}
	return firstNames[rand.Intn(len(firstNames))]
}

func randomLastName() string {
	lastNames := []string{
		"Smith",
		"Johnson",
		"Williams",
		"Jones",
		"Brown",
		"Davis",
		"Miller",
		"Wilson",
	}
	return lastNames[rand.Intn(len(lastNames))]
}

func randomString(n int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
