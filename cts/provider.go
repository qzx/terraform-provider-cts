package cts

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	client "github.com/qzx/cts-client-go"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CTS_HOST", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cts_task": resourceTask(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cts_task": dataSourceTask(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := client.NewClient(host, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create CTS client",
			Detail:   "Unable to create anonymous CTS client",
		})
		return nil, diags
	}

	return c, diags
}

func buildChecksumID(v []string) string {
	sort.Strings(v)

	h := md5.New()
	// Hash.Write never returns an error. See https://pkg.go.dev/hash#Hash
	_, _ = h.Write([]byte(strings.Join(v, "")))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func expandNestedSet(m interface{}, target string) []string {
	var res []string
	vL := m.(*schema.Set).List()
	for _, v := range vL {
		res = append(res, v.(string))
	}
	return res
}
