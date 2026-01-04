package models

import "time"

// Campaign representa una campaña en la base de datos
type Campaign struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CampaignStatus representa los estados válidos de una campaña
const (
	StatusDraft     = "draft"
	StatusActive    = "active"
	StatusPaused    = "paused"
	StatusCompleted = "completed"
	StatusArchived  = "archived"
)

// CreateCampaignRequest representa la solicitud para crear una campaña
type CreateCampaignRequest struct {
	Name   string `json:"name" validate:"required,min=3,max=255"`
	Status string `json:"status,omitempty" validate:"omitempty,oneof=draft active paused completed archived"`
}

// UpdateCampaignRequest representa la solicitud para actualizar una campaña
type UpdateCampaignRequest struct {
	Name   string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Status string `json:"status,omitempty" validate:"omitempty,oneof=draft active paused completed archived"`
}

// CampaignResponse representa la respuesta con información de la campaña
type CampaignResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Data    *Campaign `json:"data,omitempty"`
}

// CampaignsResponse representa la respuesta con lista de campañas
type CampaignsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    []*Campaign `json:"data,omitempty"`
	Total   int         `json:"total"`
}

// CampaignStats representa estadísticas de campañas
type CampaignStats struct {
	TotalCampaigns int `json:"total_campaigns"`
	Draft          int `json:"draft"`
	Active         int `json:"active"`
	Paused         int `json:"paused"`
	Completed      int `json:"completed"`
	Archived       int `json:"archived"`
}

// CampaignStatsResponse representa la respuesta con estadísticas
type CampaignStatsResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    *CampaignStats `json:"data,omitempty"`
}

