package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Data struct {
	MonthlyIncome      float64
	IncomeTax          string
	NetPayAfterTax     string
	SSS                string
	PhilHealth         string
	PagIbig            string
	TotalContributions string
	TotalDeductions    string
	NetPay             string
}

type PHP int64

func toPHP(f float64) PHP {
	return PHP((f * 100) + 0.5)
}

func (m PHP) float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

func (m PHP) multiply(f float64) PHP {
	x := (float64(m) * f) + 0.5
	return PHP(x)
}

func getSSS(monthlyIncome float64) float64 {
	var SSS float64

	table := [][5]float64{
		{0, 4250, 390.00, 180.00, 570.00},
		{4250, 4749.99, 437.50, 202.50, 640.00},
		{4750, 5249.99, 485.00, 225.00, 710.00},
		{5250, 5749.99, 532.50, 247.50, 780.00},
		{5750, 6249.99, 580.00, 270.00, 850.00},
		{6250, 6749.99, 627.50, 292.50, 920.00},
		{6750, 7249.99, 675.00, 315.00, 990.00},
		{7250, 7749.99, 722.50, 337.50, 1060.00},
		{7750, 8249.99, 770.00, 360.00, 1130.00},
		{8250, 8749.99, 12817.50, 382.50, 1200.00},
		{8750, 9249.99, 865.00, 405.00, 1270.00},
		{9250, 9749.99, 912.50, 427.50, 1340.00},
		{9750, 10249.99, 960.00, 450.00, 1410.00},
		{10250, 10749.99, 1007.50, 472.50, 1480.00},
		{10750, 11249.99, 1055.00, 495.00, 1550.00},
		{11250, 11749.99, 1102.50, 517.50, 1620.00},
		{11750, 12249.99, 1150.00, 540.00, 1690.00},
		{12250, 12749.99, 1197.50, 562.50, 1760.00},
		{12750, 13249.99, 1245.00, 585.00, 1830.00},
		{13250, 13749.99, 1292.50, 607.50, 1900.00},
		{13750, 14249.99, 1340.00, 630.00, 1970.00},
		{14250, 14749.99, 1387.50, 652.50, 2040.00},
		{14750, 15249.99, 1455.00, 675.00, 2130.00},
		{15250, 15749.99, 1502.50, 697.50, 2200.00},
		{15750, 16249.99, 1550.00, 720.00, 2270.00},
		{16250, 16749.99, 1597.50, 742.50, 2340.00},
		{16750, 17249.99, 1645.00, 765.00, 2410.00},
		{17250, 17749.99, 1692.50, 787.50, 2480.00},
		{17750, 18249.99, 1740.00, 810.00, 2550.00},
		{18250, 18749.99, 1787.50, 832.50, 2620.00},
		{18750, 19249.99, 1835.00, 855.00, 2690.00},
		{19250, 19749.99, 1882.50, 877.50, 2760.00},
		{19750, 20249.99, 1930.00, 900.00, 2830.00},
		{20250, 20749.99, 1977.50, 922.50, 2900.00},
		{20750, 21249.99, 2025.00, 945.00, 2970.00},
		{21250, 21749.99, 2072.50, 967.50, 3040.00},
		{21750, 22249.99, 2120.00, 990.00, 3110.00},
		{22250, 22749.99, 2167.50, 1012.50, 3180.00},
		{22270, 23249.99, 2215.00, 1035.00, 3250.00},
		{23250, 23749.99, 2262.50, 1057.50, 3320.00},
		{23750, 24249.99, 2310.00, 1080.00, 3390.00},
		{24250, 24279.99, 2357.50, 1102.50, 3460.00},
		{24750, 25249.99, 2405.00, 1125.00, 3530.00},
		{25250, 25749.99, 2452.50, 1147.50, 3600.00},
		{25750, 26249.99, 2500.00, 1170.00, 3670.00},
		{26250, 26749.99, 2547.50, 1192.50, 3740.00},
		{26750, 27249.99, 2595.00, 1215.00, 3810.00},
		{27250, 27749.99, 2642.50, 1237.50, 3880.00},
		{27750, 28249.99, 2690.00, 1260.00, 3950.00},
		{28250, 28749.99, 2737.50, 1282.50, 4020.00},
		{28750, 29249.99, 2785.00, 1305.00, 4090.00},
		{29250, 29749.99, 2832.50, 1327.50, 4160.00},
		{29750, 999999999999999, 2880.00, 1350.00, 4230.00},
	}

	for _, data := range table {
		if monthlyIncome < data[1] {
			SSS = data[3]
			break
		}
		if monthlyIncome >= data[0] && monthlyIncome <= data[1] {
			SSS = data[3]
			break
		}
	}

	return SSS
}

func getPhilHealth(monthlyIncome float64) float64 {
	var philHealth float64

	switch {
	case monthlyIncome <= 10000:
		philHealth = 500.00
	case monthlyIncome >= 100000:
		philHealth = 5000.00
	default:
		philHealth = toPHP(monthlyIncome).multiply(0.05).float64() / 2
	}

	return philHealth
}

func getPagIBIG(monthlyIncome float64) float64 {
	var pagIBIG float64
	var maxContribution float64

	switch {
	case monthlyIncome > 5000:
		maxContribution = 5000
	default:
		maxContribution = monthlyIncome
	}

	maxContributionPHP := toPHP(maxContribution)

	switch {
	case maxContribution > 1500:
		pagIBIG = maxContributionPHP.multiply(0.02).float64()
	default:
		pagIBIG = maxContributionPHP.multiply(0.01).float64()
	}

	return pagIBIG
}

func getTotalContributions(SSS float64, philHealth float64, pagIbig float64) float64 {
	SSSPHP := toPHP(SSS)
	PhilHealthPHP := toPHP(philHealth)
	pagIBIGPHP := toPHP(pagIbig)

	return (SSSPHP + PhilHealthPHP + pagIBIGPHP).float64()
}

func getIncomeTax(monthlyIncome float64, totalContributions float64) float64 {
	var taxableIncome float64 = monthlyIncome - totalContributions
	var incomeTax float64

	switch {
	case taxableIncome < 20833:
		incomeTax = 0
	case taxableIncome >= 20833 && taxableIncome < 33333:
		incomeTax = 0.00 + (taxableIncome-20833)*0.15
	case taxableIncome >= 33333 && taxableIncome < 66667:
		incomeTax = 1875 + (taxableIncome-33333)*0.20
	case taxableIncome >= 66667 && taxableIncome < 166667:
		incomeTax = 8541.8 + (taxableIncome-66667)*0.25
	case taxableIncome >= 166667 && taxableIncome < 666667:
		incomeTax = 33541.8 + (taxableIncome-166667)*0.30
	case taxableIncome >= 666667:
		incomeTax = 183541.8 + (taxableIncome-666667)*0.35
	}

	return incomeTax
}

func toPHPCurrency(amount float64) string {
	amountStr := strconv.FormatFloat(amount, 'f', 2, 64)

	parts := strings.Split(amountStr, ".")

	wholeNumber := addCommas(parts[0])

	var formattedAmount string
	if len(parts) > 1 {
		formattedAmount = fmt.Sprintf("₱ %s.%s", wholeNumber, parts[1])
	} else {
		formattedAmount = fmt.Sprintf("₱ %s", wholeNumber)
	}

	return formattedAmount
}

func addCommas(amount string) string {
	n := len(amount)
	if n <= 3 {
		return amount
	}
	return addCommas(amount[:n-3]) + "," + amount[n-3:]
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("index.html")

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	monthlyIncome, _ := strconv.ParseFloat(strings.TrimSpace(r.FormValue("monthly-income")), 64)

	// Calculate values
	SSS := getSSS(monthlyIncome)
	philHealth := getPhilHealth(monthlyIncome)
	pagIbig := getPagIBIG(monthlyIncome)
	totalContributions := getTotalContributions(SSS, philHealth, pagIbig)
	incomeTax := getIncomeTax(monthlyIncome, totalContributions)

	// PHP Computations
	netPayAfterTaxPHP := monthlyIncome - incomeTax
	totalDeductionsPHP := incomeTax + totalContributions
	netPayPHP := monthlyIncome - totalDeductionsPHP

	// Data to pass to template
	data := Data{
		MonthlyIncome:      monthlyIncome,
		IncomeTax:          toPHPCurrency(incomeTax),
		NetPayAfterTax:     toPHPCurrency(netPayAfterTaxPHP),
		SSS:                toPHPCurrency(SSS),
		PhilHealth:         toPHPCurrency(philHealth),
		PagIbig:            toPHPCurrency(pagIbig),
		TotalContributions: toPHPCurrency(totalContributions),
		TotalDeductions:    toPHPCurrency(totalDeductionsPHP),
		NetPay:             toPHPCurrency(netPayPHP),
	}

	tmpl.Execute(w, data)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", handler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
