package cmd

import "testing"

func TestIsTursoDB(t *testing.T) {
	tests := []struct {
		name  string
		uuid  string
		turso bool
	}{
		{
			name:  "type byte is 0x10",
			uuid:  "019db7a2-9210-79e5-afed-0b1755901d50",
			turso: true,
		},
		{
			name:  "type byte is not 0x10",
			uuid:  "00000000-1100-0000-0000-000000000000",
			turso: false,
		},
		{
			name:  "invalid uuid",
			uuid:  "not-a-uuid",
			turso: false,
		},
		{
			name:  "empty uuid",
			uuid:  "",
			turso: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTursoDB(tt.uuid)
			if got != tt.turso {
				t.Fatalf("isTursoDB(%q) = %v, want %v", tt.uuid, got, tt.turso)
			}
		})
	}
}

func TestDatabaseType(t *testing.T) {
	tests := []struct {
		uuid string
		want string
	}{
		{
			uuid: "019db7a2-9210-79e5-afed-0b1755901d50",
			want: "Turso",
		},
		{
			uuid: "00000000-1100-0000-0000-000000000000",
			want: "SQLite",
		},
		{
			uuid: "",
			want: "SQLite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := databaseType(tt.uuid)
			if got != tt.want {
				t.Fatalf("databaseType(%q) = %q, want %q", tt.uuid, got, tt.want)
			}
		})
	}
}
