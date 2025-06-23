import apiClient from "./apiClient";

export interface TestKit {
  id: string;
  name: string;
  description: string;
  price: number;
  image_url: string;
  instructions: string;
  category: string;
  in_stock: boolean;
  requires_analysis?: boolean;
  analysis_type?: string;
  reference_range?: string;
  usage_notes?: string;
  expiry_date?: string;
  manufacturer?: string;
  created_at: string;
  updated_at: string;
}

export interface ListResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export const testKitService = {
  async getAll(page = 1, limit = 10): Promise<ListResponse<TestKit>> {
    const response = await apiClient.get<TestKit[]>("/test-kits", {
      params: { page, limit },
    });
    const testKits = response.data || [];
    return {
      data: testKits,
      total: testKits.length,
      page,
      limit,
    };
  },

  async getById(id: string): Promise<TestKit> {
    const response = await apiClient.get<TestKit>(`/test-kits/${id}`);
    return response.data;
  },
};

export interface TestKitOrder {
  id: string;
  user_id: string;
  status: string;
  total_amount: number;
  shipping_address: string;
  payment_status: string;
  items: TestKitOrderItem[];
  created_at: string;
  updated_at: string;
}

interface TestKitOrderItem {
  id: string;
  order_id: string;
  test_kit_id: string;
  test_kit: TestKit;
  quantity: number;
  unit_price: number;
  subtotal: number;
}

interface CreateOrderRequest {
  test_kit_ids: { id: string; quantity: number }[];
  shipping_address: string;
  payment_method: string;
}

export const orderService = {
  async create(orderData: CreateOrderRequest): Promise<TestKitOrder> {
    const response = await apiClient.post<TestKitOrder>("/orders", orderData);
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<TestKitOrder>> {
    const response = await apiClient.get<ListResponse<TestKitOrder>>(
      "/orders",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async getById(id: string): Promise<TestKitOrder> {
    const response = await apiClient.get<TestKitOrder>(`/orders/${id}`);
    return response.data;
  },

  async validateInventory(items: { id: string; quantity: number }[]): Promise<{
    valid: boolean;
    errors?: Array<{
      id: string;
      name: string;
      requested: number;
      available: number;
    }>;
  }> {
    const response = await apiClient.post<{
      valid: boolean;
      errors?: Array<{
        id: string;
        name: string;
        requested: number;
        available: number;
      }>;
    }>("/orders/validate-inventory", { items });
    return response.data;
  },

  async cancelOrder(id: string, reason: string): Promise<TestKitOrder> {
    const response = await apiClient.post<TestKitOrder>(
      `/orders/${id}/cancel`,
      { reason }
    );
    return response.data;
  },
};

export interface HealthArticle {
  id: string;
  title: string;
  content: string;
  category: string;
  image_url: string;
  author?: {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
  };
  author_id?: string;
  summary?: string;
  tags?: string[];
  video_url?: string;
  read_time?: number;
  published?: boolean;
  view_count?: number;
  published_at?: string;
  created_at: string;
  updated_at: string;
}

export const healthEducationService = {
  async getAllArticles(
    page = 1,
    limit = 10,
    category = ""
  ): Promise<ListResponse<HealthArticle>> {
    const params: any = { page, limit };
    if (category) {
      params.category = category;
    }

    const response = await apiClient.get<{
      articles: HealthArticle[];
      page: number;
      total: number;
      total_pages: number;
    }>("/health-education/articles", {
      params,
    });

    return {
      data: response.data.articles || [],
      total: response.data.total || 0,
      page: response.data.page || page,
      limit: limit,
    };
  },

  async getArticleById(id: string): Promise<HealthArticle> {
    const response = await apiClient.get<HealthArticle>(
      `/health-education/articles/${id}`
    );
    return response.data;
  },
};

export interface MedicalRecord {
  id: string;
  user_id: string;
  record_type: string;
  title: string;
  description: string;
  file_url: string;
  date: string;
  created_at: string;
  updated_at: string;
}

interface CreateMedicalRecordRequest {
  record_type: string;
  title: string;
  description: string;
  file_url?: string;
  date: string;
}

export const medicalRecordService = {
  async create(recordData: CreateMedicalRecordRequest): Promise<MedicalRecord> {
    const response = await apiClient.post<MedicalRecord>(
      "/medical-records",
      recordData
    );
    return response.data;
  },

  async getAll(page = 1, limit = 10): Promise<ListResponse<MedicalRecord>> {
    const response = await apiClient.get<ListResponse<MedicalRecord>>(
      "/medical-records",
      {
        params: { page, limit },
      }
    );
    return response.data;
  },

  async getById(id: string): Promise<MedicalRecord> {
    const response = await apiClient.get<MedicalRecord>(
      `/medical-records/${id}`
    );
    return response.data;
  },

  async update(
    id: string,
    recordData: Partial<CreateMedicalRecordRequest>
  ): Promise<MedicalRecord> {
    const response = await apiClient.put<MedicalRecord>(
      `/medical-records/${id}`,
      recordData
    );
    return response.data;
  },
};
