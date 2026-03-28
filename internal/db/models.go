package db

type BaseModel interface {
	Id() int
}

type Title struct {
	Role string   `json:"role"`
	Tags []string `json:"tags"`
}

type Socials struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	LinkedIn string `json:"linkedin"`
	Github   string `json:"github"`
}

type Expertise struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Skills []string `json:"skills"`
}

type Experience struct {
	ID        string   `json:"id"`
	Company   string   `json:"company"`
	JobTitle  string   `json:"jobTitle"`
	Location  string   `json:"location"`
	WorkType  string   `json:"workType"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	IsCurrent bool     `json:"isCurrent"`
	Bullets   []string `json:"bullets"`
}

type Education struct {
	ID        string `json:"id"`
	UniName   string `json:"uniName"`
	Degree    string `json:"degree"`
	Location  string `json:"location"`
	StartYear string `json:"startYear"`
	EndYear   string `json:"endYear"`
}

type Resume struct {
	BaseModel
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Title      Title        `json:"title"`
	Socials    Socials      `json:"socials"`
	Expertise  []Expertise  `json:"expertise"`
	Experience []Experience `json:"experience"`
	Education  []Education  `json:"education"`
}

func (r *Resume) Id() int {
	return r.ID
}

type JD struct {
	BaseModel
	id int
}

func (j *JD) Id() int {
	return j.id
}
