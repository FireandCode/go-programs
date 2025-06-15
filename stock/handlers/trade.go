package handlers

import (
	"net/http"
	"strconv"

	"stock-trades-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TradeHandler struct {
	db *gorm.DB
}

func NewTradeHandler(db *gorm.DB) *TradeHandler {
	return &TradeHandler{db: db}
}

type CreateTradeRequest struct {
	Type      string  `json:"type" binding:"required"`
	Symbol    string  `json:"symbol" binding:"required"`
	Shares    int     `json:"shares" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Timestamp int64   `json:"timestamp"`
}

func (h *TradeHandler) CreateTrade(c *gin.Context) {
	var req CreateTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	trade := &models.Trade{
		Type:      models.TradeType(req.Type),
		UserID:    userID.(uint),
		Symbol:    req.Symbol,
		Shares:    req.Shares,
		Price:     req.Price,
		Timestamp: req.Timestamp,
	}

	if err := trade.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(trade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trade"})
		return
	}

	c.JSON(http.StatusCreated, trade)
}

func (h *TradeHandler) GetTrades(c *gin.Context) {
	var trades []models.Trade
	query := h.db.Order("id asc")

	if tradeType := c.Query("type"); tradeType != "" {
		query = query.Where("type = ?", tradeType)
	}

	if userID := c.Query("user_id"); userID != "" {
		uid, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		query = query.Where("user_id = ?", uint(uid))
	}

	if err := query.Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trades"})
		return
	}

	c.JSON(http.StatusOK, trades)
}

func (h *TradeHandler) GetTradeByID(c *gin.Context) {
	id := c.Param("id")
	var trade models.Trade

	if err := h.db.First(&trade, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "trade not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trade"})
		return
	}

	c.JSON(http.StatusOK, trade)
}

func (h *TradeHandler) HandleUnsupportedMethods(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
} 