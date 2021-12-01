package expensestore

import (
	"fmt"
	"sync"
	"testsplitwise/expense"
)

type expenseStore struct {
	expenses sync.Map
}

func NewExpenseStore() *expenseStore {
	return &expenseStore{sync.Map{}}
}

func (es *expenseStore) AddExpense(exp *expense.Expense) error {
	es.expenses.Store(exp.GetId(), exp)
	return nil
}

func (es *expenseStore) GetExpense(id uint64) (*expense.Expense, error) {
	value, ok := es.expenses.Load(id)
	if !ok {
		return nil, fmt.Errorf("expense %d does not exists", id)
	}
	exp, _ := value.(*expense.Expense)
	return exp, nil
}
