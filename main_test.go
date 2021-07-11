package gojsonparser

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_detect(t *testing.T) {
	s := "      \"this_key\":3"
	buf := bytes.NewBuffer([]byte(s))
	r := bufio.NewReader(buf)
	buf.WriteString(s)
	type args struct {
		r     *bufio.Reader
		delim rune
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{r: r, delim: '"'},
			want:    "this_key",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detect(tt.args.r, tt.args.delim)
			if (err != nil) != tt.wantErr {
				t.Errorf("detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_detectSimpleValue(t *testing.T) {
	s := "\"x\":3"
	buf := bytes.NewBuffer([]byte(s))
	r := bufio.NewReader(buf)

	s2 := "\"x\":300"
	buf2 := bytes.NewBuffer([]byte(s2))
	r2 := bufio.NewReader(buf2)

	s3 := "\"x\":\"something\""
	buf3 := bytes.NewBuffer([]byte(s3))
	r3 := bufio.NewReader(buf3)

	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{r: r},
			want:    "3",
			wantErr: false,
		},
		{
			name:    "",
			args:    args{r: r2},
			want:    "300",
			wantErr: false,
		},
		{
			name:    "",
			args:    args{r: r3},
			want:    "something",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectSimpleValue(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectSimpleValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("detectSimpleValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
