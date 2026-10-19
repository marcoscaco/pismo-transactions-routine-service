// Entidade de negocio para representar um tipo de operação no sistema.
package entity

// OperationType representa um tipo de operação no sistema.
type OperationType struct {
	ID               int64
	Description      string
	AmountMultiplier int32
}
