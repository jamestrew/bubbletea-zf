package files

import (
	"testing"
)

func BenchmarkGetWithCSlice(b *testing.B) {
	benchmarks := []struct {
		name string
		dir  string
	}{
		{"current", "."},
		{"small-plenary~125", "/home/jt/projects/plenary.nvim"},
		{"smallish-crush~625", "/home/jt/projects/crush"},
		{"medium-neovim~3500", "/home/jt/projects/neovim"},
		{"large-linux~9000", "/home/jt/projects/linux"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_, err := GetWithCSlice(bm.dir)
				if err != nil {
					b.Fatalf("GetWithCSlice failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkGetWithChannel(b *testing.B) {
	benchmarks := []struct {
		name string
		dir  string
	}{
		{"current", "."},
		{"small-plenary~125", "/home/jt/projects/plenary.nvim"},
		{"smallish-crush~625", "/home/jt/projects/crush"},
		{"medium-neovim~3500", "/home/jt/projects/neovim"},
		{"large-linux~9000", "/home/jt/projects/linux"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_, err := GetWithChannel(bm.dir)
				if err != nil {
					b.Fatalf("GetWithChannel failed: %v", err)
				}
			}
		})
	}
}

/*

  Time Performance:
  Repo      CSlice         Chan          Difference
  current   1,514,251      1,565,650     +3% (slight edge CSlice)
  plenary   3,902,207      4,206,110     +8% (slight edge CSlice)
  crush     10,887,949     10,023,648    -8% (Chan faster!)
  neovim    25,290,002     24,169,349    -4% (Chan faster!)
  linux     221,734,199    208,510,015   -6% (Chan faster!)

  Memory Usage:
  Repo      CSlice         Chan          Difference
  current   346,054        412,880       +19%
  plenary   425,206        480,841       +13%
  crush     842,002        842,845       ~0%
  neovim    3,316,412      2,723,738     -18% (Chan better!)
  linux     36,814,582     26,126,958    -31% (Chan better!)
*/
