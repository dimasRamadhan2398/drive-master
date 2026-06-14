import type { ApiResponse } from "~/composables/useApiClients";

export interface MemberCertificateResponse {
  id: string;
  memberId: string;
  memberName: string;
  memberEmail: string;
  packageName: string;
  certNumber: string;
  issuedDate: string;
  completedAt: string;
  status: "eligible" | "issued" | "expired";
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

  async downloadCertificate(userId: string, certId: string): Promise<void> {
    const { user } = useApiClients();
    const { useAuthToken } = useAuth();

    const token = useAuthToken();
    const baseURL = useRuntimeConfig().public.apiBaseUrl || "";

    const response = await fetch(`${baseURL}/api/v1/members/${userId}/certificates/${certId}/download`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token.value}`,
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