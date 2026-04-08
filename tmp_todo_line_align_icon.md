TODO: align line x-endpoints with icon adjusted positions

## Problem
Icons shift right when overlapping (overlap avoidance logic in main.py).
Lines are drawn at the true data x (Month), so the line endpoint and the
icon center are misaligned horizontally. Vertical (Percent) stays correct.

## Goal
The line should visually connect to the icon center, appearing behind the icon
(icon rendered on top of the line endpoint).

## Approach: two passes

### Pass 1: compute all adjusted icon positions (no drawing)
Run the existing month loop as-is, including overlap detection and
`previous_positions` logic, but skip all draw calls (add_artist, plt.plot lines/dots).
Store results in:
    icon_adjusted_x: Dict[Tuple[str, pd.Timestamp], pd.Timestamp] = {}
    key: (deck, month), value: adjusted x position

### Pass 2: draw icons then lines
1. Draw icons using stored positions from pass 1.
2. For each deck that qualifies for a line, build a cloned x-series:
       adjusted_months = [icon_adjusted_x.get((deck, m), m) for m in deck_data["Month"]]
   Draw the line with adjusted_months (y unchanged, original deck_data unmodified).
3. Since lines are drawn after icons, line endpoints appear on top at the connection
   point, but the icon body covers the line behind it — giving the visual of the
   line connecting into the icon center.

## Key constraint
Original deck_data is never modified. Adjusted positions are stored as a separate
clone, so all filtering/sorting logic on actual data is unaffected.

## Location
visualization/main.py — restructure the main plotting loop.
