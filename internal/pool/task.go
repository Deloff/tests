package pool

type Task struct {
	ID      int
	Created string // время создания
	Handled string // время выполнения
	Result  []byte
}
