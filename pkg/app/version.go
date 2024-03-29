package app

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/womat/debug"
)

// VERSION holds the version information with the following logic in mind
//  1 ... fixed
//  0 ... year 2020, 1->year 2021, etc.
//  7 ... month of year (7=July)
//  the date format after the + is always the first of the month
//
// VERSION differs from semantic versioning as described in https://semver.org/
// but we keep the correct syntax.
//TODO: increase version number to 1.0.1+2020xxyy
const (
	VERSION = "3.0.2+20210909"
	MODULE  = "s0counter"
)

// HandleVersion is the get application version web handler.
func (app *App) HandleVersion() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		debug.InfoLog.Print("web request version")

		return ctx.JSON(fiber.Map{
			"version":     VERSION,
			"description": MODULE,
			"about":       Version(),
		})
	}
}

// Version is the get application version as string.
func Version() string {
	return strings.TrimSpace(MODULE + " V" + strings.Split(VERSION, "+")[0])
}
