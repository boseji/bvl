// doc.go - Part of the `inventory` Package
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

// Package inventory provides a complete inventory management layer
// for use with SQLite databases.
//
// The package defines the InventoryDB type as a convenient wrapper
// around *sql.DB and exposes transaction-safe methods for adding,
// editing, deleting, listing, and querying inventory records.
//
// Inventory records include:
//
//   - id (auto-increment primary key)
//   - description (text)
//   - location (text)
//   - status (text)
//   - remarks (append-only log)
//
// The remarks field acts as an audit log with timestamped entries.
//
// Usage:
//
//	inv := inventory.NewInventoryDB("inventory.db")
//	defer inv.Close()
//
//	err := inv.AppendItem(item)
//	item, err := inv.GetItemByID(1234)
//	items, err := inv.ListAll()
//
// Transactions:
//
// For grouped writes, use:
//
//	err := inv.WithTransaction(func(tx inventory.Execer) error {
//	    return inventory.AppendItem(tx, item)
//	})
//
// Iteration:
//
//	iter, err := inv.NewItemIterator("WHERE status = ?", "Operational")
//	defer iter.Close()
//	for {
//	    item, ok, err := iter.Next()
//	    if err != nil || !ok { break }
//	    fmt.Println(item.ID, item.Description)
//	}
//
// Conventions:
//
// - All methods are documented with usage examples
// - All error messages are lowercase
// - Line width is limited to 80 characters for clarity
// - Remarks field is append-only for audit trail
//
// License:
//
// This package is GPL-2.0-only.
//
// bvl - Boseji's Inventory Management Program.
// Copyright (C) 2025 by Abhijit Bose (aka. Boseji).
//
// SPDX-License-Identifier: GPL-2.0-only
// Full Name: GNU General Public License v2.0 only
// Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
//
// Sources:
// https://github.com/boseji/bvl
package inventory
