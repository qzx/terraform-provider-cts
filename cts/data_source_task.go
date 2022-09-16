package cts

import (
	"context"

	//"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/qzx/cts-client-go"
)

func dataSourceTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaskRead,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"module": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"condition": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kv": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datacenter": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"recurse": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
									"use_as_module_input": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"service": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cts_user_defined_meta": &schema.Schema{
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"datacenter": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"regexp": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"filter": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
									"names": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"use_as_module_input": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"catalog_service": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cts_user_defined_meta": &schema.Schema{
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"datacenter": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"regexp": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_meta": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"use_as_module_input": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"schedule": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cron": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"providers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice Type
	var diags diag.Diagnostics

	taskName := d.Get("name").(string)

	task, err := c.GetTask(taskName)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(buildChecksumID([]string{task.Task.Name}))

	if err := d.Set("description", task.Task.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("module", task.Task.Module); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", task.Task.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("providers", task.Task.Providers); err != nil {
		return diag.FromErr(err)
	}

	condition := flattenConditions(&task.Task.Condition)
	if err := d.Set("condition", condition); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func flattenConditions(conditions *client.Condition) []interface{} {
	c := make(map[string]interface{})
	if conditions.ConsulKv != nil {
		c["kv"] = flattenConditionKv(conditions.ConsulKv)
	}
	if conditions.Services != nil {
		c["service"] = flattenConditionServices(conditions.Services)
	}
	if len(c) > 0 {
		cs := make([]interface{}, 1, 1)
		cs[0] = c
		return cs
	} else {
		return make([]interface{}, 0)
	}

}

func flattenConditionKv(condition *client.ConsulKv) []interface{} {
	if condition != nil {
		cs := make([]interface{}, 1, 1)
		c := make(map[string]interface{})

		c["datacenter"] = condition.Datacenter
		c["namespace"] = condition.Namespace
		c["path"] = condition.Path
		c["recurse"] = condition.Recurse
		c["use_as_module_input"] = condition.UseAsModuleInput

		cs[0] = c
		return cs
	}
	return make([]interface{}, 0)
}

func flattenConditionServices(condition *client.Services) []interface{} {
	if condition != nil {
		cs := make([]interface{}, 1, 1)
		c := make(map[string]interface{})

		c["datacenter"] = condition.Datacenter
		c["namespace"] = condition.Namespace
		c["regexp"] = condition.Regexp
		c["filter"] = condition.Filter
		c["names"] = condition.Names
		c["use_as_module_input"] = condition.UseAsModuleInput

		cs[0] = c
		return cs
	}
	return make([]interface{}, 0)
}
