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
// for use with SQLite databases. It has with full CRUD support,
// with CSV/JSON import/export, and ready-to-use CLI and web features.
//
//	बी.वी.एल - बोसजी के द्वारा रचित भंडार लेखांकन हेतु तन्त्राक्ष्।
//
// =============================================
//
// एक सुगम एवं उपयोगी भंडार संचालन हेतु तन्त्राक्ष्।
//
// एक रचनात्मक भारतीय उत्पाद ।
//
// bvl - Boseji's Inventory Management Program
//
// # Package inventory
//
// Core Features:
//
// - InventoryDB wrapper: safe, transactional DB access
// - In-memory / file SQLite support
// - Configurable sequence start (IndexStart)
//
// Data Model:
//
//   - Item struct:
//     ID, Description, Location, Status, Remarks (with FormatRemarks)
//   - FormatRemarks(): consistent timestamped remarks
//   - JSON tags for web/app/API compatibility
//
// Database Operations:
//
// - AddItem()
// - EditItem()
// - DeleteItem()
// - AppendItem()
// - AppendRemarksEntry()
// - GetItemByID()
// - ListAll()
// - ListItemsPaged() with pagination
// - NewItemIterator() with streaming Next()
// - ResetSequence()
//
// CSV Support:
//
// - ExportCSV()
// - ImportCSV()
// - ViewCSV()
// - InventoryDB wrappers
// - CLI-friendly and Excel-friendly format
//
// JSON Support:
//
// - ExportJSON()
// - ImportJSON()
// - ViewJSON()
// - ExportJSONToString()
// - ImportJSONFromString()
// - InventoryDB wrappers
// - CLI-friendly, Web API-ready format
// - jq / Web / Electron compatibility
//
// Item JSON helpers:
//
// - Item.ToJSON()
// - Item.FromJSON()
//
// Test Suite:
//
// - db_test.go: core DB functions
// - inventorydb_test.go: InventoryDB methods
// - csv_test.go: CSV
// - json_test.go: JSON
// - All error paths covered
// - Rollback scenarios covered
//
// Ready for:
//
// - CLI tools
// - Web frontends
// - Electron apps
// - REST APIs
// - jq pipelines
// - Automation scripts
// - CI/CD integration
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
