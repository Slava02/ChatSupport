package messagesrepo

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

var _ gomock.Matcher = CursorMatcher{}

// CursorMatcher is intended to be used only in tests.
type CursorMatcher struct {
	c Cursor
}

func NewCursorMatcher(c Cursor) CursorMatcher {
	return CursorMatcher{c: c}
}

func (cm CursorMatcher) Matches(x any) bool {
	v, ok := x.(*Cursor)
	if !ok {
		return false
	}

	return (v.PageSize == cm.c.PageSize) && (v.LastCreatedAt.Equal(cm.c.LastCreatedAt))
}

func (cm CursorMatcher) String() string {
	return fmt.Sprintf("{ps=%d, last_created_at=%d}", cm.c.PageSize, cm.c.LastCreatedAt.UnixNano())
}
