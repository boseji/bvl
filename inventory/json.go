// json.go - Part of the `inventory` Package
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

//
// JSON import/export/view functions for Inventory CLI
// Provides both raw DB and InventoryDB methods
// Ready for Web, Electron, CLI, jq integration
//

package inventory

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

// ExportJSON writes all inventory records to a JSON file.
//
// The output JSON is an array of Item objects:
//
//	[
//	  { "id": 1001, "description": "UPS", "location": "Rack 1", ... },
//	  { ... },
//	  ...
//	]
//
// Usage:
//
//	err := ExportJSON(db, "inventory.json")
//
// Example:
//
//	// Export inventory to file
//	err := inventory.ExportJSON(db, "export.json")
//
// Errors:
//   - returns error if database query fails
//   - returns error if JSON marshal fails
//   - returns error if file cannot be written (permission, path)
func ExportJSON(db *sql.DB, filename string) error {
	items, err := ListAll(db)
	if err != nil {
		return fmt.Errorf("export json failed: %v", err)
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json failed: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("write json failed: %v", err)
	}

	return nil
}

// ImportJSON reads inventory records from a JSON file and imports them.
//
// Existing records with matching IDs will be replaced.
//
// Usage:
//
//	err := ImportJSON(exec, "inventory.json")
//
// Example:
//
//	// Import from JSON
//	err := inventory.ImportJSON(inv, "import.json")
//
// Errors:
//   - returns error if file read fails
//   - returns error if JSON unmarshal fails
//   - returns error if individual Insert/Replace fails
func ImportJSON(exec Execer, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read json failed: %v", err)
	}

	return ImportJSONFromBytes(exec, data)
}

// ViewJSON pretty prints the content of a JSON file to stdout.
//
// Usage:
//
//	err := ViewJSON("inventory.json")
//
// Example:
//
//	err := inventory.ViewJSON("export.json")
//
// Useful for:
//   - CLI debugging
//   - Inspection
//   - jq pipelines
//
// Errors:
//   - returns error if file read fails
//   - returns error if JSON is invalid
func ViewJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read json failed: %v", err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return fmt.Errorf("format json failed: %v", err)
	}

	fmt.Println(out.String())
	return nil
}

// ExportJSONToString returns all inventory records as a JSON string.
//
// Usage:
//
//	jsonStr, err := ExportJSONToString(db)
//
// Example:
//
//	jsonStr, err := inventory.ExportJSONToString(db)
//
// Useful for:
//   - Web API response
//   - Electron UI
//   - jq processing
//   - CLI --json flag
func ExportJSONToString(db *sql.DB) (string, error) {
	items, err := ListAll(db)
	if err != nil {
		return "", fmt.Errorf("export json string failed: %v", err)
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json failed: %v", err)
	}

	return string(data), nil
}

// ImportJSONFromString reads inventory records from a JSON string.
//
// Existing records with matching IDs will be replaced.
//
// Usage:
//
//	err := ImportJSONFromString(exec, jsonString)
//
// Example:
//
//	err := inventory.ImportJSONFromString(inv, jsonPayload)
//
// Errors:
//   - returns error if JSON is invalid
//   - returns error if DB insert fails
func ImportJSONFromString(exec Execer, jsonString string) error {
	return ImportJSONFromBytes(exec, []byte(jsonString))
}

// ImportJSONFromBytes helper
func ImportJSONFromBytes(exec Execer, data []byte) error {
	var items []Item

	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("unmarshal json failed: %v", err)
	}

	for i, item := range items {
		if err := AppendItem(exec, item); err != nil {
			return fmt.Errorf("import item %d failed: %v", i, err)
		}
	}

	return nil
}

// ToJSON returns this Item as a JSON string.
//
// Usage:
//
//	jsonStr, err := item.ToJSON()
//
// Example output:
//
//	{
//	  "id": 1001,
//	  "description": "UPS",
//	  "location": "Rack 1",
//	  "status": "Operational",
//	  "remarks": "[2025-06-21 15:00] installed UPS"
//	}
//
// Useful for:
//   - Logging
//   - CLI --json
//   - Websocket events
//   - API single-item
func (item *Item) ToJSON() (string, error) {
	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal item json failed: %v", err)
	}
	return string(data), nil
}

// FromJSON parses a JSON string into this Item.
//
// Usage:
//
//	var item Item
//	err := item.FromJSON(jsonStr)
//
// Example:
//
//	var item Item
//	err := item.FromJSON(`{"id":1001,"description":"UPS"}`)
//
// Errors:
//   - returns error if JSON is invalid
func (item *Item) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), item)
}

// InventoryDB method: ExportJSON
//
// Usage:
//
//	err := inv.ExportJSON("inventory.json")
func (inv *InventoryDB) ExportJSON(filename string) error {
	return ExportJSON(inv.db, filename)
}

// InventoryDB method: ImportJSON
//
// Usage:
//
//	err := inv.ImportJSON("inventory.json")
//
// Runs inside transaction.
func (inv *InventoryDB) ImportJSON(filename string) error {
	return inv.WithTransaction(func(tx Execer) error {
		return ImportJSON(tx, filename)
	})
}

// InventoryDB method: ExportJSONToString
//
// Usage:
//
//	jsonStr, err := inv.ExportJSONToString()
func (inv *InventoryDB) ExportJSONToString() (string, error) {
	return ExportJSONToString(inv.db)
}

// InventoryDB method: ImportJSONFromString
//
// Usage:
//
//	err := inv.ImportJSONFromString(jsonString)
//
// Runs inside transaction.
func (inv *InventoryDB) ImportJSONFromString(jsonString string) error {
	return inv.WithTransaction(func(tx Execer) error {
		return ImportJSONFromString(tx, jsonString)
	})
}
