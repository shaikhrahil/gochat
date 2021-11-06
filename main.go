package main

import (
	"chatterbox/registration"
	"log"

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
	// app.Static("/", "static/public")
	app.Get("/chat", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Get("/ws/:id", registration.Register(redisDB))

	app.Listen(":3000")
}
