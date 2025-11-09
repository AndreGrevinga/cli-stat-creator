# Sample Data

This directory contains sample input data for testing and demonstrating the CLI statistics creator.

## File Structure

### input.json

Sample game score data in JSON format.

**Structure:**
```json
[
  {
    "player": "string",  // Player name
    "score": number,     // Game score (integer)
    "level": number      // Player level (integer)
  }
]
```

**Field Definitions:**
- `player`: Unique player identifier/name
- `score`: Game score value (no strict min/max, typically 0-100 range)
- `level`: Player's current level (positive integer)

**Usage:**
This sample data can be used to test statistical calculations such as:
- Average scores by level
- Min/max score analysis
- Player performance comparisons
- Level distribution statistics

**Data Characteristics:**
- Contains 10 sample records
- Score range: 67-98
- Level range: 2-5
- Balanced distribution across levels
