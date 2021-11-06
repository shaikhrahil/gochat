package main

import (
	"chatterbox/registration"
	"log"
	"os"

	"github.com/gofiber/template/html"
	"github.com/tkanos/gonfig"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
)

type Configuration struct {
	REDIS_URL string
	REDIS_PWD string
}

func main() {
	configuration := Configuration{}
	if err := gonfig.GetConf("env.json", &configuration); err != nil {
		// log.Fatalln(err.Error())
		log.Println("env.json not found")
	}

	if configuration.REDIS_URL == "" {
		configuration.REDIS_URL = os.Getenv("REDIS_URL")
	}

	if configuration.REDIS_PWD == "" {
		configuration.REDIS_PWD = os.Getenv("REDIS_PWD")
	}

	redisDB := redis.NewClient(&redis.Options{
		Addr:     configuration.REDIS_URL,
		Password: configuration.REDIS_PWD,
	})

	log.Println("Connecting to redis ..")

	connected := redisDB.Set("Connected", true, 0)

	if connected.Err() != nil {
		log.Fatal("Failed to connect to redis due to following error : \n", connected.Err())
	}

	log.Println("Starting server ..")
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/ws/:id", registration.Register(redisDB))
	// app.Static("/", "static/public")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Get("/chat", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Listen(":3000")
}
