package benchmark_test

import (
	"testing"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/stretchr/testify/assert"
)

func TestNewBenchmark(t *testing.T) {
	b := benchmark.NewBenchmark(46.5, 18)
	assert.NotNil(t, b)
}

func TestBenchmark_GetAvgTime(t *testing.T) {
	type fields struct {
		AvgTime    float64
		Difficulty uint8
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			"GetAvgTime",
			fields{
				AvgTime:    19.0,
				Difficulty: 10,
			},
			19.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := benchmark.Benchmark{
				AvgTime:    tt.fields.AvgTime,
				Difficulty: tt.fields.Difficulty,
			}
			if got := b.GetAvgTime(); got != tt.want {
				t.Errorf("Benchmark.GetAvgTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBenchmark_SetAvgTime(t *testing.T) {
	type fields struct {
		AvgTime    float64
		Difficulty uint8
	}
	type args struct {
		avg float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"SetAvgTime",
			fields{Difficulty: 10},
			args{19.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &benchmark.Benchmark{
				AvgTime:    tt.fields.AvgTime,
				Difficulty: tt.fields.Difficulty,
			}
			b.SetAvgTime(tt.args.avg)
		})
	}
}
