package oxr

type usageParams struct {
	prettyPrint bool
}

// UsageOption allows the client to specify values for a usage request.
type UsageOption func(params *usageParams)

// UsageWithPrettyPrint sets whether to minify the response.
func UsageWithPrettyPrint(active bool) UsageOption {
	return func(p *usageParams) {
		p.prettyPrint = active
	}
}
