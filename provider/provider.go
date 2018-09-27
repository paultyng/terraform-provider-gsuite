package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	directory "google.golang.org/api/admin/directory/v1"
	calendar "google.golang.org/api/calendar/v3"
	sheets "google.golang.org/api/sheets/v4"
)

var defaultScopes = []string{
	calendar.CalendarScope,
	directory.AdminDirectoryGroupScope,
	directory.AdminDirectoryUserScope,
	sheets.SpreadsheetsReadonlyScope,
}

// New returns the actual provider instance.
func New() *schema.Provider {
	var p *schema.Provider

	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			// this comes from the google provider, 0790c96cac61142042b6cc9033c74bb4cd1e562e
			"credentials": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				ValidateFunc: validateCredentials,
			},
			"scopes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
		// DatasourcesMap: map[string]*schema.Resource{
		// 	"gsuite_sheets_"
		// },
		ResourcesMap: map[string]*schema.Resource{
			"gsuite_calendar_event": resourceCalendarEvent(),

			// TODO: look at combining with these from DeviaVir/terraform-provider-gsuite
			// "gsuite_directory_group":         resourceGroup(),
			// "gsuite_directory_user":          resourceUser(),
			// "gsuite_directory_group_member":  resourceGroupMember(),
			// "gsuite_directory_group_members": resourceGroupMembers(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func scopesFromConfigOrDefault(scopesSet *schema.Set) []string {
	scopes := convertStringSet(scopesSet)
	if len(scopes) == 0 {
		log.Printf("[INFO] No scopes provided. Using default scopes.")
		scopes = defaultScopes
	}
	return scopes
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		ctx := context.Background()
		scopes := scopesFromConfigOrDefault(d.Get("scopes").(*schema.Set))
		config := &config{
			Credentials: d.Get("credentials").(string),
			Scopes:      scopes,

			StopContext: p.StopContext(),
		}

		if err := config.loadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}

		return config, nil
	}
}

// this comes from the google provider, 0790c96cac61142042b6cc9033c74bb4cd1e562e
func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err == nil {
		return
	}
	var account accountFile
	if err := json.Unmarshal([]byte(creds), &account); err != nil {
		errors = append(errors,
			fmt.Errorf("credentials are not valid JSON '%s': %s", creds, err))
	}

	return
}

func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	return s
}
