package snowflake

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateSnowflakeIDs(t *testing.T) {
	gen, err := NewGenerator(0)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n=== Пример генерации Snowflake ID ===")
	fmt.Println("nodeID: 0")
	fmt.Println("Эпоха: 2024-01-01 00:00:00 UTC")
	fmt.Println()

	var prevID int64
	for i := 0; i < 10; i++ {
		id, err := gen.Next()
		if err != nil {
			t.Fatal(err)
		}

		ts := ParseTime(id)
		nodeID := ParseNodeID(id)
		seq := ParseSequence(id)

		fmt.Printf("ID #%d:\n", i+1)
		fmt.Printf("  Decimal:    %d\n", id)
		fmt.Printf("  Hex:        0x%X\n", id)
		fmt.Printf("  Time:       %s\n", ts.Format("2006-01-02 15:04:05.000"))
		fmt.Printf("  NodeID:     %d\n", nodeID)
		fmt.Printf("  Sequence:   %d\n", seq)

		if prevID != 0 {
			fmt.Printf("  Diff:       +%d\n", id-prevID)
		}
		fmt.Println()

		prevID = id

		// Небольшая задержка для изменения timestamp
		if i == 5 {
			time.Sleep(2 * time.Millisecond)
		}
	}

	// Проверка уникальности
	fmt.Println("=== Проверка уникальности (10000 ID) ===")
	ids := make(map[int64]bool)
	for i := 0; i < 10000; i++ {
		id, err := gen.Next()
		if err != nil {
			t.Fatal(err)
		}
		if ids[id] {
			t.Fatalf("Duplicate ID found: %d", id)
		}
		ids[id] = true
	}
	fmt.Printf("Сгенерировано 10000 уникальных ID ✓\n")
	fmt.Printf("Первый: %d\n", getFirstKey(ids))
	fmt.Printf("Последний: %d\n", getFirstKey(ids)+9999)
}

func getFirstKey(m map[int64]bool) int64 {
	var first int64
	for k := range m {
		if first == 0 || k < first {
			first = k
		}
	}
	return first
}
