package ccmp

import (
	"onij/util/boost/exp"
	"testing"
)

func TestPtrEqual(t *testing.T) {
	type args[T comparable] struct {
		a1 *T
		a2 *T
	}
	type Case[T comparable] struct {
		name string
		args args[T]
		want bool
	}

	testsString := []Case[string]{
		{
			name: "case1",
			args: args[string]{a1: nil, a2: nil},
			want: true,
		},
		{
			name: "case2",
			args: args[string]{a1: nil, a2: exp.Ptr(exp.Zero[string]())},
			want: false,
		},
		{
			name: "case3",
			args: args[string]{a1: exp.Ptr(exp.Zero[string]()), a2: nil},
			want: false,
		},
		{
			name: "case4",
			args: args[string]{a1: exp.Ptr(exp.Zero[string]()), a2: exp.Ptr(exp.Zero[string]())},
			want: true,
		},
		{
			name: "case5",
			args: args[string]{a1: exp.Ptr("123"), a2: exp.Ptr("123")},
			want: true,
		},
		{
			name: "case6",
			args: args[string]{a1: exp.Ptr("123"), a2: exp.Ptr("321")},
			want: false,
		},
	}
	testsNumber := []Case[int]{
		{
			name: "case6",
			args: args[int]{a1: nil, a2: nil},
			want: true,
		},
		{
			name: "case7",
			args: args[int]{a1: nil, a2: exp.Ptr(exp.Zero[int]())},
			want: false,
		},
		{
			name: "case8",
			args: args[int]{a1: exp.Ptr(exp.Zero[int]()), a2: nil},
			want: false,
		},
		{
			name: "case9",
			args: args[int]{a1: exp.Ptr(exp.Zero[int]()), a2: exp.Ptr(exp.Zero[int]())},
			want: true,
		},
		{
			name: "case10",
			args: args[int]{a1: exp.Ptr(123), a2: exp.Ptr(123)},
			want: true,
		},
		{
			name: "case11",
			args: args[int]{a1: exp.Ptr(123), a2: exp.Ptr(321)},
			want: false,
		},
	}
	for _, test := range testsString {
		if equals := PtrEqual(test.args.a1, test.args.a2); equals != test.want {
			t.Fail()
		}
	}
	for _, test := range testsNumber {
		if equals := PtrEqual(test.args.a1, test.args.a2); equals != test.want {
			t.Fail()
		}
	}
}

func TestPtrValueOrZeroEqual(t *testing.T) {
	type args[T comparable] struct {
		a1 *T
		a2 *T
	}
	type Case[T comparable] struct {
		name string
		args args[T]
		want bool
	}

	testsString := []Case[string]{
		{
			name: "case1",
			args: args[string]{a1: nil, a2: nil},
			want: true,
		},
		{
			name: "case2",
			args: args[string]{a1: nil, a2: exp.Ptr(exp.Zero[string]())},
			want: true,
		},
		{
			name: "case3",
			args: args[string]{a1: exp.Ptr(exp.Zero[string]()), a2: nil},
			want: true,
		},
		{
			name: "case4",
			args: args[string]{a1: exp.Ptr(exp.Zero[string]()), a2: exp.Ptr(exp.Zero[string]())},
			want: true,
		},
		{
			name: "case5",
			args: args[string]{a1: exp.Ptr("123"), a2: exp.Ptr("123")},
			want: true,
		},
		{
			name: "case6",
			args: args[string]{a1: exp.Ptr("123"), a2: exp.Ptr("321")},
			want: false,
		},
	}
	testsNumber := []Case[int]{
		{
			name: "case6",
			args: args[int]{a1: nil, a2: nil},
			want: true,
		},
		{
			name: "case7",
			args: args[int]{a1: nil, a2: exp.Ptr(exp.Zero[int]())},
			want: true,
		},
		{
			name: "case8",
			args: args[int]{a1: exp.Ptr(exp.Zero[int]()), a2: nil},
			want: true,
		},
		{
			name: "case9",
			args: args[int]{a1: exp.Ptr(exp.Zero[int]()), a2: exp.Ptr(exp.Zero[int]())},
			want: true,
		},
		{
			name: "case10",
			args: args[int]{a1: exp.Ptr(123), a2: exp.Ptr(123)},
			want: true,
		},
		{
			name: "case11",
			args: args[int]{a1: exp.Ptr(123), a2: exp.Ptr(321)},
			want: false,
		},
	}
	for _, test := range testsString {
		if equals := PtrValueOrZeroEqual(test.args.a1, test.args.a2); equals != test.want {
			t.Fail()
		}
	}
	for _, test := range testsNumber {
		if equals := PtrValueOrZeroEqual(test.args.a1, test.args.a2); equals != test.want {
			t.Fail()
		}
	}
}

func TestRichTextBasicEqual(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "case1",
			args: []string{
				`
					<div>
						<p>Hello, world!</p>
						<img src="image1.jpg" width="100" height="200" alt="Image1">
						<p>This is a test.</p>
						<img src="image2.jpg" height="150">
						<img src="image3.jpg" width="300">
						<img src="image4.jpg" width="150" height="250">
					</div>
				`,
				`
					<div>
						<p>Hello, world!</p>
						<img src="image1.jpg" width="100" height="200" alt="Image1">
						<p>This is a test.</p>
						<img src="image2.jpg" height="150">
						<img src="image3.jpg" width="300">
						<img src="image4.jpg" width="150" height="250">
					</div>
				`,
			},
			want: true,
		},
		{
			name: "case2",
			args: []string{
				`
					<p>．已知<tex data-type="0">(ax-2)(x+\frac {2}{x})^{5}</tex>的展开式中的常数项为240，则 <answer data-type="3"></answer>
						<tex data-type="0">a=</tex>
						<answer data-type="3"></answer> ．
					</p>
				`,
				`
					<p>若<img src="https://amp-dev.tos-cn-beijing.volces.com/muse/img/882b42a662b4710cb9688e2fa7a28a5e.png" data-ppi="96"
							width="11" height="9">的角平分线交<em>BC</em>于<img
							src="https://amp-dev.tos-cn-beijing.volces.com/muse/img/740b8598c59da4969706f071e6c4da58.png" data-ppi="96"
							width="12" height="11">，且<tex data-type="0">c=2</tex>，求<tex data-type="0">\triangle ABD</tex>面积的取值范围． <answer
							data-type="4"></answer>
					</p>
				`,
			},
			want: false,
		},
	}
	for _, v := range tests {
		if got := RichTextBasicEqual(v.args[0], v.args[1]); got != v.want {
			t.Fail()
		}
	}
}
