package taxcalculator

var beijing map[int64]*TaxCfg = map[int64]*TaxCfg{
	2025: &TaxCfg{
		City:              "北京",
		CityEn:            "beijing",
		Year:              2025,
		ActiveAt:          7,
		UpperLimit:        35283,
		LowerLimit:        6821,
		Insurances:        beijingDefaultIns,
		IncomeTax:         beijingDefaultIncomeTax,
		Threshold:         5000,
		SpecialDeductions: beijingDefaultSDS,
	},
	2024: &TaxCfg{
		City:              "北京",
		CityEn:            "beijing",
		Year:              2024,
		ActiveAt:          7,
		UpperLimit:        35283,
		LowerLimit:        6821,
		Insurances:        beijingDefaultIns,
		IncomeTax:         beijingDefaultIncomeTax,
		Threshold:         5000,
		SpecialDeductions: beijingDefaultSDS,
	},
	2023: &TaxCfg{
		City:              "北京",
		CityEn:            "beijing",
		Year:              2023,
		ActiveAt:          7,
		UpperLimit:        33891,
		LowerLimit:        6326,
		Insurances:        beijingDefaultIns,
		IncomeTax:         beijingDefaultIncomeTax,
		Threshold:         5000,
		SpecialDeductions: beijingDefaultSDS,
	},
}

var beijingDefaultIns = Ins{
	Med:          Insurance{Company: 0.1, Individual: 0.02},
	Retirement:   Insurance{Company: 0.16, Individual: 0.08},
	Unemployment: Insurance{Company: 0.008, Individual: 0.002}, // ?
	Injury:       Insurance{Company: 0.019, Individual: 0},
	Birth:        Insurance{Company: 0.008, Individual: 0},
	HouseFund:    Insurance{Company: 0.12, Individual: 0.12},
}

var beijingDefaultIncomeTax = []IncomeTax{
	IncomeTax{Left: 0, Right: 36000, Rate: 0.03, FastFix: 0},
	IncomeTax{Left: 36000, Right: 144000, Rate: 0.1, FastFix: 2520},
	IncomeTax{Left: 144000, Right: 300000, Rate: 0.2, FastFix: 16920},
	IncomeTax{Left: 300000, Right: 420000, Rate: 0.25, FastFix: 31920},
	IncomeTax{Left: 420000, Right: 660000, Rate: 0.30, FastFix: 52920},
	IncomeTax{Left: 660000, Right: 960000, Rate: 0.35, FastFix: 85920},
	IncomeTax{Left: 960000, Right: 9999999999, Rate: 0.45, FastFix: 181920},
}

var beijingDefaultSDS = SDs{
	Child:       SpecialDeduction{Enable: false, Amount: 2000, Des: "子女教育"},
	Adult:       SpecialDeduction{Enable: false, Amount: 400, Des: "成人继续教育"},
	LifeOrDeath: SpecialDeduction{Enable: false, Amount: 0, Des: "大病医疗"},
	House:       SpecialDeduction{Enable: false, Amount: 1000, Des: "住房贷款"},
	Ranting:     SpecialDeduction{Enable: false, Amount: 1500, Des: "租房"},
	Elders:      SpecialDeduction{Enable: false, Amount: 3000, Des: "赡养老人"},
	Infant:      SpecialDeduction{Enable: false, Amount: 2000, Des: "3岁以下婴幼儿照顾"},
}
