package randutils

import "testing"

func TestRandomString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "one",
			args: args{1},
			want: 1,
		},
		{
			name: "ten",
			args: args{10},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomString(tt.args.length); len(got) != tt.want {
				t.Errorf("RandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
