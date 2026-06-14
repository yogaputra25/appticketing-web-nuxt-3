package model

import (
	"time"

	"gorm.io/gorm"
)

// =====================================================
// User
// =====================================================
type User struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	FullName     string         `gorm:"size:255;not null" json:"full_name"`
	Phone        *string        `gorm:"size:32" json:"phone,omitempty"`
	Role         string         `gorm:"size:32;not null;default:user" json:"role"` // user | admin
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) IsAdmin() bool { return u.Role == "admin" }

func (User) TableName() string { return "users" }

// =====================================================
// Event
// =====================================================
const (
	EventStatusDraft     = "draft"
	EventStatusPublished = "published"
	EventStatusCancelled = "cancelled"
	EventStatusFinished  = "finished"
)

type Event struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Venue       string         `gorm:"size:255;not null" json:"venue"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	BannerURL   *string        `gorm:"type:text" json:"banner_url,omitempty"`
	Status      string         `gorm:"size:32;not null;default:draft" json:"status"`
	Categories  []TicketCategory `gorm:"foreignKey:EventID" json:"categories,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Event) TableName() string { return "events" }

// =====================================================
// Ticket Category
// =====================================================
type TicketCategory struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	EventID        uint64         `gorm:"not null;index" json:"event_id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Description    *string        `gorm:"type:text" json:"description,omitempty"`
	Price          float64        `gorm:"type:numeric(15,2);not null" json:"price"`
	TotalStock     int            `gorm:"not null" json:"total_stock"`
	AvailableStock int            `gorm:"not null" json:"available_stock"`
	MaxPerUser     int            `gorm:"not null;default:4" json:"max_per_user"`
	SaleStartAt    *time.Time     `json:"sale_start_at,omitempty"`
	SaleEndAt      *time.Time     `json:"sale_end_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TicketCategory) TableName() string { return "ticket_categories" }

func (tc *TicketCategory) IsSoldOut() bool { return tc.AvailableStock <= 0 }

// =====================================================
// Booking
// =====================================================
const (
	BookingStatusPending = "pending_payment"
	BookingStatusPaid    = "paid"
	BookingStatusCancelled = "cancelled"
	BookingStatusExpired = "expired"
	BookingStatusRefunded = "refunded"
)

type Booking struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	BookingCode     string         `gorm:"size:32;uniqueIndex;not null" json:"booking_code"`
	UserID          uint64         `gorm:"not null;index" json:"user_id"`
	EventID         uint64         `gorm:"not null;index" json:"event_id"`
	User            *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Event           *Event         `gorm:"foreignKey:EventID" json:"event,omitempty"`
	TotalAmount     float64        `gorm:"type:numeric(15,2);not null" json:"total_amount"`
	Status          string         `gorm:"size:32;not null;default:pending_payment" json:"status"`
	ExpiresAt       *time.Time     `json:"expires_at,omitempty"`
	ETicketCodes    JSONStringList `gorm:"type:jsonb;not null;default:'[]'" json:"e_ticket_codes"`
	CancelledReason *string        `gorm:"size:255" json:"cancelled_reason,omitempty"`
	Items           []BookingItem  `gorm:"foreignKey:BookingID" json:"items,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Booking) TableName() string { return "bookings" }

// =====================================================
// Booking Item
// =====================================================
type BookingItem struct {
	ID               uint64    `gorm:"primaryKey" json:"id"`
	BookingID        uint64    `gorm:"not null;index" json:"booking_id"`
	TicketCategoryID uint64    `gorm:"not null;index" json:"ticket_category_id"`
	Quantity         int       `gorm:"not null" json:"quantity"`
	UnitPrice        float64   `gorm:"type:numeric(15,2);not null" json:"unit_price"`
	Subtotal         float64   `gorm:"type:numeric(15,2);not null" json:"subtotal"`
	CreatedAt        time.Time `json:"created_at"`
}

func (BookingItem) TableName() string { return "booking_items" }

// =====================================================
// Payment
// =====================================================
const (
	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusExpired  = "expired"
	PaymentStatusRefunded = "refunded"
)

type Payment struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	PaymentCode     string         `gorm:"size:64;uniqueIndex;not null" json:"payment_code"`
	BookingID       uint64         `gorm:"not null;index" json:"booking_id"`
	UserID          uint64         `gorm:"not null;index" json:"user_id"`
	Amount          float64        `gorm:"type:numeric(15,2);not null" json:"amount"`
	Status          string         `gorm:"size:32;not null;default:pending" json:"status"`
	PaymentMethod   string         `gorm:"size:32;not null;default:simulation" json:"payment_method"`
	PaidAt          *time.Time     `json:"paid_at,omitempty"`
	ExpiredAt       *time.Time     `json:"expired_at,omitempty"`
	GatewayResponse JSONRawMessage `gorm:"type:jsonb" json:"gateway_response,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Payment) TableName() string { return "payments" }

// =====================================================
// Queue Token
// =====================================================
const (
	QueueStatusWaiting = "waiting"
	QueueStatusReady   = "ready"
	QueueStatusUsed    = "used"
	QueueStatusExpired = "expired"
)

type QueueToken struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	UserID        uint64     `gorm:"not null;index" json:"user_id"`
	EventID       uint64     `gorm:"not null" json:"event_id"`
	Token         string     `gorm:"size:128;uniqueIndex;not null" json:"token"`
	Status        string     `gorm:"size:32;not null;default:waiting" json:"status"`
	Position      *int       `json:"position,omitempty"`
	JoinedAt      time.Time  `gorm:"not null;default:NOW()" json:"joined_at"`
	ReadyAt       *time.Time `json:"ready_at,omitempty"`
	LastActiveAt  time.Time  `gorm:"not null;default:NOW()" json:"last_active_at"`
	ExpiresAt     time.Time  `gorm:"not null" json:"expires_at"`
}

func (QueueToken) TableName() string { return "queue_tokens" }

// =====================================================
// Ticket
// =====================================================
const (
	TicketStatusActive   = "active"
	TicketStatusUsed     = "used"
	TicketStatusRefunded = "refunded"
)

type Ticket struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	BookingID    uint64     `gorm:"not null;index" json:"booking_id"`
	TicketCode   string     `gorm:"size:64;uniqueIndex;not null" json:"ticket_code"`
	CategoryName string     `gorm:"size:100;not null" json:"category_name"`
	Status       string     `gorm:"size:32;not null;default:active" json:"status"`
	ScannedAt    *time.Time `json:"scanned_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (Ticket) TableName() string { return "tickets" }

// =====================================================
// JSON helper types
// =====================================================
// JSONStringList is a []string stored as JSONB in Postgres.
// GORM's default scanner doesn't handle []string<->jsonb, so we
// provide custom Scanner/Valuer.
type JSONStringList []string

func (j JSONStringList) Value() (interface{}, error) {
	if j == nil {
		return "[]", nil
	}
	return jsonMarshal(j), nil
}

func (j *JSONStringList) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}
	if err := jsonUnmarshal(bytes, j); err != nil {
		var s string
		if jsonUnmarshal(bytes, &s) == nil {
			*j = JSONStringList{s}
			return nil
		}
		*j = JSONStringList{}
		return nil
	}
	return nil
}

// JSONRawMessage is a json.RawMessage-like for gateway_response.
type JSONRawMessage []byte

func (j JSONRawMessage) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONRawMessage) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}
	*j = bytes
	return nil
}
