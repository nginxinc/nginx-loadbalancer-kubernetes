package pointer_test

import (
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/pkg/pointer"
	"github.com/stretchr/testify/require"
)

func TestTo(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"", "hello"} {
		require.Equal(t, v, *pointer.To(v))
	}
	for _, v := range []int{0, 123456, -123456} {
		require.Equal(t, v, *pointer.To(v))
	}
	for _, v := range []int64{0, 123456, -123456} {
		require.Equal(t, v, *pointer.To(v))
	}
}

func TestFrom(t *testing.T) {
	t.Parallel()

	sv := "s"
	sd := "default"
	require.Equal(t, sd, pointer.From(nil, sd))
	require.Equal(t, sv, pointer.From(&sv, sd))

	iv := 1
	id := 2
	require.Equal(t, id, pointer.From(nil, id))
	require.Equal(t, iv, pointer.From(&iv, id))

	i64v := int64(1)
	i64d := int64(2)
	require.Equal(t, i64d, pointer.From(nil, i64d))
	require.Equal(t, i64v, pointer.From(&i64v, i64d))
}

func TestToSlice_FromSlice(t *testing.T) {
	t.Parallel()

	v := []int{1, 2, 3}
	require.Equal(t, v, pointer.FromSlice(pointer.ToSlice(v)))
	require.Nil(t, pointer.ToSlice([]string{}))
	require.Nil(t, pointer.FromSlice([]*string{}))
	require.Equal(t, []string{"A", "B"}, pointer.FromSlice([]*string{pointer.To("A"), nil, pointer.To("B")}))
}

func TestEqual(t *testing.T) {
	t.Parallel()

	require.True(t, pointer.Equal(pointer.To(1), 1))
	require.False(t, pointer.Equal(nil, 1))
	require.False(t, pointer.Equal(pointer.To(1), 2))

	s := new(struct{})
	require.False(t, pointer.Equal(&s, nil))
}
