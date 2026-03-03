# Quickstart: Responsive Results View

## Goal

Validate the results view layout in normal, wide, and narrow terminals.

## Steps

1. Run the app and complete a typing test to reach the results screen.
2. Use a normal terminal width and confirm that all results are stacked in a single column.
3. Widen the terminal until the WPM chart appears in a separate right-side column.
4. Shrink the terminal so the chart no longer fits; confirm the chart is hidden.
5. Reduce terminal height to a very small size and confirm only essential stats display.

## Expected Results

- Single-column layout is the default.
- WPM chart moves to the right only when the terminal is wide enough.
- WPM chart hides entirely when the split layout cannot fit.
- Layout updates immediately on resize without visual corruption.
