// inventory.go - Part of the `inventory` Package
//
//     ॐ भूर्भुवः स्वः
//     तत्स॑वि॒तुर्वरे॑ण्यं॒
//    भर्गो॑ दे॒वस्य॑ धीमहि।
//   धियो॒ यो नः॑ प्रचो॒दया॑त्॥
//
//
//  बी.वी.एल - बोसजी के द्वारा रचित भंडार लेखांकन हेतु तन्त्राक्ष्।
// =============================================
//
// एक सुगम एवं उपयोगी भंडार संचालन हेतु तन्त्राक्ष्।
//
// एक रचनात्मक भारतीय उत्पाद ।
//
// bvl - Boseji's Inventory Management Program
//
// Easy to use and useful stock, goods and materials handling software.
//
// Sources
// -------
// https://github.com/boseji/bvl
//
// License
// -------
//
//   bvl - Boseji's Inventory Management Program.
//   Copyright (C) 2025 by Abhijit Bose (aka. Boseji)
//
//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License version 2 only
//   as published by the Free Software Foundation.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
//
//   You should have received a copy of the GNU General Public License
//   along with this program. If not, see <https://www.gnu.org/licenses/>.
//
//  SPDX-License-Identifier: GPL-2.0-only
//  Full Name: GNU General Public License v2.0 only
//  Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
//

package inventory

import (
	"database/sql"
	"fmt"
)

// IndexStart defines the starting value for auto-incremented IDs.
const (
	// Starting value of the Index
	IndexStart = 1000
)

// Execer defines something that can Exec SQL.
// Both *sql.DB and *sql.Tx implement this.
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// InventoryDB wraps *sql.DB and provides safe transaction helpers.
//
// Users do not need to work with *sql.DB directly.
type InventoryDB struct {
	db *sql.DB
}

// NewInventoryDB opens or creates the database and returns InventoryDB.
//
// Ensures the table exists, sequence is initialized.
// Returns a ready-to-use InventoryDB wrapper.
//
// Usage:
//
//	inv := NewInventoryDB("inventory.db")
//
// Notes:
// - Underlying connection is stored in inv.db
// - Close() must be called when finished
// - Table creation is idempotent
func NewInventoryDB(dbFile string) *InventoryDB {
	db := OpenDB(dbFile)
	return &InventoryDB{db: db}
}

// WithTransaction executes the given function inside a transaction.
//
// Usage:
//
//	err := inv.WithTransaction(func(tx Execer) error {
//	    err := AddItem(tx, item)
//	    if err != nil {
//	        return err
//	    }
//	    return AppendRemarksEntry(tx, item.ID, "added new")
//	})
//
// If fn() returns error:
// - Transaction is rolled back
//
// If fn() returns nil:
// - Transaction is committed
//
// Notes:
// - Use for any group of changes that must be atomic
// - If the DB fails, returns error
func (inv *InventoryDB) WithTransaction(
	fn func(tx Execer) error) error {

	tx, err := inv.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %v", err)
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx failed: %v", err)
	}

	return nil
}

// DB returns the underlying *sql.DB (for read-only queries).
// Use only when needed, e.g. for GetItemByID.
func (inv *InventoryDB) DB() *sql.DB {
	return inv.db
}

// Close closes the underlying database connection.
func (inv *InventoryDB) Close() error {
	return inv.db.Close()
}

// AppendItem wraps AppendItem with automatic transaction.
//
// Usage:
//
//	err := inv.AppendItem(item)
func (inv *InventoryDB) AppendItem(item Item) error {
	return inv.WithTransaction(func(tx Execer) error {
		return AppendItem(tx, item)
	})
}

// AddItem wraps AddItem with automatic transaction.
//
// Usage:
//
//	err := inv.AddItem(item)
func (inv *InventoryDB) AddItem(item Item) error {
	return inv.WithTransaction(func(tx Execer) error {
		return AddItem(tx, item)
	})
}

// EditItem wraps EditItem with automatic transaction.
//
// Usage:
//
//	err := inv.EditItem(item)
func (inv *InventoryDB) EditItem(item Item) error {
	return inv.WithTransaction(func(tx Execer) error {
		return EditItem(tx, item)
	})
}

// AppendRemarksEntry wraps AppendRemarksEntry with automatic transaction.
//
// Usage:
//
//	err := inv.AppendRemarksEntry(id, "log message")
func (inv *InventoryDB) AppendRemarksEntry(id int, message string) error {
	return inv.WithTransaction(func(tx Execer) error {
		return AppendRemarksEntry(tx, id, message)
	})
}

// DeleteItem wraps DeleteItem with automatic transaction.
//
// Usage:
//
//	err := inv.DeleteItem(id)
func (inv *InventoryDB) DeleteItem(id int) error {
	return inv.WithTransaction(func(tx Execer) error {
		return DeleteItem(tx, id)
	})
}

// ResetSequence wraps ResetSequence with automatic transaction.
//
// Usage:
//
//	err := inv.ResetSequence()
func (inv *InventoryDB) ResetSequence() error {
	return inv.WithTransaction(func(tx Execer) error {
		return ResetSequence(tx)
	})
}

// GetItemByID wraps GetItemByID.
//
// Usage:
//
//	item, err := inv.GetItemByID(id)
func (inv *InventoryDB) GetItemByID(id int) (Item, error) {
	return GetItemByID(inv.db, id)
}

// ListAll wraps ListAll.
//
// Usage:
//
//	items, err := inv.ListAll()
func (inv *InventoryDB) ListAll() ([]Item, error) {
	return ListAll(inv.db)
}

// ListItemsPaged wraps ListItemsPaged.
//
// Usage:
//
//	items, err := inv.ListItemsPaged(afterID, limit)
func (inv *InventoryDB) ListItemsPaged(afterID int, limit int) ([]Item, error) {
	return ListItemsPaged(inv.db, afterID, limit)
}

// NewItemIterator returns an ItemIterator for scanning records
// with an optional WHERE clause.
//
// Usage:
//
//	iter, err := inv.NewItemIterator("WHERE status = ?", "Operational")
//	if err != nil {
//	    // handle error
//	}
//	defer iter.Close()
//
//	for {
//	    item, ok, err := iter.Next()
//	    if err != nil {
//	        // handle error
//	    }
//	    if !ok {
//	        break // end of results
//	    }
//	    fmt.Println(item.ID, item.Description)
//	}
func (inv *InventoryDB) NewItemIterator(
	whereClause string, args ...interface{},
) (*ItemIterator, error) {
	return NewItemIterator(inv.db, whereClause, args...)
}
