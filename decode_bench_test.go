package toml

import (
	"testing"
	"time"
)

func BenchmarkUnmarshal(b *testing.B) {
	var v struct {
		Title string
		Owner struct {
			Name         string
			Organization string
			Bio          string
			Dob          time.Time
		}
		Database struct {
			Server        string
			Ports         []int
			ConnectionMax int
			Enabled       bool
		}
		Servers struct {
			Alpha struct {
				IP string
				DC string
			}
			Beta struct {
				IP string
				DC string
			}
		}
		Clients struct {
			Data  []interface{}
			Hosts []string
		}
	}
	data := loadTestData("test.toml")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &v); err != nil {
			b.Fatal(err)
		}
	}
}
