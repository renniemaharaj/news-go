package router

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

var (
	portForwardPage = "http://192.168.0.1/?firewall_virt"
	routerUser      string
	routerPass      string

	l = createLogger()
)

func Initialize() {
	if !validateEnv() {
		l.Error("Missing required environment variables. Please check your .env file.")
		os.Exit(1)
	}
}

// Validate environment variables
func validateEnv() bool {
	if err := godotenv.Load(); err != nil {
		l.Error("No .env file found, or error loading it: " + err.Error())
	}

	missing := []string{}
	portForwardPage = os.Getenv("PORT_FORWARDING_PAGE")
	if portForwardPage == "" {
		missing = append(missing, "PORT_FORWARDING_PAGE")
	}
	routerUser = os.Getenv("ROUTER_USER")
	if routerUser == "" {
		missing = append(missing, "ROUTER_USER")
	}
	routerPass = os.Getenv("ROUTER_PASS")
	if routerPass == "" {
		missing = append(missing, "ROUTER_PASS")
	}
	if len(missing) > 0 {
		l.Error("Missing required environment variables: " + strings.Join(missing, ", "))
		return false
	}
	return true
}

func LoginToRouter() {
	// Launch a visible browser for debugging (headless if desired)
	url := launcher.New().
		// Set headless to false if you want to see what's going on
		Headless(false).
		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	browser.Context(ctx)
	defer cancel()

	page := browser.MustPage(portForwardPage)
	page.MustWaitLoad()

	// Optional wait for router page to fully render
	time.Sleep(2 * time.Second)

	// Check if the login form is present
	if page.MustElementR("h3", "Login") != nil {
		l.Success("Login page detected")
		time.Sleep(2 * time.Second)
		page.MustElement("#UserName").MustInput(routerUser)
		page.MustElement("#Password").MustInput(routerPass)
		page.MustElement("#ApplyButton").MustClick()

		l.Success("Submitted login credentials")
	} else {
		l.Success("Already authenticated or login not required")
	}

	// Wait for redirect and ensure we're on the port forwarding page
	page.MustWaitLoad()
	page.MustElementR("h3", "Virtual Servers / Port Forwarding")

	l.Success("Successfully reached the Port Forwarding configuration page")
}
