package zlog_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/zlog"
	"github.com/m-mizutani/zlog/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type allFieldFilter struct{}

func (x *allFieldFilter) ReplaceString(s string) string {
	return s
}

func (x *allFieldFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return fieldName != ""
}

func TestClone(t *testing.T) {
	c := zlog.NewMasking(zlog.Filters{
		filter.Value("blue"),
	})

	t.Run("string", func(t *testing.T) {
		v := c.Clone("blue is blue")
		v, ok := v.(string)
		require.True(t, ok)
		assert.Equal(t, zlog.FilteredLabel+" is "+zlog.FilteredLabel, v)
	})

	t.Run("struct", func(t *testing.T) {
		type testData struct {
			ID    int
			Name  string
			Label string
		}

		t.Run("original data is not modified when filtered", func(t *testing.T) {
			data := &testData{
				ID:    100,
				Name:  "blue",
				Label: "five",
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(*testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, zlog.FilteredLabel, copied.Name)
			assert.Equal(t, "blue", data.Name)
			assert.Equal(t, "five", data.Label)
			assert.Equal(t, "five", copied.Label)
			assert.Equal(t, 100, copied.ID)
		})

		t.Run("non-ptr struct can be modified", func(t *testing.T) {
			data := testData{
				Name:  "blue",
				Label: "five",
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, zlog.FilteredLabel, copied.Name)
			assert.Equal(t, "five", copied.Label)
		})

		t.Run("nested structure can be modified", func(t *testing.T) {
			type testDataParent struct {
				Child testData
			}

			data := &testDataParent{
				Child: testData{
					Name:  "blue",
					Label: "five",
				},
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(*testDataParent)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, zlog.FilteredLabel, copied.Child.Name)
			assert.Equal(t, "five", copied.Child.Label)
		})

		t.Run("map data", func(t *testing.T) {
			data := map[string]*testData{
				"xyz": {
					Name:  "blue",
					Label: "five",
				},
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(map[string]*testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, zlog.FilteredLabel, copied["xyz"].Name)
			assert.Equal(t, "five", copied["xyz"].Label)
		})

		t.Run("array data", func(t *testing.T) {
			data := []testData{
				{
					Name:  "orange",
					Label: "five",
				},
				{
					Name:  "blue",
					Label: "five",
				},
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.([]testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "orange", copied[0].Name)
			assert.Equal(t, zlog.FilteredLabel, copied[1].Name)
			assert.Equal(t, "five", copied[1].Label)
		})

		t.Run("array data with ptr", func(t *testing.T) {
			data := []*testData{
				{
					Name:  "orange",
					Label: "five",
				},
				{
					Name:  "blue",
					Label: "five",
				},
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.([]*testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "orange", copied[0].Name)
			assert.Equal(t, zlog.FilteredLabel, copied[1].Name)
			assert.Equal(t, "five", copied[1].Label)
		})

		t.Run("original type", func(t *testing.T) {
			type myType string
			type myData struct {
				Name myType
			}
			data := &myData{
				Name: "miss blue",
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(*myData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, myType("miss "+zlog.FilteredLabel), copied.Name)
		})

		t.Run("unexported field is also copied", func(t *testing.T) {
			type myStruct struct {
				unexported string
				Exported   string
			}

			data := &myStruct{
				unexported: "red",
				Exported:   "orange",
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(*myStruct)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "red", data.unexported)
			assert.Equal(t, "orange", data.Exported)
		})

		t.Run("various field", func(t *testing.T) {
			type child struct{}
			type myStruct struct {
				Func      func() time.Time
				Chan      chan int
				Bool      bool
				Bytes     []byte
				Interface interface{}
				Child     *child
			}
			data := &myStruct{
				Func:  time.Now,
				Chan:  make(chan int),
				Bool:  true,
				Bytes: []byte("timeless"),
				Child: nil,
			}
			v := c.Clone(data)
			require.NotNil(t, v)
			copied, ok := v.(*myStruct)
			require.True(t, ok)
			require.NotNil(t, copied)

			// function type is not compareable, but it's ok if not nil
			assert.NotNil(t, copied.Func)
			assert.Equal(t, data.Chan, copied.Chan)
			assert.Equal(t, data.Bool, copied.Bool)
			assert.Equal(t, data.Bytes, copied.Bytes)
		})
	})

	t.Run("filter various type", func(t *testing.T) {
		mask := zlog.NewMasking(zlog.Filters{
			&allFieldFilter{},
		})
		s := "test"

		type child struct {
			Data string
		}
		type myStruct struct {
			Func      func() time.Time
			Chan      chan int
			Bool      bool
			Bytes     []byte
			Strs      []string
			StrsPtr   []*string
			Interface interface{}
			Child     child
			ChildPtr  *child
		}
		data := &myStruct{
			Func:      time.Now,
			Chan:      make(chan int),
			Bool:      true,
			Bytes:     []byte("timeless"),
			Strs:      []string{"aa"},
			StrsPtr:   []*string{&s},
			Interface: &s,
			Child:     child{Data: "x"},
			ChildPtr:  &child{Data: "y"},
		}

		v := mask.Clone(data)
		require.NotNil(t, v)
		copied, ok := v.(*myStruct)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Nil(t, copied.Func)
		assert.Nil(t, copied.Chan)
		assert.Nil(t, copied.Bytes)
		assert.Nil(t, copied.Strs)
		assert.Nil(t, copied.StrsPtr)
		assert.Nil(t, copied.Interface)
		assert.Empty(t, copied.Child.Data)
		assert.Empty(t, copied.ChildPtr.Data)
	})
}
