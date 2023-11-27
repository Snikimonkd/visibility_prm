package model

import "testing"

func TestPolygon_Inside(t *testing.T) {
	type fields struct {
		Points []Point
	}
	type args struct {
		in Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "1. Inside.",
			fields: fields{
				Points: []Point{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
					{X: 1, Y: 0},
				},
			},
			args: args{
				Point{X: 0.5, Y: 0.5},
			},
			want: true,
		},
		{
			name: "2. Outside.",
			fields: fields{
				Points: []Point{
					{X: 0, Y: 0},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
					{X: 1, Y: 0},
				},
			},
			args: args{
				Point{X: 1.5, Y: 0.5},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Polygon{
				Points: tt.fields.Points,
			}
			if got := p.Inside(tt.args.in); got != tt.want {
				t.Errorf("Polygon.Inside() = %v, want %v", got, tt.want)
			}
		})
	}
}
