import axios from 'axios';

// The backend URL. In development on a physical device, replace with your
// machine's LAN IP. The EXPO_PUBLIC_ prefix exposes it via process.env.
export const BASE_URL =
  process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080';

const api = axios.create({
  baseURL: BASE_URL,
  timeout: 15_000,
  headers: { 'Content-Type': 'application/json' },
});

export interface Log {
  id: string;
  odometer_km: number;
  fuel_level?: number | null;
  notes: string;
  image_path?: string;
  parsed_by_ocr: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateLogPayload {
  odometer_km: number;
  fuel_level?: number | null;
  notes?: string;
}

export interface OCRResult {
  odometer_km?: number | null;
  fuel_level?: number | null;
  confidence: number;
}

export interface UploadImageResponse {
  log: Log;
  ocr: OCRResult;
}

export const logsApi = {
  list: (limit = 20, offset = 0) =>
    api.get<Log[]>('/api/logs', { params: { limit, offset } }).then(r => r.data),

  get: (id: string) => api.get<Log>(`/api/logs/${id}`).then(r => r.data),

  create: (payload: CreateLogPayload) =>
    api.post<Log>('/api/logs', payload).then(r => r.data),

  delete: (id: string) => api.delete(`/api/logs/${id}`),

  uploadImage: (id: string, imageUri: string) => {
    const form = new FormData();
    // React Native's FormData accepts an object with uri/type/name.
    form.append('image', {
      uri: imageUri,
      type: 'image/jpeg',
      name: 'dashboard.jpg',
    } as any);
    return api
      .post<UploadImageResponse>(`/api/logs/${id}/image`, form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then(r => r.data);
  },
};
