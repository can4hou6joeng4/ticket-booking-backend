package utils

import (
	"context"
	"time"
)

// DefaultTimeout 默认请求超时时间
const DefaultTimeout = 5 * time.Second

// CreateTimeoutContext 创建一个带超时的上下文
func CreateTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return context.WithTimeout(context.Background(), timeout)
}
