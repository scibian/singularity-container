// Copyright (c) 2021, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package slice

import "testing"

func TestContainsString(t *testing.T) {
	type args struct {
		s     []string
		match string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "NoMatchSingle",
			args: args{[]string{"a"}, "1"},
			want: false,
		},
		{
			name: "NoMatchMulti",
			args: args{[]string{"a", "b", "c"}, "1"},
			want: false,
		},
		{
			name: "NoMatchEmpty",
			args: args{[]string{}, "1"},
			want: false,
		},
		{
			name: "MatchSingle",
			args: args{[]string{"a"}, "a"},
			want: true,
		},
		{
			name: "MatchMulti",
			args: args{[]string{"a", "b", "c"}, "a"},
			want: true,
		},
		{
			name: "EmptyMatch",
			args: args{[]string{"a", "b", "c"}, ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.s, tt.args.match); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsAnyString(t *testing.T) {
	type args struct {
		s       []string
		matches []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "NoMatchSingle",
			args: args{[]string{"a"}, []string{"1"}},
			want: false,
		},
		{
			name: "NoMatchMulti",
			args: args{[]string{"a", "b", "c"}, []string{"1"}},
			want: false,
		},
		{
			name: "NoMatchEmpty",
			args: args{[]string{}, []string{"1"}},
			want: false,
		},
		{
			name: "NoMatchesSingle",
			args: args{[]string{}, []string{"1", "2", "3"}},
			want: false,
		},
		{
			name: "NoMatchesMulti",
			args: args{[]string{}, []string{"1", "2", "3"}},
			want: false,
		},
		{
			name: "NoMatchesEmpty",
			args: args{[]string{}, []string{"1", "2", "3"}},
			want: false,
		},
		{
			name: "MatchSingle",
			args: args{[]string{"a"}, []string{"a"}},
			want: true,
		},
		{
			name: "MatchMulti",
			args: args{[]string{"a", "b", "c"}, []string{"a"}},
			want: true,
		},
		{
			name: "MatchesSingle",
			args: args{[]string{"a"}, []string{"1", "a", "b"}},
			want: true,
		},
		{
			name: "MatchesMulti",
			args: args{[]string{"a", "b", "c"}, []string{"1", "a", "b"}},
			want: true,
		},
		{
			name: "EmptyMatch",
			args: args{[]string{"a", "b", "c"}, []string{""}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAnyString(tt.args.s, tt.args.matches); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsInt(t *testing.T) {
	type args struct {
		s     []int
		match int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "NoMatchSingle",
			args: args{[]int{1}, 0},
			want: false,
		},
		{
			name: "NoMatchMulti",
			args: args{[]int{1, 2, 3}, 0},
			want: false,
		},
		{
			name: "NoMatchEmpty",
			args: args{[]int{}, 0},
			want: false,
		},
		{
			name: "MatchSingle",
			args: args{[]int{1}, 1},
			want: true,
		},
		{
			name: "MatchMultiStart",
			args: args{[]int{1, 2, 3}, 1},
			want: true,
		},
		{
			name: "MatchMultiMid",
			args: args{[]int{1, 2, 3}, 2},
			want: true,
		},
		{
			name: "MatchMultiEnd",
			args: args{[]int{1, 2, 3}, 2},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsInt(tt.args.s, tt.args.match); got != tt.want {
				t.Errorf("ContainsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
