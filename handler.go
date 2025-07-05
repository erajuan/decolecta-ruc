package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the SUNAT RUC API Microservice!\n" +
		"Available endpoints:\n" +
		"- GET /sunat/ruc?numero=<RUC> - Get basic company information by RUC\n" +
		"- GET /sunat/ruc/full?numero=<RUC> - Get detailed company information by RUC\n")
}

func GetCompanyHandler(c *fiber.Ctx) error {
	var ruc string = c.Query("numero")
	if err := IsValidRuc(ruc); err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "ruc no valido"})
	}
	company, err := GetCompanyService(ruc)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(company)
}

func GetCompanyAdvanceHandler(c *fiber.Ctx) error {
	var ruc string = c.Query("numero")
	if err := IsValidRuc(ruc); err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "ruc no valido"})
	}
	company, err := GetCompanyAdvanceService(ruc)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(company)
}

func Pro5RucHandler(c *fiber.Ctx) error {
	var ruc string = c.Params("numero")
	if err := IsValidRuc(ruc); err != nil {
		return c.JSON(
			fiber.Map{
				"success": false,
				"message": "No encontrado",
			},
		)
	}
	company, err := GetCompanyService(ruc)
	if err != nil {
		log.Println(err)
		return c.JSON(
			fiber.Map{
				"success": false,
				"message": "No encontrado",
			},
		)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    company.ToPro5(),
	})

}
