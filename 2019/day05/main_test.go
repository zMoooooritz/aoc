package main

import (
	"testing"
)

var example = ``

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "add",
			input: "1,0,0,0,99",
			want:  2,
		},
		{
			name:  "mul",
			input: "2,1,0,0,99",
			want:  2,
		},
		{
			name:  "advanced",
			input: "1,1,1,4,99,5,6,0,99",
			want:  30,
		},
		{
			name:  "immediate",
			input: "1002,4,3,4,33",
			want:  1002,
		},
		{
			name:  "neg",
			input: "1101,100,-1,4,0",
			want:  1101,
		},
		{
			name:  "equal",
			input: "3,9,8,9,10,9,4,9,99,-1,8",
			want:  3,
		},
		{
			name:  "jump",
			input: "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			want:  3,
		},
		{
			name:  "full",
			input: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			want:  3,
		},

		// {
		// 	name:  "actual",
		// 	input: input,
		// 	want:  0,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		// {
		// 	name:  "example",
		// 	input: example,
		// 	want:  0,
		// },
		// {
		// 	name:  "actual",
		// 	input: input,
		// 	want:  0,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
