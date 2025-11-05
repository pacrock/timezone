package timezone

import "testing"

func sec(h, m, s int) int {
	return h*3600 + m*60 + s
}

func TestParseOffset(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		// --- Valid cases ---
		{name: "Z", input: "Z", want: 0, wantErr: false},
		{name: "UTC", input: "UTC", want: 0, wantErr: false},
		{name: "GMT", input: "GMT", want: 0, wantErr: false},
		{name: "UTC+0", input: "UTC+0", want: 0, wantErr: false},
		{name: "UTC-0", input: "UTC-0", want: 0, wantErr: false},
		{name: "UTC+00:00", input: "UTC+00:00", want: 0, wantErr: false},
		{name: "GMT+0000", input: "GMT+0000", want: 0, wantErr: false},

		{name: "UTC+5", input: "UTC+5", want: sec(5, 0, 0), wantErr: false},
		{name: "UTC+05", input: "UTC+05", want: sec(5, 0, 0), wantErr: false},
		{name: "GMT-07:00", input: "GMT-07:00", want: sec(-7, 0, 0), wantErr: false},
		{name: "GMT+0530", input: "GMT+0530", want: sec(5, 30, 0), wantErr: false},
		{name: "GMT+530", input: "GMT+530", want: sec(5, 30, 0), wantErr: false},

		{name: "+H", input: "+5", want: sec(5, 0, 0), wantErr: false},
		{name: "+HH", input: "+05", want: sec(5, 0, 0), wantErr: false},
		{name: "+HMM", input: "+530", want: sec(5, 30, 0), wantErr: false},
		{name: "+HHMM", input: "+0530", want: sec(5, 30, 0), wantErr: false},
		{name: "+HH:MM", input: "+05:30", want: sec(5, 30, 0), wantErr: false},

		{name: "-H", input: "-7", want: sec(-7, 0, 0), wantErr: false},
		{name: "-HH", input: "-07", want: sec(-7, 0, 0), wantErr: false},
		{name: "-HMM", input: "-700", want: sec(-7, 0, 0), wantErr: false},
		{name: "-HHMM", input: "-0700", want: sec(-7, 0, 0), wantErr: false},
		{name: "-HH:MM", input: "-07:00", want: sec(-7, 0, 0), wantErr: false},
		{name: "Num -09:45", input: "-09:45", want: sec(-9, -45, 0), wantErr: false},

		{name: "Max +14:00", input: "+14:00", want: sec(14, 0, 0), wantErr: false},
		{name: "Max +1400", input: "+1400", want: sec(14, 0, 0), wantErr: false},
		{name: "Max +14", input: "+14", want: sec(14, 0, 0), wantErr: false},
		{name: "Min -14:00", input: "-14:00", want: sec(-14, 0, 0), wantErr: false},
		{name: "Min -1400", input: "-1400", want: sec(-14, 0, 0), wantErr: false},
		{name: "Min -14", input: "-14", want: sec(-14, 0, 0), wantErr: false},

		// --- Invalid cases ---
		{name: "Empty", input: "", wantErr: true},
		{name: "Just +", input: "+", wantErr: true},
		{name: "Just -", input: "-", wantErr: true},
		{name: "UTC+", input: "UTC+", wantErr: true},
		{name: "GMT-", input: "GMT-", wantErr: true},

		{name: "No sign 05:00", input: "05:00", wantErr: true},
		{name: "No sign 0500", input: "0500", wantErr: true},
		{name: "No sign 5", input: "5", wantErr: true},

		{name: "Location PST", input: "PST", wantErr: true},
		{name: "Location EST", input: "EST", wantErr: true},
		{name: "Location Full", input: "America/New_York", wantErr: true},
		{name: "Invalid chars", input: "+05:XX", wantErr: true},
		{name: "Invalid trailer", input: "+05:00Z", wantErr: true},
		{name: "Invalid trailer 2", input: "Z+01:00", wantErr: true},
		{name: "Invalid format H:MM", input: "+5:00", wantErr: true},
		{name: "Invalid format HH:M", input: "+05:0", wantErr: true},
		{name: "Invalid prefix", input: "UTC+05:00abc", wantErr: true},
		{name: "Invalid format ::", input: "+05::00", wantErr: true},

		{name: "Hour > 14", input: "+15:00", wantErr: true},
		{name: "Hour > 14 num", input: "+1500", wantErr: true},
		{name: "Hour > 14 short", input: "+15", wantErr: true},
		{name: "Hour < -14", input: "-15:00", wantErr: true},
		{name: "Minute > 59", input: "+05:60", wantErr: true},
		{name: "Minute > 59 num", input: "+0560", wantErr: true},
		{name: "Minute > 59 HMM", input: "+560", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOffset(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOffset(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseOffset(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkParseOffset(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{name: "Z", input: "Z"},
		{name: "UTC", input: "UTC"},
		{name: "GMT", input: "GMT"},
		{name: "HH:MM", input: "+05:30"},
		{name: "HHMM", input: "-0400"},
		{name: "UTC+H", input: "UTC+5"},
		{name: "GMT-HH:MM", input: "GMT-07:00"},
		{name: "Error", input: "PST"},
		{name: "ErrorLong", input: "America/New_York"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				// Results must be used to prevent the compiler
				// from optimizing away the function call.
				_, _ = ParseOffset(bm.input)
			}
		})
	}
}
