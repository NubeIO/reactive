package reactive

func (n *BaseNode) setOptions(opts *Options) {
	if opts != nil {
		n.options = opts
		n.setMeta(opts)
	}
}
