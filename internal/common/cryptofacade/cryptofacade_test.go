package cryptofacade

import "testing"

func TestHash(t *testing.T) {
	body := "mock body"
	key := "mock key"
	want := "bW9jayBrZXlxedG6tJtiKDMIoX5mm6zATbSs16R8zavcQb+a5ei51g=="
	if got := Hash(&body, &key); got != want {
		t.Errorf("Hash() = %v, want %v", got, want)
	}
}
