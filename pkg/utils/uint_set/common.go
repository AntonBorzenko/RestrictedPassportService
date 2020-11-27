package uint_set

type UintSet interface {
	Has(number uint64) (bool, error)
	Insert(number uint64) error
	InsertMultiple(numbers chan uint64, ignoreErrors bool) error
}

func insertMultiple(set UintSet, numbers chan uint64, ignoreErrors bool) error {
	for number := range numbers {
		err := set.Insert(number)
		if !ignoreErrors && err != nil {
			return err
		}
	}
	return nil
}
