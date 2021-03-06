package main

import (
	"fmt"
	"strings"

	contentful "github.com/contentful-labs/contentful-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContentfulAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateAPIKey,
		Read:   resourceReadAPIKey,
		Update: resourceUpdateAPIKey,
		Delete: resourceDeleteAPIKey,
		Importer: &schema.ResourceImporter{
			State: resourceAPIKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"space_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)

	apiKey := &contentful.APIKey{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	err = client.APIKeys.Upsert(d.Get("space_id").(string), apiKey)
	if err != nil {
		return err
	}

	if err := setAPIKeyProperties(d, apiKey); err != nil {
		return err
	}

	d.SetId(apiKey.Sys.ID)

	return nil
}

func resourceUpdateAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)
	apiKeyID := d.Id()
	apiKey, err := client.APIKeys.Get(spaceID, apiKeyID)
	if err != nil {
		return err
	}

	apiKey.Name = d.Get("name").(string)
	apiKey.Description = d.Get("description").(string)

	err = client.APIKeys.Upsert(spaceID, apiKey)
	if err != nil {
		return err
	}

	if err := setAPIKeyProperties(d, apiKey); err != nil {
		return err
	}

	d.SetId(apiKey.Sys.ID)

	return nil
}

func resourceReadAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)
	apiKeyID := d.Id()

	apiKey, err := client.APIKeys.Get(spaceID, apiKeyID)
	if _, ok := err.(contentful.NotFoundError); ok {
		d.SetId("")
		return nil
	}

	return setAPIKeyProperties(d, apiKey)
}

func resourceDeleteAPIKey(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*contentful.Contentful)
	spaceID := d.Get("space_id").(string)
	apiKeyID := d.Id()

	apiKey, err := client.APIKeys.Get(spaceID, apiKeyID)
	if err != nil {
		return err
	}

	return client.APIKeys.Delete(spaceID, apiKey)
}

func setAPIKeyProperties(d *schema.ResourceData, apiKey *contentful.APIKey) error {
	if err := d.Set("space_id", apiKey.Sys.Space.Sys.ID); err != nil {
		return err
	}

	if err := d.Set("version", apiKey.Sys.Version); err != nil {
		return err
	}

	if err := d.Set("name", apiKey.Name); err != nil {
		return err
	}

	if err := d.Set("description", apiKey.Description); err != nil {
		return err
	}

	if err := d.Set("access_token", apiKey.AccessToken); err != nil {
		return err
	}

	return nil
}

func resourceAPIKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idAttr := strings.SplitN(d.Id(), "/", 2)
	if len(idAttr) == 2 {
		d.Set("space_id", idAttr[0])
		d.Set("name", idAttr[1])
	} else {
		return nil, fmt.Errorf("invalid id %q specified, should be in format \"spaceId/keyId\" for import", d.Id())
	}

	if err := resourceReadAPIKey(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
