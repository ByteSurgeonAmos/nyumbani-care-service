package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
	"github.com/sashabaranov/go-openai"
)

type AIService struct {
	config     *config.ExternalConfig
	openaiClient *openai.Client
}

func NewAIService(cfg *config.ExternalConfig) *AIService {
	client := openai.NewClient(cfg.ChatGPTAPIKey)
	return &AIService{
		config: cfg,
		openaiClient: client,
	}
}

type SymptomCheckRequest struct {
	Symptoms []string `json:"symptoms"`
	Age      int      `json:"age"`
	Gender   string   `json:"gender"`
}

type SymptomCheckResponse struct {
	PossibleConditions []PossibleCondition `json:"possible_conditions"`
	Recommendations    []string           `json:"recommendations"`
	Urgency           string             `json:"urgency"`
	Confidence        float64            `json:"confidence"`
}

type PossibleCondition struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
	Description string  `json:"description"`
}

type AnalyticsRequest struct {
	UserID       string                 `json:"user_id"`
	DataType     string                 `json:"data_type"`
	TimeRange    string                 `json:"time_range"`
	Data         map[string]interface{} `json:"data"`
}

type AnalyticsResponse struct {
	Trends    []Trend  `json:"trends"`
	Patterns  []string `json:"patterns"`
	Insights  []string `json:"insights"`
	RiskScore float64  `json:"risk_score"`
}

type Trend struct {
	Metric string    `json:"metric"`
	Values []float64 `json:"values"`
	Dates  []string  `json:"dates"`
}

// TestKitResultRequest represents a request to analyze a test kit result image
type TestKitResultRequest struct {
	TestKitType string `json:"test_kit_type"`
	ImageURL    string `json:"image_url"`
}

// TestKitResultResponse represents the analysis of a test kit result
type TestKitResultResponse struct {
	Result           string   `json:"result"`           // positive, negative, inconclusive
	Confidence       float64  `json:"confidence"`       // 0-1 confidence level
	DetectedMarkers  []string `json:"detected_markers"` // Any markers detected in the test
	RecommendedSteps []string `json:"recommended_steps"`
	Notes            string   `json:"notes"` // Any additional notes or observations
}

// Legacy function - replaced by the new one above

func (ai *AIService) AnalyzeSymptoms(request SymptomCheckRequest) (*SymptomCheckResponse, error) {
	if ai.config.ChatGPTAPIKey == "" {
		// Return mock response for development
		return ai.mockSymptomAnalysis(request), nil
	}

	// Create a prompt for GPT
	symptomsJSON, _ := json.Marshal(request.Symptoms)
	prompt := fmt.Sprintf(
		"I need a medical symptom analysis based on these symptoms: %s. "+
		"The patient is a %d year old %s. "+
		"Please provide possible conditions, recommendations, and urgency level. "+
		"Format your response as JSON with the following structure: "+
		"{\"possible_conditions\": [{\"name\": string, \"probability\": float, \"description\": string}], "+
		"\"recommendations\": [string], \"urgency\": string, \"confidence\": float}",
		string(symptomsJSON),
		request.Age,
		request.Gender,
	)

	// Create ChatGPT request
	ctx := context.Background()
	resp, err := ai.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: ai.config.ChatGPTModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a healthcare AI assistant. Respond only with valid JSON following the specified structure.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2, // Lower temperature for more consistent outputs
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatGPT API error: %v", err)
	}
		// Parse the response
	var response SymptomCheckResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse ChatGPT response: %v", err)
	}

	return &response, nil
}

func (ai *AIService) GenerateHealthAnalytics(request AnalyticsRequest) (*AnalyticsResponse, error) {
	if ai.config.ChatGPTAPIKey == "" {
		// Return mock response for development
		return ai.mockHealthAnalytics(request), nil
	}

	// Create data representation for GPT
	dataJSON, _ := json.Marshal(request.Data)
	prompt := fmt.Sprintf(
		"Analyze the following health data for user %s over %s. Data type: %s. "+
		"Data: %s "+
		"Generate health insights, trends, patterns, and risk assessment. "+
		"Format your response as JSON with the following structure: "+
		"{\"trends\": [{\"metric\": string, \"values\": [float], \"dates\": [string]}], "+
		"\"patterns\": [string], \"insights\": [string], \"risk_score\": float}",
		request.UserID,
		request.TimeRange,
		request.DataType,
		string(dataJSON),
	)

	// Create ChatGPT request
	ctx := context.Background()
	resp, err := ai.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: ai.config.ChatGPTModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a healthcare analytics AI. Respond only with valid JSON following the specified structure.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatGPT API error: %v", err)
	}
		// Parse the response
	var response AnalyticsResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse ChatGPT response: %v", err)
	}

	return &response, nil
}

// Mock implementations for development
func (ai *AIService) mockSymptomAnalysis(request SymptomCheckRequest) *SymptomCheckResponse {
	// Simple rule-based mock analysis
	conditions := []PossibleCondition{}
	urgency := "low"
		// Basic symptom matching	
	if contains(request.Symptoms, "fever") && contains(request.Symptoms, "cough") {
		conditions = append(conditions, PossibleCondition{
			Name:        "Upper Respiratory Infection",
			Probability: 0.75,
			Description: "Common viral infection affecting the upper respiratory tract",
		})
		urgency = "medium"
	}
	
	if contains(request.Symptoms, "chest pain") {
		conditions = append(conditions, PossibleCondition{
			Name:        "Chest Pain Syndrome",
			Probability: 0.6,
			Description: "Various conditions that can cause chest discomfort",
		})
		urgency = "high"
	}
	
	if len(conditions) == 0 {
		conditions = append(conditions, PossibleCondition{
			Name:        "General Symptoms",
			Probability: 0.4,
			Description: "Common symptoms that may indicate various conditions",
		})
	}

	recommendations := []string{
		"Monitor symptoms for 24-48 hours",
		"Stay hydrated and get adequate rest",
		"Consider consulting a healthcare provider if symptoms worsen",
	}

	if urgency == "high" {
		recommendations = []string{
			"Seek immediate medical attention",
			"Consider visiting an emergency room",
			"Do not delay seeking professional help",
		}
	}

	return &SymptomCheckResponse{
		PossibleConditions: conditions,
		Recommendations:    recommendations,
		Urgency:           urgency,
		Confidence:        0.7,
	}
}

func (ai *AIService) mockHealthAnalytics(request AnalyticsRequest) *AnalyticsResponse {
	return &AnalyticsResponse{
		Trends: []Trend{
			{
				Metric: "blood_pressure",
				Values: []float64{120, 125, 118, 122, 119},
				Dates:  []string{"2025-05-01", "2025-05-08", "2025-05-15", "2025-05-22", "2025-05-29"},
			},
			{
				Metric: "heart_rate",
				Values: []float64{72, 75, 68, 70, 74},
				Dates:  []string{"2025-05-01", "2025-05-08", "2025-05-15", "2025-05-22", "2025-05-29"},
			},
		},
		Patterns: []string{
			"Blood pressure shows stable trend within normal range",
			"Heart rate variability indicates good cardiovascular health",
			"Sleep patterns show improvement over the past month",
		},
		Insights: []string{
			"Your health metrics show positive trends",
			"Consider maintaining current lifestyle habits",
			"Regular monitoring is recommended",
		},
		RiskScore: 0.2, // Low risk
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Update the symptom check handler to use AI service
func (ai *AIService) ProcessSymptomCheck(symptomCheck *models.SymptomCheck) error {
	request := SymptomCheckRequest{
		Symptoms: symptomCheck.Symptoms,
		Age:      25, // Default age if not provided
		Gender:   "unknown",
	}

	response, err := ai.AnalyzeSymptoms(request)
	if err != nil {
		return err
	}

	// Update symptom check with AI results
	conditionsJson, _ := json.Marshal(response.PossibleConditions)
	
	// Store the conditions as a JSON string
	symptomCheck.Results = string(conditionsJson)
	// Store first recommendation only - API compatibility
	if len(response.Recommendations) > 0 {
		symptomCheck.Recommendations = response.Recommendations[0]
	}
	symptomCheck.UrgencyLevel = response.Urgency
	// Store additional recommendations in Results field as a JSON string
	if len(response.Recommendations) > 1 {
		// Append additional recommendations to the Results field
		additionalRecommendations := strings.Join(response.Recommendations[1:], "; ")
		symptomCheck.Results += fmt.Sprintf(", \"additional_recommendations\": \"%s\"", additionalRecommendations)
	}

	return nil
}

// AnalyzeTestKitResult analyzes a test kit result image
func (s *AIService) AnalyzeTestKitResult(req *TestKitResultRequest) (*TestKitResultResponse, error) {
	if s.openaiClient == nil {
		return nil, fmt.Errorf("OpenAI client not initialized")
	}

	// Create the ChatGPT prompt for analyzing the test kit result image
	systemPrompt := "You are a medical diagnostic assistant specialized in analyzing test kit results from images. " +
		"Provide a detailed analysis of the test kit image including: " +
		"1. Whether the result is positive, negative, or inconclusive " +
		"2. Your confidence level in the interpretation (0-1) " +
		"3. Any markers or indicators you can detect " +
		"4. Recommended next steps based on the result " +
		"5. Any additional notes or observations. " +
		"Return your analysis in JSON format that matches the TestKitResultResponse structure."

	userPrompt := fmt.Sprintf("I'm analyzing a %s test kit result from this image: %s. Please interpret the result.", 
		req.TestKitType, req.ImageURL)

	completion, err := s.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.config.ChatGPTModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0.2, // Lower temperature for more deterministic results
			MaxTokens:   1000,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to analyze test kit result: %w", err)
	}

	if len(completion.Choices) == 0 || completion.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("no response from AI service")
	}

	// Parse the response
	content := completion.Choices[0].Message.Content
	
	// Extract JSON from the response
	jsonStr := extractJSON(content)
	
	var response TestKitResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		// If unmarshalling fails, return a default response with the raw content
		return &TestKitResultResponse{
			Result:           "inconclusive",
			Confidence:       0.0,
			DetectedMarkers:  []string{},
			RecommendedSteps: []string{"Consult a healthcare professional"},
			Notes:            fmt.Sprintf("Unable to parse result automatically. Raw response: %s", content),
		}, nil
	}

	return &response, nil
}

// extractJSON extracts JSON content from a string that might contain other text
func extractJSON(content string) string {
	// Find the start of the JSON object
	start := strings.Index(content, "{")
	if start == -1 {
		return "{}"
	}

	// Find the end of the JSON object
	end := strings.LastIndex(content, "}")
	if end == -1 || end < start {
		return "{}"
	}

	// Extract the JSON part
	return content[start : end+1]
}
