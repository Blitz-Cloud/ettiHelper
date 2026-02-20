package middleware

import (
	"strings"

	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/Blitz-Cloud/ettiHelper/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func IsValidTenant(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok {
		utils.Log.Fatal("Error accessing db con from admin route")
	}
	// spew.Dump(c.GetReqHeaders())
	fullHostname := c.GetReqHeaders()["Host"][0]
	hostname := fullHostname
	if strings.Contains(fullHostname, "api.") {
		hostname = strings.Split(fullHostname, "api.")[1]
	}

	if ok, tenant := IsValidDomain(db, hostname); ok == true {
		c.Locals("tenant", tenant)
		return c.Next()
	}
	return c.Status(fiber.StatusNotFound).SendString("Tenant not found")
}

func IsValidDomain(db *gorm.DB, origin string) (bool, types.Namespace) {

	domain := types.Domain{}
	if strings.Contains(origin, "http://") {
		origin = strings.Split(origin, "http://")[1]
	}

	if strings.Contains(origin, "https://") {
		origin = strings.Split(origin, "https://")[1]
	}
	if err := db.First(&domain, "value = ?", origin).Error; err != nil {
		return false, types.Namespace{}
	}

	namespace := types.Namespace{}
	if err := db.First(&namespace, "id = ?", domain.NamespaceId).Error; err != nil {
		return false, types.Namespace{}
	}
	return true, namespace
}

// func TenantRouteProtector(c *fiber.Ctx) error {
// }
