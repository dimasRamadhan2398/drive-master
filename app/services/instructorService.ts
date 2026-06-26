import type {
  ApiResponse,
  PaginatedResponse,
} from "~/composables/useApiClients";

// Instructor profile details (nested in API response)
export interface InstructorDetails {
  userId: string;
  bnspCertificateNumber?: string;
  numberOfStudents: number;
  sessionsCompleted: number;
  averageRating: number;
  description?: string;
  specialization?: string;
  licenseNumber?: string;
  yearsOfExperience: number;
  licenseExpiry?: string;
  isActive: boolean;
  photoURL?: string;
  bio?: string;
}

// Instructor user data from API
export interface InstructorApiResponse {
  userId: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  phoneNumber?: string;
  image?: string;
  dateOfBirth?: string;
  address?: string;
  roleId: number;
  role: {
    id: number;
    name: string;
  };
  instructorProfile?: InstructorDetails;
}

export interface Instructor {
  userId: string;
  name: string;
  email: string;
  phone: string;
  image?: string;
  bio?: string;
  status: "active" | "inactive" | "pending";
  joinedDate: string;
  totalStudents: number;
  rating: number;
  certifications: string[];
  yearsOfExperience: number;
}

export interface MediaMetadata {
  filename: string;
  mimeType: string;
  size: number;
  url: string;
  uploadedAt: string;
}

export interface UpdateInstructorData {
  firstName?: string;
  lastName?: string;
  phoneNumber?: string;
  address?: string;
  dateOfBirth?: string;
  bio?: string;
}

// Request data for creating a new instructor (combined user + profile)
export interface CreateInstructorRequest {
  // User fields
  firstName: string;
  lastName: string;
  username: string;
  password: string;
  email: string;
  phoneNumber: string;
  dateOfBirth?: string;
  address?: string;

  // Instructor profile fields
  licenseNumber?: string;
  licenseExpiry?: string;
  bnspCertificateNumber?: string;
  yearsOfExperience?: number;
  specialization?: string;
  description?: string;
}

// Legacy interfaces for backwards compatibility
export interface CreateUserRequest {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  phoneNumber: string;
  password: string;
  roleId: number;
}

export interface CreateInstructorProfileRequest {
  specialization?: string;
  description?: string;
  licenseNumber?: string;
  licenseExpiry?: string;
  bnspCertificateNumber?: string;
  yearsOfExperience?: number;
}

export interface PaginatedInstructorsResult {
  instructors: Instructor[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// Map API response to Instructor model
export const mapApiToInstructor = (item: InstructorApiResponse): Instructor => {
  const profile = item.instructorProfile;
  return {
    userId: item.userId,
    name: `${item.firstName} ${item.lastName}`.trim(),
    email: item.email,
    phone: item.phoneNumber || "",
    image: profile?.photoURL || item.image,
    bio: profile?.bio || "",
    status: profile?.isActive ? "active" : "inactive",
    joinedDate: profile?.yearsOfExperience
      ? `${profile.yearsOfExperience} years`
      : new Date().toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
    totalStudents: profile?.numberOfStudents || 0,
    rating: profile?.averageRating || 0,
    certifications: profile?.bnspCertificateNumber
      ? [`BNSP: ${profile.bnspCertificateNumber}`]
      : [],
    yearsOfExperience: profile?.yearsOfExperience || 0,
  };
};

// Generate numeric ID from UUID
function hashCode(str: string): number {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash;
  }
  return Math.abs(hash);
}

export const instructorService = {
  // GET /instructors - Get instructor lists
  async fetchAll(
    params: { page?: number; limit?: number; search?: string } = {},
  ): Promise<PaginatedInstructorsResult> {
    const { user, extractPaginatedData } = useApiClients();

    const queryParams = new URLSearchParams();
    if (params.page) queryParams.set("page", String(params.page));
    if (params.limit) queryParams.set("limit", String(params.limit));
    if (params.search) queryParams.set("search", params.search);

    const queryString = queryParams.toString();
    const url = `/instructors/all${queryString ? `?${queryString}` : ""}`;

    const response = await user<PaginatedResponse<InstructorApiResponse>>(url, {
      method: "GET",
    });

    const { data, pagination } = extractPaginatedData(response);
    return {
      instructors: Array.isArray(data) ? data.map(mapApiToInstructor) : [],
      pagination,
    };
  },

  // GET /instructors/:id/profile - Get instructor profile
  async fetchById(userId: string): Promise<Instructor | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<InstructorApiResponse>>(
        `/instructors/${userId}/profile`,
        {
          method: "GET",
        },
      );
      return mapApiToInstructor(extractData(response));
    } catch {
      return null;
    }
  },

  // PUT /instructors/:id/profile - Update instructor profile
  async update(
    userId: string,
    data: UpdateInstructorData,
  ): Promise<Instructor | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<InstructorApiResponse>>(
        `/instructors/${userId}/profile`,
        {
          method: "PUT",
          body: data,
        },
      );
      return mapApiToInstructor(extractData(response));
    } catch {
      return null;
    }
  },

  // POST /instructors/:id/media/upload - Upload profile picture
  async uploadProfilePic(
    userId: string,
    file: File,
  ): Promise<{ url: string; metadata: MediaMetadata } | null> {
    const { user, extractData } = useApiClients();
    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("folder", "instructors/profiles");
      formData.append("filename", `${userId}_${Date.now()}`);

      const response = await user<ApiResponse<{ url: string }>>(
        `/instructors/${userId}/media/upload`,
        {
          method: "POST",
          body: formData,
          headers: {
            "Content-Type": "multipart/form-data",
          },
        },
      );
      return { url: extractData(response).url, metadata: {} as MediaMetadata };
    } catch {
      return null;
    }
  },

  // POST /instructors/:id/media/upload-base64 - Upload base64 media
  async uploadBase64Media(
    userId: string,
    base64Data: string,
    mimeType: string,
  ): Promise<{ url: string; metadata: MediaMetadata } | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<{ url: string }>>(
        `/instructors/${userId}/media/upload-base64`,
        {
          method: "POST",
          body: { data: base64Data, mimeType },
        },
      );
      return { url: extractData(response).url, metadata: {} as MediaMetadata };
    } catch {
      return null;
    }
  },

  // DELETE /instructors/:id/media - Delete profile picture
  async deleteProfilePic(userId: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user(`/instructors/${userId}/media`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // DELETE /instructors/:id - Delete instructor fully (user + profile)
  async delete(userId: string): Promise<boolean> {
    const { user } = useApiClients();
    try {
      await user(`/instructors/${userId}`, { method: "DELETE" });
      return true;
    } catch {
      return false;
    }
  },

  // GET /instructors/:id/media/metadata - Get media metadata
  async getMediaMetadata(userId: string): Promise<MediaMetadata | null> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<MediaMetadata>>(
        `/instructors/${userId}/media/metadata`,
        {
          method: "GET",
        },
      );
      return extractData(response);
    } catch {
      return null;
    }
  },

  // POST /instructors/register - Create new instructor (user + profile in one call)
  async create(data: CreateInstructorRequest): Promise<Instructor | null> {
    const { user, extractData } = useApiClients();
    console.log("[instructorService.create] Sending request with data:", data);

    try {
      const response = await user<ApiResponse<InstructorApiResponse>>(
        "/instructors/register",
        {
          method: "POST",
          body: data,
        },
      );

      console.log("[instructorService.create] Raw response:", response);

      const extracted = extractData(response);
      console.log("[instructorService.create] Extracted data:", extracted);

      const mapped = mapApiToInstructor(extracted);
      console.log("[instructorService.create] Mapped instructor:", mapped);

      return mapped;
    } catch (err) {
      console.error("[instructorService.create] Error:", err);
      return null;
    }
  },
};
