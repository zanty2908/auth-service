package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePhone(t *testing.T) {
	phoneList := []string{
		"+84909000999",
		"+61416624855",
		"+49 176 1234 5678",
		"+5681.8383-7370",
		"+35585 22 34 10",
	}

	for _, value := range phoneList {
		phone, err := ValidatePhone(value)
		require.NoError(t, err, "Phone: "+value)
		require.NotEmpty(t, phone)
	}

}
