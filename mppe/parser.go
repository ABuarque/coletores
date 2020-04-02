package main

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dadosjusbr/storage"
)

const (
	// page name to process
	sheetName = "Sheet"

	// index of unique register code on the row
	registerCodeIndex = 0

	// index of name on the row
	nameIndex = 1

	// index of role on the row
	roleIndex = 2

	// index of total discount
	totalDiscountIndex = 14

	// index of ceil retention
	ceilRetentionIndex = 13

	// index of income tax
	incomeTaxIndex = 12

	// index of prev contribution
	prevContributionIndex = 11

	// index of income details
	totalIncomeDetailsIndex = 10

	// index of wage
	wageIndex = 4

	// index of total perks
	totalPerkIndex = 6
)

// Parse parses the xlsx tables
func Parse(paths []string) ([]storage.Employee, error) {
	var employees []storage.Employee
	for _, path := range paths {
		documentIdentification := getFileDocumentation(path)
		file, err := excelize.OpenFile(path)
		if err != nil {
			return nil, fmt.Errorf("error opening document %s for parse: %q", documentIdentification, err)
		}
		rows := file.GetRows(sheetName)
		numberOfRows := len(rows)
		var employee storage.Employee
		for index, row := range rows {
			if index == 0 || index == 1 || index == 2 || index == numberOfRows-1 {
				continue
			}
			discounts, err := getDiscounts(row, documentIdentification)
			if err != nil {
				return nil, err
			}
			incomeDetails, err := getIncome(row, documentIdentification)
			if err != nil {
				return nil, err
			}
			employee = storage.Employee{
				Reg:       row[registerCodeIndex],
				Name:      row[nameIndex],
				Role:      row[roleIndex],
				Type:      getType(documentIdentification),
				Workplace: "mppe",
				Active:    isActive(documentIdentification),
				Income:    incomeDetails,
				Discounts: discounts,
			}
			employees = append(employees, employee)
		}
	}
	return employees, nil
}

// it returns the total discounts sumary
func getDiscounts(row []string, documentIdentification string) (*storage.Discount, error) {
	totalDiscount, err := strconv.ParseFloat(row[totalDiscountIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing total discount from string to float64 for document %s: %q", documentIdentification, err)
	}
	ceilRetention, err := strconv.ParseFloat(row[ceilRetentionIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing ceil retention from string to float64 for document %s: %q", documentIdentification, err)
	}
	incomeTax, err := strconv.ParseFloat(row[incomeTaxIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing income tax from string to float64 for document %s: %q", documentIdentification, err)
	}
	prevContribution, err := strconv.ParseFloat(row[prevContributionIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing prev contribution from string to float64 for document %s: %q", documentIdentification, err)
	}
	return &storage.Discount{
		Total:            totalDiscount,
		CeilRetention:    &ceilRetention,
		IncomeTax:        &incomeTax,
		PrevContribution: &prevContribution,
	}, nil
}

// it returns the incomes sumary
func getIncome(row []string, documentIdentification string) (*storage.IncomeDetails, error) {
	total, err := strconv.ParseFloat(row[totalIncomeDetailsIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing total of income details from string to float64 for document %s: %q", documentIdentification, err)
	}
	wage, err := strconv.ParseFloat(row[wageIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing wage of income details from string to float64 for document %s: %q", documentIdentification, err)
	}
	perks, err := getPerks(row, documentIdentification)
	if err != nil {
		return nil, err
	}
	return &storage.IncomeDetails{
		Total: total,
		Wage:  &wage,
		Perks: perks,
		Other: nil,
	}, nil
}

// it retrieves employee perks
func getPerks(row []string, documentIdentification string) (*storage.Perks, error) {
	total, err := strconv.ParseFloat(row[totalPerkIndex], 64)
	if err != nil {
		return nil, fmt.Errorf("error on parsing total of perks from string to float64 for document %s: %q", documentIdentification, err)
	}
	return &storage.Perks{
		Total: total,
	}, nil
}

// it returns the employee type
func getType(documentIdentification string) string {
	switch documentIdentification {
	case "proventos-de-todos-os-membros-inativos":
		return "membro"
	case "proventos-de-todos-os-servidores-inativos":
		return "servidor"
	case "remuneracao-de-todos-os-membros-ativos":
		return "membro"
	case "remuneracao-de-todos-os-servidores-atuvos":
		return "servidor"
	case "valores-percebidos-por-todos-os-colaboradores":
		return "colaborador"
	case "valores-percebidos-por-todos-os-pensionistas":
		return "pensionista"
	default:
		return "indefinido"
	}
}

// it checks if the document is of active members or not
func isActive(documentIdentification string) bool {
	switch documentIdentification {
	case "proventos-de-todos-os-membros-inativos":
		return false
	case "proventos-de-todos-os-servidores-inativos":
		return false
	case "verbas-referentes-a-exercicios-anteriores":
		return false
	case "verbas-indenizatorias-e-outras-remuneracoes-temporarias":
		return false
	case "valores-percebidos-por-todos-os-pensionistas":
		return false
	default:
		return true
	}
}

// it cuts off month, year and extension from name
// to get file name
func getFileDocumentation(fileName string) string {
	fileSize := len(fileName)
	name := fileName[0 : fileSize-13]
	return name
}
