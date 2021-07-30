package compose

import (
	"reflect"
	"testing"
	"time"
)

func Test_getNearestInterval(t *testing.T) {
	type args struct {
		reference time.Time
		interval  time.Duration
		now       time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "0",
			args: args{
				reference: time.Time{},
				interval:  24 * time.Hour,
				now:       time.Time{}.Add(36 * time.Hour),
			},
			want:    time.Time{},
			wantErr: false,
		},
		{
			name: "1",
			args: args{
				reference: time.Time{},
				interval:  24 * time.Hour,
				now:       time.Time{}.Add(50 * time.Hour),
			},
			want:    time.Time{}.Add(24 * time.Hour),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNearestInterval(tt.args.reference, tt.args.interval, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNearestInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNearestInterval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
