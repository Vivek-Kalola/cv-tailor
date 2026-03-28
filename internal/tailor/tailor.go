package tailor

import (
	"cv-tailoring/internal/db"
)

type Tailor struct {
}

func NewTailor() *Tailor {
	return &Tailor{}
}

func (t Tailor) TailorResume(text string, resume *db.Resume, skills []string) (string, error) {
	//TODO: Implement this

	result := "Tailored Resume"
	for _, skill := range skills {
		result += skill
	}

	return result, nil
}
