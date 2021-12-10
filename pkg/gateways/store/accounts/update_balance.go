package accounts

func (s AccountStorage) UpdateBalance(originID, destinationID string, amount int) error {

	originAccount := s.storage[originID]
	destinationAccount := s.storage[destinationID]

	originAccount.Balance -= amount
	destinationAccount.Balance += amount

	s.Upsert(originID, originAccount)
	s.Upsert(destinationID, destinationAccount)
	return nil
}