package expense

import "sync/atomic"

var seqNo uint64

type ExpenseType int

var (
	EXACT ExpenseType = 1
	EQUAL ExpenseType = 2
)

type Expense struct {
	id      uint64
	paidBy  uint64
	forUser uint64
	amount  float64
	_type   ExpenseType
}

func (e Expense) GetId() uint64 {
	return e.id
}

func NewExpense(paidBy, forUser uint64, amount float64, expType ExpenseType) Expense {
	atomic.AddUint64(&seqNo, 1)
	return Expense{seqNo, paidBy, forUser, amount, expType}
}

func (e *Expense) GetPaidBy() uint64 {
	return e.paidBy
}

func (e *Expense) PaidForUser() uint64 {
	return e.forUser
}

func (e *Expense) GetAmount() float64 {
	return e.amount
}
