# Forgot Password - Password Reset Guide

This guide explains how to use the forgot password feature to change a user's password. OTP verification is now optional - you can directly reset the password with just email and new password.

## Overview

The forgot password flow can be used in two ways:

**Simple Flow (Recommended):**
- Directly reset password with email and new password (no OTP required)

**Full Flow (Optional):**
1. **Request Password Reset** - User requests a password reset by providing their email (sends OTP via email)
2. **Verify OTP** (optional) - User verifies the OTP received via email
3. **Reset Password** - User sets a new password (no OTP required in this step)

## API Endpoints

All endpoints are public (no authentication required) and are under the `/password` route group:

- `POST /password/forgot` - Request password reset (sends OTP via email)
- `POST /password/verify-otp` - Verify OTP code (optional step)
- `POST /password/reset` - Reset password with new password

## Step-by-Step Flow

### Step 1: Request Password Reset

**Endpoint:** `POST /password/forgot`

**Description:** 
- Sends a 6-digit OTP code to the user's email
- Requires the user to be verified (email must be verified)
- Generates OTP and stores it in the database with an expiration time
- OTP expires after a certain period (configured in utils)

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:3000/password/forgot \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset OTP has been sent to your email"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Failed to process password reset request",
  "error": "user not found or not verified"
}
```

**Requirements:**
- Email must be registered and verified in the system
- User must exist and be verified

---

### Step 2: Verify OTP (Optional)

**Endpoint:** `POST /password/verify-otp`

**Description:**
- Validates the OTP code entered by the user
- Checks if OTP is valid, not expired, and not already used
- This step is optional - you can skip directly to Step 3 if preferred

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:3000/password/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP verified successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Invalid or expired OTP"
}
```

**Requirements:**
- Email must match the one used in Step 1
- OTP must be the 6-digit code received via email
- OTP must not be expired
- OTP must not have been used already

---

### Step 3: Reset Password

**Endpoint:** `POST /password/reset`

**Description:**
- Changes the user's password to the new password provided
- No OTP required - only email and new password needed
- Hashes the new password using bcrypt before storing

**Request Body:**
```json
{
  "email": "user@example.com",
  "new_password": "newpassword123"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:3000/password/reset \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "new_password": "newpassword123"
  }'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Password has been reset successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Failed to reset password",
  "error": "user not found"
}
```

**Requirements:**
- Email must be registered in the system
- New password must be between 6-20 characters (as per validation)

---

## Complete Flow Example

Here's a complete example of resetting a password. **Note:** Steps 1 and 2 are now optional since OTP verification is no longer required for password reset.

### Simple Flow (Recommended):
```bash
# Direct password reset (no OTP required)
curl -X POST http://localhost:3000/password/reset \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "new_password": "myNewPassword123"
  }'
```

### Full Flow (with OTP verification):
```bash
# Step 1: Request password reset (optional - sends OTP via email)
curl -X POST http://localhost:3000/password/forgot \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'

# Check email for OTP (e.g., "123456")

# Step 2 (Optional): Verify OTP
curl -X POST http://localhost:3000/password/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'

# Step 3: Reset password (no OTP needed)
curl -X POST http://localhost:3000/password/reset \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "new_password": "myNewPassword123"
  }'
```

## Implementation Details

### OTP Generation
- **Length:** 6 digits
- **Format:** Numeric only
- **Expiration:** Configured in `utils.GetExpiryTime()` (typically 10-15 minutes)
- **Storage:** Stored in `forgot_passwords` table with:
  - User ID
  - Email
  - OTP code
  - Expiration timestamp
  - Used flag (prevents reuse)

### Password Security
- Passwords are hashed using bcrypt with cost factor 10
- Original password is never stored
- OTP codes are single-use (marked as used after successful reset)

### Email Sending
- OTP is sent via email using the configured SMTP settings
- Email service is handled by `utils.SendPasswordResetEmail()`
- Requires proper SMTP configuration in environment variables

### Validation Rules
- **Email:** Must be a valid email format
- **New Password:** Required, must be 6-20 characters long
- **OTP:** Optional - only needed if using the OTP verification flow

## Error Scenarios

1. **User not found**
   - Error: "user not found"
   - Solution: Ensure user exists in the system

2. **Email format invalid**
   - Error: "Invalid request" with validation error
   - Solution: Provide a valid email address

3. **Password length invalid**
   - Error: "Invalid request" with validation error
   - Solution: Password must be 6-20 characters

## Security Considerations

1. **Password Hashing:** New passwords are hashed using bcrypt before storage
2. **No Authentication Required:** Reset endpoint is public - anyone with the email can reset the password
3. **OTP Verification:** OTP endpoints are still available for additional security verification if needed
4. **Email Verification:** The forgot password request still verifies the user exists

## Notes

- **OTP is no longer required** for password reset - you can directly reset with email and new password
- Step 1 (Request Password Reset) and Step 2 (Verify OTP) are now optional
- You can directly call `/password/reset` with just email and new password
- After successful password reset, the user can login with the new password
- The old password becomes invalid immediately after reset
- If you want additional security, you can still use the OTP verification endpoints before resetting

