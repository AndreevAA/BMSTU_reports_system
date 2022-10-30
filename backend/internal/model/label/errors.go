package label

type CanNotCreateLabelErr struct{}

func (a *CanNotCreateLabelErr) Error() string {
	return "can't create label"
}

type LabelNotFoundErr struct{}

func (a *LabelNotFoundErr) Error() string {
	return "label does not exist or does not belong to user"
}
