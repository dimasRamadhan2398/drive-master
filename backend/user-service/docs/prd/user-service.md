# Product Requirement Document (PRD)
## Driving School Website — Backend System

**Version:** 1.0  
**Date:** April 17, 2026  
**Status:** Draft  

---

## 1. Overview

### 1.1 Purpose
This document defines the backend requirements for a driving school website platform. The backend will serve as the core engine managing course scheduling, instructor operations, and payment processing — enabling a seamless experience for students, instructors, and administrators.

### 1.2 Goals
- Provide a robust, scalable API to support the driving school's core operations.
- Ensure data consistency and integrity across courses, instructors, and financial transactions.
- Enable real-time availability management for scheduling.
- Support secure payment processing and financial reporting.

### 1.3 Stakeholders
| Role | Responsibility |
|---|---|
| Admin | Manage courses, instructors, schedules, and view financial reports |
| Instructor | View and manage assigned schedules |
| Student | Browse, book courses, and make payments |
| Developer | Build and maintain the backend system |

---

## 2. Tech Stack (Recommended)

| Layer | Technology |
|---|---|
| Language | Node.js (TypeScript) or Python (FastAPI) |
| Database | PostgreSQL (relational data) + Redis (caching & sessions) |
| Auth | JWT + OAuth 2.0 |
| Payment Gateway | Midtrans / Xendit (Indonesia) |
| API Style | RESTful API with OpenAPI/Swagger documentation |
| Hosting | AWS / GCP / Railway |
| Storage | AWS S3 (for documents, certificates) |

---

## 3. Modules

---

### Module A: User Profile Data

#### 3.1 Overview
Manages all course types offered by the driving school, session availability, and bookings made by students.

#### 3.2 Data Entities

**Course**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `name` | String | Course name (e.g., "Beginner Car License") |
| `description` | Text | Course details |
| `category` | Enum | `car`, `motorcycle`, `truck` |
| `duration_hours` | Integer | Total hours of training |
| `max_students` | Integer | Maximum students per session |
| `price` | Decimal | Course price |
| `is_active` | Boolean | Whether the course is available |
| `certificate_issuer` | String | The issuer of the certificate |
| `created_at` | Timestamp | |
| `updated_at` | Timestamp | |

**Schedule (Session)**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `course_id` | UUID | FK → Course |
| `instructor_id` | UUID | FK → Instructor |
| `start_datetime` | Timestamp | Session start time |
| `end_datetime` | Timestamp | Session end time |
| `location` | String | Practice location / track |
| `available_slots` | Integer | Remaining bookable slots |
| `status` | Enum | `open`, `full`, `cancelled`, `completed` |

**Booking**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `student_id` | UUID | FK → User |
| `schedule_id` | UUID | FK → Schedule |
| `status` | Enum | `pending`, `confirmed`, `cancelled`, `completed` |
| `booked_at` | Timestamp | |
| `notes` | Text | Optional student notes |

#### 3.3 API Endpoints

| Method | Endpoint | Description | Access |
|---|---|---|---|
| GET | `/api/courses` | List all active courses | Public |
| GET | `/api/courses/:id` | Get course detail | Public |
| POST | `/api/courses` | Create a new course | Admin |
| PUT | `/api/courses/:id` | Update course | Admin |
| DELETE | `/api/courses/:id` | Deactivate course | Admin |
| GET | `/api/schedules` | List schedules (with filters: date, course, instructor) | Public |
| GET | `/api/schedules/:id` | Get schedule detail | Public |
| POST | `/api/schedules` | Create a new schedule | Admin |
| PUT | `/api/schedules/:id` | Update schedule | Admin |
| DELETE | `/api/schedules/:id` | Cancel schedule | Admin |
| POST | `/api/bookings` | Book a session | Student |
| GET | `/api/bookings/me` | Get current student's bookings | Student |
| PUT | `/api/bookings/:id/cancel` | Cancel a booking | Student / Admin |
| GET | `/api/bookings` | List all bookings | Admin |

#### 3.4 Business Rules
- A schedule cannot be booked if `available_slots = 0`.
- Cancellation by a student is only allowed up to **24 hours** before the session start.
- When a booking is cancelled, `available_slots` must be incremented back.
- Admin can cancel any schedule; all affected bookings must be updated to `cancelled` and students notified.
- A student cannot book the same schedule twice.

---

### Module B: Instructor Management

#### 3.5 Overview
Manages instructor profiles, availability, and their assignment to schedules.

#### 3.6 Data Entities

**Instructor**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `user_id` | UUID | FK → User (for login) |
| `full_name` | String | Full name |
| `phone` | String | Contact number |
| `license_number` | String | Teaching license ID |
| `license_expiry` | Date | License expiry date |
| `specializations` | Array | `["car", "motorcycle"]` |
| `bio` | Text | Short instructor bio |
| `photo_url` | String | Profile photo |
| `is_active` | Boolean | Employment status |
| `created_at` | Timestamp | |

**InstructorAvailability**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `instructor_id` | UUID | FK → Instructor |
| `day_of_week` | Enum | `monday` … `sunday` |
| `start_time` | Time | Availability start |
| `end_time` | Time | Availability end |

#### 3.7 API Endpoints

| Method | Endpoint | Description | Access |
|---|---|---|---|
| GET | `/api/instructors` | List all active instructors | Public |
| GET | `/api/instructors/:id` | Get instructor profile | Public |
| POST | `/api/instructors` | Create instructor profile | Admin |
| PUT | `/api/instructors/:id` | Update instructor profile | Admin |
| DELETE | `/api/instructors/:id` | Deactivate instructor | Admin |
| GET | `/api/instructors/:id/schedules` | Get schedules assigned to instructor | Admin / Instructor |
| POST | `/api/instructors/:id/availability` | Set availability | Admin / Instructor |
| GET | `/api/instructors/:id/availability` | Get availability | Admin / Instructor |
| PUT | `/api/instructors/:id/availability/:avail_id` | Update availability slot | Admin / Instructor |

#### 3.8 Business Rules
- An instructor cannot be assigned to overlapping schedules.
- The system must validate that the instructor's `license_expiry` is not past before assigning to a future schedule.
- Deactivating an instructor must unassign them from all future `open` schedules. Admin must be alerted to reassign.
- Availability slots are used to suggest valid schedule times when creating sessions.

---

### Module C: Payment & Billing

#### 3.9 Overview
Handles payment for course bookings, transaction history, and basic financial reporting.

#### 3.10 Data Entities

**Invoice**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `booking_id` | UUID | FK → Booking |
| `student_id` | UUID | FK → User |
| `amount` | Decimal | Total amount charged |
| `discount` | Decimal | Discount applied |
| `tax` | Decimal | Tax amount |
| `total` | Decimal | Final payable amount |
| `status` | Enum | `unpaid`, `paid`, `refunded`, `failed` |
| `due_date` | Timestamp | Payment deadline |
| `issued_at` | Timestamp | |

**Payment**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `invoice_id` | UUID | FK → Invoice |
| `gateway` | String | `midtrans`, `xendit`, `manual` |
| `gateway_transaction_id` | String | ID from payment gateway |
| `method` | String | `bank_transfer`, `qris`, `credit_card`, `cash` |
| `amount_paid` | Decimal | Amount received |
| `paid_at` | Timestamp | Confirmed payment time |
| `status` | Enum | `pending`, `success`, `failed`, `refunded` |

#### 3.11 API Endpoints

| Method | Endpoint | Description | Access |
|---|---|---|---|
| POST | `/api/payments/initiate` | Create payment session for a booking | Student |
| GET | `/api/payments/:id` | Get payment status | Student / Admin |
| POST | `/api/payments/webhook` | Receive callback from payment gateway | Gateway (Public) |
| GET | `/api/invoices/me` | List current student's invoices | Student |
| GET | `/api/invoices/:id` | Get invoice detail | Student / Admin |
| GET | `/api/invoices` | List all invoices | Admin |
| POST | `/api/invoices/:id/refund` | Issue a refund | Admin |
| GET | `/api/reports/revenue` | Revenue report (filter by date range) | Admin |

#### 3.12 Business Rules
- A booking is only set to `confirmed` after payment reaches `success` status.
- Invoices are auto-generated when a booking is created, with a **24-hour payment due date**.
- If payment is not completed within the due date, the booking is auto-cancelled and the slot is released.
- Refunds are only eligible for bookings cancelled **before** the session starts.
- Webhook endpoint must validate gateway signature before processing.
- All monetary values must be stored in IDR (Indonesian Rupiah) as integers (in cents/sen).

### Module D: Services & Syllabus

#### 3.13 Overview
Write all the package options, services and syllabus on every courses

#### 3.14 Data Entities
---
**Syllabus**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `course_id` | UUID | FK → User (for login) |
| `title` | String | Full name |
| `phone` | String | Contact number |
| `license_number` | String | Teaching license ID |
| `license_expiry` | Date | License expiry date |
| `specializations` | Array | `["car", "motorcycle"]` |
| `bio` | Text | Short instructor bio |
| `photo_url` | String | Profile photo |
| `is_active` | Boolean | Employment status |
| `created_at` | Timestamp | |

**Syllabus Section**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `syllabus_topic_id` | UUID | FK → Existing Syllabus Topic (for login) |
| `title` | String | Syllabus Title |
| `description` | String | Syllabus Section Description / Opening Paragraph |

**Syllabus Subsection**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `syllabus_section_id` | UUID | FK → All points of syllabus specific section (for login) |
| `title` | String | Syllabus  |
| `description` | String | Syllabus Subsection  |

**Service Area**
| Field | Type | Description |
|---|---|---|
| `id` | UUID | Primary key |
| `area_name` | String | Area name |
| `address` | String | Address |
| `phone_number` | String | Phone number |
| `province` | String | Location province |
| `city` | Date | Location city |
| `district` | Array | Location district |
| `village` | Text | Short instructor bio |
| `photo_url` | String | Profile photo |
| `is_active` | Boolean | Employment status |
| `created_at` | Timestamp | |

## 4. Cross-Cutting Concerns

### 4.1 Authentication & Authorization

- JWT-based authentication for all protected routes.
- Role-based access control (RBAC) with roles: `admin`, `instructor`, `student`.
- Token expiry: Access token 15 minutes, Refresh token 7 days.

### 4.2 Notifications

| Trigger | Recipient | Channel |
|---|---|---|
| Booking confirmed | Student | Email / WhatsApp |
| Payment due reminder | Student | Email |
| Schedule cancelled | Student + Instructor | Email |
| Payment received | Student | Email |
| Upcoming session (24h before) | Student + Instructor | Email / WhatsApp |

### 4.3 Error Handling

- All API responses follow a standard format:
```json
{
  "success": true | false,
  "message": "...",
  "data": { ... } | null,
  "errors": [ ... ] | null
}
```
- HTTP status codes must be used correctly (200, 201, 400, 401, 403, 404, 422, 500).

### 4.4 Data Validation
- All inputs must be validated server-side.
- Phone numbers must follow Indonesian format (`+62xxxxxxxxxx`).
- Date/time must be stored in UTC; displayed in `Asia/Jakarta` (WIB).

### 4.5 Logging & Monitoring
- All API requests must be logged (method, endpoint, user, timestamp, response status).
- Payment transactions require audit logs that are immutable.
- System errors must be reported to a monitoring tool (e.g., Sentry).

---

## 5. Non-Functional Requirements

| Category | Requirement |
|---|---|
| Performance | API response time < 300ms for 95th percentile |
| Scalability | System should handle up to 1,000 concurrent users |
| Availability | 99.5% uptime SLA |
| Security | HTTPS enforced; PII data encrypted at rest |
| Compliance | Payment data must not be stored raw; comply with PCI-DSS via gateway tokenization |
| Backup | Daily automated database backups; 30-day retention |

---

## 6. Out of Scope (v1.0)

- Mobile application backend (separate project)
- Online theory exam/quiz system
- Certificate generation
- Multi-branch / franchise management
- Loyalty or referral program

---

## 7. Open Questions

| # | Question | Owner | Status |
|---|---|---|---|
| 1 | Which payment gateway will be used (Midtrans or Xendit)? | Product | Open |
| 2 | Is WhatsApp notification required at launch or post-MVP? | Product | Open |
| 3 | Should instructors have a separate login portal or use the same interface? | Product | Open |
| 4 | What is the refund policy timeline (hours/days before session)? | Business | Open |

---

## 8. Revision History

| Version | Date | Author | Notes |
|---|---|---|---|
| 1.0 | April 17, 2026 | — | Initial draft |