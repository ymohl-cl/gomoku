package alignment

import (
	"testing"

	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

func BenchmarkAnalyzeThree_firstCase(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnalyzeThree(&mask)
	}

	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeThree_lastCase(b *testing.B) {
	mask := [11]uint8{0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnalyzeThree(&mask)
	}

	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeThree_noThree(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnalyzeThree(&mask)
	}

	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeWin_noWin(b *testing.B) {
	mask := [11]uint8{0, 0, 2, 1, 0, 1, 1, 1, 2, 0, 0}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnalyzeWin(&mask)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeWin_Win(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		AnalyzeWin(&mask)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkSetStyleByMask_flanked(b *testing.B) {
	mask := [11]uint8{3, 3, 0, 0, 1, 1, 1, 2, 0, 0, 0}
	a := Alignment{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		a.setStyleByMask(&mask, 2, rdef.Player1)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkSetStyleByMask_free(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0}
	a := Alignment{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		a.setStyleByMask(&mask, 2, rdef.Player1)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkSetStyleByMask_half(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 0, 1, 1, 1, 2, 0, 0, 0}
	a := Alignment{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		a.setStyleByMask(&mask, 2, rdef.Player1)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkIsBetter_betterBySize(b *testing.B) {
	var one, two Alignment

	one = Alignment{Size: 2, Style: rdef.AlignFree}
	two = Alignment{Size: 1, Style: rdef.AlignFree}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		one.IsBetter(&two)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkIsBetter_betterByStyle(b *testing.B) {
	var one, two Alignment

	one = Alignment{Size: 2, Style: rdef.AlignFree}
	two = Alignment{Size: 2, Style: rdef.AlignHalf}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		one.IsBetter(&two)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkIsBetter_lowerByStyle(b *testing.B) {
	var one, two Alignment

	one = Alignment{Size: 3, Style: rdef.AlignHalf}
	two = Alignment{Size: 3, Style: rdef.AlignFree}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		one.IsBetter(&two)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkNew_freethree(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 0, 1, 1, 1, 0, 3, 3, 3}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		New(&mask, rdef.Player1)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkNew_noAvailableSpots(b *testing.B) {
	mask := [11]uint8{0, 1, 1, 2, 0, 1, 0, 0, 2, 3, 3}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		New(&mask, rdef.Player1)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkNew_fourSpotsFlanked(b *testing.B) {
	mask := [11]uint8{3, 3, 3, 1, 0, 1, 1, 1, 2, 0, 0}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		New(&mask, rdef.Player1)
	}
	b.StopTimer()
	b.ReportAllocs()
}
