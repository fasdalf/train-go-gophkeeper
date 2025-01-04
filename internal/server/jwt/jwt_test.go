package jwt

import (
	"fmt"
	"time"
)

func Example() {
	startID := uint64(100)
	key := "mock key"
	exp := 3 * time.Hour

	gotToken, err := BuildJWTString(startID, &key, exp)
	if err != nil {
		fmt.Println(err)
	}

	gotID, err := GetUserID(&gotToken, &key)
	if err != nil {
		fmt.Println(err)
	}

	if gotID != startID {
		fmt.Println("got:", gotID, "want:", startID)
	} else {
		fmt.Println("IDs match")
	}

	// Output:
	// IDs match
}
