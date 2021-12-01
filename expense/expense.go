package expense

import "sync/atomic"

var seqNo int64

type ExpenseType int

var (
	EXACT ExpenseType = 1
	EQUAL ExpenseType = 2
)

type Expense struct {
	id      int64
	paidBy  int64
	forUser int64
	amount  float64
	_type   ExpenseType
}

func (e Expense) GetId() int64 {
	return e.id
}

func NewExpense(paidBy, forUser int64, amount float64, expType ExpenseType) Expense {
	atomic.AddInt64(&seqNo, 1)
	return Expense{seqNo, paidBy, forUser, amount, expType}
}

func (e *Expense) GetPaidBy() int64 {
	return e.paidBy
}

func (e *Expense) PaidForUser() int64 {
	return e.forUser
}

func (e *Expense) GetAmount() float64 {
	return e.amount
}
