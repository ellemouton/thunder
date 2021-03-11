// +build !prod

package passwd

func Protected() (bool, string) {
	return false, ""
}
