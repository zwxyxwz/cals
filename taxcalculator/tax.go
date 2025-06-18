package taxcalculator

type Insurance struct {
	Company    float64 `json:"company"`
	Individual float64 `json:"individual"`
}
type Ins struct {
	Med          Insurance `json:"med"`
	Retirement   Insurance `json:"retirement"`
	Unemployment Insurance `json:"unemployment"`
	Injury       Insurance `json:"injury"`
	Birth        Insurance `json:"birth"`
	HouseFund    Insurance `json:"house_fund"`
}

type IncomeTax struct {
	Left    float64 `json:"left"`
	Right   float64 `json:"right"`
	Rate    float64 `json:"rate"`
	FastFix float64 `json:"fast_fix"`
}

type SpecialDeduction struct {
	Enable bool    `json:"enable"`
	Amount float64 `json:"amount"`
	Des    string  `json:"des"`
}

type SDs struct {
	Child       SpecialDeduction `json:"child"`
	Adult       SpecialDeduction `json:"adult"`
	LifeOrDeath SpecialDeduction `json:"life_or_death"`
	House       SpecialDeduction `json:"house"`
	Ranting     SpecialDeduction `json:"ranting"`
	Elders      SpecialDeduction `json:"elders"`
	Infant      SpecialDeduction `json:"infant"`
}

var (
	ChildTag = 1 << 1
	AdultTag = 1 << 2
	// LifeOrDeathTag = 1 << 3 // not able to cal
	HouseTag   = 1 << 4
	RantingTag = 1 << 5
	EldersTag  = 1 << 6
	InfantTag  = 1 << 7

	taxCfgOpMap = map[int]taxCfgOp{
		ChildTag:   enableChildDeduction(),
		AdultTag:   enableAdultDeduction(),
		HouseTag:   enableHouseDeduction(),
		RantingTag: enableRantingDeduction(),
		EldersTag:  enableEldersDeduction(),
		InfantTag:  enableInfantDeduction(),
	}
)

type TaxCfg struct {
	City              string      `json:"city"`
	CityEn            string      `json:"city_en"`
	Year              int         `json:"year"`      //
	ActiveAt          int         `json:"active_at"` // 生效月份
	UpperLimit        float64     `json:"upper_limit"`
	LowerLimit        float64     `json:"lower_limit"`
	Insurances        Ins         `json:"insurances"`         // 五险一金
	IncomeTax         []IncomeTax `json:"income_tax"`         // 个税计算阶梯配置
	Threshold         float64     `json:"threshold"`          // 个税起征点
	SpecialDeductions SDs         `json:"special_deductions"` // 专项附加扣除
}

type Salary struct {
	Original       float64 `json:"original"`
	AfterInsurance float64 `json:"after_insurance"`
	AfterThreshold float64 `json:"after_threshold"`
	Income         float64 `json:"income"`
	Tax            float64 `json:"tax"`
}

type salaryOp func(*Salary)

func doInsurance(tax *TaxCfg) salaryOp {
	return salaryOp(func(s *Salary) {
		insurances := []Insurance{tax.Insurances.Med, tax.Insurances.Retirement, tax.Insurances.Unemployment, tax.Insurances.Injury, tax.Insurances.Birth}
		var i float64
		base := s.Original
		if base > tax.UpperLimit {
			base = tax.UpperLimit
			i += base * tax.Insurances.HouseFund.Individual
		}
		if base < tax.LowerLimit {
			base = tax.LowerLimit
			i += s.Original * tax.Insurances.HouseFund.Individual
		}
		for _, insurance := range insurances {
			i += base * insurance.Individual
		}
		s.AfterInsurance = s.Original - i
	})
}

func doThreshold(threshold float64, sds SDs) salaryOp {
	return salaryOp(func(s *Salary) {
		if s.AfterInsurance >= 0 {
			ss := []SpecialDeduction{sds.Child, sds.Adult, sds.LifeOrDeath, sds.House, sds.Ranting, sds.Elders, sds.Infant}
			s.AfterThreshold = s.AfterInsurance - threshold
			for _, sd := range ss {
				if sd.Enable {
					s.AfterThreshold -= sd.Amount
				}
			}
			if s.AfterThreshold <= 0 {
				s.AfterThreshold = 0
			}
		}
	})
}

func (s *Salary) Do(options ...salaryOp) {
	for _, option := range options {
		option(s)
	}
}

type taxCfgOp func(*TaxCfg)

// SetHouseFund 仅有公积金在不同公司可能有区别，因此支持单独调整
func (t *TaxCfg) SetHouseFund(houseFund Insurance) *TaxCfg {
	// 默认是一起调整
	if houseFund.Company > 0.12 || houseFund.Individual > 0.12 {
		houseFund.Company = 0.12
		houseFund.Individual = 0.12
	}
	if houseFund.Company < 0.5 || houseFund.Individual < 0.5 {
		houseFund.Company = 0.5
		houseFund.Individual = 0.5
	}
	t.Insurances.HouseFund = houseFund
	return t
}

func enableChildDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.Child.Enable = true
	}
}

func enableAdultDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.Adult.Enable = true
	}
}

func enableLifeOrDeathDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.LifeOrDeath.Enable = true
	}
}

func enableHouseDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.House.Enable = true
	}
}

func enableRantingDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.Ranting.Enable = true
	}
}

func enableEldersDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.Elders.Enable = true
	}
}

func enableInfantDeduction() taxCfgOp {
	return func(t *TaxCfg) {
		t.SpecialDeductions.Infant.Enable = true
	}
}
