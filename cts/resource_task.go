package cts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/qzx/cts-client-go"
)

func resourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskCreate,
		ReadContext:   resourceTaskRead,
		UpdateContext: resourceTaskUpdate,
		DeleteContext: resourceTaskDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"module": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"condition": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kv": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datacenter": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"namespace": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"path": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"recurse": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"use_as_module_input": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"service": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cts_user_defined_meta": &schema.Schema{
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"datacenter": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"namespace": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"regexp": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"filter": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"names": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										ForceNew: true,
									},
									"use_as_module_input": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"catalog_service": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cts_user_defined_meta": &schema.Schema{
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"datacenter": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"namespace": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"regexp": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"node_meta": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"use_as_module_input": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"schedule": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cron": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"providers": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTaskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	data, err := taskResourceData(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	task := client.Task{
		Task: data,
	}
	o, err := c.CreateTask(task)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildChecksumID([]string{o.Task.Name}))

	return resourceTaskRead(ctx, d, m)
}

func resourceTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice Type
	var diags diag.Diagnostics

	taskName := d.Get("name").(string)

	_, err := c.GetTask(taskName)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceTaskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	if d.HasChange("enabled") {
		name := d.Get("name").(string)

		v := d.Get("enabled")
		if err := c.UpdateTaskEnable(name, v.(bool)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaskRead(ctx, d, m)
}

func resourceTaskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice Type
	var diags diag.Diagnostics

	taskName := d.Get("name").(string)

	err := c.DeleteTask(taskName)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func taskResourceData(d *schema.ResourceData, meta interface{}) (client.TaskItem, error) {
	data := client.TaskItem{}

	if v, ok := d.GetOk("name"); ok {
		data.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		data.Description = v.(string)
	}
	if v, ok := d.GetOk("module"); ok {
		data.Module = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		data.Enabled = v.(bool)
	}
	if v, ok := d.GetOk("providers"); ok {
		vL := v.(*schema.Set).List()
		providers := make([]string, 0, len(vL))
		for _, v := range vL {
			providers = append(providers, v.(string))
		}
		data.Providers = providers
	}
	if v, ok := d.GetOk("condition"); ok {
		condition := v.([]interface{})[0].(map[string]interface{})

		if v, ok := condition["kv"]; ok && len(v.([]interface{})) > 0 {
			ckv := &client.ConsulKv{}
			kv := v.([]interface{})[0].(map[string]interface{})
			if v, ok := kv["datacenter"]; ok {
				ckv.Datacenter = v.(string)
			}
			if v, ok := kv["namespace"]; ok {
				ckv.Namespace = v.(string)
			}
			if v, ok := kv["path"]; ok {
				ckv.Path = v.(string)
			}
			if v, ok := kv["Recurse"]; ok {
				ckv.Recurse = v.(bool)
			}
			if v, ok := kv["use_as_module_input"]; ok {
				bv := v.(bool)
				ckv.UseAsModuleInput = &bv
			}
			data.Condition.ConsulKv = ckv
		}
		if v, ok := condition["service"]; ok && len(v.([]interface{})) > 0 {
			cs := &client.Services{}
			s := v.([]interface{})[0].(map[string]interface{})
			if v, ok := s["datacenter"]; ok {
				cs.Datacenter = v.(string)
			}
			if v, ok := s["namespace"]; ok {
				cs.Namespace = v.(string)
			}
			if v, ok := s["filter"]; ok {
				cs.Filter = v.(string)
			}
			if v, ok := s["regexp"]; ok {
				cs.Regexp = v.(string)
			}
			if v, ok := s["names"]; ok {
				cs.Names = v.([]string)
			}
			if v, ok := s["use_as_module_input"]; ok {
				bv := v.(bool)
				cs.UseAsModuleInput = &bv
			}
			data.Condition.Services = cs
		}
	}

	return data, nil
}
