package vm

import "testing"

func TestY(t *testing.T) {
	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 100} {
		n1 := lazyFactorial(NewNumber(n))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		if n1 != n2 {
			t.Fail()
		}
	}

	for _, ts := range [][]*Thunk{
		{NewNumber(7)},
		{NewNumber(13), NewString("foobarbaz")},
		{NewNumber(42), Nil, Nil},
	} {
		t.Log(lazyFactorial(ts...))
	}
}

func strictFactorial(n float64) float64 {
	if n == 0 {
		return 1
	}

	return n * strictFactorial(n-1)
}

func lazyFactorial(ts ...*Thunk) float64 {
	return float64(App(App(Y, NewLazyFunction(lazyFactorialImpl)), ts...).Eval().(numberType))
}

func lazyFactorialImpl(ts ...*Thunk) Object {
	// fmt.Println(len(ts))

	return App(If,
		App(Equal, ts[1], NewNumber(0)),
		NewNumber(1),
		App(Mult,
			ts[1],
			App(ts[0], append([]*Thunk{App(Sub, ts[1], NewNumber(1))}, ts[2:]...)...)))
}

func TestYsSingleF(t *testing.T) {
	fs := App(Ys, NewLazyFunction(lazyFactorialImpl)).Eval().([]*Thunk)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 100} {
		n1 := float64(App(fs[0], NewNumber(n)).Eval().(numberType))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		if n1 != n2 {
			t.Fail()
		}
	}
}

func TestYsMultipleFs(t *testing.T) {
	even := NewLazyFunction(func(ts ...*Thunk) Object {
		n := ts[2]

		return App(If,
			App(Equal, n, NewNumber(0)),
			True,
			App(ts[1], App(Sub, n, NewNumber(1))))
	})

	odd := NewLazyFunction(func(ts ...*Thunk) Object {
		n := ts[2]

		return App(If,
			App(Equal, n, NewNumber(0)),
			False,
			App(ts[0], App(Sub, n, NewNumber(1))))
	})

	fs := App(Ys, even, odd).Eval().([]*Thunk)

	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 100, 121, 256, 1023} {
		b1 := bool(App(fs[0], NewNumber(n)).Eval().(boolType))
		b2 := bool(App(fs[1], NewNumber(n)).Eval().(boolType))

		t.Logf("n = %v, even? %v, odd? %v\n", n, b1, b2)

		if b1 != (int(n)%2 == 0) {
			t.Fail()
		}

		if b2 != (int(n)%2 != 0) {
			t.Fail()
		}
	}
}
