package conv

import (
	"testing"
)

func TestBooleanToInt(t *testing.T) {
	type args struct {
		b1 bool
		b2 *bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{b1: true},
			want: 1,
		},
		{
			name: "case2",
			args: args{b1: false},
			want: 0,
		},
		{
			name: "case3",
			args: args{b2: exp.Ptr(true)},
			want: 1,
		},
		{
			name: "case4",
			args: args{b2: exp.Ptr(false)},
			want: 0,
		},
	}
	for _, v := range tests {
		if v.args.b2 != nil {
			if got := BooleanToInt[int](v.args.b2); got != v.want {
				t.Fail()
			}
			continue
		}

		if got := BooleanToInt[int](v.args.b1); got != v.want {
			t.Fail()
		}
	}
}

func TestStringToFloat64(t *testing.T) {
	type args struct {
		b1 string
		b2 *string
	}
	tests := []struct {
		name string
		args args
		want *float64
	}{
		{
			name: "case1",
			args: args{b1: "1.1"},
			want: exp.Ptr(1.1),
		},
		{
			name: "case2",
			args: args{b1: "b"},
			want: nil,
		},
		{
			name: "case3",
			args: args{b2: exp.Ptr("1")},
			want: exp.Ptr(float64(1)),
		},
		{
			name: "case4",
			args: args{b2: exp.Ptr("")},
			want: nil,
		},
	}
	for _, v := range tests {
		var got *float64
		if v.args.b2 != nil {
			got = StringToFloat64(v.args.b2)
		} else {
			got = StringToFloat64(v.args.b1)
		}
		if got == nil && v.want != nil {
			t.Fail()
		}
		if got != nil && v.want == nil {
			t.Fail()
		}
		if got == nil && v.want == nil {
			continue
		}
		if *got != *v.want {
			t.Fail()
		}
	}
}
