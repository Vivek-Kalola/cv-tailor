package handler

import (
	"cv-tailoring/internal/ai"
	"cv-tailoring/internal/analyzer"
	"cv-tailoring/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AnalyzeRequest struct {
	JdText string `json:"jd_text"`
}

type AnalyzerHandling struct {
	db       *db.DB
	analyzer *analyzer.Analyzer
	logger   *logrus.Logger
}

func AnalyzeHandler(db *db.DB, api *ai.GeminiAPI, logger *logrus.Logger) *AnalyzerHandling {
	return &AnalyzerHandling{
		db:       db,
		analyzer: analyzer.NewAnalyzer(api, logger),
		logger:   logger,
	}
}

func (h *AnalyzerHandling) Analyze(context *gin.Context) {

	body := &AnalyzeRequest{}

	if err := context.ShouldBindBodyWithJSON(body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	resume, ok := h.db.GetDefaultResume()

	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "resume not found"})
		return
	}

	result, err := h.analyzer.Analyze(body.JdText, &resume)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, result)
}
