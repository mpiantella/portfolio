package domain

import "time"

// Project represents a portfolio project/case study
type Project struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	Category            string    `json:"category"`
	ShortDescription    string    `json:"short_description"`
	Challenge           string    `json:"challenge"`
	Solution            string    `json:"solution"`
	Impact              Impact    `json:"impact"`
	Technologies        []string  `json:"technologies"`
	Timeline            string    `json:"timeline"`
	Featured            bool      `json:"featured"`
	Images              []string  `json:"images"`
	ArchitectureDiagram string    `json:"architecture_diagram,omitempty"`
	Outcomes            []string  `json:"outcomes"`
	Lessons             []string  `json:"lessons"`
	CreatedAt           time.Time `json:"created_at"`
}

// Impact describes the measurable impact of a project
type Impact struct {
	Metric  string `json:"metric"`
	Details string `json:"details"`
}

// Patent represents an intellectual property patent
type Patent struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	PatentNumber     string   `json:"patent_number,omitempty"`
	Status           string   `json:"status"` // Granted, Pending
	Year             int      `json:"year"`
	Description      string   `json:"description"`
	Impact           string   `json:"impact"`
	TechnicalDetails string   `json:"technical_details,omitempty"`
	Link             string   `json:"link,omitempty"`
	CoInventors      []string `json:"co_inventors,omitempty"`
}

// SpeakingEngagement represents a conference talk, webinar, or presentation
type SpeakingEngagement struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Event        string    `json:"event"`
	Date         time.Time `json:"date"`
	Location     string    `json:"location"`
	Type         string    `json:"type"` // Conference, Webinar, Internal
	VideoURL     string    `json:"video_url,omitempty"`
	SlidesURL    string    `json:"slides_url,omitempty"`
	Description  string    `json:"description"`
	Topics       []string  `json:"topics"`
	AudienceSize int       `json:"audience_size,omitempty"`
	KeyTakeaways []string  `json:"key_takeaways,omitempty"`
}

// Certification represents a professional certification
type Certification struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Issuer        string    `json:"issuer"`
	IssueDate     time.Time `json:"issue_date"`
	ExpiryDate    time.Time `json:"expiry_date,omitempty"`
	CredentialID  string    `json:"credential_id,omitempty"`
	CredentialURL string    `json:"credential_url,omitempty"`
	BadgeURL      string    `json:"badge_url,omitempty"`
	Skills        []string  `json:"skills,omitempty"`
}

// Testimonial represents a recommendation or testimonial
type Testimonial struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Role        string    `json:"role"`
	Company     string    `json:"company"`
	Content     string    `json:"content"`
	ImageURL    string    `json:"image_url,omitempty"`
	LinkedInURL string    `json:"linkedin_url,omitempty"`
	Date        time.Time `json:"date"`
	Featured    bool      `json:"featured"`
}

// ContactMessage represents a contact form submission
type ContactMessage struct {
	Name             string    `json:"name" validate:"required"`
	Email            string    `json:"email" validate:"required,email"`
	Company          string    `json:"company,omitempty"`
	Message          string    `json:"message" validate:"required"`
	PreferredContact string    `json:"preferred_contact"` // Email, LinkedIn, Phone
	Topic            string    `json:"topic"`             // Consulting, Speaking, Mentorship, Collaboration
	CreatedAt        time.Time `json:"created_at"`
}

// Skill represents a technical or soft skill
type Skill struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"` // Architecture, Languages, Cloud, Leadership, etc.
	Level       string   `json:"level"`    // Expert, Advanced, Intermediate
	YearsExp    int      `json:"years_experience"`
	Keywords    []string `json:"keywords,omitempty"`
	Highlighted bool     `json:"highlighted"`
}

// Experience represents a work experience entry
type Experience struct {
	Company      string    `json:"company"`
	Position     string    `json:"position"`
	Location     string    `json:"location"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date,omitempty"`
	Current      bool      `json:"current"`
	Description  string    `json:"description"`
	Highlights   []string  `json:"highlights"`
	Technologies []string  `json:"technologies"`
}

// Education represents an educational credential
type Education struct {
	Institution string    `json:"institution"`
	Degree      string    `json:"degree"`
	Field       string    `json:"field"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date"`
	GPA         string    `json:"gpa,omitempty"`
	Honors      []string  `json:"honors,omitempty"`
	Description string    `json:"description,omitempty"`
}

// SocialLink represents a social media or professional network link
type SocialLink struct {
	Platform string `json:"platform"` // LinkedIn, GitHub, Twitter, etc.
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
	Icon     string `json:"icon"` // CSS class or icon name
}

// ProjectsData represents the structure of projects.json
type ProjectsData struct {
	Projects []Project `json:"projects"`
}

// PatentsData represents the structure of patents.json
type PatentsData struct {
	Patents []Patent `json:"patents"`
}

// SpeakingData represents the structure of speaking.json
type SpeakingData struct {
	SpeakingEngagements []SpeakingEngagement `json:"speaking_engagements"`
}

// CertificationsData represents the structure of certifications.json
type CertificationsData struct {
	Certifications []Certification `json:"certifications"`
}

// ProfileData represents complete portfolio profile data
type ProfileData struct {
	FullName        string       `json:"full_name"`
	Title           string       `json:"title"`
	Tagline         string       `json:"tagline"`
	Bio             string       `json:"bio"`
	Email           string       `json:"email"`
	Phone           string       `json:"phone,omitempty"`
	Location        string       `json:"location"`
	LinkedInURL     string       `json:"linkedin_url"`
	ProfileImageURL string       `json:"profile_image_url"`
	ResumeURL       string       `json:"resume_url"`
	YearsExperience int          `json:"years_experience"`
	Languages       []Language   `json:"languages"`
	SocialLinks     []SocialLink `json:"social_links"`
	Skills          []Skill      `json:"skills"`
	Experiences     []Experience `json:"experiences"`
	Education       []Education  `json:"education"`
}

// Language represents a spoken language
type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"` // Native, Fluent, Professional, Conversational
}

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
	Uptime  string `json:"uptime,omitempty"`
}
