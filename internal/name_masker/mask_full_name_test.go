package name_masker

import "testing"

func TestMaskFullName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "it should connect two names",
			args: args{
				name: "Celal Şengör",
			},
			want: "Ce*** ***ör",
		},
		{
			name: "it should connect three names",
			args: args{
				name: "Ali Celal Şengör",
			},
			want: "Al* Ce*** ***ör",
		},
		{
			name: "it should connect four names",
			args: args{
				name: "Ali Mehmet Celal Şengör",
			},
			want: "Al* Me*** Ce*** ***ör",
		},
		{
			name: "it should handle cyrillic alphabet",
			args: args{
				name: "Пётр Кропо́ткин",
			},
			want: "Пё** ***ин",
		},
		{
			name: "it should handle long cyrillic  alphabet",
			args: args{
				name: "Пётр Алексе́евич Кропо́ткин",
			},
			want: "Пё** Ал*** ***ин",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskFullName(tt.args.name); got != tt.want {
				t.Errorf("MaskFullName() = %v, want %v", got, tt.want)
			}
		})
	}
}
