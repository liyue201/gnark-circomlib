package circuits

import "testing"

func TestPOSEIDON_C(t *testing.T) {
	for i := 2; i <= 17; i++ {
		a := POSEIDON_C(i)
		t.Logf("len(a) = %v", len(a))
		t.Logf("a[0] = %v", a[0].Text(16))
	}
}

func TestPOSEIDON_M(t *testing.T) {
	for i := 2; i <= 17; i++ {
		a := POSEIDON_M(i)
		t.Logf("len(a) = %v", len(a))
		for i := 0; i < len(a); i++ {
			t.Logf("a[i] = %v", len(a[i]))
		}
	}
}

func TestPOSEIDON_S(t *testing.T) {
	for i := 2; i <= 17; i++ {
		a := POSEIDON_S(i)
		t.Logf("len(a) = %v", len(a))
		t.Logf("a[0] = %v", a[0].Text(16))
	}
}

func TestPOSEIDON_P(t *testing.T) {
	for i := 2; i <= 17; i++ {
		a := POSEIDON_P(i)
		t.Logf("len(a) = %v", len(a))
		for i := 0; i < len(a); i++ {
			t.Logf("a[i] = %v", len(a[i]))
		}
	}
}
