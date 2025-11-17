# Sample Data

This directory contains sample input data for testing and demonstrating the CLI statistics creator.

## File Structure

### input.json

Sample game score data in JSON format.

**Structure:**
```json
[
  {
    "player": {
      "name": "string",  // Player name
      "id": number       // Player ID (integer)
    },
    "score": number,     // Game score (integer)
    "level": number      // Player level (integer)
  }
]
```

**Field Definitions:**
- `player`: Player object containing:
  - `name`: Player's name (string)
  - `id`: Unique player identifier (positive integer)
- `score`: Game score value (no strict min/max, typically 0-100 range)
- `level`: Player's current level (positive integer)

**Usage:**
This sample data can be used to test statistical calculations such as:
- Average scores by level
- Min/max score analysis
- Player performance comparisons
- Level distribution statistics

**Data Characteristics:**
- Contains 40 sample records
- Score range: 67-98
- Level range: 2-6
- Balanced distribution across levels and players
