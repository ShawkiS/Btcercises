package HashSignatureSchemes

import (
	"fmt"
	"testing"

	utils "github.com/Btcercises/HashSignatureSchemes/utils"
)

func TestGoodSig(t *testing.T) {

	msg := utils.GetMessageFromString("good")

	sec, pub, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	sig := Sign(msg, sec)

	worked := Verify(msg, pub, sig)

	if !worked {
		t.Fatalf("Verify returned false, expected true")
	}
}

func TestBadSig(t *testing.T) {

	msg := utils.GetMessageFromString("bad")

	sec, pub, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	sig := Sign(msg, sec)

	sig.Preimage[16] = sig.Preimage[26].Hash()

	worked := Verify(msg, pub, sig)

	if worked {
		t.Fatalf("Verify returned true, expected false")
	}

	msg = utils.GetMessageFromString("worse")
	worked = Verify(msg, pub, sig)

	if worked {
		t.Fatalf("Verify returned true, expected false")
	}
}

func TestGoodMany(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := fmt.Sprintf("good %d", i)
		msg := utils.GetMessageFromString(s)
		sec, pub, err := GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		sig := Sign(msg, sec)

		worked := Verify(msg, pub, sig)
		if !worked {
			t.Fatalf("Verify returned false, expected true")
		}
	}
}

func TestBadMany(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := fmt.Sprintf("bad %d", i)
		msg := utils.GetMessageFromString(s)

		sec, pub, err := GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		sig := Sign(msg, sec)
		sig.Preimage[i%10] = sig.Preimage[i%11].Hash()

		worked := Verify(msg, pub, sig)
		if worked {
			t.Fatalf("Verify returned true, expected false")
		}
	}
}
