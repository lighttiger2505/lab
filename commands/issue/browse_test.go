package issue

import (
	"testing"

	"github.com/lighttiger2505/lab/cmd"
)

func Test_browseMethod_Process(t *testing.T) {
	type fields struct {
		opener cmd.URLOpener
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
				opener: &cmd.MockOpener{
					MockOpen: func(url string) error {
						got := url
						want := "https://domain/group/repository/issues"
						if got != want {
							t.Errorf("invalid url, \ngot:%#v\nwant:%#v", got, want)

						}
						return nil
					},
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
				opener: &cmd.MockOpener{
					MockOpen: func(url string) error {
						got := url
						want := "https://domain/group/repository/issues/12"
						if got != want {
							t.Errorf("invalid url, \ngot:%#v\nwant:%#v", got, want)

						}
						return nil
					},
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
			m := &browseMethod{
				opener: tt.fields.opener,
				url:    tt.fields.url,
				id:     tt.fields.id,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("browseMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("browseMethod.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
