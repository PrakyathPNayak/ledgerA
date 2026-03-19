package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler handles category and subcategory endpoints.
type CategoryHandler struct {
	categoryService    service.CategoryService
	subcategoryService service.SubcategoryService
	userService        service.UserService
}

// NewCategoryHandler creates CategoryHandler.
func NewCategoryHandler(categoryService service.CategoryService, subcategoryService service.SubcategoryService, userService service.UserService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService, subcategoryService: subcategoryService, userService: userService}
}

func (h *CategoryHandler) resolveUserID(c *gin.Context) (uuid.UUID, bool) {
	uid, err := middleware.GetFirebaseUID(c)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", "unauthorized")
		return uuid.Nil, false
	}
	user, err := h.userService.GetMe(c.Request.Context(), uid)
	if err != nil {
		dto.Error(c, 404, "ERR_NOT_FOUND", "user not found")
		return uuid.Nil, false
	}
	return user.ID, true
}

// List handles GET /api/v1/categories.
func (h *CategoryHandler) List(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	categories, _, err := h.categoryService.List(c.Request.Context(), userID)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to list categories")
		return
	}
	result := make([]dto.CategoryWithSubsResponse, 0, len(categories))
	for _, category := range categories {
		subs, _, subErr := h.subcategoryService.ListByCategory(c.Request.Context(), userID, category.ID)
		if subErr != nil {
			dto.Error(c, 500, "ERR_INTERNAL", "failed to list subcategories")
			return
		}
		subResponses := make([]dto.SubcategoryResponse, 0, len(subs))
		for _, sub := range subs {
			subResponses = append(subResponses, dto.NewSubcategoryResponse(sub))
		}
		result = append(result, dto.CategoryWithSubsResponse{CategoryResponse: dto.NewCategoryResponse(category), Subcategories: subResponses})
	}
	dto.Success(c, result)
}

// Create handles POST /api/v1/categories.
func (h *CategoryHandler) Create(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	category, err := h.categoryService.Create(c.Request.Context(), userID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to create category")
		return
	}
	dto.Created(c, dto.NewCategoryResponse(*category))
}

// Update handles PATCH /api/v1/categories/:id.
func (h *CategoryHandler) Update(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid category id")
		return
	}
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	category, err := h.categoryService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to update category")
		return
	}
	dto.Success(c, dto.NewCategoryResponse(*category))
}

// Delete handles DELETE /api/v1/categories/:id.
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid category id")
		return
	}
	if err := h.categoryService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to delete category")
		return
	}
	c.Status(204)
}

// CreateSubcategory handles POST /api/v1/categories/:id/subcategories.
func (h *CategoryHandler) CreateSubcategory(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid category id")
		return
	}
	var req dto.CreateSubcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	sub, err := h.subcategoryService.Create(c.Request.Context(), userID, categoryID, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to create subcategory")
		return
	}
	dto.Created(c, dto.NewSubcategoryResponse(*sub))
}

// UpdateSubcategory handles PATCH /api/v1/subcategories/:id.
func (h *CategoryHandler) UpdateSubcategory(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid subcategory id")
		return
	}
	var req dto.UpdateSubcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", err.Error())
		return
	}
	sub, err := h.subcategoryService.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to update subcategory")
		return
	}
	dto.Success(c, dto.NewSubcategoryResponse(*sub))
}

// DeleteSubcategory handles DELETE /api/v1/subcategories/:id.
func (h *CategoryHandler) DeleteSubcategory(c *gin.Context) {
	userID, ok := h.resolveUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Error(c, 400, "ERR_BAD_REQUEST", "invalid subcategory id")
		return
	}
	if err := h.subcategoryService.Delete(c.Request.Context(), userID, id); err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to delete subcategory")
		return
	}
	c.Status(204)
}
