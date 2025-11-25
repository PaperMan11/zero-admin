package syncx

import (
	"testing"
)

func TestNewLimit(t *testing.T) {
	limit := NewLimit(2)
	if cap(limit.pool) != 2 {
		t.Errorf("expected capacity 2, got %d", cap(limit.pool))
	}
}

func TestLimit_Borrow_Return(t *testing.T) {
	limit := NewLimit(1)

	// 正常借出并归还
	limit.Borrow()
	if err := limit.Return(); err != nil {
		t.Errorf("unexpected error on return: %v", err)
	}

	// 归还多于借出的情况
	if err := limit.Return(); err != ErrLimitReturn {
		t.Errorf("expected ErrLimitReturn, got %v", err)
	}
}

func TestLimit_TryBorrow(t *testing.T) {
	limit := NewLimit(1)

	// 成功尝试借用
	if !limit.TryBorrow() {
		t.Error("expected TryBorrow to succeed")
	}

	// 失败尝试借用（资源已满）
	if limit.TryBorrow() {
		t.Error("expected TryBorrow to fail")
	}
}
