package main

import "fmt"

func main() {
	fmt.Println("Reconciliation Service")

	// input
	internalData := os.Open()

	var uc ReconciliationUsecase
	uc = new ReconciliationUsecaseImpl()
	uc.ReconcileData()
}
