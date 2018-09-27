package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/version"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	calendar "google.golang.org/api/calendar/v3"
	sheets "google.golang.org/api/sheets/v4"
)

// Config is the structure used to instantiate the Google Calendar provider.
type config struct {
	Credentials string
	Scopes      []string

	Calendar *calendar.Service
	Sheets   *sheets.Service
	//Directory *directory.Service

	StopContext context.Context

	client      *http.Client
	userAgent   string
	tokenSource oauth2.TokenSource
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *config) loadAndValidate(ctx context.Context) error {
	var account accountFile
	var client *http.Client
	var tokenSource oauth2.TokenSource

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return fmt.Errorf("Error loading credentials: %s", err)
		}

		// Assume account_file is a JSON string
		if err := parseJSON(&account, contents); err != nil {
			return fmt.Errorf("Error parsing credentials '%s': %s", contents, err)
		}

		// Get the token for use in our requests
		log.Printf("[INFO] Requesting Google token...")
		log.Printf("[INFO]   -- Email: %s", account.ClientEmail)
		log.Printf("[INFO]   -- Scopes: %s", c.Scopes)
		log.Printf("[INFO]   -- Private Key Length: %d", len(account.PrivateKey))

		conf := jwt.Config{
			Email:      account.ClientEmail,
			PrivateKey: []byte(account.PrivateKey),
			Scopes:     c.Scopes,
			TokenURL:   "https://accounts.google.com/o/oauth2/token",
		}

		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client = conf.Client(ctx)

		tokenSource = conf.TokenSource(ctx)
	} else {
		log.Printf("[INFO] Authenticating using DefaultClient")
		err := error(nil)
		client, err = google.DefaultClient(ctx, c.Scopes...)
		if err != nil {
			return err
		}

		tokenSource, err = google.DefaultTokenSource(ctx, c.Scopes...)
		if err != nil {
			return err
		}
	}

	c.tokenSource = tokenSource

	client.Transport = logging.NewTransport("Google", client.Transport)

	projectURL := "https://www.terraform.io"
	userAgent := fmt.Sprintf("Terraform/%s (+%s)",
		version.String(), projectURL)

	c.client = client
	c.userAgent = userAgent

	// Create the calendar service.
	calendarSvc, err := calendar.New(client)
	if err != nil {
		return nil
	}
	calendarSvc.UserAgent = userAgent
	c.Calendar = calendarSvc

	// Create the sheets service
	sheetsSvc, err := sheets.New(client)
	if err != nil {
		return nil
	}
	sheetsSvc.UserAgent = userAgent
	c.Sheets = sheetsSvc

	// // Create the directory service.
	// directorySvc, err := directory.New(client)
	// if err != nil {
	// 	return nil
	// }
	// directorySvc.UserAgent = userAgent
	// c.directory = directorySvc

	return nil
}

// accountFile represents the structure of the account file JSON file.
type accountFile struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
}

func parseJSON(result interface{}, contents string) error {
	r := strings.NewReader(contents)
	dec := json.NewDecoder(r)

	return dec.Decode(result)
}
