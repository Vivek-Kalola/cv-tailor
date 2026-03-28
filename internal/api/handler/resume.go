package handler

import (
	"cv-tailoring/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResumeHandler struct {
	db *db.DB
}

func NewResumeHandler(db *db.DB) *ResumeHandler {
	return &ResumeHandler{db: db}
}

func (h *ResumeHandler) GetDefaultResume(context *gin.Context) {
	resume, ok := h.db.GetDefaultResume()
	if !ok {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	context.JSON(http.StatusOK, resume)
}

func (h *ResumeHandler) DeleteDefaultResume(context *gin.Context) {
	err := h.db.DeleteResume(db.DefaultID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Done()
}

func (h *ResumeHandler) GetResumeByID(context *gin.Context) {

	stringID := context.Param("id")

	id, err := strconv.Atoi(stringID)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
	}

	resume, ok := h.db.GetResume(id)
	if !ok {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	context.JSON(http.StatusOK, resume)
}

func (h *ResumeHandler) CreateResume(context *gin.Context) {

	resume := &db.Resume{}

	err := context.ShouldBindBodyWithJSON(resume)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.db.UpsertResume(*resume)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Resume uploaded successfully", "id": id})
}

func (h *ResumeHandler) DeleteResumesByID(context *gin.Context) {
	stringID := context.Param("id")

	id, err := strconv.Atoi(stringID)

	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.db.DeleteResume(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Done()
}
