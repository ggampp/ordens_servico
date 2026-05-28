package model

// DashboardSummary aggregates the headline indicators.
type DashboardSummary struct {
	OpenOrders       int64              `json:"open_orders"`
	InProgressOrders int64              `json:"in_progress_orders"`
	CompletedOrders  int64              `json:"completed_orders"`
	AssignedOrders   int64              `json:"assigned_orders"`
	CancelledOrders  int64              `json:"cancelled_orders"`
	ActiveEmployees  int64              `json:"active_employees"`
	OrdersByEmployee []OrdersByEmployee `json:"orders_by_employee"`
	OrdersByPriority map[string]int64   `json:"orders_by_priority"`
}

// OrdersByEmployee counts orders grouped by responsible employee.
type OrdersByEmployee struct {
	EmployeeID   *int64 `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	Count        int64  `json:"count"`
}
