# Integration Checklist & API Templates

## ✅ What's Ready to Use

- [x] All 7 payment & authentication pages created
- [x] OTP verification flow
- [x] Free trial session management
- [x] Dynamic promo pricing system
- [x] Credit card Luhn validation
- [x] Onboarding modal component
- [x] Form validation with Zod
- [x] Responsive UI with Nuxt UI
- [x] Complete documentation

---

## ⏳ What Needs Backend Implementation

### 1. Midtrans Integration

Create file: `/app/server/api/midtrans/create-transaction.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { plan, paymentMethod, email } = await readBody(event)
  
  // 1. Validate input
  // 2. Get pricing from database (apply promo if active)
  // 3. Call Midtrans API to create transaction
  // 4. Store order in database
  // 5. Return snap token to client
  
  return {
    snapToken: 'snap_token_here',
    orderId: 'order_123',
    amount: 500000
  }
})
```

Create file: `/app/server/api/midtrans/verify.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { orderId } = await readBody(event)
  
  // 1. Call Midtrans API to verify payment status
  // 2. Update transaction status in database
  // 3. Create user session/subscription if successful
  // 4. Grant free trial if applicable
  
  return {
    status: 'success' | 'failed',
    message: 'Payment verified'
  }
})
```

Create file: `/app/server/api/midtrans/callback.ts`
```typescript
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  
  // 1. Verify webhook signature with Midtrans
  // 2. Update transaction status based on notification
  // 3. Grant access/subscription if payment confirmed
  // 4. Log webhook for audit trail
  
  return { status: 'ok' }
})
```

### 2. OTP Delivery Service

Create file: `/app/server/api/auth/send-otp.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { email } = await readBody(event)
  
  // 1. Generate 6-digit OTP code
  // 2. Store OTP in database with expiry (5 minutes)
  // 3. Send OTP via email service (SendGrid, Mailgun, etc.)
  // 4. Return success confirmation
  
  return { success: true, message: 'OTP sent to email' }
})
```

Create file: `/app/server/api/auth/verify-otp.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { email, code } = await readBody(event)
  
  // 1. Retrieve stored OTP for email
  // 2. Compare code with stored OTP
  // 3. Check if OTP expired
  // 4. Increment attempt counter
  // 5. If valid: Mark email as verified, create user session
  // 6. Return success/error
  
  return {
    success: true,
    userId: 'user_123',
    sessionToken: 'token_here'
  }
})
```

### 3. Free Trial Management

Create file: `/app/server/api/free-trial/claim.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { userId } = await readUserSession(event)
  
  // 1. Check if user eligible for free trial
  // 2. Check if already used free trial
  // 3. Create free trial session record
  // 4. Set expiry time (15 minutes from now)
  // 5. Return session details
  
  return {
    sessionId: 'session_123',
    expiresAt: new Date(),
    instructor: { name: 'Pak Ahmad', phone: '08123456789' },
    vehicle: { plate: 'B 1234 XYZ', model: 'Honda Civic' }
  }
})
```

Create file: `/app/server/api/free-trial/complete.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { sessionId } = await readBody(event)
  
  // 1. Retrieve session details
  // 2. Mark session as completed
  // 3. Store completion timestamp
  // 4. Issue certificate (if applicable)
  // 5. Update user free trial status (marked as used)
  
  return {
    success: true,
    certificateUrl: 'https://...'
  }
})
```

### 4. User Profile Setup

Create file: `/app/server/api/user/setup-profile.ts`
```typescript
export default defineEventHandler(async (event) => {
  const { userId } = await readUserSession(event)
  const { password, fullName } = await readBody(event)
  
  // 1. Hash password with bcrypt
  // 2. Validate full name
  // 3. Update user profile in database
  // 4. Mark onboarding as complete
  // 5. Verify email match with government ID (if required)
  
  return { success: true }
})
```

---

## Database Schema (SQL)

### Users Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR UNIQUE NOT NULL,
  password_hash VARCHAR,
  full_name VARCHAR,
  phone VARCHAR,
  created_at TIMESTAMP DEFAULT NOW(),
  onboarding_completed BOOLEAN DEFAULT FALSE,
  email_verified BOOLEAN DEFAULT FALSE
);
```

### Free Trials Table
```sql
CREATE TABLE free_trials (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  status VARCHAR DEFAULT 'available', -- available, used, expired
  used_at TIMESTAMP,
  duration_minutes INTEGER DEFAULT 15,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id) -- One free trial per user
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  amount DECIMAL(10, 2),
  status VARCHAR DEFAULT 'pending', -- pending, success, failed
  payment_method VARCHAR,
  midtrans_order_id VARCHAR UNIQUE,
  plan VARCHAR,
  created_at TIMESTAMP DEFAULT NOW(),
  completed_at TIMESTAMP
);
```

### Sessions Table
```sql
CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  is_free_trial BOOLEAN DEFAULT FALSE,
  instructor_name VARCHAR,
  vehicle_info JSONB,
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  status VARCHAR DEFAULT 'scheduled', -- scheduled, in_progress, completed
  created_at TIMESTAMP DEFAULT NOW()
);
```

### OTP Table (Temporary)
```sql
CREATE TABLE otp_codes (
  id UUID PRIMARY KEY,
  email VARCHAR NOT NULL,
  code VARCHAR(6) NOT NULL,
  attempts INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '5 minutes',
  verified_at TIMESTAMP
);
```

---

## Environment Variables to Add

```bash
# Email Service
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your_sendgrid_api_key
SMTP_FROM=noreply@drivingmaster.com

# Database
DATABASE_URL=postgresql://...

# Midtrans
NUXT_PUBLIC_MIDTRANS_CLIENT_KEY=your_client_key
MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_SANDBOX_MODE=true

# Security
JWT_SECRET=your_jwt_secret_key
SESSION_SECRET=your_session_secret_key

# Promo
NUXT_PUBLIC_PROMO_START_TIME=1704067200
NUXT_PUBLIC_PROMO_END_TIME=1706745600
NUXT_PUBLIC_PROMO_DISCOUNT_PERCENTAGE=20

# Free Trial
NUXT_PUBLIC_FREE_TRIAL_DURATION_MINUTES=15
```

---

## Testing API Endpoints

### Test OTP Flow
```bash
# Send OTP
curl -X POST http://localhost:3000/api/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# Verify OTP (any 6-digit code works in dev)
curl -X POST http://localhost:3000/api/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","code":"123456"}'
```

### Test Payment Flow
```bash
# Create transaction
curl -X POST http://localhost:3000/api/midtrans/create-transaction \
  -H "Content-Type: application/json" \
  -d '{
    "plan":"standard",
    "paymentMethod":"va",
    "email":"user@example.com"
  }'

# Verify payment
curl -X POST http://localhost:3000/api/midtrans/verify \
  -H "Content-Type: application/json" \
  -d '{"orderId":"order_123"}'
```

### Test Free Trial
```bash
# Claim free trial
curl -X POST http://localhost:3000/api/free-trial/claim \
  -H "Authorization: Bearer token_here"

# Complete free trial
curl -X POST http://localhost:3000/api/free-trial/complete \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"session_123"}'
```

---

## Frontend Hooks & Middleware

### Create Authentication Middleware
File: `/app/middleware/auth.ts`
```typescript
export default defineRouteMiddleware((to, from) => {
  const { data: session } = await useAsyncData('session', () => $fetch('/api/auth/session'))
  
  if (!session.value && to.path.startsWith('/dashboard')) {
    return navigateTo('/login')
  }
})
```

### Create Composable for Free Trial
File: `/app/composables/useFreeTrial.ts`
```typescript
export const useFreeTrial = () => {
  const isTrialActive = ref(false)
  const remainingTime = ref(0)
  
  const claimFreeTrial = async () => {
    const { data } = await $fetch('/api/free-trial/claim')
    isTrialActive.value = true
    return data
  }
  
  const completeFreeTrial = async (sessionId) => {
    await $fetch('/api/free-trial/complete', {
      method: 'POST',
      body: { sessionId }
    })
  }
  
  return { isTrialActive, remainingTime, claimFreeTrial, completeFreeTrial }
}
```

---

## Production Deployment Checklist

- [ ] Get Midtrans API keys (sandbox first, then production)
- [ ] Setup email service (SendGrid, Mailgun, AWS SES)
- [ ] Configure database (Supabase, Neon, AWS RDS)
- [ ] Implement all backend API routes
- [ ] Setup environment variables in Vercel
- [ ] Test payment flow in sandbox mode
- [ ] Setup error logging (Sentry, LogRocket)
- [ ] Configure rate limiting
- [ ] Add CSRF protection
- [ ] Setup monitoring & alerts
- [ ] Test end-to-end user flows
- [ ] Deploy to Vercel

---

## Support & Documentation

- **Payment Docs**: See `PAYMENT_FLOW.md`
- **Flow Diagrams**: See `FLOW_DIAGRAMS.md`
- **Quick Start**: See `QUICKSTART.md`
- **All Features**: See `FEATURES_IMPLEMENTED.md`

Ready to build! 🚀
