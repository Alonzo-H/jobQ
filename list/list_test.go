package list_test

import (
	"jobQ/list"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyListPop(t *testing.T) {
	l := list.New()
	_, err := l.Pop()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestEmptyListHead(t *testing.T) {
	l := list.New()
	_, err := l.Head()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestListPushHead(t *testing.T) {
	l := list.New()
	l.Push(1)
	i, err := l.Head()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
}

func TestListPushPop(t *testing.T) {
	l := list.New()
	l.Push(1)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
}

func TestListPushPopHead(t *testing.T) {
	l := list.New()
	l.Push(1)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	_, err = l.Head()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestListPushPopPop(t *testing.T) {
	l := list.New()
	l.Push(1)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	_, err = l.Pop()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestListPushHeadPop(t *testing.T) {
	l := list.New()
	l.Push(1)
	i, err := l.Head()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	i, err = l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
}

func TestListPushHeadHead(t *testing.T) {
	l := list.New()
	l.Push(1)
	i, err := l.Head()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	i, err = l.Head()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
}

func TestListPushPushHead(t *testing.T) {
	l := list.New()
	l.Push(1)
	l.Push(2)
	i, err := l.Head()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
}

func TestListPushPushPop(t *testing.T) {
	l := list.New()
	l.Push(1)
	l.Push(2)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
}

func TestListPushPushPopHead(t *testing.T) {
	l := list.New()
	l.Push(1)
	l.Push(2)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	i, err = l.Head()
	require.Nil(t, err)
	require.EqualValues(t, 2, i)
}

func TestListPushPushPopPop(t *testing.T) {
	l := list.New()
	l.Push(1)
	l.Push(2)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	i, err = l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 2, i)
}

func TestListPushPushPopPopHead(t *testing.T) {
	l := list.New()
	l.Push(1)
	l.Push(2)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)
	i, err = l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 2, i)
	_, err = l.Head()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestListPushPopPushPopPushPopPop(t *testing.T) {
	l := list.New()

	l.Push(1)
	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)

	l.Push(2)
	i, err = l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 2, i)

	l.Push(3)
	i, err = l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 3, i)

	_, err = l.Pop()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestSingleNodeRemove(t *testing.T) {
	l := list.New()

	n := l.Push(1)
	l.Remove(n)

	_, err := l.Pop()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestDoubleNodeRemoveHead(t *testing.T) {
	l := list.New()

	n := l.Push(1)
	l.Push(2)

	l.Remove(n)

	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, i, 2)

	_, err = l.Pop()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestDoubleNodeRemoveTail(t *testing.T) {
	l := list.New()

	l.Push(1)
	n := l.Push(2)

	l.Remove(n)

	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, i, 1)

	_, err = l.Pop()
	require.ErrorIs(t, err, list.ErrEmptyList)
}

func TestTripleNodeRemoveMiddle(t *testing.T) {
	l := list.New()

	l.Push(1)
	n := l.Push(2)
	l.Push(3)

	l.Remove(n)

	i, err := l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 1, i)

	i, err = l.Pop()
	require.Nil(t, err)
	require.EqualValues(t, 3, i)
}
