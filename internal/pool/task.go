package pool

type Task struct {
	Id      int
	Created string // время создания
	Handled string // время выполнения
	Result  []byte
}
