package provider

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	sheets "google.golang.org/api/sheets/v4"
)

const (
	calendarScope       = calendar.CalendarScope
	sheetsReadonlyScope = sheets.SpreadsheetsReadonlyScope
)

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
	sheets   *sheets.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate() error {
	// TODO: dynamically determine scopes to use
	oauthScopes := []string{
		calendarScope,
		sheetsReadonlyScope,
	}

	log.Printf("[INFO] authenticating with local client")
	client, err := google.DefaultClient(context.Background(), oauthScopes...)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	// Use a custom user-agent string. This helps google with analytics and it's
	// just a nice thing to do.
	client.Transport = logging.NewTransport("Google", client.Transport)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraform.VersionString())

	// Create the calendar service.
	calendarSvc, err := calendar.New(client)
	if err != nil {
		return nil
	}
	calendarSvc.UserAgent = userAgent
	c.calendar = calendarSvc

	// Create the sheets service
	sheetsSvc, err := sheets.New(client)
	if err != nil {
		return nil
	}
	sheetsSvc.UserAgent = userAgent
	c.sheets = sheetsSvc

	return nil
}
