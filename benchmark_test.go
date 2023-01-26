package gclog

import (
	"errors"
	"io"
	"math"
	"testing"
	"time"
)

var ints = []int{30, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
var uints = []uint{30, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
var err = errors.New("hi")

func BenchmarkStructuredLogging(b *testing.B) {
	// f, e := os.Create("l.log")
	// if e != nil {
	// 	fmt.Println(e.Error())
	// }
	// log := New(f, true)
	log := New(io.Discard, true)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cl := log.With().Str("name", "User").Str("id", "1").Bool("verified", true).Logger()

			cl.StartJson().
				Str("name", "user").
				Strs("name", []string{"user1", "user1", "user1", "user1", "user1", "user1", "user1", "user1", "user1", "user1"}).
				Int("age", 30).
				Ints("age", ints).
				Uints("age", uints).
				Bool("married", true).
				Float32("k1", 1.0).
				Float64("k2", 2.0).
				Floats32("k3", []float32{1.0, 2.0}).
				Floats64("k4", []float64{1.0, 2.0}).
				Floats32("k5", []float32{float32(math.NaN()), float32(math.Inf(-1)), float32(math.Inf(1))}).
				Floats64("k6", []float64{math.NaN(), math.Inf(-1), math.Inf(1)}).
				Float32("k7", float32(math.NaN())).
				Float64("k8", math.Inf(-1)).
				Float64("k9", math.Inf(1)).
				Bytes("bytes", []byte("hello bytes")).
				Int8("k_int", 10).
				Ints8("k_int", []int8{1, 2, 3}).
				Int16("k_int", 10).
				Ints16("k_int", []int16{1, 2, 3}).
				Int32("k_int", 10).
				Ints32("k_int", []int32{1, 2, 3}).
				Int64("k_int", 10).
				Ints64("k_int", []int64{1, 2, 3}).
				Uint8("k_uint", 10).
				Uints8("k_uint", []uint8{1, 2, 3}).
				Uint16("k_uint", 10).
				Uints16("k_uint", []uint16{1, 2, 3}).
				Uint32("k_uint", 10).
				Uints32("k_uint", []uint32{1, 2, 3}).
				Uint64("k_uint", 10).
				Uints64("k_uint", []uint64{1, 2, 3}).
				Uint64("k_uint", 10).
				Bools("jj", []bool{true, false}).
				Err(nil).
				Err(err).
				Time("tt", time.Now()).
				Dur("dur", 1*time.Hour). // Allocates: 24 B
				Interface("iface", nil). // Allocates:  8 B
				Msgf("%s", "a")
			cl.EndWith()
		}
	})
}
