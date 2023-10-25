package services_test

import (
	"fmt"
	"go-unit-test/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckGrade(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected string
	}{
		{"A", 80, "A"},
		{"B", 70, "B"},
		{"C", 60, "C"},
		{"D", 50, "D"},
		{"F", 40, "F"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test CheckGrade for grade %s", tt.name), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, services.CheckGrade(tt.score))
		})
	}
}

func BenchmarkCheckGrade(b *testing.B) {
	for i := 0; i < b.N; i++ {
		services.CheckGrade(80)
	}
}

func ExampleCheckGrade() {
	fmt.Println(services.CheckGrade(80))
	// Output: A
}
