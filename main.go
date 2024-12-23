package main

import (
	"os"

	"github.com/atotto/clipboard"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	groq "github.com/hasitpbhatt/groq-go"
	"github.com/shreeyash-ugale/go-cheat-ml/pkg"
)

func main() {
	engine := html.New(".", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/send", htmlForm)

	app.Post("/form", formHandler)

	app.Get("/clip", clipHandler)

	app.Get("/ai", aiHandler)
	app.Listen(":3000")
}

func htmlForm(c *fiber.Ctx) error {
	return c.SendFile("form.html")
}

func aiHandler(c *fiber.Ctx) error {
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		return c.Status(500).SendString("Failed to get clipboard content")
	}
	client := groq.NewClient(groq.WithAPIKey(os.Getenv("GROQ_API_KEY")))
	resp, err := client.ChatCompletion([]groq.Message{
		{
			Content: pkg.Prompt(),
			Role:    "system",
		},
		{
			Content: clipboardContent,
			Role:    "user",
		},
	})
	if err != nil {
		return c.Status(500).SendString("Failed to chat with groq cloud")
	}
	data := ""
	for _, c := range resp.Choices {
		data += (c.Message.Content)
	}
	return c.Render("disp", fiber.Map{
		"clipboardContent": data,
	})

}

func formHandler(c *fiber.Ctx) error {
	data := c.FormValue("data")
	err := clipboard.WriteAll(data)
	if err != nil {
		return c.Status(500).SendString("Failed to write to clipboard")

	}
	return c.SendString("Received data: " + data)
}

func clipHandler(c *fiber.Ctx) error {
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		return c.Status(500).SendString("Failed to get clipboard content")
	}
	return c.Render("disp", fiber.Map{
		"clipboardContent": clipboardContent,
	})
}
