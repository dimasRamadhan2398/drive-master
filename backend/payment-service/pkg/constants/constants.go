package constants

const (
	// HTTP Status
	StatusOK           = 200
	StatusCreated      = 201
	StatusAccepted     = 202
	StatusNoContent    = 204
	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusForbidden    = 403
	StatusNotFound     = 404
	StatusConflict     = 409
	StatusGone         = 410
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
	StatusInternalServerError = 500
	StatusPaymentRequired     = 402

	// Pagination
	DefaultPageSize = 20
	MaxPageSize     = 100

	// Payment Status
	PaymentStatusPending   = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed    = "failed"
	PaymentStatusExpired   = "expired"
	PaymentStatusCancelled = "cancelled"
	PaymentStatusRefunded  = "refunded"

	// Transaction Status
	TransactionStatusPending   = "pending"
	TransactionStatusSuccess  = "success"
	TransactionStatusFailed   = "failed"
	TransactionStatusReversed = "reversed"

	// Payment Method Types
	PaymentMethodCreditCard  = "credit_card"
	PaymentMethodDebitCard   = "debit_card"
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodEWallet     = "ewallet"
	PaymentMethodQRIS        = "qris"
	PaymentMethodCOD         = "cod"
	PaymentMethodVA          = "virtual_account"

	// Currency
	DefaultCurrency = "IDR"

	// Payment Gateway
	GatewayMidtrans = "midtrans"
	GatewayXendit   = "xendit"
)

var PaymentStatuses = []string{
	PaymentStatusPending,
	PaymentStatusSuccess,
	PaymentStatusFailed,
	PaymentStatusExpired,
	PaymentStatusCancelled,
	PaymentStatusRefunded,
}

var TransactionStatuses = []string{
	TransactionStatusPending,
	TransactionStatusSuccess,
	TransactionStatusFailed,
	TransactionStatusReversed,
}

var PaymentMethods = []string{
	PaymentMethodCreditCard,
	PaymentMethodDebitCard,
	PaymentMethodBankTransfer,
	PaymentMethodEWallet,
	PaymentMethodQRIS,
	PaymentMethodCOD,
	PaymentMethodVA,
}