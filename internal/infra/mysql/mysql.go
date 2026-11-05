package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

func New() (*sql.DB, error) {
	loc, err := time.LoadLocation(time.Local.String())

	if err != nil {
		return nil, err
	}

	cfg := mysql.Config{
		User:   "root",
		Passwd: "1234",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "example",
		// MySQL 8.0 미만 구형 인증 방식 호환이 필요한 경우 사용
		AllowNativePasswords: true,
		// MySQL의 DATETIME을 Go의 time.Time으로 자동 변환
		ParseTime: true,
		// 타임존 세팅 (time.Local 또는 time.UTC)
		Loc: loc,
		// 이모지 지원을 위한 멀티바이트 문자셋
		Collation: "utf8mb4_0900_ai_ci",
		// 서버 초기 연결 제한 시간
		Timeout: 5 * time.Second,
		// 쿼리 읽기 제한 시간 (Slow 쿼리 방어)
		ReadTimeout: 10 * time.Second,
		// 쿼리 쓰기 제한 시간
		WriteTimeout: 10 * time.Second,
		// Go가 SQL을 미리 조립해서 단 1번만 통신함 -> 3번에서 1번으로 줄어듦
		InterpolateParams: true,
	}

	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)

	// 동시에 열 수 있는 최대 커넥션 수
	db.SetMaxOpenConns(10)
	// 풀(Pool)에 유휴 상태로 대기시킬 최대 커넥션 수 (MaxOpen과 동일하게 맞추는 게 정석)
	db.SetMaxIdleConns(10)
	// 커넥션 재사용 최대 재사용 시간 (MySQL의 wait_timeout보다 무조건 짧아야 함)
	db.SetConnMaxLifetime(5 * time.Minute)
	// 유휴 커넥션이 풀에 머물 수 있는 최대 시간
	db.SetConnMaxIdleTime(3 * time.Minute)

	ctx, ctxCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer ctxCancel()

	// 여기서 처음으로 DB 서버로 TCP 연결을 시도
	conn, err := db.Conn(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "SELECT 1"); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
