# Payment & Authentication Flow Documentation

## User Journey Overview

### Flow 1: Guest Checkout (No Registration Required)
1. User views packages on `/packages`
2. User clicks "Book Now" or "Get Started"
3. Redirects to `/payment-method` → Select payment method & enter email
4. Redirects to `/payment` → Display payment instructions for selected method
5. Redirects to `/payment-status` → Show success/failure with CTA to dashboard

### Flow 2: Standard Registration + Plan Selection + Payment
1. User fills `/register` form (3-step flow)
2. On submission → Redirects to `/verify` for OTP verification
3. After OTP verification → Redirects to `/onboarding` (promo banner display)
4. From onboarding → Redirects to `/select-plan` (plan selection with free trial info)
5. Plan selection → Redirects to `/payment-method` (payment method selection)
6. After payment → Redirects to `/payment-status`
7. On success → Redirects to `/dashboard/free-trial` or main dashboard

---

## Payment Methods Supported

All methods managed through `/payment` page:
- **Virtual Account (VA)** - BCA, BNI, Mandiri, Permata, Bank Transfer
- **QRIS** - Quick Response Code Indonesian Standard
- **Bank Transfer** - Direct bank deposit
- **Credit/Debit Card** - Full Luhn validation on front-end

---

## Free Trial Feature

### Eligibility
- **Free registered users**: 1 free trial (one-time)
- **Paid registered users**: 1 free trial (one-time, after payment)

### Free Trial Details
- **Duration**: 10-15 minutes (configurable via env vars)
- **Location**: `/dashboard/free-trial`
- **Countdown Timer**: Real-time countdown to session end
- **Features**:
  - Instructor assignment (mock)
  - Vehicle assignment (mock)
  - Session start/end buttons
  - Session history tracking

### Accessing Free Trial
1. New registered user → After onboarding completes
2. After payment → Redirect to free-trial page directly
3. From dashboard → Button to access if eligible

---

## OTP Verification System

### Implementation
- **Mock OTP generation**: 6-digit code
- **Verification time**: 5 minutes
- **Resend cooldown**: 30 seconds
- **Max attempts**: 5 attempts

### OTP Format
- Code: 6 digits (000000-999999)
- Sent to email in registration form
- Mock code for testing: Any 6-digit number

---

## Payment Status Handling

### Success Flow
- Payment verified → Redirect to `/payment-status?status=success`
- Auto-redirect to dashboard after 3 seconds
- User can manually proceed to dashboard

### Failure Flow
- Payment failed → Redirect to `/payment-status?status=failed`
- User can retry payment or return to select-plan

---

## Promo System

### Dynamic Promo Timing
- **Start/End times**: Unix timestamps (env vars)
- **Discount percentage**: Configurable (env vars)
- **Display location**: `/select-plan` page
- **Visibility**: Only shows when current time is within promo window

### Promo Banner
- Shows countdown to promo end
- Highlights discount percentage
- Works for both free and paid plans
- Auto-hides when promo expires

---

## Environment Variables Required

Copy from `.env.example`:
```
NUXT_PUBLIC_MIDTRANS_CLIENT_KEY=your_key
MIDTRANS_SERVER_KEY=your_key
NUXT_PUBLIC_PROMO_START_TIME=unix_timestamp
NUXT_PUBLIC_PROMO_END_TIME=unix_timestamp
NUXT_PUBLIC_PROMO_DISCOUNT_PERCENTAGE=20
NUXT_PUBLIC_FREE_TRIAL_DURATION_MINUTES=15
```

---

## Midtrans Integration Points

### Client-Side Integration
- Snap.js script loaded in `/payment` page
- Midtrans client key from env vars
- Payment method selection linked to Snap token

### Server-Side Implementation Needed
- `/api/midtrans/create-transaction` - Create Midtrans transaction
- `/api/midtrans/verify` - Verify payment status
- `/api/midtrans/callback` - Handle Midtrans webhooks

---

## Testing Payment Methods

### VA (Virtual Account)
- Customer can transfer to generated account number
- Payment auto-confirms within 1 hour

### QRIS
- Scan QR with supporting bank app
- Payment instantly confirms

### Credit Card (Mock)
- Card number: 4111 1111 1111 1111
- Expiry: Any future date
- CVV: Any 3 digits

---

## Onboarding Modal (Dashboard)

### Trigger
- Appears on first dashboard login
- Modal blocks dashboard access until completed
- Modal cannot be dismissed without completing

### Content
1. Password creation field (min 8 chars)
2. Full name field (must match KTP for verification)
3. Submit button

### On Complete
- User data saved
- Modal closes
- Dashboard fully accessible
- Free trial access enabled (if eligible)

---

## Validation Rules

### Credit Card
- **Luhn algorithm**: Implemented on front-end
- **Format**: 16 digits (spaces optional)
- **Expiry**: MM/YY format, future date
- **CVV**: 3-4 digits

### Email
- **Format**: Valid email pattern
- **Usage**: Used for OTP delivery

### Password (Onboarding)
- **Min length**: 8 characters
- **Requirements**: Letters + numbers (enforced via Zod schema)

---

## Route Protection

### Protected Routes
- `/dashboard/*` - Requires authentication
- `/select-plan` - Requires registration
- `/dashboard/free-trial` - Requires completed free trial eligibility check

### Public Routes
- `/login`, `/register`, `/packages` - No auth required
- `/payment-method`, `/payment` - Can be public or protected based on business logic
- `/verify`, `/onboarding` - Part of registration flow

---

## Notes for Backend Integration

When integrating backend services:
1. Replace mock data with real API calls
2. Store free trial usage count in database
3. Implement proper OTP delivery (email service)
4. Connect Midtrans API for real transactions
5. Add Row Level Security (RLS) policies if using Supabase
6. Implement proper session management
