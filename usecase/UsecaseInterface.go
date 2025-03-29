package usecase

type ReconciliationUsecase interface {
	ReconcileData(internalTrxData multipart.File, bankStatement []multipart.File)
}
