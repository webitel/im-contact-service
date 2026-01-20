package utils

// Map applies the given function to each element of the given slice and returns a new slice with the results.
// The function f is applied to each element of the slice s and the results are collected in a new slice of type U.
// The resulting slice is of the same length as the original slice s.
func Map[T, U any](s []T, f func(T) U) []U {
    res := make([]U, len(s))
    for i, v := range s {
        res[i] = f(v)
    }
    return res
}