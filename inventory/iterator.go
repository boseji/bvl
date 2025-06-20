// iterator.go - Part of the `inventory` Package
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

// ItemIterator provides a streaming interface to iterate over
// inventory records, one item at a time.
//
// Internally uses sql.Rows and rows.Next().
// Allows for processing large result sets with low memory usage.
//
// Usage:
//
//	iter, err := NewItemIterator(db, "WHERE status = ?", "Operational")
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
//
// Use cases:
//
// - To process large inventories without loading all into memory
// - To stream records to external systems
// - To filter results with WHERE clause
//
// Notes:
//
// - You must call Close() when done to release database resources
// - The iterator must be used in a single goroutine
// - If WHERE clause is empty (""), all records are returned
// - Always check for error on Next() even if ok == false
type ItemIterator struct {
	rows *sql.Rows
}

// NewItemIterator returns an ItemIterator for scanning records
// in the inventory table with an optional WHERE clause.
//
// The iterator streams results one at a time and uses minimal memory.
//
// Usage:
//
//	iter, err := NewItemIterator(inv.DB(), "WHERE status = ?", "Operational")
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
//
// Use cases:
//
// - To process large inventories without loading entire table
// - To filter items with a dynamic WHERE clause
// - To support streaming export to CSV, JSON, etc.
//
// Notes:
//
// - WHERE clause must begin with "WHERE ..." or be empty string ""
// - Use parameter substitution for arguments (? placeholders)
// - Must call Close() when done to release database resources
func NewItemIterator(
	db *sql.DB, whereClause string, args ...interface{},
) (*ItemIterator, error) {

	query := `
        SELECT id, description, location, status, remarks
        FROM inventory `
	if whereClause != "" {
		query += whereClause
	}
	query += " ORDER BY id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("iterator query failed: %v", err)
	}

	return &ItemIterator{rows: rows}, nil
}

// Next returns the next item from the iterator.
//
// Usage:
//
//	for {
//	    item, ok, err := iter.Next()
//	    if err != nil {
//	        // handle error (scan error)
//	    }
//	    if !ok {
//	        break // no more rows
//	    }
//	    fmt.Println(item.ID, item.Description)
//	}
//
// Return values:
//
// - item: the next Item in the result set (if ok == true)
// - ok: true if a row was returned, false if at end of result set
// - err: non-nil if scan failed
//
// Use cases:
//
// - To process items one at a time from a filtered query
// - For streaming export (CSV, JSON)
// - For large inventories without loading all items into memory
//
// Notes:
//
// - Must check err even if ok == false
// - Must call Close() on the iterator after use
// - Each call advances the cursor (forward-only)
// - This is not thread-safe: use only in single goroutine
func (it *ItemIterator) Next() (Item, bool, error) {
	var item Item
	if it.rows.Next() {
		err := it.rows.Scan(
			&item.ID, &item.Description, &item.Location,
			&item.Status, &item.Remarks)
		if err != nil {
			return item, false, fmt.Errorf("iterator scan failed: %v", err)
		}
		return item, true, nil
	}
	return item, false, nil
}

// Close releases the database resources held by this iterator.
//
// You must call Close() when you are finished iterating,
// otherwise database connections may be leaked.
//
// Usage:
//
//	iter, err := NewItemIterator(inv.DB(), "WHERE status = ?", "Operational")
//	if err != nil {
//	    // handle error
//	}
//	defer iter.Close() // always defer Close!
//
//	for {
//	    item, ok, err := iter.Next()
//	    if err != nil {
//	        // handle error
//	    }
//	    if !ok {
//	        break
//	    }
//	    fmt.Println(item.ID, item.Description)
//	}
//
// Use cases:
//
// - Always used after calling NewItemIterator()
// - Should be deferred immediately after iterator creation
//
// Notes:
//
// - Safe to call even if no rows were read
// - Safe to call multiple times (subsequent calls will do nothing)
// - Does not affect the underlying database connection
func (it *ItemIterator) Close() error {
	return it.rows.Close()
}
