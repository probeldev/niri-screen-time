package cache

import (
	"log"
	"sync"
	"time"

	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
)

// ScreenTimeCache - буфер между сбором данных и их сохранением в БД
type ScreenTimeCache struct {
	db          *db.ScreenTimeDB
	buffer      []model.ScreenTime
	bufferMutex sync.Mutex
	flushPeriod time.Duration
	maxBuffer   int
	stopChan    chan struct{}
}

// NewScreenTimeCache создает новый кэш
func NewScreenTimeCache(db *db.ScreenTimeDB, flushPeriod time.Duration, maxBuffer int) *ScreenTimeCache {
	return &ScreenTimeCache{
		db:          db,
		buffer:      make([]model.ScreenTime, 0, maxBuffer),
		flushPeriod: flushPeriod,
		maxBuffer:   maxBuffer,
		stopChan:    make(chan struct{}),
	}
}

// Start запускает фоновую горутину для периодического сброса буфера
func (stc *ScreenTimeCache) Start() {
	go stc.flushWorker()
}

// Stop останавливает фоновую горутину и сбрасывает оставшиеся данные
func (stc *ScreenTimeCache) Stop() {
	close(stc.stopChan)
	stc.flushBuffer() // Сброс оставшихся данных
}

// Add добавляет запись в буфер
func (stc *ScreenTimeCache) Add(st model.ScreenTime) {
	stc.bufferMutex.Lock()
	defer stc.bufferMutex.Unlock()

	stc.buffer = append(stc.buffer, st)

	// Если буфер заполнен, сбрасываем его
	if len(stc.buffer) >= stc.maxBuffer {
		stc.flushBuffer()
	}
}

// flushWorker периодически сбрасывает буфер в БД
func (stc *ScreenTimeCache) flushWorker() {
	ticker := time.NewTicker(stc.flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stc.flushBuffer()
		case <-stc.stopChan:
			return
		}
	}
}

// flushBuffer сбрасывает содержимое буфера в БД
func (stc *ScreenTimeCache) flushBuffer() {
	stc.bufferMutex.Lock()
	defer stc.bufferMutex.Unlock()

	if len(stc.buffer) == 0 {
		return
	}

	// Копируем буфер, чтобы минимизировать время блокировки
	records := make([]model.ScreenTime, len(stc.buffer))
	copy(records, stc.buffer)

	// Очищаем буфер
	stc.buffer = stc.buffer[:0]

	// Сохраняем в БД в отдельной горутине, чтобы не блокировать основной поток
	go func() {
		if err := stc.db.BulkInsert(records); err != nil {
			log.Printf("Failed to bulk insert records: %v", err)
			// При ошибке можно добавить записи обратно в буфер
			stc.bufferMutex.Lock()
			stc.buffer = append(stc.buffer, records...)
			stc.bufferMutex.Unlock()
		}
	}()
}
