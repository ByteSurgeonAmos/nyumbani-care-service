package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Prescription models
type Prescription struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	DoctorID     uuid.UUID      `gorm:"type:uuid;not null" json:"doctor_id"`
	ImageURL     string         `json:"image_url"`
	Status       string         `json:"status"` // uploaded, processing, approved, rejected, dispensed
	Notes        string         `json:"notes"`
	PharmacyNotes string        `json:"pharmacy_notes"`
	Medications  []PrescriptionMedication `gorm:"foreignKey:PrescriptionID" json:"medications"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Doctor       User           `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
}

type PrescriptionMedication struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PrescriptionID uuid.UUID      `gorm:"type:uuid;not null" json:"prescription_id"`
	Name           string         `json:"name"`
	Dosage         string         `json:"dosage"`
	Frequency      string         `json:"frequency"`
	Duration       string         `json:"duration"`
	Instructions   string         `json:"instructions"`
	Price          float64        `json:"price"`
	Available      bool           `gorm:"default:true" json:"available"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Lab Work models
type LabTest struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Category    string         `json:"category"`
	SampleType  string         `json:"sample_type"` // blood, urine, stool, etc.
	PreparationInstructions string `json:"preparation_instructions"`
	TurnaroundTime string    `json:"turnaround_time"`
	Available   bool          `gorm:"default:true" json:"available"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type LabBooking struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	LabTestID       uuid.UUID      `gorm:"type:uuid;not null" json:"lab_test_id"`
	BookingDate     time.Time      `json:"booking_date"`
	PreferredTime   string         `json:"preferred_time"`
	Status          string         `json:"status"` // booked, confirmed, sample_collected, processing, completed, cancelled
	SampleCollectionMethod string `json:"sample_collection_method"` // home_collection, lab_visit
	Address         string         `json:"address"`
	ContactNumber   string         `json:"contact_number"`
	SpecialInstructions string    `json:"special_instructions"`
	PaymentStatus   string         `json:"payment_status"` // pending, paid, refunded
	TotalPrice      float64        `json:"total_price"`
	TrackingNumber  string         `json:"tracking_number"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LabTest         LabTest        `gorm:"foreignKey:LabTestID" json:"lab_test,omitempty"`
}

type LabResult struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	LabBookingID  uuid.UUID      `gorm:"type:uuid;not null" json:"lab_booking_id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	ResultData    string         `json:"result_data"` // JSON string with test values
	ReferenceRanges string       `json:"reference_ranges"`
	Interpretation string        `json:"interpretation"`
	DoctorComments string        `json:"doctor_comments"`
	ResultDate    time.Time      `json:"result_date"`
	VerifiedBy    string         `json:"verified_by"`
	ReportURL     string         `json:"report_url"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	LabBooking    LabBooking     `gorm:"foreignKey:LabBookingID" json:"lab_booking,omitempty"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Health Education models
type HealthArticle struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	Summary     string         `json:"summary"`
	Category    string         `json:"category"`
	Tags        []string       `gorm:"type:text[]" json:"tags"`
	AuthorID    uuid.UUID      `gorm:"type:uuid" json:"author_id"`
	ImageURL    string         `json:"image_url"`
	VideoURL    string         `json:"video_url"`
	ReadTime    int            `json:"read_time"` // in minutes
	Published   bool           `gorm:"default:false" json:"published"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Author      User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

type HealthQuiz struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Questions   []QuizQuestion `gorm:"foreignKey:QuizID" json:"questions"`
	TimeLimit   int            `json:"time_limit"` // in minutes
	PassingScore int           `json:"passing_score"`
	Published   bool           `gorm:"default:false" json:"published"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type QuizQuestion struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	QuizID    uuid.UUID      `gorm:"type:uuid;not null" json:"quiz_id"`
	Question  string         `json:"question"`
	Options   []string       `gorm:"type:text[]" json:"options"`
	CorrectAnswer int        `json:"correct_answer"` // index of correct option
	Explanation string       `json:"explanation"`
	Points    int            `gorm:"default:1" json:"points"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Telehealth models
type TelehealthSession struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PatientID     uuid.UUID      `gorm:"type:uuid;not null" json:"patient_id"`
	ProviderID    uuid.UUID      `gorm:"type:uuid;not null" json:"provider_id"`
	ProviderType  string         `json:"provider_type"` // doctor, nurse, pharmacist
	SessionType   string         `json:"session_type"` // video, audio, chat
	ScheduledAt   time.Time      `json:"scheduled_at"`
	Duration      int            `json:"duration"` // in minutes
	Status        string         `json:"status"` // scheduled, active, completed, cancelled, no_show
	MeetingURL    string         `json:"meeting_url"`
	Notes         string         `json:"notes"`
	Diagnosis     string         `json:"diagnosis"`
	Treatment     string         `json:"treatment"`
	Price         float64        `json:"price"`
	PaymentStatus string         `json:"payment_status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Patient       User           `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Provider      User           `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
}

// AI Symptom Checker models
type SymptomCheck struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Symptoms    []string       `gorm:"type:text[]" json:"symptoms"`
	Severity    string         `json:"severity"` // mild, moderate, severe
	Duration    string         `json:"duration"`
	Age         int            `json:"age"`
	Gender      string         `json:"gender"`
	Results     string         `json:"results"` // JSON string with AI analysis
	Recommendations string     `json:"recommendations"`
	UrgencyLevel string        `json:"urgency_level"` // low, medium, high, emergency
	FollowUpRequired bool      `json:"follow_up_required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CareSense Analytics models
type CareSenseAnalytics struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	AnalysisType    string         `json:"analysis_type"` // health_trends, risk_assessment, wellness_score
	DataSource      string         `json:"data_source"` // test_results, consultations, symptoms, etc.
	AnalysisData    string         `json:"analysis_data"` // JSON string with analysis results
	Insights        []string       `gorm:"type:text[]" json:"insights"`
	Recommendations []string       `gorm:"type:text[]" json:"recommendations"`
	Score           float64        `json:"score"` // wellness score or risk score
	GeneratedAt     time.Time      `json:"generated_at"`
	ExpiresAt       time.Time      `json:"expires_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Payment models
type Payment struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	OrderID       uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Amount        float64        `json:"amount"`
	Currency      string         `json:"currency"`
	Method        string         `json:"method"` // mpesa, paystack, bank_transfer
	Status        string         `json:"status"` // pending, processing, completed, failed, refunded
	TransactionID string         `json:"transaction_id"`
	PhoneNumber   string         `json:"phone_number,omitempty"`
	FailureReason string         `json:"failure_reason,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TestKitResult models
type TestKitResult struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	OrderID         uuid.UUID      `gorm:"type:uuid" json:"order_id"`
	TestKitID       uuid.UUID      `gorm:"type:uuid" json:"test_kit_id"`
	ImageURL        string         `json:"image_url"`
	Result          string         `json:"result"`          // positive, negative, inconclusive
	AIConfidence    float64        `json:"ai_confidence"`   // 0-1 confidence level
	DetectedMarkers []string       `gorm:"type:text[]" json:"detected_markers"`
	RecommendedSteps []string      `gorm:"type:text[]" json:"recommended_steps"`
	Notes           string         `json:"notes"`
	ReviewedBy      *uuid.UUID     `gorm:"type:uuid" json:"reviewed_by"` // Healthcare professional who reviewed
	ReviewNotes     string         `json:"review_notes"`
	Status          string         `json:"status"` // pending, reviewed, confirmed
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order           TestKitOrder   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	TestKit         TestKit        `gorm:"foreignKey:TestKitID" json:"test_kit,omitempty"`
	Reviewer        *User          `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

func (p *Prescription) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (pm *PrescriptionMedication) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}

func (lt *LabTest) BeforeCreate(tx *gorm.DB) error {
	if lt.ID == uuid.Nil {
		lt.ID = uuid.New()
	}
	return nil
}

func (lb *LabBooking) BeforeCreate(tx *gorm.DB) error {
	if lb.ID == uuid.Nil {
		lb.ID = uuid.New()
	}
	return nil
}

func (lr *LabResult) BeforeCreate(tx *gorm.DB) error {
	if lr.ID == uuid.Nil {
		lr.ID = uuid.New()
	}
	return nil
}

func (ha *HealthArticle) BeforeCreate(tx *gorm.DB) error {
	if ha.ID == uuid.Nil {
		ha.ID = uuid.New()
	}
	return nil
}

func (hq *HealthQuiz) BeforeCreate(tx *gorm.DB) error {
	if hq.ID == uuid.Nil {
		hq.ID = uuid.New()
	}
	return nil
}

func (qq *QuizQuestion) BeforeCreate(tx *gorm.DB) error {
	if qq.ID == uuid.Nil {
		qq.ID = uuid.New()
	}
	return nil
}

func (ts *TelehealthSession) BeforeCreate(tx *gorm.DB) error {
	if ts.ID == uuid.Nil {
		ts.ID = uuid.New()
	}
	return nil
}

func (sc *SymptomCheck) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}

func (ca *CareSenseAnalytics) BeforeCreate(tx *gorm.DB) error {
	if ca.ID == uuid.Nil {
		ca.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (tr *TestKitResult) BeforeCreate(tx *gorm.DB) error {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	return nil
}
