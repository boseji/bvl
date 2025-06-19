-- Migration Script for Inventory Database
--
-- bvl - Boseji's Inventory Management Program
--
-- Sources
-- -------
-- https://github.com/boseji/bvl
--
-- License
-- -------
--
--   bvl - Boseji's Inventory Management Program
--   Copyright (C) 2025 by Abhijit Bose (aka. Boseji)
--
--   This program is free software: you can redistribute it and/or modify
--   it under the terms of the GNU General Public License version 2 only
--   as published by the Free Software Foundation.
--
--   This program is distributed in the hope that it will be useful,
--   but WITHOUT ANY WARRANTY; without even the implied warranty of
--   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. 
--
--   You should have received a copy of the GNU General Public License
--   along with this program. If not, see <https://www.gnu.org/licenses/>.
--
--  SPDX-License-Identifier: GPL-2.0-only
--  Full Name: GNU General Public License v2.0 only
--  Please visit <https://spdx.org/licenses/GPL-2.0-only.html> for details.
--

PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

-- Create inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT,
    location TEXT,
    status TEXT,
    remarks TEXT
);

-- Initialize AUTOINCREMENT sequence to start at 1001
DELETE FROM sqlite_sequence WHERE name = 'inventory';
INSERT INTO sqlite_sequence (name, seq) VALUES ('inventory', 1000);

COMMIT;
