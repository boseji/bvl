// db_test.go - Part of Tests for the `inventory` Package
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

package inventory_test

import (
	"database/sql"
	"testing"

	"github.com/boseji/bvl/inventory"
)

func setupTestDB(t *testing.T) *sql.DB {
	db := inventory.OpenDB(":memory:")
	if db == nil {
		t.Fatal("failed to create test DB")
	}
	return db
}

func TestOpenDB(t *testing.T) {
	db := inventory.OpenDB(":memory:")
	if db == nil {
		t.Fatal("expected non-nil db")
	}
	db.Close()
}

func TestResetSequence(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := inventory.ResetSequence(db)
	if err != nil {
		t.Fatalf("ResetSequence failed: %v", err)
	}
}

func TestAddAndGetItem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	item := inventory.Item{
		Description: "UPS", Location: "Rack 1",
		Status: "Operational", Remarks: "installed",
	}
	err := inventory.AddItem(db, item)
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	items, err := inventory.ListAll(db)
	if err != nil || len(items) != 1 {
		t.Fatalf("ListAll failed: %v", err)
	}

	got, err := inventory.GetItemByID(db, items[0].ID)
	if err != nil {
		t.Fatalf("GetItemByID failed: %v", err)
	}
	if got.Description != "UPS" {
		t.Errorf("unexpected Description: %s", got.Description)
	}
}

func TestGetItemByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := inventory.GetItemByID(db, 9999)
	if err == nil {
		t.Errorf("expected error for missing ID")
	}
}

func TestAppendItem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	item := inventory.Item{
		ID:          1005,
		Description: "Battery", Location: "Shelf 1",
		Status: "New", Remarks: "received",
	}
	err := inventory.AppendItem(db, item)
	if err != nil {
		t.Fatalf("AppendItem failed: %v", err)
	}

	got, err := inventory.GetItemByID(db, 1005)
	if err != nil {
		t.Fatalf("GetItemByID failed: %v", err)
	}
	if got.Description != "Battery" {
		t.Errorf("unexpected Description: %s", got.Description)
	}
}

func TestEditItem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	item := inventory.Item{
		Description: "Switch", Location: "Rack 2",
		Status: "Installed", Remarks: "new",
	}
	err := inventory.AddItem(db, item)
	if err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	items, _ := inventory.ListAll(db)
	item = items[0]
	item.Status = "Active"
	item.Remarks = "activated"

	err = inventory.EditItem(db, item)
	if err != nil {
		t.Fatalf("EditItem failed: %v", err)
	}

	got, _ := inventory.GetItemByID(db, item.ID)
	if got.Status != "Active" {
		t.Errorf("unexpected Status: %s", got.Status)
	}
}

func TestEditItem_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	item := inventory.Item{
		ID: 9999, Description: "Ghost", Location: "Void",
		Status: "Lost", Remarks: "none",
	}
	err := inventory.EditItem(db, item)
	if err != nil {
		t.Fatalf("EditItem on missing ID should not fail: %v", err)
	}
}

func TestAppendRemarksEntry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	item := inventory.Item{
		Description: "Router", Location: "Rack 3",
		Status: "Operational", Remarks: "setup",
	}
	_ = inventory.AddItem(db, item)
	items, _ := inventory.ListAll(db)

	err := inventory.AppendRemarksEntry(db, items[0].ID, "tested")
	if err != nil {
		t.Fatalf("AppendRemarksEntry failed: %v", err)
	}
}

func TestAppendRemarksEntry_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := inventory.AppendRemarksEntry(db, 9999, "should fail")
	if err == nil {
		t.Errorf("expected error for missing ID")
	}
}

func TestDeleteItem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	item := inventory.Item{
		Description: "Server", Location: "Rack 4",
		Status: "Provisioning", Remarks: "image install",
	}
	_ = inventory.AddItem(db, item)
	items, _ := inventory.ListAll(db)

	err := inventory.DeleteItem(db, items[0].ID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

	_, err = inventory.GetItemByID(db, items[0].ID)
	if err == nil {
		t.Errorf("expected error for deleted item")
	}
}

func TestDeleteItem_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	err := inventory.DeleteItem(db, 9999)
	if err != nil {
		t.Fatalf("DeleteItem on missing ID should not fail: %v", err)
	}
}

func TestListItemsPaged(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	for i := 0; i < 5; i++ {
		item := inventory.Item{
			Description: "Device", Location: "Loc",
			Status: "OK", Remarks: "added",
		}
		_ = inventory.AddItem(db, item)
	}

	items, err := inventory.ListItemsPaged(db, 0, 3)
	if err != nil || len(items) != 3 {
		t.Fatalf("ListItemsPaged failed: %v", err)
	}

	// edge case afterID too large
	empty, err := inventory.ListItemsPaged(db, 9999, 5)
	if err != nil {
		t.Fatalf("ListItemsPaged failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected empty slice for large afterID")
	}
}

func TestItemIterator(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	for i := 0; i < 3; i++ {
		item := inventory.Item{
			Description: "FW", Location: "Loc",
			Status: "Ready", Remarks: "fw added",
		}
		_ = inventory.AddItem(db, item)
	}

	iter, err := inventory.NewItemIterator(db, "WHERE status = ?", "Ready")
	if err != nil {
		t.Fatalf("NewItemIterator failed: %v", err)
	}
	defer iter.Close()

	count := 0
	for {
		item, ok, err := iter.Next()
		if err != nil {
			t.Fatalf("Iterator Next failed: %v", err)
		}
		if !ok {
			break
		}
		if item.Status != "Ready" {
			t.Errorf("unexpected status: %s", item.Status)
		}
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 items, got %d", count)
	}
}

func TestItemIterator_BadWhereClause(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := inventory.NewItemIterator(db, "WHERE no_such_field = 1")
	if err == nil {
		t.Errorf("expected error for bad WHERE clause")
	}
}
