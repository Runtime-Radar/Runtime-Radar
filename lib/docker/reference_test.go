package docker

import (
	"testing"
)

func TestParseReference(t *testing.T) {
	tests := []struct {
		name        string
		ref         string
		expected    *Reference
		expectError bool
	}{
		{
			name:        "Simple image name",
			ref:         "nginx",
			expected:    &Reference{"nginx:latest", ""},
			expectError: false,
		},
		{
			name:        "Image with tag",
			ref:         "nginx:1.20",
			expected:    &Reference{"nginx:1.20", ""},
			expectError: false,
		},
		{
			name:        "Image with digest",
			ref:         "nginx@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa",
			expected:    &Reference{"nginx@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa", ""},
			expectError: false,
		},
		{
			name:        "Image with invalid digest",
			ref:         "nginx@sha256:abc123",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Image with tag and digest",
			ref:         "nginx:1.20@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa",
			expected:    &Reference{"nginx:1.20@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa", ""},
			expectError: false,
		},
		{
			name:        "docker.io library image",
			ref:         "docker.io/library/nginx",
			expected:    &Reference{"nginx:latest", ""},
			expectError: false,
		},
		{
			name:        "docker.io library image with tag",
			ref:         "docker.io/library/nginx:1.20",
			expected:    &Reference{"nginx:1.20", ""},
			expectError: false,
		},
		{
			name:        "docker.io user image",
			ref:         "docker.io/user/nginx",
			expected:    &Reference{"user/nginx:latest", ""},
			expectError: false,
		},
		{
			name:        "docker.io user image with tag",
			ref:         "docker.io/user/nginx:1.20",
			expected:    &Reference{"user/nginx:1.20", ""},
			expectError: false,
		},
		{
			name:        "Custom registry",
			ref:         "gcr.io/project/nginx",
			expected:    &Reference{"project/nginx:latest", "gcr.io"},
			expectError: false,
		},
		{
			name:        "Custom registry with tag",
			ref:         "gcr.io/project/nginx:1.20",
			expected:    &Reference{"project/nginx:1.20", "gcr.io"},
			expectError: false,
		},
		{
			name:        "Custom registry with digest",
			ref:         "gcr.io/project/nginx@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa",
			expected:    &Reference{"project/nginx@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa", "gcr.io"},
			expectError: false,
		},
		{
			name:        "Custom registry with tag and digest",
			ref:         "gcr.io/project/nginx:1.20@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa",
			expected:    &Reference{"project/nginx:1.20@sha256:7cc4b5aefd1d0cadf8d97d4350462ba51c694ebca145b08d7d41b41acc8db5aa", "gcr.io"},
			expectError: false,
		},
		{
			name:        "localhost registry with port",
			ref:         "localhost:5000/nginx",
			expected:    &Reference{"nginx:latest", "localhost:5000"},
			expectError: false,
		},
		{
			name:        "localhost registry with port and tag",
			ref:         "localhost:5000/nginx:1.20",
			expected:    &Reference{"nginx:1.20", "localhost:5000"},
			expectError: false,
		},
		{
			name:        "Complex registry with nested path",
			ref:         "registry.example.com/team/project/service:v1.0.0",
			expected:    &Reference{"team/project/service:v1.0.0", "registry.example.com"},
			expectError: false,
		},
		{
			name:        "Image with underscores and dashes",
			ref:         "my-registry.com/my_team/my-service_v2:latest",
			expected:    &Reference{"my_team/my-service_v2:latest", "my-registry.com"},
			expectError: false,
		},
		{
			name:        "Image with version tag",
			ref:         "alpine:3.14.2",
			expected:    &Reference{"alpine:3.14.2", ""},
			expectError: false,
		},
		{
			name:        "Image with semantic version tag",
			ref:         "node:16.14.0-alpine",
			expected:    &Reference{"node:16.14.0-alpine", ""},
			expectError: false,
		},
		{
			name:        "Empty reference",
			ref:         "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid reference with invalid characters",
			ref:         "nginx:tag with spaces",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Reference with uppercase (should be valid)",
			ref:         "NGINX:LATEST",
			expected:    nil,
			expectError: true, // docker references are case-sensitive and uppercase is invalid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := ParseReference(tt.ref)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected error to be nil, got %v", err)
			}

			if ref.Image != tt.expected.Image {
				t.Errorf("Expected image to be %q, got %q", tt.expected.Image, ref.Image)
			}

			if ref.Registry != tt.expected.Registry {
				t.Errorf("Expected registry to be %q, got %q", tt.expected.Registry, ref.Registry)
			}
		})
	}
}
