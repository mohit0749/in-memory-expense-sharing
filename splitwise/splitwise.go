package splitwise

import (
	"errors"
	"fmt"
	"testsplitwise/expense/storage/expensestore"
	"testsplitwise/user"
	"testsplitwise/user/storage/userstore"

	"testsplitwise/expense"
	expensestorage "testsplitwise/expense/storage"
	userstorage "testsplitwise/user/storage"
)

type Splitwise struct {
	expenses expensestorage.ExpenseStore
	users    userstorage.UserStore
}

func NewSplitwise() Splitwise {
	expensestorage.InitExpenseStore(expensestore.NewExpenseStore())
	userstorage.InitUserStore(userstore.NewUserStore())
	return Splitwise{expensestorage.GetExpenseStore(), userstorage.GetUserStore()}
}

func (sw *Splitwise) CreateUser(name, email, contact string) uint64 {
	u := user.NewUser(name, email, contact)
	sw.users.AddUser(&u)
	return u.GetId()
}

func (sw *Splitwise) CreateExpense(userId uint64, forUser []uint64, amount []float64, expType expense.ExpenseType) error {
	var exps []*expense.Expense
	var err error
	if len(forUser) == 0 || len(amount) == 0 {
		return errors.New("invalid users/amount")
	}
	switch expType {
	case expense.EQUAL:
		exps, err = sw.splitEqually(userId, forUser, amount[0])
	case expense.EXACT:
		exps, err = sw.splitExact(userId, forUser, amount)
	}
	if err != nil {
		return err
	}
	for _, exp := range exps {
		sw.expenses.AddExpense(exp)
	}
	return nil
}

func (sw Splitwise) ShowExpense(userId uint64) {
	user, err := sw.users.GetUser(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	owes := user.GetOwes()
	owesMp := make(map[uint64]float64)
	for _, owe := range owes {
		exp, err := sw.expenses.GetExpense(owe)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if exp.GetPaidBy() == userId {
			continue
		}
		owesMp[exp.GetPaidBy()] += exp.GetAmount()
	}
	for k, v := range owesMp {
		fmt.Printf("%d owes %d: %.2f\n", userId, k, v)
	}
	lends := user.GetLend()
	lendMp := make(map[uint64]float64)
	for _, l := range lends {
		exp, err := sw.expenses.GetExpense(l)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if userId == exp.PaidForUser() {
			continue
		}
		lendMp[exp.PaidForUser()] += exp.GetAmount()
	}

	for k, v := range lendMp {
		fmt.Printf("%d owes %d: %.2f\n", k, userId, v)
	}
}

func (sw Splitwise) splitEqually(paidBy uint64, forUsers []uint64, amount float64) ([]*expense.Expense, error) {
	exps := make([]*expense.Expense, 0)
	if len(forUsers) == 0 {
		return exps, errors.New("zero users mentioned")
	}
	equalAmount := amount / float64(len(forUsers))
	payer, err := sw.users.GetUser(paidBy)
	if err != nil {
		return nil, fmt.Errorf("user %d does not exists", paidBy)
	}
	for i := 0; i < len(forUsers); i++ {
		exp := expense.NewExpense(paidBy, forUsers[i], equalAmount, expense.EQUAL)
		user, err := sw.users.GetUser(forUsers[i])
		if err != nil {
			return nil, fmt.Errorf("user %d does not exists", forUsers[i])
		}
		user.Owe(exp.GetId())
		payer.Lend(exp.GetId())
		exps = append(exps, &exp)
	}
	return exps, nil
}

func (sw Splitwise) splitExact(paidBy uint64, forUsers []uint64, amount []float64) ([]*expense.Expense, error) {
	exps := make([]*expense.Expense, 0)
	if len(forUsers) != len(amount) {
		return exps, errors.New("invalid amount/users")
	}
	payer, err := sw.users.GetUser(paidBy)
	if err != nil {
		return nil, fmt.Errorf("user %d does not exists", paidBy)
	}
	for i := 0; i < len(forUsers); i++ {
		exp := expense.NewExpense(paidBy, forUsers[i], amount[i], expense.EXACT)
		user, err := sw.users.GetUser(forUsers[i])
		if err != nil {
			return nil, fmt.Errorf("user %d does not exists", forUsers[i])
		}
		user.Owe(exp.GetId())
		payer.Lend(exp.GetId())
		exps = append(exps, &exp)
	}
	return exps, nil
}
