package amqp

// DeliveryGuarantee controls how a message is published.
type DeliveryGuarantee uint8

const (
	// AtLeastOnce enables publisher confirms and persistent delivery (default).
	AtLeastOnce DeliveryGuarantee = iota
	// AtMostOnce is fire-and-forget publishing.
	AtMostOnce
)
