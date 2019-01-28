package internal

import (
	"testing"

	"github.com/lighttiger2505/lab/internal/browse"
)

func Test_BrowseMethod_Process(t *testing.T) {
	type fields struct {
		opener browse.URLOpener
		opt    *BrowseOption
		url    string
		id     int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "browse issue list page",
			fields: fields{
				opener: &browse.MockOpener{
					MockOpen: func(url string) error {
						got := url
						want := "https://domain/group/repository/issues"
						if got != want {
							t.Errorf("invalid url, \ngot:%#v\nwant:%#v", got, want)

						}
						return nil
					},
				},
				opt: &BrowseOption{
					Browse: true,
					URL:    false,
					Copy:   false,
				},
				url: "https://domain/group/repository/issues",
				id:  0,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "browse issue detail page",
			fields: fields{
				opener: &browse.MockOpener{
					MockOpen: func(url string) error {
						got := url
						want := "https://domain/group/repository/issues/12"
						if got != want {
							t.Errorf("invalid url, \ngot:%#v\nwant:%#v", got, want)

						}
						return nil
					},
				},
				opt: &BrowseOption{
					Browse: true,
					URL:    false,
					Copy:   false,
				},
				url: "https://domain/group/repository/issues",
				id:  12,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "show issue detail page url",
			fields: fields{
				opener: &browse.MockOpener{},
				opt: &BrowseOption{
					Browse: false,
					URL:    true,
					Copy:   false,
				},
				url: "https://domain/group/repository/issues",
				id:  12,
			},
			want:    "https://domain/group/repository/issues/12",
			wantErr: false,
		},
		{
			name: "copy issue detail page url",
			fields: fields{
				opener: &browse.MockOpener{},
				opt: &BrowseOption{
					Browse: false,
					URL:    false,
					Copy:   true,
				},
				url: "https://domain/group/repository/issues",
				id:  12,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BrowseMethod{
				Opener: tt.fields.opener,
				Opt:    tt.fields.opt,
				URL:    tt.fields.url,
				ID:     tt.fields.id,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("BrowseOption.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BrowseOption.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
