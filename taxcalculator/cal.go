package taxcalculator

import "fmt"

/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/29 -- 14:34
 @Author  : bishop ❤️ MONEY
 @Description: 工资个税计算器
*/

type TaxCalculator struct {
	year             int64
	defaultCfg       *TaxCfg
	fixedCfg         *TaxCfg
	specialDeduction int
}

func NewTaxCalculator(year int64, specialDeduction int) (*TaxCalculator, error) {
	if specialDeduction|HouseTag == 1 && specialDeduction|RantingTag == 1 {
		// 二者只能生效一个
		return nil, fmt.Errorf("住房贷款和租房见面只能生效一个，请任选其一")
	}
	tc := &TaxCalculator{year: year, specialDeduction: specialDeduction}
	tc.defaultCfg = beijing[year]
	tc.fixedCfg = beijing[year-1]

	for _, tag := range []int{ChildTag, AdultTag, HouseTag, RantingTag, EldersTag, InfantTag} {
		if tc.specialDeduction|tag == 1 {
			if op, ok := taxCfgOpMap[tag]; ok {
				op(tc.defaultCfg)
				op(tc.fixedCfg)
			}
		}
	}
	return tc, nil
}

func (t *TaxCalculator) Cal(originalSalary float64, month int64) (original, pureIncome, tax string) {
	res := t.calOneMonth(originalSalary, month)
	return fmt.Sprintf("%.4f", res.Original), fmt.Sprintf("%.4f", res.Income), fmt.Sprintf("%.4f", res.Tax)
}

func (t *TaxCalculator) CalThisYear(originalSalarys []float64) (pureIncome []string) {
	res := t.calEveryMonth(originalSalarys)
	for _, r := range res {
		pureIncome = append(pureIncome, fmt.Sprintf("%.4f", r.Income))
	}
	return
}

func (t *TaxCalculator) calOneMonth(original float64, month int64) *Salary {
	var originals []float64
	for i := int64(0); i < month; i++ {
		originals = append(originals, original)
	}
	res := t.calEveryMonth(originals)
	return res[len(res)-1]
}

func (t *TaxCalculator) calEveryMonth(originals []float64) []*Salary {
	var totalTaxFreeIncome float64
	var totalTax float64
	var res []*Salary

	for i := int64(1); i <= int64(len(originals)); i++ {
		cfg := t.fixCfgWithMonth(i)
		s := &Salary{Original: originals[i-1]}
		// 预缴五险一金 + 调整起税点，计算待缴税部分
		s.Do(doInsurance(cfg), doThreshold(cfg.Threshold, cfg.SpecialDeductions))
		// 统计当月全年税前实发
		totalTaxFreeIncome += s.AfterThreshold
		// 确定当月税率基准
		taxCfg := t.getTaxCfg(totalTaxFreeIncome, cfg.IncomeTax)
		// 计算当月税额
		tax := totalTaxFreeIncome*taxCfg.Rate - totalTax - taxCfg.FastFix
		// 统计全年税额
		totalTax += tax
		// 计算当月真实收入
		s.Tax = tax
		s.Income = s.AfterInsurance - tax
		// 加入结果集
		res = append(res, s)
	}
	return res
}

// CalOneMonthFreelance 自由职业者每个月薪水不同 or 本人当年有调薪；提供综合计算结果
func (t *TaxCalculator) calOneMonthFreelance(originals []float64) []*Salary {
	return t.calEveryMonth(originals)
}

func (t *TaxCalculator) getTaxCfg(totalTaxFreeIncome float64, incomeTax []IncomeTax) IncomeTax {
	for _, t := range incomeTax {
		if totalTaxFreeIncome >= t.Left && totalTaxFreeIncome < t.Right {
			return t
		}
	}
	return IncomeTax{}
}

func (t *TaxCalculator) fixCfgWithMonth(month int64) *TaxCfg {
	switch month {
	case 7, 8, 9, 10, 11, 12:
		return t.defaultCfg
	case 1, 2, 3, 4, 5, 6:
		return t.fixedCfg
	default:
		return t.defaultCfg
	}
}
