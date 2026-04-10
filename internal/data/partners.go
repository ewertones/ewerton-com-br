package data

type Partner struct {
	Slug           string
	Name           string
	Logo           string
	URL            string
	DisplayURL     string
	HasDescription bool
}

type PartnerGroup struct {
	Slug     string
	Partners []Partner
}

var PartnerGroups = []PartnerGroup{
	{
		Slug: "financial",
		Partners: []Partner{
			{Slug: "dashfi", Name: "DashFi", Logo: "/static/images/partners/dashfi.webp", URL: "https://dash.fi", DisplayURL: "dash.fi", HasDescription: true},
			{Slug: "cerc", Name: "CERC", Logo: "/static/images/partners/cerc.webp", URL: "https://cerc.com", DisplayURL: "cerc.com", HasDescription: true},
		},
	},
	{
		Slug: "marketing",
		Partners: []Partner{
			{Slug: "stagwell", Name: "Stagwell", Logo: "/static/images/partners/stagwell.webp", URL: "https://stagwellglobal.com", DisplayURL: "stagwellglobal.com", HasDescription: true},
			{Slug: "begrowth", Name: "BeGrowth", Logo: "/static/images/partners/begrowth.webp", URL: "https://begrowth.com.br", DisplayURL: "begrowth.com.br", HasDescription: true},
		},
	},
	{
		Slug: "tax",
		Partners: []Partner{
			{Slug: "ey", Name: "EY", Logo: "/static/images/partners/ey.webp", URL: "https://ey.com", DisplayURL: "ey.com"},
		},
	},
	{
		Slug: "food",
		Partners: []Partner{
			{Slug: "morsum", Name: "Morsum", Logo: "/static/images/partners/morsum.webp", URL: "https://www.dotfoods.com", DisplayURL: "dotfoods.com", HasDescription: true},
		},
	},
	{
		Slug: "outsourcing",
		Partners: []Partner{
			{Slug: "softensity", Name: "Softensity", Logo: "/static/images/partners/softensity.webp", URL: "https://softensity.com", DisplayURL: "softensity.com"},
			{Slug: "truelogic", Name: "Truelogic", Logo: "/static/images/partners/truelogic.webp", URL: "https://truelogic.io", DisplayURL: "truelogic.io"},
			{Slug: "blueshift", Name: "BlueShift", Logo: "/static/images/partners/blueshift.webp", URL: "https://blueshift.com.br", DisplayURL: "blueshift.com.br"},
		},
	},
}
