/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2019-10-29 11:26:52
**/

package db

import "errors"

//Model ...
type Model struct {
	tx *Tx
}

//GetTx ...
func (m *Model) GetTx() (*Tx, error) {
	if m.tx == nil {
		if lastCn.db == nil {
			return nil, errors.New(`no connection available`)
		}
		m.tx = &Tx{db: lastCn.db}
	}
	return m.tx, nil
}

//SetTx ...
func (m *Model) SetTx(tx *Tx) {
	m.tx = tx
}
