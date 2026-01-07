# Web Interface Enhancement and Charts Design Document

**Date:** 2025-12-09
**Status:** Planning
**Author:** Architecture Design Session

## Overview

This document describes the design for enhancing the web interface with improved UI/UX and adding interactive statistical charts. The enhancements will include visual data representations (histograms, distribution charts, standard deviation visualizations) while maintaining the existing htmx-based architecture.

## Goals

1. **Improve UI/UX**: Modernize the interface with better visual hierarchy, colors, and responsive design
2. **Add Data Visualizations**: Implement interactive charts for statistical analysis
3. **Maintain Architecture**: Keep the existing htmx + Go template approach
4. **Enhance User Experience**: Make statistics easier to understand through visual representations

## Technology Stack

### Chart Library Selection

**Recommended: Chart.js v4.x**

**Rationale:**
- Lightweight (~200KB minified)
- Excellent documentation and community support
- Native support for histograms, box plots, and distribution charts
- Easy integration with htmx and server-rendered HTML
- No complex dependencies
- Good accessibility features
- Responsive by default

**CDN Link:**
```html
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
```

**Alternative: Plotly.js** (if more advanced statistical features are needed later)
- More feature-rich for scientific visualizations
- Built-in statistical analysis tools
- Trade-off: Larger bundle size (~3MB)

### UI Framework

**Continue with:**
- Vanilla CSS with CSS custom properties (CSS variables)
- No CSS framework to keep the bundle small
- Progressive enhancement approach

## High-Level Architecture

### Component Structure

```
web/
├── static/
│   ├── index.html           (existing - enhanced)
│   ├── styles.css           (existing - enhanced)
│   ├── app.js              (existing - enhanced)
│   └── charts.js           (new - chart management)
└── templates/
    ├── results.html        (existing - enhanced)
    ├── overall-stats.html  (existing - enhanced with chart container)
    ├── level-stats.html    (existing - enhanced with chart container)
    └── player-stats.html   (existing - enhanced with chart container)
```

### Data Flow

1. User uploads file → htmx POST to `/api/stats`
2. Server processes data and renders templates with statistics
3. Templates include chart data in `data-*` attributes or `<script>` tags
4. JavaScript detects new content and initializes charts
5. Charts render using statistics data from the DOM

## UI/UX Enhancements

### Visual Design Improvements

#### Color Scheme
```css
:root {
  /* Primary colors */
  --primary-color: #3b82f6;      /* Blue - main actions */
  --primary-hover: #2563eb;
  --primary-light: #dbeafe;

  /* Secondary colors */
  --secondary-color: #8b5cf6;    /* Purple - accents */
  --success-color: #10b981;      /* Green - positive stats */
  --warning-color: #f59e0b;      /* Orange - warnings */
  --error-color: #ef4444;        /* Red - errors */

  /* Neutral colors */
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --text-tertiary: #9ca3af;
  --border-color: #e5e7eb;

  /* Chart colors */
  --chart-blue: #3b82f6;
  --chart-green: #10b981;
  --chart-purple: #8b5cf6;
  --chart-orange: #f59e0b;
  --chart-red: #ef4444;
  --chart-teal: #14b8a6;
  --chart-pink: #ec4899;
  --chart-yellow: #eab308;
}

/* Dark mode support (future enhancement) */
@media (prefers-color-scheme: dark) {
  :root {
    --bg-primary: #1f2937;
    --bg-secondary: #111827;
    --text-primary: #f9fafb;
    /* ... more dark mode variables */
  }
}
```

#### Typography
```css
/* Font stack */
:root {
  --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
               'Helvetica Neue', Arial, sans-serif;
  --font-mono: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono',
               Consolas, monospace;
}

/* Type scale */
--text-xs: 0.75rem;    /* 12px */
--text-sm: 0.875rem;   /* 14px */
--text-base: 1rem;     /* 16px */
--text-lg: 1.125rem;   /* 18px */
--text-xl: 1.25rem;    /* 20px */
--text-2xl: 1.5rem;    /* 24px */
--text-3xl: 1.875rem;  /* 30px */
--text-4xl: 2.25rem;   /* 36px */
```

#### Spacing System
```css
:root {
  --spacing-xs: 0.25rem;   /* 4px */
  --spacing-sm: 0.5rem;    /* 8px */
  --spacing-md: 1rem;      /* 16px */
  --spacing-lg: 1.5rem;    /* 24px */
  --spacing-xl: 2rem;      /* 32px */
  --spacing-2xl: 3rem;     /* 48px */
  --spacing-3xl: 4rem;     /* 64px */
}
```

### Layout Enhancements

#### Header Section
- Add visual hierarchy with better typography
- Include breadcrumb or status indicator
- Add optional filters toggle button
- Improve mobile responsiveness

#### Upload Section
- Better visual feedback for drag-and-drop
- File preview/validation before upload
- Progress indicator during processing
- Collapsible filters section to reduce clutter

#### Results Section
- Tabbed interface for Overall/Levels/Players
- Expandable/collapsible cards for each statistic group
- Side-by-side view: tables on left, charts on right
- Print-friendly view option

### Component Enhancements

#### Statistics Cards
```html
<div class="stats-card">
  <div class="stats-card-header">
    <h3 class="stats-card-title">Overall Statistics</h3>
    <div class="stats-card-actions">
      <button class="icon-btn" data-action="toggle-chart">
        <svg><!-- chart icon --></svg>
      </button>
      <button class="icon-btn" data-action="export">
        <svg><!-- download icon --></svg>
      </button>
    </div>
  </div>

  <div class="stats-card-body">
    <div class="stats-table-container">
      <!-- Existing table -->
    </div>

    <div class="stats-chart-container" id="chart-overall">
      <canvas id="chart-overall-canvas"></canvas>
    </div>
  </div>
</div>
```

#### Improved Tables
- Alternating row colors for better readability
- Sticky headers for long tables
- Sort functionality (client-side)
- Highlight min/max values
- Color-coded values based on thresholds

## Chart Implementations

### Chart Types and Use Cases

#### 1. Score Distribution Histogram
**Location:** Overall Statistics section

**Purpose:** Show the distribution of all scores across all games

**Chart Configuration:**
```javascript
{
  type: 'bar',
  data: {
    labels: ['0-10', '11-20', '21-30', ...], // Score ranges
    datasets: [{
      label: 'Number of Games',
      data: [5, 12, 23, ...], // Count of games in each range
      backgroundColor: 'var(--chart-blue)',
      borderColor: 'var(--primary-color)',
      borderWidth: 1
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      title: {
        display: true,
        text: 'Score Distribution'
      },
      legend: {
        display: false
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Frequency'
        }
      },
      x: {
        title: {
          display: true,
          text: 'Score Range'
        }
      }
    }
  }
}
```

#### 2. Level Performance Bar Chart
**Location:** Level Statistics section

**Purpose:** Compare average scores across different levels

**Chart Configuration:**
```javascript
{
  type: 'bar',
  data: {
    labels: ['Level 1', 'Level 2', 'Level 3', ...],
    datasets: [
      {
        label: 'Average Score',
        data: [85.5, 78.2, 92.1, ...],
        backgroundColor: 'var(--chart-blue)',
        borderColor: 'var(--primary-color)',
        borderWidth: 1
      },
      {
        label: 'Median Score',
        data: [84.0, 77.0, 91.0, ...],
        backgroundColor: 'var(--chart-green)',
        borderColor: 'var(--success-color)',
        borderWidth: 1
      }
    ]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      title: {
        display: true,
        text: 'Performance by Level'
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Score'
        }
      }
    }
  }
}
```

#### 3. Player Performance Comparison
**Location:** Player Statistics section

**Purpose:** Compare players' average scores with error bars for standard deviation

**Chart Configuration:**
```javascript
{
  type: 'bar',
  data: {
    labels: ['Player 1', 'Player 2', 'Player 3', ...],
    datasets: [{
      label: 'Average Score',
      data: [88.5, 75.3, 92.1, ...],
      backgroundColor: players.map((_, i) => chartColors[i % chartColors.length]),
      borderWidth: 1
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    indexAxis: 'y', // Horizontal bars
    plugins: {
      title: {
        display: true,
        text: 'Player Performance Comparison'
      }
    },
    scales: {
      x: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Average Score'
        }
      }
    }
  }
}
```

#### 4. Box Plot for Score Distribution (Advanced)
**Location:** Overall Statistics section (optional detailed view)

**Purpose:** Show quartiles, median, and outliers

**Implementation:** Use Chart.js boxplot plugin
```bash
npm install @sgratzl/chartjs-chart-boxplot
```

#### 5. Trend Line Chart (Future Enhancement)
**Location:** New "Trends" section

**Purpose:** Show score progression over time (requires timestamp data)

**Note:** Would require adding timestamp field to GameScore struct

### Chart Data Preparation

#### Backend Changes (Go Templates)

**Strategy 1: Data Attributes (Recommended)**
```html
<!-- In overall-stats.html -->
<div class="stats-chart-container"
     id="chart-overall"
     data-chart-type="histogram"
     data-chart-data='{"bins": [0,10,20,30], "frequencies": [5,12,23,15]}'>
  <canvas id="chart-overall-canvas"></canvas>
</div>
```

**Strategy 2: Inline Script Tags**
```html
<script type="application/json" id="chart-overall-data">
{
  "scores": [85, 92, 78, ...],
  "labels": ["Game 1", "Game 2", ...],
  "statistics": {
    "mean": 85.5,
    "median": 84.0,
    "stdDev": 8.3
  }
}
</script>
```

#### Frontend Chart Initialization

**charts.js Structure:**
```javascript
// Chart management module
const ChartManager = {
  charts: new Map(),

  // Initialize all charts in a container
  initializeCharts(container = document) {
    const chartContainers = container.querySelectorAll('[data-chart-type]');
    chartContainers.forEach(container => {
      const chartType = container.dataset.chartType;
      const chartData = JSON.parse(container.dataset.chartData || '{}');

      this.createChart(container, chartType, chartData);
    });
  },

  // Create specific chart types
  createChart(container, type, data) {
    const canvas = container.querySelector('canvas');
    if (!canvas) return;

    const chartId = container.id;

    // Destroy existing chart if present
    if (this.charts.has(chartId)) {
      this.charts.get(chartId).destroy();
    }

    let config;
    switch(type) {
      case 'histogram':
        config = this.createHistogram(data);
        break;
      case 'bar':
        config = this.createBarChart(data);
        break;
      case 'boxplot':
        config = this.createBoxPlot(data);
        break;
      default:
        console.warn(`Unknown chart type: ${type}`);
        return;
    }

    const chart = new Chart(canvas, config);
    this.charts.set(chartId, chart);
  },

  // Create histogram configuration
  createHistogram(data) {
    return {
      type: 'bar',
      data: {
        labels: data.bins,
        datasets: [{
          label: 'Frequency',
          data: data.frequencies,
          backgroundColor: 'rgba(59, 130, 246, 0.8)',
          borderColor: 'rgb(59, 130, 246)',
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: 'Score Distribution'
          },
          legend: {
            display: false
          }
        }
      }
    };
  },

  // Create bar chart configuration
  createBarChart(data) {
    // Implementation
  },

  // Destroy all charts
  destroyAll() {
    this.charts.forEach(chart => chart.destroy());
    this.charts.clear();
  },

  // Export chart as image
  exportChart(chartId, filename = 'chart.png') {
    const chart = this.charts.get(chartId);
    if (!chart) return;

    const url = chart.toBase64Image();
    const link = document.createElement('a');
    link.download = filename;
    link.href = url;
    link.click();
  }
};

// Initialize charts when htmx loads new content
document.addEventListener('htmx:afterSwap', (event) => {
  ChartManager.initializeCharts(event.detail.target);
});

// Initialize charts on page load
document.addEventListener('DOMContentLoaded', () => {
  ChartManager.initializeCharts();
});
```

### Statistical Calculations for Charts

#### Histogram Bin Calculation
Add utility function in `charts.js`:
```javascript
function calculateHistogramBins(scores, binCount = 10) {
  const min = Math.min(...scores);
  const max = Math.max(...scores);
  const binWidth = (max - min) / binCount;

  const bins = [];
  const frequencies = new Array(binCount).fill(0);

  for (let i = 0; i < binCount; i++) {
    const binStart = min + (i * binWidth);
    const binEnd = binStart + binWidth;
    bins.push(`${Math.round(binStart)}-${Math.round(binEnd)}`);

    // Count scores in this bin
    frequencies[i] = scores.filter(score =>
      score >= binStart && (i === binCount - 1 ? score <= binEnd : score < binEnd)
    ).length;
  }

  return { bins, frequencies };
}
```

#### Standard Deviation Calculation
If not already available, add to Go backend:
```go
// In internal/stats/stats.go
func CalculateStandardDeviation(scores []int, mean float64) float64 {
    if len(scores) == 0 {
        return 0
    }

    var sumSquaredDiff float64
    for _, score := range scores {
        diff := float64(score) - mean
        sumSquaredDiff += diff * diff
    }

    variance := sumSquaredDiff / float64(len(scores))
    return math.Sqrt(variance)
}
```

Update Statistics struct:
```go
type Statistics struct {
    // ... existing fields ...
    StandardDeviation float64 `json:"standardDeviation"`
}
```

## Backend Changes

### Template Data Structure Updates

**internal/handlers/stats.go** - Update template data:
```go
type TemplateData struct {
    Overall      OverallStatsData
    LevelStats   []LevelStatsData
    PlayerStats  []PlayerStatsData
    ShowDetailed bool
    ChartData    ChartData // New
}

type ChartData struct {
    ScoreDistribution HistogramData       `json:"scoreDistribution"`
    LevelPerformance  BarChartData        `json:"levelPerformance"`
    PlayerComparison  BarChartData        `json:"playerComparison"`
    AllScores         []int               `json:"allScores"` // For client-side calculations
}

type HistogramData struct {
    Bins        []string `json:"bins"`
    Frequencies []int    `json:"frequencies"`
}

type BarChartData struct {
    Labels  []string  `json:"labels"`
    Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
    Label           string   `json:"label"`
    Data            []float64 `json:"data"`
    BackgroundColor string   `json:"backgroundColor"`
}
```

### Helper Functions

Add chart data preparation functions:
```go
// In internal/handlers/charts.go (new file)
package handlers

import (
    "cli-stat-creator/internal/stats"
    "math"
)

// PrepareHistogramData creates histogram data from scores
func PrepareHistogramData(scores []stats.GameScore, binCount int) HistogramData {
    if len(scores) == 0 {
        return HistogramData{}
    }

    // Extract score values
    values := make([]int, len(scores))
    for i, gs := range scores {
        values[i] = gs.Score
    }

    // Find min and max
    min, max := values[0], values[0]
    for _, v := range values {
        if v < min { min = v }
        if v > max { max = v }
    }

    // Calculate bins
    binWidth := float64(max - min) / float64(binCount)
    bins := make([]string, binCount)
    frequencies := make([]int, binCount)

    for i := 0; i < binCount; i++ {
        binStart := float64(min) + float64(i) * binWidth
        binEnd := binStart + binWidth
        bins[i] = fmt.Sprintf("%.0f-%.0f", binStart, binEnd)

        // Count scores in bin
        for _, score := range values {
            if float64(score) >= binStart &&
               (i == binCount-1 || float64(score) < binEnd) {
                frequencies[i]++
            }
        }
    }

    return HistogramData{
        Bins: bins,
        Frequencies: frequencies,
    }
}

// PrepareLevelChartData creates bar chart data for levels
func PrepareLevelChartData(levelStats map[int]stats.Statistics) BarChartData {
    // Implementation
}

// PreparePlayerChartData creates bar chart data for players
func PreparePlayerChartData(playerStats map[stats.Player]stats.Statistics) BarChartData {
    // Implementation
}
```

## Responsive Design

### Breakpoints
```css
:root {
  --breakpoint-sm: 640px;   /* Mobile landscape */
  --breakpoint-md: 768px;   /* Tablet */
  --breakpoint-lg: 1024px;  /* Desktop */
  --breakpoint-xl: 1280px;  /* Large desktop */
}
```

### Mobile Optimizations
- Stack charts below tables on mobile
- Collapsible sections by default
- Simplified chart types on small screens
- Touch-friendly controls (44px minimum)
- Horizontal scroll for wide tables

### Tablet Layout
- Side-by-side table and chart in landscape
- Tabs for different statistic sections
- Larger touch targets

### Desktop Layout
- Two-column layout: tables left, charts right
- Multiple charts visible simultaneously
- Hover interactions for additional detail

## Accessibility

### ARIA Labels and Roles
```html
<div class="stats-chart-container"
     role="img"
     aria-label="Bar chart showing score distribution">
  <canvas id="chart-overall-canvas"></canvas>
</div>
```

### Keyboard Navigation
- Tab through interactive elements
- Enter/Space to toggle chart visibility
- Arrow keys for chart data point navigation (Chart.js built-in)

### Screen Reader Support
- Provide text alternative for charts
- Announce dynamic content changes
- Descriptive link text

### Color Contrast
- Ensure WCAG AA compliance (4.5:1 for normal text)
- Use patterns in addition to colors for charts
- Test with color blindness simulators

## Performance Considerations

### Chart Rendering
- Lazy load Chart.js library (defer attribute)
- Initialize charts only when visible (Intersection Observer)
- Throttle chart updates during window resize
- Destroy charts before creating new ones

### Data Optimization
- Limit histogram bins based on data size
- Sample large datasets (>1000 points) for scatter plots
- Use appropriate canvas size (max 2x display size)

### Bundle Size
- Chart.js core only: ~180KB
- Use tree-shaking if using module bundler
- Cache Chart.js CDN with long expiry

## Testing Strategy

### Visual Regression Testing
- Screenshot comparison for UI changes
- Test across different browsers
- Mobile and desktop viewports

### Chart Functionality Testing
- Verify correct data rendering
- Test interaction (hover, click, zoom)
- Export functionality
- Responsive behavior

### Accessibility Testing
- Keyboard navigation
- Screen reader compatibility
- Color contrast validation

### Browser Compatibility
- Chrome/Edge (latest 2 versions)
- Firefox (latest 2 versions)
- Safari (latest 2 versions)
- Mobile browsers (iOS Safari, Chrome Mobile)

## Implementation Phases

### Phase 1: UI/UX Foundation
**Goal:** Modernize existing interface without charts

**Tasks:**
1. Create CSS custom properties system
2. Enhance typography and spacing
3. Improve color scheme and visual hierarchy
4. Update card components with better styling
5. Add responsive breakpoints
6. Implement collapsible sections
7. Add loading states and animations
8. Test across browsers and devices

**Deliverables:**
- Enhanced `styles.css` with design system
- Updated HTML templates with better structure
- Improved mobile responsiveness

### Phase 2: Chart Infrastructure
**Goal:** Set up Chart.js and basic charting capabilities

**Tasks:**
1. Add Chart.js to project (CDN link in index.html)
2. Create `charts.js` module with ChartManager
3. Add chart containers to templates
4. Implement htmx integration for dynamic charts
5. Add chart initialization on page load and htmx swap
6. Create basic histogram chart
7. Test chart rendering and responsiveness

**Deliverables:**
- `web/static/charts.js` with chart management
- Updated templates with chart containers
- Working score distribution histogram

### Phase 3: Statistical Charts
**Goal:** Implement all chart types with real data

**Tasks:**
1. Add histogram bin calculation in `internal/handlers/charts.go`
2. Update Statistics struct with StandardDeviation field
3. Implement level performance bar chart
4. Implement player comparison bar chart
5. Add chart data to template context
6. Create chart configuration for each type
7. Test with various datasets
8. Add error handling for edge cases (no data, single data point)

**Deliverables:**
- `internal/handlers/charts.go` with data preparation
- All three main chart types implemented
- Updated templates with chart data attributes

### Phase 4: Advanced Features
**Goal:** Add polish and advanced functionality

**Tasks:**
1. Implement chart export (download as PNG)
2. Add chart visibility toggles
3. Implement chart tooltips with detailed info
4. Add animation and transitions
5. Create print-friendly view
6. Optimize chart performance
7. Add client-side chart filtering
8. Implement dark mode support (optional)

**Deliverables:**
- Export functionality
- Enhanced user interactions
- Performance optimizations

### Phase 5: Documentation and Polish
**Goal:** Document changes and ensure quality

**Tasks:**
1. Update CLAUDE.md with new components
2. Update README.md with screenshots
3. Add inline code documentation
4. Write user guide for chart features
5. Accessibility audit and fixes
6. Cross-browser testing
7. Performance testing
8. User acceptance testing

**Deliverables:**
- Updated documentation
- Accessibility report
- Browser compatibility matrix
- Performance benchmarks

## Future Enhancements

### Additional Chart Types
- Scatter plot for score vs. level correlation
- Radar chart for multi-dimensional player performance
- Heat map for level difficulty by player
- Timeline/trend charts (requires timestamp data)

### Advanced Features
- Chart comparison mode (overlay multiple datasets)
- Custom chart color themes
- Chart annotations and markers
- Zoom and pan for large datasets
- 3D visualizations (using Chart.js 3D plugin)

### Data Export
- Export charts as SVG
- Export complete report as PDF
- Export data as CSV
- Share chart links

### Interactivity
- Click chart elements to filter data
- Drag to select score ranges
- Real-time chart updates
- Chart legend filtering

### Analytics
- Track which charts users interact with most
- A/B test chart types
- Gather feedback on chart usefulness

## Migration Strategy

### Backward Compatibility
- All changes are additive
- Existing functionality remains unchanged
- JavaScript gracefully degrades without Chart.js
- CSS enhancements use progressive enhancement

### Rollout Plan
1. Deploy Phase 1 (UI/UX) first for user feedback
2. Beta test Phase 2-3 (charts) with small user group
3. Gather feedback and iterate
4. Full rollout of all phases
5. Monitor performance and user engagement

## Success Metrics

### User Experience
- Reduced time to understand statistics (measured via user testing)
- Increased user satisfaction (surveys)
- Reduced support requests about interpreting data

### Technical
- Page load time < 2 seconds
- Chart render time < 500ms
- Lighthouse score > 90
- Zero accessibility violations

### Engagement
- Users spend more time exploring data
- Increased usage of detailed statistics
- More files uploaded per session

## Summary

This design enhances the cli-stat-creator web interface with modern UI/UX improvements and interactive statistical visualizations. By using Chart.js, we add powerful data visualization capabilities while keeping the bundle size manageable and maintaining the existing htmx-based architecture. The phased implementation approach allows for iterative development and user feedback, ensuring the enhancements truly improve the user experience.

The combination of better visual design, clearer hierarchy, and interactive charts will make statistical analysis more accessible and insightful for users, transforming raw numbers into meaningful visual insights.
