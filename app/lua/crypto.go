package lua

import (
	"github.com/yuin/gopher-lua"
	"crypto/sha1"
	"fmt"
)

// Sha1Hash returns the sha1 hash of the given string
func Sha1Hash(L *lua.LState) int {
	// Get string to be hashed
	str := L.Get(2)

	// Check for valid string type
	if str.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Hash string using sha1
	data := sha1.Sum([]byte(str.String()))

	// Convert byte array to string and push tu stack
	L.Push(
		lua.LString(
			fmt.Sprintf("%x", data),
		),
	)

	return 1
}
