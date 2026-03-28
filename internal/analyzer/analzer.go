package analyzer

import (
	"cv-tailoring/internal/ai"
	"cv-tailoring/internal/db"
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

type Analyzer struct {
	geminiAPI *ai.GeminiAPI
	logger    *logrus.Logger
}

type AnalyzeResult struct {
	MatchScore    float64  `json:"matchScore"`
	CandidateName string   `json:"candidateName"`
	Skills        []*Skill `json:"skills"`
}

type Skill struct {
	Skill           string `json:"skill"`
	Category        string `json:"category"`
	Impact          string `json:"impact"`
	Required        bool   `json:"required"`
	Matched         bool   `json:"matched"`
	ConfidenceScore int    `json:"confidenceScore"`
}

func (a *Analyzer) Analyze(jd string, resume *db.Resume) (*AnalyzeResult, error) {

	if true {
		file, err := os.ReadFile("./db/dummy.analyze")

		if err != nil {
			return nil, err
		}

		analyzeResult := &AnalyzeResult{}

		err = json.Unmarshal(file, &analyzeResult)
		if err != nil {
			return nil, err
		}
		return analyzeResult, nil
	}

	result, err := a.geminiAPI.AnalyzeAgainstJD(resume, jd)

	if err != nil {
		return nil, err
	}

	analyzeResult := &AnalyzeResult{}

	err = json.Unmarshal([]byte(result), analyzeResult)

	if err != nil {
		return nil, err
	}

	return analyzeResult, nil
}

func NewAnalyzer(geminiAPI *ai.GeminiAPI, logger *logrus.Logger) *Analyzer {
	return &Analyzer{
		geminiAPI: geminiAPI,
		logger:    logger,
	}
}
