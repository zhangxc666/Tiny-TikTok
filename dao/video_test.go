package dao

import (
	"testing"
)

func TestVideoDao_UpdateVideoInfo(t *testing.T) {
	DBInit()
	type args struct {
		video *Video
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: testing.CoverMode(), args: args{video: &Video{
			ID:        1,
			TimeStamp: 4,
		}}, wantErr: false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vi := VideoDao{}
			if err := vi.UpdateVideoInfo(tt.args.video); (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideoInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
