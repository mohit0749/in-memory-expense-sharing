package splitwise

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testsplitwise/expense"
	"testsplitwise/user"
)

type Splitwise struct {
	lock     *sync.Mutex
	expenses map[int64]*expense.Expense
	users    map[int64]*user.User
}

func (sw *Splitwise) CreateUser(name, email, contact string) int64 {
	u := user.NewUser(name, email, contact)
	sw.lock.Lock()
	defer sw.lock.Unlock()
	sw.users[u.GetId()] = &u
	return u.GetId()
}

func NewSplitwise() Splitwise {
	return Splitwise{&sync.Mutex{}, make(map[int64]*expense.Expense, 0), make(map[int64]*user.User, 0)}
}

func (sw *Splitwise) CreateExpense(userId int64, forUser []int64, amount []float64, _type string) error {
	var exps []expense.Expense
	var err error
	if len(forUser) == 0 || len(amount) == 0 {
		return errors.New("invalid users/amount")
	}
	switch strings.ToLower(_type) {
	case "equal":
		exps, err = sw.splitEqually(userId, forUser, amount[0])
	case "exact":
		exps, err = sw.splitExact(userId, forUser, amount)
	}
	if err != nil {
		return err
	}
	for _, exp := range exps {
		sw.lock.Lock()
		sw.expenses[exp.GetId()] = &exp
		sw.lock.Unlock()
	}
	return nil
}

func (sw *Splitwise) ShowAllExpenses() {
	for k, _ := range sw.users {
		sw.ShowExpense(k)
	}
}

func (sw Splitwise) ShowExpense(userId int64) {
	user, err := sw.getUser(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	owes := user.GetOwes()
	owesMp := make(map[int64]float64)
	for _, owe := range owes {
		exp, err := sw.getExpense(owe)
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
	lendMp := make(map[int64]float64)
	for _, l := range lends {
		exp, err := sw.getExpense(l)
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

func (sw *Splitwise) getExpense(id int64) (*expense.Expense, error) {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	exp, ok := sw.expenses[id]
	if !ok {
		return nil, errors.New("exp does not exists")
	}
	return exp, nil
}

func (sw *Splitwise) getUser(id int64) (*user.User, error) {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	user, ok := sw.users[id]
	if !ok {
		return nil, errors.New("user does not exists")
	}
	return user, nil
}

func (sw Splitwise) splitEqually(paidBy int64, forUsers []int64, amount float64) ([]expense.Expense, error) {
	exps := make([]expense.Expense, 0)
	if len(forUsers) == 0 {
		return exps, errors.New("zero users mentioned")
	}
	equalAmount := amount / float64(len(forUsers))
	payer, err := sw.getUser(paidBy)
	if err != nil {
		return nil, fmt.Errorf("user %d does not exists", paidBy)
	}
	for i := 0; i < len(forUsers); i++ {
		exp := expense.NewExpense(paidBy, forUsers[i], equalAmount, expense.EXACT)
		user, err := sw.getUser(forUsers[i])
		if err != nil {
			return nil, fmt.Errorf("user %d does not exists", forUsers[i])
		}
		user.Owe(exp.GetId())
		payer.Lend(exp.GetId())
		fmt.Println(exp)
		exps = append(exps, exp)
	}
	fmt.Println(exps)
	return exps, nil
}

func (sw Splitwise) splitExact(paidBy int64, forUsers []int64, amount []float64) ([]expense.Expense, error) {
	exps := make([]expense.Expense, 0)
	if len(forUsers) != len(amount) {
		return exps, errors.New("invalid amount/users")
	}
	payer, err := sw.getUser(paidBy)
	if err != nil {
		return nil, fmt.Errorf("user %d does not exists", paidBy)
	}
	for i := 0; i < len(forUsers); i++ {
		exp := expense.NewExpense(paidBy, forUsers[i], amount[i], expense.EXACT)
		user, err := sw.getUser(forUsers[i])
		if err != nil {
			return nil, fmt.Errorf("user %d does not exists", forUsers[i])
		}
		user.Owe(exp.GetId())
		payer.Lend(exp.GetId())
		exps = append(exps, exp)
	}
	return exps, nil
}
