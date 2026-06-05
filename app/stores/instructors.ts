import { defineStore } from "pinia";
import { instructorService } from "~/services/instructorService";
import type {
  Instructor,
  UpdateInstructorData,
  CreateInstructorRequest,
} from "~/services/instructorService";

interface InstructorsState {
  instructors: Instructor[];
  currentInstructor: Instructor | null;
  isLoading: boolean;
  error: string | null;
  activeId: string | null;
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
  searchQuery: string;
  useServerPagination: boolean;
}

const initialInstructors: Instructor[] = [
  {
    userId: "inst-001",
    name: "Ahmad Wijaya",
    email: "ahmad@driveschool.com",
    phone: "081234567890",
    image: "/images/instructors/ahmad.jpg",
    bio: "Professional driving instructor with 5 years experience.",
    status: "active",
    joinedDate: "Jan 15, 2026",
    totalStudents: 45,
    rating: 4.8,
    certifications: ["Certified Driving Instructor", "Safety Expert"],
    yearsOfExperience: 5,
  },
  {
    userId: "inst-002",
    name: "Dewi Kartika",
    email: "dewi@driveschool.com",
    phone: "081234567891",
    image: "/images/instructors/dewi.jpg",
    bio: "Patient and friendly instructor specializing in beginners.",
    status: "active",
    joinedDate: "Feb 20, 2026",
    totalStudents: 32,
    rating: 4.9,
    certifications: ["Certified Driving Instructor", "Beginner Specialist"],
    yearsOfExperience: 3,
  },
  {
    userId: "inst-003",
    name: "Budi Santoso",
    email: "budi@driveschool.com",
    phone: "081234567892",
    bio: "Experienced instructor with background in defensive driving.",
    status: "active",
    joinedDate: "Dec 5, 2025",
    totalStudents: 58,
    rating: 4.7,
    certifications: ["Certified Driving Instructor", "Defensive Driving"],
    yearsOfExperience: 7,
  },
];

export const useInstructorsStore = defineStore("instructors", {
  state: (): InstructorsState => ({
    instructors: [],
    currentInstructor: null,
    isLoading: false,
    activeId: null,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      totalPages: 0,
    },
    searchQuery: "",
    useServerPagination: true,
  }),

  getters: {
    activeInstructors: (state) =>
      state.instructors.filter((i) => i.status === "active"),
    inactiveInstructors: (state) =>
      state.instructors.filter((i) => i.status === "inactive"),
    pendingInstructors: (state) =>
      state.instructors.filter((i) => i.status === "pending"),
    totalInstructors: (state) => state.instructors.length,
    topRatedInstructors: (state) =>
      [...state.instructors].sort((a, b) => b.rating - a.rating).slice(0, 5),
    getInstructorByUserId: (state) => (userId: string) =>
      state.instructors.find((i) => i.userId === userId),
  },

  actions: {
    setActiveID(id: string | null) {
      this.activeId = id;
    },

    getActiveID() {
      return this.activeId;
    },

    async fetchInstructors(page = 1, resetPage = true) {
      this.isLoading = true;
      this.error = null;

      try {
        if (resetPage) {
          this.pagination.page = page;
        }

        const params: { page: number; limit: number; search?: string } = {
          page: this.pagination.page,
          limit: this.pagination.limit,
        };

        if (this.searchQuery) {
          params.search = this.searchQuery;
        }

        const result = await instructorService.fetchAll(params);

        this.instructors = result.instructors;
        this.pagination = result.pagination;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch instructors";
        console.error("Error fetching instructors:", err);
        this.instructors = [...initialInstructors];
      } finally {
        this.isLoading = false;
      }
    },

    async fetchInstructorById(userId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const instructor = await instructorService.fetchById(userId);
        if (instructor) {
          this.currentInstructor = instructor;
          return instructor;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to fetch instructor";
        console.error("Error fetching instructor:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    setSearchQuery(query: string) {
      this.searchQuery = query;
      if (this.useServerPagination) {
        this.fetchInstructors(1, true);
      }
    },

    setPage(page: number) {
      this.pagination.page = page;
      this.fetchInstructors(page, false);
    },

    async updateInstructor(userId: string, data: UpdateInstructorData) {
      try {
        const updatedInstructor = await instructorService.update(userId, data);
        if (updatedInstructor) {
          const index = this.instructors.findIndex((i) => i.userId === userId);
          if (index !== -1) {
            this.instructors[index] = {
              ...this.instructors[index],
              ...updatedInstructor,
            };
          }
          if (this.currentInstructor?.userId === userId) {
            this.currentInstructor = {
              ...this.currentInstructor,
              ...updatedInstructor,
            };
          }
          return updatedInstructor;
        }
        return null;
      } catch {
        // Fallback to local update if API fails
        return this.updateInstructorLocal(userId, data);
      }
    },

    updateInstructorLocal(
      userId: string,
      data: UpdateInstructorData,
    ): Instructor | null {
      const index = this.instructors.findIndex((i) => i.userId === userId);
      if (index !== -1) {
        const instructor = this.instructors[index]!;
        this.instructors[index] = {
          ...instructor,
          ...data,
          name: data.firstName
            ? `${data.firstName} ${data.lastName ?? instructor.name.split(" ").slice(1).join(" ")}`
            : instructor.name,
          phone: data.phoneNumber ?? instructor.phone,
        } as Instructor;
        return this.instructors[index]!;
      }
      return null;
    },

    async uploadProfilePic(userId: string, file: File) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await instructorService.uploadProfilePic(userId, file);
        if (result) {
          const index = this.instructors.findIndex((i) => i.userId === userId);
          if (index !== -1) {
            this.instructors[index]!.image = result.url;
          }
          if (this.currentInstructor?.userId === userId) {
            this.currentInstructor!.image = result.url;
          }
          return result;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error
            ? err.message
            : "Failed to upload profile picture";
        console.error("Error uploading profile picture:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async uploadBase64Media(
      userId: string,
      base64Data: string,
      mimeType: string,
    ) {
      this.isLoading = true;
      this.error = null;

      try {
        const result = await instructorService.uploadBase64Media(
          userId,
          base64Data,
          mimeType,
        );
        if (result) {
          const index = this.instructors.findIndex((i) => i.userId === userId);
          if (index !== -1) {
            this.instructors[index]!.image = result.url;
          }
          if (this.currentInstructor?.userId === userId) {
            this.currentInstructor!.image = result.url;
          }
          return result;
        }
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to upload media";
        console.error("Error uploading base64 media:", err);
        return null;
      } finally {
        this.isLoading = false;
      }
    },

    async deleteProfilePic(userId: string) {
      this.isLoading = true;
      this.error = null;

      try {
        const success = await instructorService.deleteProfilePic(userId);
        if (success) {
          const index = this.instructors.findIndex((i) => i.userId === userId);
          if (index !== -1) {
            this.instructors[index]!.image = undefined;
          }
          if (this.currentInstructor?.userId === userId) {
            this.currentInstructor!.image = undefined;
          }
        }
        return success;
      } catch (err) {
        this.error =
          err instanceof Error
            ? err.message
            : "Failed to delete profile picture";
        console.error("Error deleting profile picture:", err);
        return false;
      } finally {
        this.isLoading = false;
      }
    },

    async getMediaMetadata(userId: string) {
      try {
        return await instructorService.getMediaMetadata(userId);
      } catch {
        return null;
      }
    },

    async deleteInstructor(userId: string) {
      try {
        const success = await instructorService.deleteProfilePic(userId);
        if (success) {
          this.instructors = this.instructors.filter(
            (i) => i.userId !== userId,
          );
        }
        return success;
      } catch {
        // Fallback to local delete if API fails
        this.instructors = this.instructors.filter((i) => i.userId !== userId);
        return true;
      }
    },

    async createInstructor(data: CreateInstructorRequest, imageFile?: File) {
      this.isLoading = true;
      this.error = null;

      console.log("[createInstructor] Starting...");
      console.log("[createInstructor] Step 1: Data received:", data);
      console.log("[createInstructor] Step 2: Image file received:", imageFile);

      try {
        console.log(
          "[createInstructor] Step 3: Calling instructorService.create...",
        );
        const instructor = await instructorService.create(data);
        console.log(
          "[createInstructor] Step 4: instructorService.create returned:",
          instructor,
        );

        if (instructor) {
          console.log(
            "[createInstructor] Step 5: Instructor created successfully, adding to list...",
          );
          this.instructors.unshift(instructor);
          console.log(
            "[createInstructor] Step 6: Instructor added to store. userId:",
            instructor.userId,
          );

          if (instructor.userId && imageFile) {
            console.log(
              "[createInstructor] Step 7: Uploading profile picture...",
            );
            const uploadResult = await instructorService.uploadProfilePic(
              instructor.userId,
              imageFile,
            );
            console.log(
              "[createInstructor] Step 8: Upload result:",
              uploadResult,
            );

            if (uploadResult) {
              const index = this.instructors.findIndex(
                (i) => i.userId === instructor.userId,
              );
              if (index !== -1) {
                this.instructors[index]!.image = uploadResult.url;
                console.log(
                  "[createInstructor] Step 9: Image URL updated in store",
                );
              }
            } else {
              console.log(
                "[createInstructor] Step 9: Upload failed, skipping image update",
              );
            }
          } else {
            console.log(
              "[createInstructor] Step 7: No image file or userId, skipping upload",
            );
          }

          console.log(
            "[createInstructor] Step 10: Returning instructor:",
            instructor,
          );
          return instructor;
        }

        console.log("[createInstructor] Step 4 ERROR: instructor is null");
        return null;
      } catch (err) {
        this.error =
          err instanceof Error ? err.message : "Failed to create instructor";
        console.error("[createInstructor] ERROR caught:", err);
        console.error("[createInstructor] Error message:", this.error);
        return null;
      } finally {
        this.isLoading = false;
        console.log("[createInstructor] Finished. isLoading set to false");
      }
    },

    filterInstructors(searchQuery: string) {
      return this.instructors.filter((instructor) => {
        const matchesSearch =
          instructor.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          instructor.email.toLowerCase().includes(searchQuery.toLowerCase());
        return matchesSearch;
      });
    },

    setCurrentInstructor(instructor: Instructor | null) {
      this.currentInstructor = instructor;
    },

    reset() {
      this.instructors = [...initialInstructors];
      this.currentInstructor = null;
      this.isLoading = false;
      this.error = null;
    },
  },
});
