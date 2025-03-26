package numbers

import "testing"

func TestJoinComma(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want string
	}{
		{
			name: "case1",
			args: []int{1, 2, 3, 4},
			want: "1,2,3,4",
		},
	}
	for _, v := range tests {
		if got := JoinComma(v.args, false); got.String() != v.want {
			t.Fail()
		}
	}
}

func TestParseCommaNumber(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []int32
	}{
		{
			name: "case1",
			args: ",1,",
			want: []int32{1},
		},
	}
	for _, v := range tests {
		if got := NewCommaNumber[int32](v.args).Parse(); got == nil {
			t.Fail()
		}
	}
}
