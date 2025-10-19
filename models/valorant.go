package models

import (
	"time"

	"gorm.io/gorm"
)

// ValorantPlayer представляет игрока в Valorant
type ValorantPlayer struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	GameName   string         `json:"game_name"`   // Игровое имя
	Tag        string         `json:"tag"`         // Тег игрока
	Region     string         `json:"region"`      // Регион (eu, na, ap, etc.)
	Rank       string         `json:"rank"`        // Текущий ранг
	RankRating int            `json:"rank_rating"` // Рейтинг ранга
	PeakRank   string         `json:"peak_rank"`   // Пиковый ранг
	PeakRating int            `json:"peak_rating"` // Пиковый рейтинг
	Level      int            `json:"level"`       // Уровень аккаунта
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ValorantMatch представляет матч в Valorant
type ValorantMatch struct {
	ID        uint                  `json:"id" gorm:"primaryKey"`
	MatchID   string                `json:"match_id" gorm:"uniqueIndex"` // ID матча от Riot
	Map       string                `json:"map"`                         // Карта
	Mode      string                `json:"mode"`                        // Режим игры
	Result    string                `json:"result"`                      // Результат (win/loss)
	Score     string                `json:"score"`                       // Счет
	Duration  int                   `json:"duration"`                    // Длительность в секундах
	Date      time.Time             `json:"date"`                        // Дата матча
	TeamID    *uint                 `json:"team_id"`
	Team      *Team                 `json:"team" gorm:"foreignKey:TeamID"`
	Players   []ValorantPlayerMatch `json:"players" gorm:"foreignKey:MatchID"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt gorm.DeletedAt        `json:"deleted_at" gorm:"index"`
}

// ValorantPlayerMatch связывает игрока с матчем
type ValorantPlayerMatch struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MatchID     uint           `json:"match_id"`
	Match       ValorantMatch  `json:"match" gorm:"foreignKey:MatchID"`
	PlayerID    uint           `json:"player_id"`
	Player      ValorantPlayer `json:"player" gorm:"foreignKey:PlayerID"`
	Agent       string         `json:"agent"`        // Агент
	Kills       int            `json:"kills"`        // Убийства
	Deaths      int            `json:"deaths"`       // Смерти
	Assists     int            `json:"assists"`      // Помощи
	Score       int            `json:"score"`        // Очки
	Damage      int            `json:"damage"`       // Урон
	Headshots   int            `json:"headshots"`    // Хедшоты
	FirstKills  int            `json:"first_kills"`  // Первые убийства
	FirstDeaths int            `json:"first_deaths"` // Первые смерти
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ValorantStats статистика игрока
type ValorantStats struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	PlayerID       uint           `json:"player_id"`
	Player         ValorantPlayer `json:"player" gorm:"foreignKey:PlayerID"`
	TotalMatches   int            `json:"total_matches"`   // Всего матчей
	Wins           int            `json:"wins"`            // Победы
	Losses         int            `json:"losses"`          // Поражения
	WinRate        float64        `json:"win_rate"`        // Процент побед
	AverageScore   float64        `json:"average_score"`   // Средний счет
	AverageKills   float64        `json:"average_kills"`   // Средние убийства
	AverageDeaths  float64        `json:"average_deaths"`  // Средние смерти
	AverageAssists float64        `json:"average_assists"` // Средние помощи
	HeadshotRate   float64        `json:"headshot_rate"`   // Процент хедшотов
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
