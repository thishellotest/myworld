package biz

import "time"

type AsanaWebhookVo struct {
	Events []AsanaWebhookEvent `json:"events"`
}

type AsanaWebhookEvent struct {
	User      AsanaWebhookUser     `json:"user"`
	CreatedAt time.Time            `json:"created_at"`
	Action    string               `json:"action"`
	Resource  AsanaWebhookResource `json:"resource"`
	Parent    interface{}          `json:"parent"`
	Change    AsanaWebhookChange   `json:"change"`
}

func (c *AsanaWebhookEvent) IsTaskWebhook() bool {
	if c.Resource.ResourceType == "task" {
		if c.Action == "changed" || c.Action == "added" || c.Action == "removed" ||
			c.Action == "undeleted" || c.Action == "deleted" {
			return true
		}
	}
	return false
}

/*
{"events":[{"user":{"gid":"1206230291638946","resource_type":"user"},"created_at":"2023-12-30T10:12:25.274Z","action":"changed","resource":{"gid":"1206234446219801","resource_type":"task","resource_subtype":"default_task"},"parent":null,"change":{"field":"name","action":"changed"}}]}
*/
func (c *AsanaWebhookEvent) IsUserWebhook() bool {
	if c.Resource.ResourceType == "task" {
		if c.Action == "changed" || c.Action == "added" {
			return true
		}
	}
	return false
}

type AsanaWebhookChange struct {
	Field    string               `json:"field"`
	Action   string               `json:"action"`
	NewValue AsanaWebhookNewValue `json:"new_value"`
}

type AsanaWebhookUser struct {
	Gid          string `json:"gid"`
	ResourceType string `json:"resource_type"`
}
type AsanaWebhookResource struct {
	Gid             string `json:"gid"`
	ResourceType    string `json:"resource_type"`
	ResourceSubtype string `json:"resource_subtype"`
}
type AsanaWebhookEnumValue struct {
	Gid          string `json:"gid"`
	ResourceType string `json:"resource_type"`
}
type AsanaWebhookNewValue struct {
	Gid             string                `json:"gid"`
	ResourceType    string                `json:"resource_type"`
	ResourceSubtype string                `json:"resource_subtype"`
	EnumValue       AsanaWebhookEnumValue `json:"enum_value"`
}
