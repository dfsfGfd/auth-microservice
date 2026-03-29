// Package snowflake предоставляет генератор Snowflake ID.
//
// Snowflake — это распределённый генератор уникальных ID,
// создающий 64-битные целые числа, сортируемые по времени.
//
// Формат ID (64 бита):
//   - 41 бит: миллисекунды с кастомной эпохи
//   - 10 бит: node ID (0-1023)
//   - 12 бит: sequence number (0-4095)
//
// Для одного инстанса используется nodeID = 0.
package snowflake

import (
	"errors"
	"sync"
	"time"
)

// Ошибки генератора.
var (
	ErrInvalidNodeID = errors.New("nodeID must be between 0 and 1023")
	ErrClockMoved    = errors.New("clock moved backwards")
)

// Константы временной эпохи.
const (
	// Epoch — кастомная эпоха (2024-01-01 00:00:00 UTC).
	Epoch int64 = 1704067200000

	// NodeBits — количество бит для node ID.
	NodeBits = 10

	// SeqBits — количество бит для sequence number.
	SeqBits = 12

	// MaxNodeID — максимальный node ID (2^10 - 1).
	MaxNodeID = 1023

	// MaxSeq — максимальный sequence number (2^12 - 1).
	MaxSeq = 4095
)

// Generator генерирует уникальные Snowflake ID.
// Потокобезопасен.
type Generator struct {
	nodeID    int64
	epoch     int64
	mu        sync.Mutex
	lastTime  int64
	sequence  int64
}

// NewGenerator создаёт новый генератор с указанным nodeID.
// nodeID должен быть в диапазоне [0, 1023].
func NewGenerator(nodeID int64) (*Generator, error) {
	if nodeID < 0 || nodeID > MaxNodeID {
		return nil, ErrInvalidNodeID
	}

	return &Generator{
		nodeID:   nodeID,
		epoch:    Epoch,
		lastTime: -1,
	}, nil
}

// Next генерирует следующий уникальный ID.
// Возвращает ошибку, если системное время ушло назад.
func (g *Generator) Next() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	ts := g.timestamp()

	// Проверка: время не должно идти назад
	if ts < g.lastTime {
		return 0, ErrClockMoved
	}

	// Сброс sequence при переходе на новую миллисекунду
	if ts > g.lastTime {
		g.sequence = 0
	}

	// Инкремент sequence с wrap-around
	g.sequence = (g.sequence + 1) & MaxSeq
	if g.sequence == 0 && ts == g.lastTime {
		// Лимит исчерпан, ждём следующую миллисекунду
		for ts <= g.lastTime {
			ts = g.timestamp()
		}
		g.sequence = 0
	}

	g.lastTime = ts

	// Формирование ID:
	// [41 бит timestamp][10 бит nodeID][12 бит sequence]
	id := (ts-g.epoch)<<22 | (g.nodeID << 12) | g.sequence
	return id, nil
}

// timestamp возвращает текущее время в миллисекундах.
func (g *Generator) timestamp() int64 {
	return time.Now().UnixMilli()
}

// ParseTime извлекает время создания из ID.
func ParseTime(id int64) time.Time {
	ts := (id >> 22) + Epoch
	return time.UnixMilli(ts)
}

// ParseNodeID извлекает nodeID из ID.
func ParseNodeID(id int64) int64 {
	return (id >> 12) & MaxNodeID
}

// ParseSequence извлекает sequence number из ID.
func ParseSequence(id int64) int64 {
	return id & MaxSeq
}
