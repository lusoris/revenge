# Questions - MFA Service Testing

**Date**: 2026-02-04
**Status**: RESOLVED

---

## Question 1: Redundant Nonce Field in DB

**Answer**: Remove it (clean approach)

**Action**: DB migration required to remove `nonce` column from `user_totp_secrets`

---

## Question 2: TOTP Re-Enrollment Flow

**Answer**: Auto-overwrite existing secret

**Action**: Implement upsert logic in GenerateSecret

---

## Question 3: Error Handling Strategy

**Answer**: Use `errors.Is(err, pgx.ErrNoRows)` for consistency

**Action**: Fix HasTOTP to properly distinguish no-rows from real errors
