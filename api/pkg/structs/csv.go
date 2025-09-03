package structs

type OptionsWithSelectedCount struct {
	Option        string `csv:"option"`
	SelectedCount int    `csv:"selected_count"`
}
