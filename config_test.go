package golivewire

import "testing"

func TestSetBaseURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://localhost:8081", "http://localhost:8081"},
		{"http://localhost:8081/", "http://localhost:8081"},
		{"http://localhost:8081/test", "http://localhost:8081/test"},
		{"http://localhost:8081/test/", "http://localhost:8081/test"},
		{"http://localhost", "http://localhost"},
		{"http://localhost/", "http://localhost"},
		{"http://localhost/test", "http://localhost/test"},
		{"http://localhost/test/", "http://localhost/test"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			SetBaseURL(test.input)
			if baseURL != test.expected {
				t.Errorf("expected %v, got %v", test.expected, baseURL)
			}
		})
	}
}

func TestSetBaseURL_Error(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"/"},
		{"/test"},
		{""},
		{"local/"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("expect error")
				}
			}()

			SetBaseURL(test.input)
		})
	}
}
