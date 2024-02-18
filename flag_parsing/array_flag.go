package flag_parsing

type ScopeFlags []string

func (i *ScopeFlags) String() string {
	// Return openid as default value
	if len(*i) == 0 {
		return "openid"
	}
	// Start with nothing
	result := ""
	for _, content := range *i {
		result += content + "+"
	}
	return result[:len(result)-1]
}

func (i *ScopeFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
