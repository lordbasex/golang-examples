package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"api/internal/models"
)

// CampaignHandler contiene las dependencias para los handlers de campañas
type CampaignHandler struct {
	DB *sql.DB
}

// NewCampaignHandler crea una nueva instancia de CampaignHandler
func NewCampaignHandler(db *sql.DB) *CampaignHandler {
	return &CampaignHandler{
		DB: db,
	}
}

// GetAll obtiene todas las campañas
// @Summary Listar campañas
// @Description Obtiene la lista de todas las campañas, opcionalmente filtradas por estado
// @Tags campaigns
// @Produce json
// @Param status query string false "Filtrar por estado (draft, active, paused, completed, archived)"
// @Success 200 {object} models.CampaignsResponse
// @Failure 500 {object} models.CampaignsResponse
// @Router /api/v1/campaigns [get]
func (h *CampaignHandler) GetAll(c *fiber.Ctx) error {
	status := c.Query("status") // Filtro opcional por estado
	
	var query string
	var rows *sql.Rows
	var err error

	if status != "" {
		query = `SELECT id, name, status, created_at FROM campaigns WHERE status = ? ORDER BY created_at DESC`
		rows, err = h.DB.Query(query, status)
	} else {
		query = `SELECT id, name, status, created_at FROM campaigns ORDER BY created_at DESC`
		rows, err = h.DB.Query(query)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignsResponse{
			Success: false,
			Message: "Error al obtener campañas: " + err.Error(),
		})
	}
	defer rows.Close()

	campaigns := make([]*models.Campaign, 0)
	for rows.Next() {
		campaign := &models.Campaign{}
		if err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.Status, &campaign.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignsResponse{
				Success: false,
				Message: "Error al leer campañas: " + err.Error(),
			})
		}
		campaigns = append(campaigns, campaign)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignsResponse{
			Success: false,
			Message: "Error al procesar campañas: " + err.Error(),
		})
	}

	return c.JSON(models.CampaignsResponse{
		Success: true,
		Data:    campaigns,
		Total:   len(campaigns),
	})
}

// GetByID obtiene una campaña por su ID
// @Summary Obtener campaña por ID
// @Description Obtiene una campaña específica por su ID
// @Tags campaigns
// @Produce json
// @Param id path int true "Campaign ID"
// @Success 200 {object} models.CampaignResponse
// @Failure 400 {object} models.CampaignResponse
// @Failure 404 {object} models.CampaignResponse
// @Failure 500 {object} models.CampaignResponse
// @Router /api/v1/campaigns/{id} [get]
func (h *CampaignHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	campaignID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "ID de campaña inválido",
		})
	}

	query := `SELECT id, name, status, created_at FROM campaigns WHERE id = ?`
	
	campaign := &models.Campaign{}
	err = h.DB.QueryRow(query, campaignID).Scan(
		&campaign.ID,
		&campaign.Name,
		&campaign.Status,
		&campaign.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(models.CampaignResponse{
			Success: false,
			Message: "Campaña no encontrada",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignResponse{
			Success: false,
			Message: "Error al obtener campaña: " + err.Error(),
		})
	}

	return c.JSON(models.CampaignResponse{
		Success: true,
		Data:    campaign,
	})
}

// Create crea una nueva campaña
// @Summary Crear campaña
// @Description Crea una nueva campaña en la base de datos
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body models.CreateCampaignRequest true "Datos de la campaña"
// @Success 201 {object} models.CampaignResponse
// @Failure 400 {object} models.CampaignResponse
// @Failure 500 {object} models.CampaignResponse
// @Router /api/v1/campaigns [post]
func (h *CampaignHandler) Create(c *fiber.Ctx) error {
	req := new(models.CreateCampaignRequest)
	
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "Datos inválidos: " + err.Error(),
		})
	}

	// Validación básica
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "El nombre de la campaña es requerido",
		})
	}

	// Si no se proporciona status, usar 'draft' por defecto
	if req.Status == "" {
		req.Status = models.StatusDraft
	}

	// Validar que el status sea válido
	validStatuses := map[string]bool{
		models.StatusDraft:     true,
		models.StatusActive:    true,
		models.StatusPaused:    true,
		models.StatusCompleted: true,
		models.StatusArchived:  true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "Estado inválido. Valores permitidos: draft, active, paused, completed, archived",
		})
	}

	query := `INSERT INTO campaigns (name, status) VALUES (?, ?)`
	
	result, err := h.DB.Exec(query, req.Name, req.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignResponse{
			Success: false,
			Message: "Error al crear campaña: " + err.Error(),
		})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignResponse{
			Success: false,
			Message: "Error al obtener ID de la campaña creada",
		})
	}

	// Obtener la campaña creada
	campaign := &models.Campaign{}
	err = h.DB.QueryRow(
		`SELECT id, name, status, created_at FROM campaigns WHERE id = ?`,
		id,
	).Scan(&campaign.ID, &campaign.Name, &campaign.Status, &campaign.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(models.CampaignResponse{
			Success: true,
			Message: "Campaña creada exitosamente",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.CampaignResponse{
		Success: true,
		Message: "Campaña creada exitosamente",
		Data:    campaign,
	})
}

// Update actualiza una campaña existente
// @Summary Actualizar campaña
// @Description Actualiza los datos de una campaña existente
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path int true "Campaign ID"
// @Param campaign body models.UpdateCampaignRequest true "Datos a actualizar"
// @Success 200 {object} models.CampaignResponse
// @Failure 400 {object} models.CampaignResponse
// @Failure 404 {object} models.CampaignResponse
// @Failure 500 {object} models.CampaignResponse
// @Router /api/v1/campaigns/{id} [put]
func (h *CampaignHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	campaignID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "ID de campaña inválido",
		})
	}

	req := new(models.UpdateCampaignRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "Datos inválidos: " + err.Error(),
		})
	}

	// Verificar que la campaña existe
	var exists bool
	err = h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM campaigns WHERE id = ?)`, campaignID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(models.CampaignResponse{
			Success: false,
			Message: "Campaña no encontrada",
		})
	}

	// Validar status si se proporciona
	if req.Status != "" {
		validStatuses := map[string]bool{
			models.StatusDraft:     true,
			models.StatusActive:    true,
			models.StatusPaused:    true,
			models.StatusCompleted: true,
			models.StatusArchived:  true,
		}
		if !validStatuses[req.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
				Success: false,
				Message: "Estado inválido. Valores permitidos: draft, active, paused, completed, archived",
			})
		}
	}

	// Construir query dinámico solo con los campos proporcionados
	query := `UPDATE campaigns SET `
	args := make([]interface{}, 0)
	updates := make([]string, 0)

	if req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, req.Name)
	}
	if req.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, req.Status)
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "No hay datos para actualizar",
		})
	}

	// Unir los updates y agregar la condición WHERE
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"
	args = append(args, campaignID)

	_, err = h.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignResponse{
			Success: false,
			Message: "Error al actualizar campaña: " + err.Error(),
		})
	}

	// Obtener la campaña actualizada
	campaign := &models.Campaign{}
	err = h.DB.QueryRow(
		`SELECT id, name, status, created_at FROM campaigns WHERE id = ?`,
		campaignID,
	).Scan(&campaign.ID, &campaign.Name, &campaign.Status, &campaign.CreatedAt)

	if err != nil {
		return c.JSON(models.CampaignResponse{
			Success: true,
			Message: "Campaña actualizada exitosamente",
		})
	}

	return c.JSON(models.CampaignResponse{
		Success: true,
		Message: "Campaña actualizada exitosamente",
		Data:    campaign,
	})
}

// Delete elimina una campaña
// @Summary Eliminar campaña
// @Description Elimina una campaña de la base de datos
// @Tags campaigns
// @Produce json
// @Param id path int true "Campaign ID"
// @Success 200 {object} models.CampaignResponse
// @Failure 400 {object} models.CampaignResponse
// @Failure 404 {object} models.CampaignResponse
// @Failure 500 {object} models.CampaignResponse
// @Router /api/v1/campaigns/{id} [delete]
func (h *CampaignHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	campaignID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.CampaignResponse{
			Success: false,
			Message: "ID de campaña inválido",
		})
	}

	// Verificar que la campaña existe antes de eliminar
	var exists bool
	err = h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM campaigns WHERE id = ?)`, campaignID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(models.CampaignResponse{
			Success: false,
			Message: "Campaña no encontrada",
		})
	}

	// Nota: Las relaciones campaign_members se eliminan automáticamente por ON DELETE CASCADE
	query := `DELETE FROM campaigns WHERE id = ?`
	_, err = h.DB.Exec(query, campaignID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignResponse{
			Success: false,
			Message: "Error al eliminar campaña: " + err.Error(),
		})
	}

	return c.JSON(models.CampaignResponse{
		Success: true,
		Message: "Campaña eliminada exitosamente (incluye relaciones con miembros)",
	})
}

// GetStats obtiene estadísticas de campañas por estado
// @Summary Estadísticas de campañas
// @Description Obtiene un resumen de campañas agrupadas por estado
// @Tags campaigns
// @Produce json
// @Success 200 {object} models.CampaignStatsResponse
// @Failure 500 {object} models.CampaignStatsResponse
// @Router /api/v1/campaigns/stats [get]
func (h *CampaignHandler) GetStats(c *fiber.Ctx) error {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END) as draft,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status = 'paused' THEN 1 ELSE 0 END) as paused,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'archived' THEN 1 ELSE 0 END) as archived
		FROM campaigns
	`

	stats := &models.CampaignStats{}
	err := h.DB.QueryRow(query).Scan(
		&stats.TotalCampaigns,
		&stats.Draft,
		&stats.Active,
		&stats.Paused,
		&stats.Completed,
		&stats.Archived,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.CampaignStatsResponse{
			Success: false,
			Message: "Error al obtener estadísticas: " + err.Error(),
		})
	}

	return c.JSON(models.CampaignStatsResponse{
		Success: true,
		Data:    stats,
	})
}

