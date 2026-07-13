import type { ApiResponse } from "~/composables/useApiClients";
import { useAuthStore } from "~/stores/auth";

export interface MemberCertificateResponse {
  id: string;
  memberId: string;
  memberName: string;
  memberEmail: string;
  packageName: string;
  certNumber: string;
  issuedDate: string;
  completedAt: string;
  status: "eligible" | "issued" | "expired" | "revoked";
}

export const certificateService = {
  async getMemberCertificates(userId: string): Promise<MemberCertificateResponse[]> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<MemberCertificateResponse[]>>(
        `/members/${userId}/certificates`,
        { method: "GET" }
      );
      return extractData(response) || [];
    } catch {
      return [];
    }
  },

  async getAllCertificates(): Promise<MemberCertificateResponse[]> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<MemberCertificateResponse[]>>(
        `/certificates`,
        { method: "GET" }
      );
      return extractData(response) || [];
    } catch (e) {
      console.error("Error fetching all certificates:", e);
      return [];
    }
  },

  async issueCertificate(input: {
    memberId: string;
    packageId: string;
    packageName: string;
    entitlementId: string;
  }): Promise<any> {
    const { user, extractData } = useApiClients();
    const response = await user<ApiResponse<any>>(`/certificates`, {
      method: "POST",
      body: input,
    });
    return extractData(response);
  },

  async revokeCertificate(certId: string): Promise<void> {
    const { user } = useApiClients();
    await user(`/certificates/${certId}`, {
      method: "DELETE",
    });
  },

  async getStats(): Promise<any> {
    const { user, extractData } = useApiClients();
    try {
      const response = await user<ApiResponse<any>>(`/certificates/stats`, {
        method: "GET",
      });
      return extractData(response);
    } catch (e) {
      console.error("Error fetching certificate stats:", e);
      return null;
    }
  },

  async getCertificatePDFBlobUrl(certId: string): Promise<string> {
    const authStore = useAuthStore();
    const token = authStore.accessToken;
    const baseURL = useRuntimeConfig().public.apiBase || "";

    const response = await fetch(`${baseURL}/certificates/${certId}/pdf`, {
      method: "GET",
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to fetch certificate PDF");
    }

    const blob = await response.blob();
    return window.URL.createObjectURL(blob);
  },

  async downloadCertificatePDF(certId: string, certNumber: string): Promise<void> {
    const authStore = useAuthStore();
    const token = authStore.accessToken;
    const baseURL = useRuntimeConfig().public.apiBase || "";

    const response = await fetch(`${baseURL}/certificates/${certId}/pdf`, {
      method: "GET",
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to download certificate PDF");
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `certificate-${certNumber}.pdf`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  },

  async downloadCertificate(userId: string, certId: string): Promise<void> {
    const authStore = useAuthStore();
    const token = authStore.accessToken;
    const baseURL = useRuntimeConfig().public.apiBase || "";

    const response = await fetch(`${baseURL}/members/${userId}/certificates/${certId}/download`, {
      method: "GET",
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to download certificate");
    }

    const blob = await response.blob();
    const contentDisposition = response.headers.get("Content-Disposition");
    let filename = `certificate-${certId}.pdf`;

    if (contentDisposition) {
      const match = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
      if (match && match[1]) {
        filename = match[1].replace(/['"]/g, "");
      }
    }

    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  },
};