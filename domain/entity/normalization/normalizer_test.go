package normalization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeMultipleWhitespaces(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "If the input string has multiple whitespaces in a sequence, they should be converted into a single one",
			args: args{
				s: "abc     def",
			},
			want: "abc def",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveMultipleWhitespaces(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
