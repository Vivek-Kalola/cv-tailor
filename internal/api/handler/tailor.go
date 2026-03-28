package handler

import (
	"cv-tailoring/internal/ai"
	"cv-tailoring/internal/db"
	"cv-tailoring/internal/tailor"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TailorRequest struct {
	JdText         string   `json:"jd_text"`
	ResumeID       int      `json:"resume_id"`
	SelectedSkills []string `json:"selected_skills"`
}

type TailorResponse struct {
	TailoredContent string `json:"tailored_content"`
	PdfPath         string `json:"pdf_path"`
}

type TailorHandling struct {
	db     *db.DB
	tailor *tailor.Tailor
}

func TailorHandler(db *db.DB, api *ai.GeminiAPI, logger *logrus.Logger) *TailorHandling {
	return &TailorHandling{
		db:     db,
		tailor: tailor.NewTailor(),
	}
}

func (h *TailorHandling) Tailor(context *gin.Context) {

	body := &TailorRequest{}

	if err := context.ShouldBindBodyWithJSON(body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resume, ok := h.db.GetResume(body.ResumeID)

	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "resume not found"})
		return
	}

	result, err := h.tailor.TailorResume(body.JdText, &resume, body.SelectedSkills)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tailorResponse := &TailorResponse{
		TailoredContent: result,
		PdfPath:         "/path/to/pdf",
	}

	context.JSON(http.StatusOK, tailorResponse)
}
