package sqlbuilder

type Locker interface {
	Build() (string, []interface{})
}

type ShareLocker struct{}

func (locker *ShareLocker) Build() (string, []interface{}) {
	return "lock in share mode", nil
}

type UpdateLocker struct{}

func (locker *UpdateLocker) Build() (string, []interface{}) {
	return "for update", nil
}
