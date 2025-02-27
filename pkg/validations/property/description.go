package property

type Description struct {
}

func (d *Description) Name() string {
	return "description"
}

func (d *Description) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		newProperty.Description = ""
		oldProperty.Description = ""
		return NewDiff(oldProperty, newProperty)
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, nil
}
