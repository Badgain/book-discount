package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Badgain/book-discount/internal/domain"
	"github.com/Badgain/book-discount/internal/handler/dto"
)

// DiscountHandler обрабатывает HTTP запросы для расчета скидок
type DiscountHandler struct {
	service domain.DiscountCalculator
}

// NewDiscountHandler создает новый экземпляр DiscountHandler
func NewDiscountHandler(service domain.DiscountCalculator) *DiscountHandler {
	return &DiscountHandler{
		service: service,
	}
}

// CalculateDiscount обрабатывает POST запрос для расчета скидки
func (h *DiscountHandler) CalculateDiscount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.DiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Расчет скидки
	response, err := h.service.Calculate(ctx, req.CustomerTypeAsDomain(), req.BooksAsDomain())
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправка ответа
	h.sendJSON(w, discountAsDTO(response), http.StatusOK)
}

// sendJSON отправляет JSON ответ
func (h *DiscountHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Default().Error("unable to encode response", slog.String("error", err.Error()))
	}
}

// sendError отправляет ошибку в формате JSON
func (h *DiscountHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	h.sendJSON(w, dto.ErrorResponse{Error: message}, statusCode)
}

func discountAsDTO(discount domain.Discount) dto.DiscountResponse {
	return dto.DiscountResponse{
		OriginalAmount:  discount.CartAmount,
		DiscountPercent: discount.DiscountPercent,
		DiscountAmount:  discount.DiscountAmount,
		FinalAmount:     discount.TotalCost,
	}
}
