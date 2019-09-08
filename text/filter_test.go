package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	okChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func TestJustAlphaNumeric(t *testing.T) {
	assert.Equal(t, okChars, okChars)
	assert.Equal(t, okChars, JustAlphaNumeric(
		"a~b-c=d_e+f[g]h{i}j\\k|l/m?n.o>p,q<r`s!t@u#v$w%x^y&z*A(B)C-D_E=F+++++G$#%&*HIJKLMNO[]P QRS'''TU,.,.V - WXYZ0_123,.456789"))
}
