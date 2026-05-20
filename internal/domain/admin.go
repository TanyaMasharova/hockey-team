// internal/domain/admin.go
package domain

// internal/domain/admin.go
type SalesByMonth struct {
	Month   string  `db:"month"`
	Tickets int     `db:"tickets"`
	Revenue float64 `db:"revenue"` // float64, не int
}
type SectorPopularity struct {
	Sector           string  `db:"sector"`
	Sold             int     `db:"sold"`
	Capacity         int     `db:"capacity"`
	OccupancyPercent float64 `db:"occupancy_percent"`
}

type TicketStatus struct {
	Status string `db:"status"`
	Count  int    `db:"count"`
}

type AvgPriceBySector struct {
	Sector      string  `db:"sector"`
	SectorType  string  `db:"sector_type"`
	AvgPrice    float64 `db:"avg_price"` // ← float64
	TicketsSold int     `db:"tickets_sold"`
}

type TopBuyer struct {
	FullName     string  `db:"full_name"`
	Email        string  `db:"email"`
	TicketsCount int     `db:"tickets_count"`
	TotalSpent   float64 `db:"total_spent"` // ← float64
}

type AdminStats struct {
	TotalUsers    int     `db:"total_users"`
	TotalTickets  int     `db:"total_tickets"`
	TotalRevenue  float64 `db:"total_revenue"` // ← float64
	ActiveTickets int     `db:"active_tickets"`
	UsedTickets   int     `db:"used_tickets"`
}
