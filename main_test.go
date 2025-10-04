package main

import (
	"testing"
)

func BenchmarkGetFilesCSlice(b *testing.B) {
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
				_, err := getFilesCSlice(bm.dir)
				if err != nil {
					b.Fatalf("getFilesCSlice failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkGetFilesChan(b *testing.B) {
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
				_, err := getFilesChan(bm.dir)
				if err != nil {
					b.Fatalf("getFilesChan failed: %v", err)
				}
			}
		})
	}
}
