package ai

import (
	"context"
	"cv-tailoring/internal/db"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type GeminiAPI struct {
	client *genai.Client
	model  *genai.GenerativeModel
	logger *logrus.Logger
}

func NewGeminiAPI(logger *logrus.Logger) (*GeminiAPI, error) {

	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.2)

	return &GeminiAPI{
		client: client,
		model:  model,
		logger: logger,
	}, nil
}

func (g *GeminiAPI) Close() {
	_ = g.client.Close()
}

func (g *GeminiAPI) AnalyzeAgainstJD(resume *db.Resume, jd string) (string, error) {

	skills, err := g.extractSkills(jd)

	marshaled, err := json.Marshal(resume)

	if err != nil {
		return "", err
	}

	auditPrompt := `Role: High-Precision ATS Auditor.
    Task: Match the "Required Skills List" against the "Candidate Resume".
    
    Rules:
    - Zero-Inference: Only matched = true if explicitly in resume.
    - Semantic Equivalence: K8s = Kubernetes, etc.
    - No Hallucinations: Use exact text.
    
    Output Format: Return raw JSON only (no markdown).
    Schema:
    {
       "candidateName": "string",
       "matchScore": "number",
       "skills": [{"skill": "string", "category": "string", "impact": "High|Medium|Low", "required": true, "matched": boolean, "confidenceScore": int}]
    }

    Required Skills List:
    ` + skills + `

    Candidate Resume:
    ` + string(marshaled)

	content, err := g.model.GenerateContent(context.Background(), genai.Text(auditPrompt))
	if err != nil {
		return "", err
	}

	if len(content.Candidates) > 0 && len(content.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%s", content.Candidates[0].Content.Parts[0]), nil
	}

	return "", errors.New("no valid response found")
}

func (g *GeminiAPI) extractSkills(jd string) (string, error) {
	extractionPrompt := `Identify every single technical skill, soft skill, tool, and qualification in the following Job Description. 
    Break down complex sentences into individual items (e.g., "Java and Spring Boot" becomes "Java", "Spring Boot").
    Return the result as a simple JSON array of strings: ["skill1", "skill2", ...]. 
    Do not skip anything. Output ONLY the JSON array.
    
    Job Description:
    ` + jd

	content, err := g.model.GenerateContent(context.Background(), genai.Text(extractionPrompt))

	if err != nil {
		return "", err
	}

	if len(content.Candidates) > 0 && len(content.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%s", content.Candidates[0].Content.Parts[0]), nil
	}

	return "", errors.New("no valid response found")
}
