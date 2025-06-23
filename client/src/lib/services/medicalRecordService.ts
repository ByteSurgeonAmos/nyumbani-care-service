import apiClient from "../api/apiClient";
import type { AnalyzeTestKitResponse } from "../api/extendedServices";

interface MedicalRecord {
  id?: string;
  title: string;
  date: string;
  recordType: string;
  notes: string;
  doctor?: string;
  symptoms?: string;
  testResultId?: string;
  testResult?: any;
}

interface TestResultForMedicalRecord {
  testType: string;
  testDate: string;
  result: string;
  interpretation: string;
  labName: string;
  doctorNotes?: string;
}

export const medicalRecordService = {
  /**
   * Creates a medical record from a test result
   */
  async createFromTestResult(testResult: AnalyzeTestKitResponse): Promise<any> {
    try {
      // First create the basic medical record
      const recordData: MedicalRecord = {
        title: `Test Result: ${testResult.test_kit?.name || "Unknown Test"}`,
        date: testResult.created_at || new Date().toISOString(),
        recordType: "test_result",
        notes: `${testResult.interpretation}\n\nAdvice: ${testResult.advice}`,
        testResultId: testResult.id,
      };

      // Make API call to create medical record
      const response = await apiClient.post("/medical-records", recordData);

      // For now we'll use a simulated response since backend integration is pending
      return {
        success: true,
        recordId: response.data?.id || "simulated-id",
        message: "Test result saved to medical records",
      };
    } catch (error) {
      console.error("Error saving test result to medical records:", error);
      throw error;
    }
  },

  /**
   * Transforms a test kit analysis result to a medical record test result format
   */
  convertAnalysisToTestResult(
    analysis: AnalyzeTestKitResponse
  ): TestResultForMedicalRecord {
    return {
      testType: analysis.test_kit?.name || "Home Test Kit",
      testDate: analysis.created_at || new Date().toISOString(),
      result: analysis.result,
      interpretation: analysis.interpretation,
      labName: "Nyumbani Care Analysis",
      doctorNotes: analysis.advice,
    };
  },

  /**
   * Lists all medical records for the current user
   */
  async getAll(): Promise<any> {
    try {
      const response = await apiClient.get("/medical-records");
      return response.data;
    } catch (error) {
      console.error("Error fetching medical records:", error);
      throw error;
    }
  },

  /**
   * Gets a specific medical record by ID
   */
  async getById(id: string): Promise<any> {
    try {
      const response = await apiClient.get(`/medical-records/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching medical record ${id}:`, error);
      throw error;
    }
  },
};
