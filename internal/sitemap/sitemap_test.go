package sitemap_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"git.sr.ht/~jamesponddotco/sitred/internal/sitemap"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileName string
		want     []string
		wantErr  error
	}{
		{
			name:     "valid sitemap",
			fileName: "valid-sitemap.xml",
			want: []string{
				"http://example.com/page1",
				"http://example.com/page2",
			},
			wantErr: nil,
		},
		{
			name:     "invalid sitemap",
			fileName: "invalid-sitemap.xml",
			want:     nil,
			wantErr:  sitemap.ErrSitemap,
		},
		{
			name:     "empty sitemap",
			fileName: "empty-sitemap.xml",
			want:     []string{},
			wantErr:  nil,
		},
		{
			name:     "sitemap with image",
			fileName: "image-included-sitemap.xml",
			want: []string{
				"http://example.com/page1",
				"http://example.com/page2",
			},
			wantErr: nil,
		},
		{
			name:     "malformed sitemap",
			fileName: "malformed-sitemap.xml",
			want:     nil,
			wantErr:  sitemap.ErrSitemap,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open("testdata/" + tt.fileName)
			if err != nil {
				t.Fatalf("could not open test file: %v", err)
			}
			defer file.Close()

			got, err := sitemap.Parse(file)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)

				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)

					return
				}
			}
		})
	}
}

func TestParseLargeSitemap(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/large-sitemap.xml")
	if err != nil {
		t.Fatalf("could not open test file: %v", err)
	}
	defer file.Close()

	got, err := sitemap.Parse(file)
	if err != nil {
		t.Errorf("Parse() error = %v", err)

		return
	}

	if len(got) < 330000 {
		t.Errorf("Expected at least 330000 URLs, got %d", len(got))
	}
}

func BenchmarkParseLargeSitemap(b *testing.B) {
	file, err := os.Open("testdata/large-sitemap.xml")
	if err != nil {
		b.Fatalf("could not open test file: %v", err)
	}
	defer file.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := sitemap.Parse(file)
		if err != nil {
			b.Fatalf("Parse() error = %v", err)
		}

		// Reset file position for the next iteration.
		_, _ = file.Seek(0, io.SeekStart)
	}
}
