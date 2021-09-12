package helper

type ArrayInt []int

func (a ArrayInt) InArray(i int) bool {
    for _, v := range a {
        if v == i {
            return true
        }
    }
    return false
}
