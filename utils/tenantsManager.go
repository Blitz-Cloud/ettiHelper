package utils

import (
	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
)

// a struct that will serve as holeder for the rest of the functions
// it is design to make easier the management of tenants(namespaces) and their associated middleware functions
type TenantsManager struct {
	fiberApp *fiber.App
	Store    map[string][]fiber.Handler
}

func InitTenant(app *fiber.App) TenantsManager {
	return TenantsManager{
		fiberApp: app,
		Store:    map[string][]fiber.Handler{},
	}
}

// this function does almost nothing important i made it in case in the feature i will want to also register routes under a specific domain or subdomain
func (tenantManager *TenantsManager) RegisterRouter(router func(*fiber.App)) {
	router(tenantManager.fiberApp)
}

func (tenantManager *TenantsManager) RegisterMiddleware(tenantName string, middleware fiber.Handler) {
	if _, ok := tenantManager.Store[tenantName]; !ok {
		tenantManager.Store[tenantName] = append(tenantManager.Store[tenantName], middleware)
	}
	Log.Info("Middleware for %s already exists", tenantName)
}

func (tenantManager *TenantsManager) TenantMiddlewareDispatcher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenant, ok := c.Locals("tenant").(types.Namespace)
		if !ok {
			Log.Fatal("Failed to access tenant inside middleware dispatcher")
		}
		tenantName := tenant.AuthFlow
		spew.Dump("Dispatcher")
		spew.Dump(tenant)
		if hanlers, ok := tenantManager.Store[tenantName]; ok {
			for _, handler := range hanlers {
				return handler(c)
			}
		}
		return c.Next()
	}
}
