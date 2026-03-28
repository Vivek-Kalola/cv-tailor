package api

import (
	"cv-tailoring/internal/ai"
	"cv-tailoring/internal/api/handler"
	"cv-tailoring/internal/db"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewServer(port int, database *db.DB, logger *logrus.Logger) error {
	r := gin.Default()

	geminiAPI, err := ai.NewGeminiAPI(logger)
	if err != nil {
		logger.Error(err)
	}

	resumeHandler := handler.NewResumeHandler(database)
	analyzeHandler := handler.AnalyzeHandler(database, geminiAPI, logger)
	tailorHandler := handler.TailorHandler(database, geminiAPI, logger)

	addRouters(r, resumeHandler, analyzeHandler, tailorHandler)

	logger.Info("Starting server on port " + strconv.Itoa(port))
	return r.Run(fmt.Sprintf(":%d", port))
}

func addRouters(r *gin.Engine, resumeHandler *handler.ResumeHandler, analyzeHandler *handler.AnalyzerHandling, tailorHandler *handler.TailorHandling) {
	r.HandleMethodNotAllowed = true
	r.HandleMethodNotAllowed = true

	r.GET("/api/resumes/default", resumeHandler.GetDefaultResume)
	r.POST("/api/resumes", resumeHandler.CreateResume)
	r.GET("/api/resumes/:id", resumeHandler.GetResumeByID)
	r.DELETE("/api/resumes/:id", resumeHandler.DeleteResumesByID)
	r.DELETE("/api/resumes/default", resumeHandler.DeleteDefaultResume)

	r.POST("/api/analyze", analyzeHandler.Analyze)
	r.POST("/api/tailor", tailorHandler.Tailor)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
