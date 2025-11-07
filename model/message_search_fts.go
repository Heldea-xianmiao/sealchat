package model

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ftsVersionCurrent = 1
	ftsRebuildTimeout = 3 * time.Second
)

var (
	ftsInitialized atomic.Bool
	ftsRebuilding  atomic.Bool
	lastFTSError   atomic.Value
)

func init() {
	lastFTSError.Store("")
}

func resetSQLiteFTSState() {
	sqliteFTSReady = false
	ftsInitialized.Store(false)
	ftsRebuilding.Store(false)
	lastFTSError.Store("")
}

func setLastFTSError(err error) {
	if err == nil {
		lastFTSError.Store("")
		return
	}
	lastFTSError.Store(strings.TrimSpace(err.Error()))
}

func disableSQLiteFTS(err error) {
	setLastFTSError(err)
	sqliteFTSReady = false
}

func ReportSQLiteFTSFailure(err error) {
	if err == nil {
		return
	}
	disableSQLiteFTS(err)
}

func sqliteFTSSchemaValid(conn *gorm.DB) (bool, error) {
	if conn == nil {
		return false, errors.New("nil connection for FTS check")
	}
	if !conn.Migrator().HasTable("messages_fts") {
		return false, nil
	}
	var triggerCount int64
	if err := conn.Raw(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'trigger' AND name IN ('messages_ai','messages_ad','messages_au')`,
	).Scan(&triggerCount).Error; err != nil {
		return false, err
	}
	return triggerCount == 3, nil
}

type ftsVersionRecord struct {
	Key       string    `gorm:"primaryKey;size:64"`
	Version   int       `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Status    string    `gorm:"size:32"`
	Message   string    `gorm:"size:255"`
}

func ensureSQLiteFTSManager(conn *gorm.DB) error {
	if !ftsInitialized.CompareAndSwap(false, true) {
		return nil
	}
	if err := conn.AutoMigrate(&ftsVersionRecord{}); err != nil {
		return err
	}
	rec, err := getFTSVersion(conn)
	if err != nil {
		return err
	}
	if rec.Version >= ftsVersionCurrent && rec.Status == "ready" {
		if ok, schemaErr := sqliteFTSSchemaValid(conn); schemaErr == nil && ok {
			sqliteFTSReady = true
			return nil
		} else if schemaErr != nil {
			log.Printf("检测到消息 FTS schema 校验错误，准备重建: %v", schemaErr)
		} else {
			log.Printf("检测到消息 FTS schema 缺失，准备重建")
		}
	}
	go rebuildFTSInBackground(conn)
	return nil
}

func getFTSVersion(conn *gorm.DB) (*ftsVersionRecord, error) {
	rec := &ftsVersionRecord{}
	err := conn.Where("key = ?", "messages_fts").
		First(rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rec.Key = "messages_fts"
		rec.Version = 0
		rec.Status = "unknown"
		return rec, nil
	}
	return rec, err
}

func rebuildFTSInBackground(conn *gorm.DB) {
	if !ftsRebuilding.CompareAndSwap(false, true) {
		return
	}
	defer ftsRebuilding.Store(false)

	start := time.Now()
	if err := markFTSStatus(conn, "building", ""); err != nil {
		log.Printf("记录 FTS 状态失败: %v", err)
	}

	if err := rebuildFTS(conn); err != nil {
		setLastFTSError(err)
		_ = markFTSStatus(conn, "error", err.Error())
		log.Printf("重建 FTS 失败: %v", err)
		return
	}
	duration := time.Since(start)
	setLastFTSError(nil)
	sqliteFTSReady = true
	_ = markFTSStatus(conn, "ready", fmt.Sprintf("rebuilt in %s", duration))
	log.Printf("消息 FTS 重建完成，用时 %s", duration)
}

func rebuildFTS(conn *gorm.DB) error {
	statements := []string{
		`DROP TRIGGER IF EXISTS messages_ai;`,
		`DROP TRIGGER IF EXISTS messages_ad;`,
		`DROP TRIGGER IF EXISTS messages_au;`,
		`DROP TABLE IF EXISTS messages_fts;`,
		`CREATE VIRTUAL TABLE messages_fts USING fts5(
			message_id UNINDEXED,
			content,
			tokenize = 'unicode61 remove_diacritics 0'
		);`,
		`CREATE TRIGGER messages_ai AFTER INSERT ON messages BEGIN
			INSERT INTO messages_fts(message_id, content) VALUES (new.id, COALESCE(new.content, ''));
		END;`,
		`CREATE TRIGGER messages_ad AFTER DELETE ON messages BEGIN
			DELETE FROM messages_fts WHERE message_id = old.id;
		END;`,
		`CREATE TRIGGER messages_au AFTER UPDATE ON messages BEGIN
			INSERT OR REPLACE INTO messages_fts(message_id, content) VALUES (new.id, COALESCE(new.content, ''));
		END;`,
	}
	for _, stmt := range statements {
		if err := conn.Exec(stmt).Error; err != nil {
			return err
		}
	}
	if err := conn.Exec(`
		INSERT INTO messages_fts(message_id, content)
		SELECT id, COALESCE(content, '')
		FROM messages;
	`).Error; err != nil {
		return err
	}
	return nil
}

func markFTSStatus(conn *gorm.DB, status, message string) error {
	record := ftsVersionRecord{
		Key: "messages_fts",
		Version: func() int {
			if status == "ready" {
				return ftsVersionCurrent
			}
			return 0
		}(),
		Status:    status,
		Message:   strings.TrimSpace(message),
		UpdatedAt: time.Now(),
	}
	return conn.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&record).Error
}

func LastFTSError() string {
	val := lastFTSError.Load()
	if val == nil {
		return ""
	}
	if msg, ok := val.(string); ok {
		return msg
	}
	return fmt.Sprintf("%v", val)
}
