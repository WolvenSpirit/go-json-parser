package gojsonparser

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func Test_detect(t *testing.T) {
	s := "\"this_key\": 3"
	buf := bytes.NewBuffer([]byte(s))
	r := bufio.NewReader(buf)
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

	s2 := "\"x\":{\"y\":3}"
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
			want:    "{\"y\":3}",
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

func Test_object(t *testing.T) {
	s := "{\"x\":3}"
	buf := bytes.NewBuffer([]byte(s))
	r := bufio.NewReader(buf)

	s2 := "{\"x\":{\"y\":300}}"
	buf2 := bytes.NewBuffer([]byte(s2))
	r2 := bufio.NewReader(buf2)

	s3 := "{\"x\":{\"z\":{\"y\":\"something\"},\"foo\":3},\"bar\":0}"
	buf3 := bytes.NewBuffer([]byte(s3))
	r3 := bufio.NewReader(buf3)

	s4 := ""
	buf4 := bytes.NewBuffer([]byte(s4))
	r4 := bufio.NewReader(buf4)

	s5 := `{\
				"x\":{
					\"z\":{
						\"y\":\"something\"
						},
					\"foo\":{
						\"j\":{
							\"i\":3
							}
						}
					},
				\"bar\":0
				}`
	buf5 := bytes.NewBuffer([]byte(s5))
	r5 := bufio.NewReader(buf5)

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
			name:    "1d",
			args:    args{r: r},
			want:    "\"x\":3",
			wantErr: false,
		},
		{
			name:    "2d",
			args:    args{r: r2},
			want:    "\"x\":{\"y\":300}",
			wantErr: false,
		},
		{
			name:    "3d",
			args:    args{r: r3},
			want:    "\"x\":{\"z\":{\"y\":\"something\"},\"foo\":3},\"bar\":0",
			wantErr: false,
		},
		{
			name:    "Empty",
			args:    args{r: r4},
			want:    "",
			wantErr: true,
		},
		{
			name:    "5d",
			args:    args{r: r5},
			want:    s5[1 : len(s5)-1],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectObject(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("object() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("object() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse1D(t *testing.T) {
	s5 := `{"x":3,"bar":0,"foo":"something"}`
	m := make(map[string]string)
	m["x"] = "3"
	m["bar"] = "0"
	m["foo"] = "something"
	buf5 := bytes.NewBuffer([]byte(s5))
	r5 := bufio.NewReader(buf5)
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "5d",
			args:    args{r: r5},
			want:    m,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse1D() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse1D() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse2Dimensional(t *testing.T) {
	s := `{"x":3,"bar":0,"foo":{"y":"something"}}`
	buf := bytes.NewBuffer([]byte(s))
	r := bufio.NewReader(buf)
	m := make(map[string]string)
	m["x"] = "3"
	m["bar"] = "0"
	m["foo"] = `{"y":"something"}`
	child := make(map[string]string)
	child["y"] = "something"
	expect := make(map[string]Value)
	expect["x"] = Value{String: "3"}
	expect["bar"] = Value{String: "0"}
	expect["foo"] = Value{String: `{"y":"something"}`, Map: child}
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]Value
		wantErr bool
	}{
		{
			name:    "",
			args:    args{r: r},
			want:    expect,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse2Dimensional(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse2Dimensional() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse2Dimensional() = \n%+v\n%+v", got, tt.want)
			}
		})
	}
}
