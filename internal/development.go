//go:build !production

package internal

func init() {
	Version = "development"
}
