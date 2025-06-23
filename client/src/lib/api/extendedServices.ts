import apiClient from "./apiClient";
import type { TestKit, ListResponse } from "./apiServices";

// Lab Tests
export interface LabTest {
  id: string;
  name: string;
  description: string;
  price: number;
  duration: string;
  requirements: string;
  category: string;
  image_url?: string;
  created_at: string;
  updated_at: string;
}

export const labTestService = {
  async getAll(page = 1, limit = 10): Promise<ListResponse<LabTest>> {
    const response = await apiClient.get<LabTest[]>("/lab-tests", {
      params: { page, limit },
    });
    // Handle the case where the API returns an array directly
    const labTests = response.data || [];
    return {
      data: labTests,
      total: labTests.length,
      page,
      limit,
    };
  },
};

// Test Results
export interface TestResult {
  id: string;
  user_id: string;
  test_kit_id: string;
  test_kit?: TestKit;
  result_value: string;
  result_interpretation: string;
  result_date: string;
  image_url?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

interface CreateTestResultRequest {
  test_kit_id: string;
  result_value: string;
  result_interpretation?: string;
  result_date: string;
  image_url?: string;
  notes?: string;
}

export const testResultService = {
  async create(resultData: CreateTestResultRequest): Promise<TestResult> {
    const response = await apiClient.post<TestResult>(
      "/test-results",
      resultData
    );
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<TestResult>> {
    const response = await apiClient.get<ListResponse<TestResult>>(
      "/test-results",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async getById(id: string): Promise<TestResult> {
    const response = await apiClient.get<TestResult>(`/test-results/${id}`);
    return response.data;
  },

  async update(
    id: string,
    resultData: Partial<CreateTestResultRequest>
  ): Promise<TestResult> {
    const response = await apiClient.put<TestResult>(
      `/test-results/${id}`,
      resultData
    );
    return response.data;
  },
};

// Test Kit Results (for image analysis)
interface AnalyzeTestKitRequest {
  test_kit_id: string;
  image_data: string; // Base64 encoded image
}

export interface AnalyzeTestKitResponse {
  id: string;
  user_id?: string;
  test_kit_id?: string;
  test_kit?: TestKit;
  order_id?: string;
  order?: any; // Order relationship
  result: string;
  confidence: number;
  interpretation: string;
  advice: string;
  image_url?: string;
  created_at?: string;
  updated_at?: string;
  detected_markers?: string[];
  status?: string;
  review_notes?: string;
  recommended_steps?: string[];
  notes?: string;
}

export const testKitResultService = {
  async analyzeImage(
    data: AnalyzeTestKitRequest
  ): Promise<AnalyzeTestKitResponse> {
    const response = await apiClient.post<AnalyzeTestKitResponse>(
      "/test-kits/results/analyze",
      data
    );
    return response.data;
  },

  async getAll(
    page = 1,
    limit = 10
  ): Promise<ListResponse<AnalyzeTestKitResponse>> {
    const response = await apiClient.get<ListResponse<AnalyzeTestKitResponse>>(
      "/test-kits/results",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async getById(id: string): Promise<AnalyzeTestKitResponse> {
    const response = await apiClient.get<AnalyzeTestKitResponse>(
      `/test-kits/results/${id}`
    );
    return response.data;
  },
};

// Consultations
export interface Consultation {
  id: string;
  user_id: string;
  healthcare_provider_id?: string;
  consultation_type: string; // e.g. "video", "in-person", "chat"
  status: string; // e.g. "scheduled", "completed", "cancelled"
  scheduled_time: string;
  end_time?: string;
  reason: string;
  notes?: string;
  follow_up_needed: boolean;
  follow_up_date?: string;
  created_at: string;
  updated_at: string;
}

interface CreateConsultationRequest {
  healthcare_provider_id?: string;
  consultation_type: string;
  scheduled_time: string;
  reason: string;
}

export const consultationService = {
  async create(data: CreateConsultationRequest): Promise<Consultation> {
    const response = await apiClient.post<Consultation>("/consultations", data);
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<Consultation>> {
    const response = await apiClient.get<ListResponse<Consultation>>(
      "/consultations",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async getById(id: string): Promise<Consultation> {
    const response = await apiClient.get<Consultation>(`/consultations/${id}`);
    return response.data;
  },

  async update(id: string, data: Partial<Consultation>): Promise<Consultation> {
    const response = await apiClient.put<Consultation>(
      `/consultations/${id}`,
      data
    );
    return response.data;
  },
};

// Prescriptions
export interface Prescription {
  id: string;
  user_id: string;
  healthcare_provider_id: string;
  consultation_id?: string;
  medication_name: string;
  dosage: string;
  frequency: string;
  duration: string;
  notes?: string;
  status: string; // e.g. "active", "completed", "cancelled"
  created_at: string;
  updated_at: string;
}

interface CreatePrescriptionRequest {
  healthcare_provider_id: string;
  consultation_id?: string;
  medication_name: string;
  dosage: string;
  frequency: string;
  duration: string;
  notes?: string;
}

export const prescriptionService = {
  async create(data: CreatePrescriptionRequest): Promise<Prescription> {
    const response = await apiClient.post<Prescription>("/prescriptions", data);
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<Prescription>> {
    const response = await apiClient.get<ListResponse<Prescription>>(
      "/prescriptions",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async updateStatus(id: string, status: string): Promise<Prescription> {
    const response = await apiClient.put<Prescription>(
      `/prescriptions/${id}/status`,
      { status }
    );
    return response.data;
  },
};

// Lab Bookings
export interface LabBooking {
  id: string;
  user_id: string;
  lab_test_id: string;
  lab_test?: LabTest;
  scheduled_date: string;
  status: string; // e.g. "pending", "confirmed", "completed", "cancelled"
  payment_status: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

interface CreateLabBookingRequest {
  lab_test_id: string;
  scheduled_date: string;
  notes?: string;
}

export const labBookingService = {
  async create(data: CreateLabBookingRequest): Promise<LabBooking> {
    const response = await apiClient.post<LabBooking>("/lab-bookings", data);
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<LabBooking>> {
    const response = await apiClient.get<ListResponse<LabBooking>>(
      "/lab-bookings",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async updateStatus(id: string, status: string): Promise<LabBooking> {
    const response = await apiClient.put<LabBooking>(
      `/lab-bookings/${id}/status`,
      { status }
    );
    return response.data;
  },
};

// Telehealth Sessions
export interface TelehealthSession {
  id: string;
  user_id: string;
  provider_id?: string;
  consultation_id?: string;
  session_token: string;
  start_time: string;
  end_time?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

interface CreateTelehealthSessionRequest {
  provider_id?: string;
  consultation_id?: string;
}

export const telehealthService = {
  async create(
    data: CreateTelehealthSessionRequest
  ): Promise<TelehealthSession> {
    const response = await apiClient.post<TelehealthSession>(
      "/telehealth/sessions",
      data
    );
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<TelehealthSession>> {
    const response = await apiClient.get<ListResponse<TelehealthSession>>(
      "/telehealth/sessions",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },
};

// Symptom Checker
export interface SymptomCheck {
  id: string;
  user_id: string;
  symptoms: string[];
  risk_level: string;
  recommendations: string[];
  created_at: string;
  updated_at: string;
}

interface CreateSymptomCheckRequest {
  symptoms: string[];
}

export const symptomCheckService = {
  async create(data: CreateSymptomCheckRequest): Promise<SymptomCheck> {
    const response = await apiClient.post<SymptomCheck>(
      "/symptoms/check",
      data
    );
    return response.data;
  },

  async getHistory(page = 1, limit = 10): Promise<ListResponse<SymptomCheck>> {
    const response = await apiClient.get<ListResponse<SymptomCheck>>(
      "/symptoms/history",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },
};

// CareSense Analytics
export interface CareSenseAnalytics {
  id: string;
  user_id: string;
  health_score: number;
  risk_factors: string[];
  recommendations: string[];
  insights: any;
  created_at: string;
  updated_at: string;
}

export const careSenseService = {
  async generateAnalytics(): Promise<CareSenseAnalytics> {
    const response = await apiClient.post<CareSenseAnalytics>(
      "/caresense/analytics"
    );
    return response.data;
  },

  async getAnalytics(): Promise<CareSenseAnalytics> {
    const response = await apiClient.get<CareSenseAnalytics>(
      "/caresense/analytics"
    );
    return response.data;
  },
};

// Payments
export interface Payment {
  id: string;
  user_id: string;
  amount: number;
  currency: string;
  payment_method: string;
  status: string;
  reference: string;
  description: string;
  created_at: string;
  updated_at: string;
}

interface ProcessPaystackPaymentRequest {
  email: string;
  amount: number;
  reference?: string;
  callback_url?: string;
  order_id?: string;
}

interface PaystackResponse {
  authorization_url: string;
  access_code: string;
  reference: string;
}

export const paymentService = {
  async processPaystackPayment(
    data: ProcessPaystackPaymentRequest
  ): Promise<PaystackResponse> {
    const response = await apiClient.post<PaystackResponse>(
      "/payments/paystack",
      data
    );
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<Payment>> {
    const response = await apiClient.get<ListResponse<Payment>>("/payments", {
      params: { page, limit },
    });
    return response.data;
  },

  async getById(id: string): Promise<Payment> {
    const response = await apiClient.get<Payment>(`/payments/${id}`);
    return response.data;
  },
};

// File Upload
interface FileUploadResponse {
  file_url: string;
}

export const uploadService = {
  async uploadFile(file: File): Promise<FileUploadResponse> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await apiClient.post<FileUploadResponse>(
      "/uploads/file",
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }
    );
    return response.data;
  },
};
