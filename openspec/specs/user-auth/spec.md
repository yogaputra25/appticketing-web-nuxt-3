# User Authentication

## Purpose

Manage user registration, login, profile management, and role-based access control using JWT authentication.

## Requirements

### Requirement: User Registration
The system SHALL allow new users to register with email, password, full name, and phone number.

#### Scenario: Successful registration
- **WHEN** user submits valid registration data (unique email, password ≥ 8 chars)
- **THEN** system creates user account with hashed password (bcrypt) and returns success

#### Scenario: Duplicate email registration
- **WHEN** user submits registration with an email that already exists
- **THEN** system returns 409 Conflict error with message "Email already registered"

#### Scenario: Invalid registration data
- **WHEN** user submits invalid data (missing fields, weak password, invalid email format)
- **THEN** system returns 400 Bad Request with field-level validation errors

### Requirement: User Login
The system SHALL authenticate users via email and password and return a JWT token.

#### Scenario: Successful login
- **WHEN** user submits valid credentials
- **THEN** system returns JWT access token (expires 24h) and user profile data

#### Scenario: Invalid credentials
- **WHEN** user submits wrong email or password
- **THEN** system returns 401 Unauthorized with generic message "Invalid credentials"

### Requirement: JWT-based Authentication
The system SHALL protect authenticated endpoints using JWT bearer tokens.

#### Scenario: Valid token access
- **WHEN** client sends request with valid JWT in Authorization header
- **THEN** system processes the request and attaches user context

#### Scenario: Missing or invalid token
- **WHEN** client sends request without token or with expired/invalid token
- **THEN** system returns 401 Unauthorized

### Requirement: Role-based Access Control
The system SHALL distinguish between `user` and `admin` roles for authorization.

#### Scenario: Admin endpoint access by admin
- **WHEN** admin user accesses admin-only endpoint
- **THEN** system processes the request successfully

#### Scenario: Non-admin access to admin endpoint
- **WHEN** regular user accesses admin-only endpoint
- **THEN** system returns 403 Forbidden

### Requirement: User Profile Management
The system SHALL allow authenticated users to view and update their profile.

#### Scenario: Get current user profile
- **WHEN** authenticated user requests `/api/auth/me`
- **THEN** system returns the user profile (id, email, name, phone, role)

#### Scenario: Update profile
- **WHEN** authenticated user submits profile updates (name, phone)
- **THEN** system updates the user record and returns updated profile
