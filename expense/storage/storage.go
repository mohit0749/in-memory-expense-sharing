package storage

import (
	"sync"
	"testsplitwise/expense"
)

type ExpenseStore interface {
	AddExpense(*expense.Expense) error
	GetExpense(uint64) (*expense.Expense, error)
}

var defaultStore ExpenseStore
var singleton sync.Once

func InitExpenseStore(store ExpenseStore) {
	singleton.Do(func() {
		defaultStore = store
	})
}

func GetExpenseStore() ExpenseStore {
	return defaultStore
}
