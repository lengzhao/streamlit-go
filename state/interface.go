package state

import (
	"time"

	"github.com/lengzhao/streamlit-go/widgets"
)

// SessionInterface 会话接口，避免循环依赖
type SessionInterface interface {
	ID() string
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Has(key string) bool
	Clear()
	LastAccessedAt() time.Time
	CreatedAt() time.Time
	AddWidget(widget widgets.Widget)
	GetWidgets() []widgets.Widget
	ClearWidgets()
}
