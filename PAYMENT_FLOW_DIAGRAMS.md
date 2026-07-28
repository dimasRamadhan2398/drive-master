# User Flow Diagrams

## Complete User Journey

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         DRIVING SCHOOL APP FLOWS                            │
└─────────────────────────────────────────────────────────────────────────────┘

╔════════════════════════════════════════════════════════════════════════════╗
║ FLOW 1: GUEST CHECKOUT (No Registration)                                  ║
╚════════════════════════════════════════════════════════════════════════════╝

  ┌──────────┐
  │ PACKAGES │
  └────┬─────┘
       │ "Book Now"
       ▼
  ┌──────────────────┐
  │ PAYMENT METHOD   │ ◄─── Select: VA / QRIS / Bank Transfer / CC
  │ Enter email      │      Enter phone number
  └────┬─────────────┘      Validate credit card (Luhn)
       │
       ▼
  ┌──────────────┐
  │ PAYMENT PAGE │ ◄─────── Show payment instructions
  │ Instructions │         Based on selected method
  └────┬─────────┘
       │
       ▼
  ┌───────────────────┐
  │ PAYMENT STATUS    │ ◄─── Success/Failure page
  └───────────────────┘


╔════════════════════════════════════════════════════════════════════════════╗
║ FLOW 2: REGISTERED USER WITH FREE TRIAL                                   ║
╚════════════════════════════════════════════════════════════════════════════╝

  ┌──────────┐
  │ REGISTER │ ◄───── 3-step form
  └────┬─────┘       Step 1: Personal info
       │             Step 2: Account details
       ▼             Step 3: Confirmation
  ┌────────────────┐
  │ VERIFY (OTP)   │ ◄───── 6-digit code
  │ 6 digit input  │        Resend after 30s
  └────┬───────────┘
       │
       ▼
  ┌─────────────┐
  │ ONBOARDING  │ ◄───── Promo banner (if active)
  │ Promo info  │        Countdown timer
  └────┬────────┘
       │
       ▼
  ┌──────────────┐
  │ SELECT PLAN  │ ◄───── Show plans with:
  │ Plan choice  │        - Free trial duration
  └────┬─────────┘        - Discounted pricing
       │                   - Plan comparison
       ▼
  ┌──────────────────┐
  │ PAYMENT METHOD   │ ◄─── Same as Flow 1
  └────┬─────────────┘
       │
       ▼
  ┌──────────────┐
  │ PAYMENT PAGE │ ◄─────── Midtrans integration
  └────┬─────────┘
       │
       ▼
  ┌───────────────────┐
  │ PAYMENT STATUS    │
  └────┬──────────────┘
       │
       ▼ On Success
  ┌─────────────┐
  │ DASHBOARD   │ ◄───── Onboarding Modal appears
  │ (First time)│       (Must complete: password + full name)
  └────┬────────┘
       │ Modal complete
       ▼
  ┌──────────────────┐
  │ DASHBOARD        │
  │ Main page        │
  └────┬─────────────┘
       │ Click "Start Free Trial"
       ▼
  ┌───────────────────────────┐
  │ FREE TRIAL SESSION        │ ◄─── Features:
  │ - 15 min countdown timer  │      - Instructor assigned
  │ - Session details         │      - Vehicle assigned
  │ - Start/End buttons       │      - Start/End controls
  │ - Session history         │      - Completion tracking
  └───────────────────────────┘


╔════════════════════════════════════════════════════════════════════════════╗
║ PAYMENT METHODS FLOW                                                       ║
╚════════════════════════════════════════════════════════════════════════════╝

  PAYMENT METHOD SELECTED
         │
    ┌────┼────┬──────────┬─────────────┐
    ▼    ▼    ▼          ▼             ▼
   VA  QRIS  Bank       Credit Card    
   │    │    Transfer       │
   │    │    │              │
   │    │    ▼              ▼
   │    │  [Show bank     [Show card
   │    │   details]       form]
   │    │                  │
   │    ▼                  ▼
   │  [Show QR]        [Luhn validate]
   │    │               │
   ▼    ▼               ▼
   └────┴───────────────┘
        │
        ▼
   [Midtrans Snap]
        │
        ▼
   PAYMENT STATUS PAGE


╔════════════════════════════════════════════════════════════════════════════╗
║ FREE TRIAL ELIGIBILITY                                                     ║
╚════════════════════════════════════════════════════════════════════════════╝

  FREE REGISTERED USER
  │
  └─ Eligible for 1x free trial
     Duration: 10-15 minutes
     Access via: /dashboard/free-trial
     Usage: Tracked in database

  PAID REGISTERED USER (After Payment)
  │
  └─ Eligible for 1x free trial
     Duration: 10-15 minutes
     Access via: /dashboard/free-trial
     Usage: Tracked in database


╔════════════════════════════════════════════════════════════════════════════╗
║ PROMO SYSTEM FLOW                                                          ║
╚════════════════════════════════════════════════════════════════════════════╝

  Current Time < PROMO_START_TIME
  │
  └─ ❌ NO PROMO BADGE

  PROMO_START_TIME ≤ Current Time ≤ PROMO_END_TIME
  │
  └─ ✅ SHOW PROMO BANNER
     │
     ├─ Discount percentage displayed
     ├─ Countdown to expiry
     ├─ Price calculation (original - discount%)
     └─ Apply to all plans

  Current Time > PROMO_END_TIME
  │
  └─ ❌ HIDE PROMO BANNER


╔════════════════════════════════════════════════════════════════════════════╗
║ ONBOARDING MODAL (Dashboard First Load)                                   ║
╚════════════════════════════════════════════════════════════════════════════╝

  Dashboard Loads
  │
  └─ Check if user completed onboarding
     │
     ├─ If NO: Show modal (cannot dismiss)
     │  │
     │  ├─ Input 1: Password (min 8 chars, letters + numbers)
     │  ├─ Input 2: Full Name (must match KTP)
     │  └─ Submit button
     │     │
     │     ▼
     │  Save to database
     │     │
     │     ▼
     │  Modal closes
     │     │
     │     ▼
     │  Dashboard unlocks
     │
     └─ If YES: Show dashboard immediately


╔════════════════════════════════════════════════════════════════════════════╗
║ OTP VERIFICATION FLOW                                                      ║
╚════════════════════════════════════════════════════════════════════════════╝

  Register form submitted
  │
  ▼
  Redirect to /verify?email=user@example.com
  │
  ▼
  Display 6-digit input
  │
  ├─ Resend button (disabled first 30 seconds)
  ├─ Max 5 attempts
  └─ Timeout: 5 minutes
  │
  ▼
  User enters code
  │
  ▼ Code validation
  ├─ If valid → Redirect to /onboarding
  └─ If invalid → Show error + decrement attempts

```

---

## State Management Flow

```
User Action → Form Input → Validation (Zod) → Submit → API Call → Redirect

Example: Payment Method Selection
─────────────────────────────────
VA selected
   ↓
Form valid (email format checked)
   ↓
Submit clicked
   ↓
Store in route query: ?method=va&email=user@email.com
   ↓
Redirect to /payment
   ↓
Payment page reads query params
   ↓
Display VA instructions
```

---

## Database Relationships (To Implement)

```
┌──────────────┐
│ users        │
├──────────────┤
│ id (PK)      │
│ email        │
│ password     │
│ full_name    │
│ created_at   │
└──────┬───────┘
       │ 1:N
       │
       ▼
┌──────────────────────┐
│ user_free_trials     │
├──────────────────────┤
│ id (PK)              │
│ user_id (FK)         │
│ duration_minutes     │
│ status (used/unused) │
│ created_at           │
│ used_at              │
└──────┬───────────────┘
       │
       └─ 1:1 relationship with sessions

┌──────────────────────────────┐
│ transactions                 │
├──────────────────────────────┤
│ id (PK)                      │
│ user_id (FK)                 │
│ amount                       │
│ status (pending/success)     │
│ payment_method               │
│ midtrans_order_id            │
│ created_at                   │
│ completed_at                 │
└──────────────────────────────┘

┌──────────────────────────────┐
│ sessions                     │
├──────────────────────────────┤
│ id (PK)                      │
│ user_id (FK)                 │
│ instructor_name              │
│ vehicle_info                 │
│ start_time                   │
│ end_time                     │
│ status (scheduled/completed) │
│ is_free_trial                │
└──────────────────────────────┘
```

---

## Component Hierarchy

```
app.vue
├── pages/
│   ├── packages.vue ─────────┐
│   ├── login.vue            │
│   ├── register.vue ────────┬┘
│   ├── verify.vue           │
│   ├── onboarding.vue       │
│   ├── select-plan.vue ─────┤
│   ├── payment-method.vue ──┼─ Payment Flow
│   ├── payment.vue ─────────┤
│   ├── payment-status.vue ──┤
│   └── dashboard/
│       ├── index.vue ───────┼─ Protected Routes
│       │   └── OnboardingModal (NEW)
│       └── free-trial.vue ──┤
│
└── components/
    └── OnboardingModal.vue
```
