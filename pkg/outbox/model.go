package outbox

import (
	"encoding/json"

	"github.com/uptrace/bun"
)

// Entry represents a single outbox row.
type Entry struct {
	bun.BaseModel `bun:"table:outbox"`

	ID        int64           `bun:"id,pk,autoincrement"`
	EventName string          `bun:"event_name,notnull"`
	Payload   json.RawMessage `bun:"payload,type:jsonb,notnull"`
	Headers   map[string]any  `bun:"headers,type:jsonb,notnull,default:'{}'"`
	CreatedAt int64           `bun:"created_at,notnull"`
	Published bool            `bun:"published,notnull,default:false"`
}
