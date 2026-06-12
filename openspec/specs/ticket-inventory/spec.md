# Ticket Inventory

## Purpose

Manage ticket categories per event, track available stock, and ensure atomic stock operations to prevent overselling.

## Requirements

### Requirement: Create Ticket Category
The system SHALL allow admins to create ticket categories (e.g., VIP, Regular, Festival) under an event with name, price, and total stock.

#### Scenario: Admin creates ticket category
- **WHEN** admin submits valid category data (name, price > 0, stock > 0)
- **THEN** system creates category and initializes available stock equal to total stock

#### Scenario: Non-admin attempts to create category
- **WHEN** regular user attempts to create ticket category
- **THEN** system returns 403 Forbidden

### Requirement: List Ticket Categories per Event
The system SHALL expose ticket categories under a published event publicly.

#### Scenario: Public list of categories
- **WHEN** client requests categories for a published event
- **THEN** system returns categories with name, price, and remaining available stock

#### Scenario: Hide sold-out categories
- **WHEN** a category has 0 available stock
- **THEN** system returns category with `is_sold_out: true` flag

### Requirement: Update Ticket Category
The system SHALL allow admins to update ticket category details, but stock changes have constraints.

#### Scenario: Admin updates price and name
- **WHEN** admin updates name and/or price
- **THEN** system updates the category and reflects changes immediately

#### Scenario: Admin attempts to reduce stock below sold count
- **WHEN** admin attempts to set total stock lower than already-sold count
- **THEN** system returns 400 with validation error

### Requirement: Atomic Stock Decrement
The system SHALL atomically decrement ticket stock on reservation to prevent overselling under concurrent requests.

#### Scenario: Concurrent reservations
- **WHEN** multiple users attempt to reserve the last ticket simultaneously
- **THEN** system uses row-level lock (`SELECT FOR UPDATE`) and processes requests serially; only one reservation succeeds, the rest receive "sold out"

#### Scenario: Stock goes to zero
- **WHEN** last ticket is reserved
- **THEN** system marks category as `is_sold_out: true` and prevents further reservations

### Requirement: Stock Restoration on Expiry/Cancellation
The system SHALL restore ticket stock when a booking expires or is cancelled.

#### Scenario: Booking expires before payment
- **WHEN** a pending booking exceeds its TTL (10 minutes)
- **THEN** background job releases the held tickets back to available stock

#### Scenario: User cancels pending booking
- **WHEN** user explicitly cancels a pending payment booking
- **THEN** system releases the held tickets back to available stock immediately
